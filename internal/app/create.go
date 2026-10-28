package app

import (
	"bytes"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/tclasen/factoryctl/internal/defaults"
	"github.com/tclasen/factoryctl/internal/knowledge"
	"github.com/tclasen/factoryctl/internal/project"
	"github.com/tclasen/factoryctl/internal/safeio"
	"github.com/tclasen/factoryctl/internal/templates"
)

func Create(projects []project.Project, p project.Project, typ, slug string, inputs map[string]any) (*knowledge.Document, error) {
	return CreateChecked(projects, p, typ, slug, inputs, "")
}

func CreateChecked(projects []project.Project, p project.Project, typ, slug string, inputs map[string]any, expectedDefinition string) (*knowledge.Document, error) {
	if !slugPattern.MatchString(slug) {
		return nil, Fail("invalid_path", "slug must use lowercase kebab-case")
	}
	snap, e := Load(projects)
	if e != nil {
		return nil, e
	}
	b := snap.Bundle(p.Checkout)
	if b == nil {
		return nil, Fail("not_found", "selected bundle is unavailable")
	}
	def, known := b.Definitions[typ]
	identity := DefinitionHash(def)
	if expectedDefinition != "" && expectedDefinition != identity {
		return nil, Fail("stale_definition", "effective definition changed; refresh")
	}
	if known && (def.Invalid || !def.Enabled) {
		return nil, Fail("schema_conflict", "type definition is invalid, duplicated or disabled")
	}
	directory := "records"
	if known {
		directory = def.Directory
	}
	dest := path.Join(directory, slug+".md")
	if def.Filename != "" {
		dest = path.Join(directory, def.Filename)
	}
	var raw []byte
	var assets map[string]string
	if templatePath := knowledge.String(def.Raw, "template"); known && templatePath != "" {
		tp := path.Clean(path.Join(path.Dir(def.Path), templatePath))
		if b.Docs[def.Path] != nil && b.Docs[tp] == nil {
			return nil, Fail("invalid_document", "explicit local template reference is missing: "+tp)
		}
		docs := map[string]*knowledge.Document{}
		for key, data := range defaults.Files(true, p.ID, p.Name) {
			docs[key] = knowledge.Parse(key, data)
		}
		for key, d := range b.Docs {
			docs[key] = d
		}
		context := map[string]any{"Project": map[string]any{"ID": p.ID, "Name": p.Name}, "Selection": []string{}, "Links": []any{}}
		assets, e = templates.DependencyVersions(p.Root, docs, tp)
		if e != nil {
			return nil, Fail("invalid_document", e.Error())
		}
		raw, e = templates.Render(p.Root, docs, tp, context, inputs)
		if e == nil {
			after, err := templates.DependencyVersions(p.Root, docs, tp)
			if err != nil || !reflect.DeepEqual(assets, after) {
				return nil, Fail("stale_revision", "template dependencies changed during rendering; refresh")
			}
		}
	} else {
		m := map[string]any{"type": typ, "status": "draft", "factory": map[string]any{"format": 1}}
		allowed := map[string]bool{"title": true, "summary": true}
		for field, spec := range def.Fields {
			allowed[field] = true
			if value, ok := knowledge.Mapping(spec)["default"]; ok {
				m[field] = value
			}
		}
		for k, v := range inputs {
			if !allowed[k] {
				return nil, Fail("invalid_document", "unknown input: "+k)
			}
			if k != "summary" {
				m[k] = v
			}
		}
		if strings.TrimSpace(knowledge.String(m, "title")) == "" {
			return nil, Fail("invalid_document", "missing required input: title")
		}
		body := "\n# " + knowledge.String(m, "title") + "\n\n" + knowledge.String(inputs, "summary") + "\n"
		raw, e = knowledge.Encode(m, body)
	}
	if e != nil {
		return nil, Fail("invalid_document", e.Error())
	}
	return Mutate(projects, p, Mutation{Kind: "create", Path: dest, Document: raw, ExpectedDefinition: identity, DefinitionType: typ, ExpectedSnapshot: snap.Revision, ExpectedAssets: assets})
}

