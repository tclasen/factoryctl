package app

import (
	"path"
	"sort"
	"strings"

	"github.com/tclasen/factoryctl/internal/knowledge"
	"github.com/tclasen/factoryctl/internal/project"
	"github.com/tclasen/factoryctl/internal/templates"
)

func actionApplies(b *Bundle, prompt *knowledge.Document, record *knowledge.Document) bool {
	if record == nil {
		return true
	}
	def := b.Definitions[record.Type]
	for _, v := range knowledge.List(def.Raw["actions"]) {
		if target, ok := v.(string); ok && resolveActionRef(def.Path, target) == prompt.Ref.Path {
			return true
		}
	}
	for _, v := range knowledge.List(prompt.Frontmatter["applies_to"]) {
		if v == "*" || v == record.Type {
			return true
		}
	}
	return false
}
func Actions(projects []project.Project, p project.Project, recordPath string) ([]map[string]any, error) {
	snap, err := Load(projects)
	if err != nil {
		return nil, err
	}
	b := snap.Bundle(p.Checkout)
	if b == nil {
		return nil, Fail("not_found", "selected bundle unavailable")
	}
	record := b.Docs[recordPath]
	if recordPath != "" && record == nil {
		return nil, Fail("not_found", "record unavailable")
	}
	result := []map[string]any{}
	for name, d := range b.Docs {
		if d.Type != "Prompt" || !actionApplies(b, d, record) {
			continue
		}
		label := knowledge.String(d.Frontmatter, "label")
		if label == "" {
			label = d.Title
		}
		reason := ""
		if record != nil {
			def, known := b.Definitions[record.Type]
			if known && (!def.Enabled || def.Invalid || def.ExtensionInvalid) {
				reason = "record type is disabled or invalid"
			}
		}
		if err := templates.Validate(p.Root, b.Docs, name); err != nil {
			reason = err.Error()
		}
		if profile := knowledge.String(d.Frontmatter, "profile"); profile != "" {
			target := resolveActionRef(name, profile)
			if b.Docs[target] == nil || b.Docs[target].Type != "Agent Profile" {
				reason = "profile is unavailable"
			}
		}
		for _, dg := range snap.Diagnostics {
			if dg.Path == name && dg.Checkout == p.Checkout && dg.Severity == "error" {
				reason = dg.Message
				break
			}
		}
		result = append(result, map[string]any{"path": name, "title": d.Title, "label": label, "inputs": d.Frontmatter["inputs"], "profile": d.Frontmatter["profile"], "disabled_reason": reason})
	}
	sort.Slice(result, func(i, j int) bool { return result[i]["path"].(string) < result[j]["path"].(string) })
	return result, nil
}
func resolveActionRef(from, target string) string {
	if strings.HasPrefix(target, "/") {
		return path.Clean(strings.TrimPrefix(target, "/"))
	}
	return path.Clean(path.Join(path.Dir(from), target))
}

// RenderPrompt is a preview only. No adapter, shell or conversation is started.
func RenderPrompt(projects []project.Project, p project.Project, promptPath, recordPath string, inputs map[string]any) (map[string]any, error) {
	actions, err := Actions(projects, p, recordPath)
	if err != nil {
		return nil, err
	}
	found := false
	for _, a := range actions {
		if a["path"] == promptPath {
			found = true
			if a["disabled_reason"] != "" {
				return nil, Fail("invalid_document", a["disabled_reason"].(string))
			}
		}
	}
	if !found {
		return nil, Fail("not_found", "applicable prompt unavailable")
	}
	snap, err := Load(projects)
	if err != nil {
		return nil, err
	}
	b := snap.Bundle(p.Checkout)
	prompt := b.Docs[promptPath]
	if prompt == nil || prompt.Type != "Prompt" {
		return nil, Fail("not_found", "prompt unavailable")
	}
	record := b.Docs[recordPath]
	if !actionApplies(b, prompt, record) {
		return nil, Fail("invalid_arguments", "prompt does not apply to record")
	}
	ctx := map[string]any{"Record": nil, "Project": map[string]any{"ID": p.ID, "Name": p.Name}, "Selection": []string{}, "Links": []any{}}
	if record != nil {
		def, known := b.Definitions[record.Type]
		if known && (def.Invalid || def.ExtensionInvalid || !def.Enabled) {
			return nil, Fail("schema_conflict", "record type is invalid or disabled")
		}
		ctx["Record"] = map[string]any{"Path": record.Ref.Path, "Title": record.Title, "Type": record.Type, "Revision": record.Revision, "Frontmatter": record.Frontmatter, "Body": record.Body}
		ctx["Selection"] = []string{recordPath}
		links := []Link{}
		for _, l := range snap.Links {
			if l.Source.Checkout == p.Checkout && l.Source.Path == recordPath || l.Resolved != nil && l.Resolved.Checkout == p.Checkout && l.Resolved.Path == recordPath {
				links = append(links, l)
			}
		}
		ctx["Links"] = links
	}
	profileBody := ""
	if profile := knowledge.String(prompt.Frontmatter, "profile"); profile != "" {
		node := b.Docs[resolveActionRef(promptPath, profile)]
		if node == nil || node.Type != "Agent Profile" {
			return nil, Fail("invalid_document", "profile unavailable")
		}
		profileBody = node.Body
	}
	rendered, err := templates.Render(p.Root, b.Docs, promptPath, ctx, inputs)
	if err != nil {
		return nil, Fail("invalid_document", err.Error())
	}
	return map[string]any{"prompt": string(rendered), "profile_instructions": profileBody, "context": ctx, "prompt_path": promptPath, "prompt_revision": prompt.Revision, "record_path": recordPath, "launched": false}, nil
}
