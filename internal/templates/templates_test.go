package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tclasen/factoryctl/internal/defaults"
	"github.com/tclasen/factoryctl/internal/knowledge"
)

func fixture(payload string) *knowledge.Document {
	return knowledge.Parse("templates/test.md", []byte("---\ntype: Record Template\nfor_type: Custom\ninputs:\n  title: {kind: string, required: true}\n---\n# Template\n```factory-template\n"+payload+"\n```\n"))
}
func TestEscapedInputsAndLiteralData(t *testing.T) {
	docs := map[string]*knowledge.Document{"templates/test.md": fixture("---\ntype: Custom\ntitle: {{ yaml .Input.title }}\n---\n{{ .Input.title }}")}
	value := "quotes: \"text\"\n{{ include \"evil\" . }}"
	b, e := Render(t.TempDir(), docs, "templates/test.md", nil, map[string]any{"title": value})
	if e != nil {
		t.Fatal(e)
	}
	d := knowledge.Parse("custom/test.md", b)
	if len(d.Diagnostics) > 0 || d.Title != value {
		t.Fatalf("did not preserve quoted data: %s %+v", b, d.Diagnostics)
	}
	for _, inputs := range []map[string]any{{}, {"title": "x", "extra": "bad"}, {"title": 5}, {"title": " "}} {
		if _, e := Render(t.TempDir(), docs, "templates/test.md", nil, inputs); e == nil {
			t.Fatalf("accepted invalid inputs: %#v", inputs)
		}
	}
}
func TestIncludesAndRestrictions(t *testing.T) {
	for _, tc := range []struct{ name, payload string }{
		{"cycle", `{{ include "templates/test.md" . }}`},
		{"escape", `{{ include "../outside.md" . }}`},
		{"builtin", `{{ printf "%s" .Input.title }}`},
		{"call", `{{ call .Input.title }}`},
		{"recursive", `{{define "loop"}}{{template "loop" .}}{{end}}{{template "loop" .}}`},
		{"missing", `{{ .Input.absent }}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			docs := map[string]*knowledge.Document{"templates/test.md": fixture(tc.payload)}
			if _, e := Render(t.TempDir(), docs, "templates/test.md", nil, map[string]any{"title": "x"}); e == nil {
				t.Fatal("accepted unsafe or invalid template")
			}
		})
	}
	root := fixture(`{{ include "templates/child.md" . }}`)
	child := fixture(`{{ yaml .Input.title }}`)
	docs := map[string]*knowledge.Document{"templates/test.md": root, "templates/child.md": child}
	b, e := Render(t.TempDir(), docs, "templates/test.md", nil, map[string]any{"title": "ok"})
	if e != nil || strings.TrimSpace(string(b)) != "\"ok\"" {
		t.Fatalf("include failed: %q %v", b, e)
	}
	child.Type = "Prompt"
	if _, e = Render(t.TempDir(), docs, "templates/test.md", nil, map[string]any{"title": "ok"}); e == nil {
		t.Fatal("accepted include type mismatch")
	}
}
func TestAssetAndPayloadSafety(t *testing.T) {
	root := t.TempDir()
	if e := os.MkdirAll(filepath.Join(root, ".factory/templates"), 0700); e != nil {
		t.Fatal(e)
	}
	outside := filepath.Join(t.TempDir(), "secret")
	if e := os.WriteFile(outside, []byte("secret"), 0600); e != nil {
		t.Fatal(e)
	}
	if e := os.Symlink(outside, filepath.Join(root, ".factory/templates/evil")); e != nil {
		t.Fatal(e)
	}
	docs := map[string]*knowledge.Document{"templates/test.md": fixture(`{{ include ".factory/templates/evil" . }}`)}
	if _, e := Render(root, docs, "templates/test.md", nil, map[string]any{"title": "x"}); e == nil {
		t.Fatal("followed symlink")
	}
	for _, d := range []*knowledge.Document{fixture("one\n```\n```factory-template\ntwo"), knowledge.Parse("x.md", []byte("---\ntype: Record Template\nasset: .factory/templates/x\n---\n# Template\n```factory-template\nx\n```")), knowledge.Parse("x.md", []byte("---\ntype: Record Template\n---\n```factory-template\nx\n```"))} {
		if _, e := Payload(d); e == nil {
			t.Fatal("accepted ambiguous/misplaced payload")
		}
	}
}
func TestOutputLimit(t *testing.T) {
	docs := map[string]*knowledge.Document{"templates/test.md": fixture(`{{ .Input.title }}{{ .Input.title }}`)}
	if _, e := Render(t.TempDir(), docs, "templates/test.md", nil, map[string]any{"title": strings.Repeat("x", knowledge.MaxDocument)}); e == nil {
		t.Fatal("unbounded output")
	}
}
func TestDefaultTemplateRenders(t *testing.T) {
	docs := map[string]*knowledge.Document{}
	for p, b := range defaults.Files(true, "demo", "Demo") {
		docs[p] = knowledge.Parse(p, b)
	}
	b, e := Render(t.TempDir(), docs, "templates/task.md", map[string]any{"Project": map[string]any{"ID": "demo", "Name": "Demo"}}, map[string]any{"title": "Try it", "summary": "Actual summary"})
	if e != nil {
		t.Fatal(e)
	}
	d := knowledge.Parse("tasks/try-it.md", b)
	if len(d.Diagnostics) > 0 || d.Type != "Task" {
		t.Fatalf("invalid rendered task: %s %+v", b, d.Diagnostics)
	}
}

func TestIncludeDepthAndPromptLimit(t *testing.T) {
	docs := map[string]*knowledge.Document{}
	for i := 0; i < 19; i++ {
		p := fmt.Sprintf("templates/p%d.md", i)
		docs[p] = fixture(fmt.Sprintf(`{{ include "templates/p%d.md" . }}`, i+1))
	}
	if _, err := Render(t.TempDir(), docs, "templates/p0.md", nil, map[string]any{"title": "ok"}); err == nil || !strings.Contains(err.Error(), "depth") {
		t.Fatalf("missing depth rejection: %v", err)
	}
	d := knowledge.Parse("prompts/test.md", []byte("---\ntype: Prompt\ninputs:\n  title: {kind: string, required: true}\n---\n# Template\n```factory-prompt\n{{ .Input.title }}\n```\n"))
	docs = map[string]*knowledge.Document{"prompts/test.md": d}
	if _, err := Render(t.TempDir(), docs, "prompts/test.md", nil, map[string]any{"title": strings.Repeat("x", MaxPrompt+1)}); err == nil {
		t.Fatal("prompt output exceeded limit")
	}
}

func TestStaticValidationFindsUnexecutedDependencies(t *testing.T) {
	for _, payload := range []string{`{{ if false }}{{ include "templates/absent.md" . }}{{ end }}`, `{{ if false }}{{ include "templates/test.md" . }}{{ end }}`, `{{ include .Input.title . }}`, `{{ if false }}{{ printf "hidden" }}{{ end }}`} {
		docs := map[string]*knowledge.Document{"templates/test.md": fixture(payload)}
		if err := Validate(t.TempDir(), docs, "templates/test.md"); err == nil {
			t.Fatalf("static validator missed dependency: %s", payload)
		}
	}
	root := fixture(`{{ if false }}{{ include "templates/child.md" . }}{{ end }}`)
	child := fixture(`{{ .Input.title }}`)
	docs := map[string]*knowledge.Document{"templates/test.md": root, "templates/child.md": child}
	if err := Validate(t.TempDir(), docs, "templates/test.md"); err != nil {
		t.Fatal(err)
	}
	child.Frontmatter["inputs"] = map[string]any{"bad": map[string]any{"kind": "integer", "default": "oops"}}
	if err := Validate(t.TempDir(), docs, "templates/test.md"); err == nil {
		t.Fatal("invalid included input default accepted")
	}
	child.Frontmatter["inputs"] = map[string]any{"bad": map[string]any{"kind": "string", "required": "true"}}
	if err := Validate(t.TempDir(), docs, "templates/test.md"); err == nil {
		t.Fatal("invalid required flag accepted")
	}
}