type InitResult struct {
	Created     []string `json:"created"`
	Unchanged   []string `json:"unchanged"`
	ProjectRoot string   `json:"project_root"`
}

func Init(root, id, starter string) (InitResult, error) {
	result := InitResult{ProjectRoot: root, Created: []string{}, Unchanged: []string{}}
	if !slugPattern.MatchString(id) {
		return result, Fail("invalid_document", "project ID must use lowercase kebab-case")
	}
	if starter != "" && starter != "software" {
		return result, Fail("invalid_document", "only the software starter is supported")
	}
	canonical, e := filepath.EvalSymlinks(root)
	if e != nil {
		return result, e
	}
	canonical, e = filepath.Abs(canonical)
	if e != nil {
		return result, e
	}
	root = canonical
	repo, e := project.FindRepository(root)
	if e != nil || repo == "" {
		return result, Fail("ambiguous_project", "init requires an explicitly selected directory within a Git working tree")
	}
	result.ProjectRoot = root
	files := defaults.Files(starter == "software", id, filepath.Base(root))
	p := project.Project{ID: id, Name: filepath.Base(root), Root: root, RepositoryRoot: repo, Checkout: root}
	b := &Bundle{Project: p, Docs: map[string]*knowledge.Document{}}
	bundlePath := filepath.Join(root, "factory")
	if st, e := os.Lstat(bundlePath); e == nil {
		if st.Mode()&os.ModeSymlink != 0 || !st.IsDir() {
			return result, Fail("invalid_path", "factory must be a real directory")
		}
		existing, e := Load([]project.Project{p})
		if e != nil {
			return result, e
		}
		for name, d := range existing.Bundles[0].Docs {
			b.Docs[name] = d
		}
	}
	for name, data := range files {
		if e := SafePath(bundlePath, name, true); e != nil {
			return result, e
		}
		current, e := safeio.ReadFile(filepath.Join(bundlePath, filepath.FromSlash(name)), knowledge.MaxDocument)
		if e == nil {
			if !bytes.Equal(current, data) {
				return result, Fail("already_exists", "init would change existing file: "+name)
			}
			result.Unchanged = append(result.Unchanged, name)
		} else if !os.IsNotExist(e) {
			return result, e
		}
		b.Docs[name] = knowledge.Parse(name, data)
	}
	proposed := &Snapshot{Bundles: []*Bundle{b}}
	v := proposed.Validate("factory")
	if v.Errors > 0 {
		return result, &Error{Code: "invalid_document", Message: "initial bundle is invalid", Diagnostics: v.Diagnostics}
	}
	projectRoot, e := os.OpenRoot(root)
	if e != nil {
		return result, e
	}
	defer projectRoot.Close()
	if e = projectRoot.MkdirAll("factory", 0755); e != nil {
		return result, e
	}
	bundle, e := projectRoot.OpenRoot("factory")
	if e != nil {
		return result, e
	}
	defer bundle.Close()
	for _, name := range sortedPaths(b.Docs) {
		data, exists := files[name]
		if !exists {
			continue
		}
		if e := SafePath(bundlePath, name, true); e != nil {
			return result, e
		}
		if e = bundle.MkdirAll(path.Dir(name), 0755); e != nil {
			return result, e
		}
		f, e := bundle.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
		if os.IsExist(e) {
			current, e := safeio.ReadRoot(bundle, name, knowledge.MaxDocument)
			if e != nil || !bytes.Equal(current, data) {
				return result, Fail("already_exists", "concurrent init conflict at "+name)
			}
			continue
		}
		if e != nil {
			return result, e
		}
		_, e = f.Write(data)
		if e == nil {
			e = f.Sync()
		}
		ce := f.Close()
		result.Created = append(result.Created, name)
		if e == nil {
			e = ce
		}
		if e != nil {
			return result, e
		}
	}
	return result, nil
}
