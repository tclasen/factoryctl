// Package templates renders bounded, deterministic document and prompt templates.
package templates

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"text/template/parse"
	"time"

	"github.com/tclasen/factoryctl/internal/knowledge"
	"github.com/tclasen/factoryctl/internal/safeio"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

const MaxPrompt = 128 << 10
const MaxIncludeDepth = 16

// Payload extracts exactly one matching fenced payload under a Template heading.
func Payload(d *knowledge.Document) (string, error) {
	kind := "factory-template"
	if d.Type == "Prompt" {
		kind = "factory-prompt"
	} else if d.Type != "Record Template" {
		return "", fmt.Errorf("%s: expected Record Template or Prompt", d.Ref.Path)
	}
	source := []byte(d.Body)
	root := goldmark.DefaultParser().Parse(text.NewReader(source))
	count := 0
	payload := ""
	inSection := false
	sections := 0
	err := ast.Walk(root, func(n ast.Node, enter bool) (ast.WalkStatus, error) {
		if !enter {
			return ast.WalkContinue, nil
		}
		if h, ok := n.(*ast.Heading); ok && h.Level == 1 {
			inSection = string(h.Text(source)) == "Template"
			if inSection {
				sections++
			}
		}
		if fence, ok := n.(*ast.FencedCodeBlock); ok && string(fence.Language(source)) == kind {
			if !inSection {
				return ast.WalkStop, fmt.Errorf("%s: payload must be under # Template", d.Ref.Path)
			}
			count++
			var b strings.Builder
			for i := 0; i < fence.Lines().Len(); i++ {
				s := fence.Lines().At(i)
				b.Write(s.Value(source))
			}
			payload = b.String()
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		return "", err
	}
	asset := knowledge.String(d.Frontmatter, "asset")
	if count > 1 || count == 1 && sections != 1 || count == 1 && asset != "" || count == 0 && asset == "" {
		return "", fmt.Errorf("%s: expected exactly one payload source", d.Ref.Path)
	}
	return payload, nil
}

// Render never mutates files or starts a runtime. The caller validates record
// output against the proposed graph before storing it.
func Render(projectRoot string, docs map[string]*knowledge.Document, templatePath string, context map[string]any, inputs map[string]any) ([]byte, error) {
	if err := Validate(projectRoot, docs, templatePath); err != nil {
		return nil, err
	}
	d := docs[templatePath]
	if d == nil {
		return nil, fmt.Errorf("template not found: %s", templatePath)
	}
	values, err := validateInputs(knowledge.Mapping(d.Frontmatter["inputs"]), inputs)
	if err != nil {
		return nil, err
	}
	// JSON round-trip ensures context consists solely of data, never callable Go methods.
	raw, err := json.Marshal(context)
	if err != nil {
		return nil, fmt.Errorf("invalid context: %w", err)
	}
	ctx := map[string]any{}
	if err = json.Unmarshal(raw, &ctx); err != nil {
		return nil, err
	}
	if ctx == nil {
		ctx = map[string]any{}
	}
	ctx["Input"] = values
	limit := knowledge.MaxDocument
	if d.Type == "Prompt" {
		limit = MaxPrompt
	}
	r := renderer{root: projectRoot, docs: docs, kind: d.Type, limit: limit, active: map[string]bool{}}
	return r.render(templatePath, ctx, 0)
}

type renderer struct {
	root   string
	docs   map[string]*knowledge.Document
	kind   string
	limit  int
	active map[string]bool
	calls  int
}

func (r *renderer) render(name string, ctx any, depth int) ([]byte, error) {
	r.calls++
	if r.calls > 1024 {
		return nil, fmt.Errorf("template include budget exceeded")
	}
	if depth > MaxIncludeDepth {
		return nil, fmt.Errorf("include depth exceeds %d", MaxIncludeDepth)
	}
	if !safePath(name) {
		return nil, fmt.Errorf("invalid template path: %s", name)
	}
	if r.active[name] {
		return nil, fmt.Errorf("include cycle at %s", name)
	}
	r.active[name] = true
	defer delete(r.active, name)
	var source string
	if strings.HasPrefix(name, ".factory/") {
		b, e := readAsset(r.root, name)
		if e != nil {
			return nil, e
		}
		source = string(b)
	} else {
		d := r.docs[name]
		if d == nil {
			return nil, fmt.Errorf("include target not found: %s", name)
		}
		if d.Type != r.kind {
			return nil, fmt.Errorf("include type mismatch: %s", name)
		}
		p, e := Payload(d)
		if e != nil {
			return nil, e
		}
		source = p
		if asset := knowledge.String(d.Frontmatter, "asset"); asset != "" {
			b, e := readAsset(r.root, asset)
			if e != nil {
				return nil, e
			}
			source = string(b)
		}
	}
	funcs := template.FuncMap{
		// JSON is a YAML 1.2 subset and safely encodes values in every scalar position.
		"yaml": func(v any) (string, error) { b, e := json.Marshal(v); return string(b), e },
		"json": func(v any) (string, error) { b, e := json.Marshal(v); return string(b), e },
		"markdownLink": func(label, target string) (string, error) {
			u, e := url.Parse(target)
			if e != nil {
				return "", e
			}
			if u.Scheme != "" && u.Scheme != "http" && u.Scheme != "https" && u.Scheme != "mailto" {
				return "", fmt.Errorf("unsafe link scheme")
			}
			label = strings.NewReplacer("\\", "\\\\", "[", "\\[", "]", "\\]", "\n", " ").Replace(label)
			target = strings.NewReplacer(" ", "%20", "(", "%28", ")", "%29", "\n", "%0A", "\r", "%0D").Replace(target)
			return "[" + label + "](" + target + ")", nil
		},
		"include": func(target string, data any) (string, error) {
			b, e := r.render(target, data, depth+1)
			return string(b), e
		},
	}
	t, e := template.New(name).Option("missingkey=error").Funcs(funcs).Parse(source)
	if e != nil {
		return nil, fmt.Errorf("%s: %w", name, e)
	}
	for _, part := range t.Templates() {
		if part.Name() != name {
			return nil, fmt.Errorf("%s: use include for named partials; inline definitions unsupported", name)
		}
		if e = allowed(part.Tree.Root); e != nil {
			return nil, fmt.Errorf("%s: %w", name, e)
		}
	}
	b := &boundedBuffer{limit: r.limit}
	if e = t.Execute(b, ctx); e != nil {
		return nil, fmt.Errorf("%s: %w", name, e)
	}
	return b.Bytes(), nil
}
func safePath(s string) bool {
	return s != "" && s != "." && !strings.ContainsAny(s, "\\\x00%") && !strings.HasPrefix(s, "/") && path.Clean(s) == s && s != ".." && !strings.HasPrefix(s, "../")
}
func readAsset(root, name string) ([]byte, error) {
	if !safePath(name) || !strings.HasPrefix(name, ".factory/") {
		return nil, fmt.Errorf("asset must be inside .factory/: %s", name)
	}
	current := root
	for _, part := range strings.Split(name, "/") {
		current = filepath.Join(current, part)
		st, e := os.Lstat(current)
		if e != nil {
			return nil, e
		}
		if st.Mode()&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("symlink asset forbidden: %s", name)
		}
	}
	// OpenRoot supplies containment even if a path component changes during open.
	base, e := os.OpenRoot(root)
	if e != nil {
		return nil, e
	}
	defer base.Close()
	f, e := safeio.OpenRegularRoot(base, name, os.O_RDONLY)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	st, e := f.Stat()
	if e != nil {
		return nil, e
	}
	if !st.Mode().IsRegular() {
		return nil, fmt.Errorf("asset is not a regular file")
	}
	b, e := io.ReadAll(io.LimitReader(f, knowledge.MaxDocument+1))
	if len(b) > knowledge.MaxDocument {
		return nil, fmt.Errorf("asset exceeds 1 MiB")
	}
	return b, e
}

type boundedBuffer struct {
	bytes.Buffer
	limit int
}

func (b *boundedBuffer) Write(p []byte) (int, error) {
	if len(p) > b.limit-b.Len() {
		return 0, fmt.Errorf("template output exceeds %d bytes", b.limit)
	}
	return b.Buffer.Write(p)
}
func allowed(n parse.Node) error {
	if n == nil {
		return nil
	}
	switch v := n.(type) {
	case *parse.ListNode:
		if v == nil {
			return nil
		}
		for _, x := range v.Nodes {
			if e := allowed(x); e != nil {
				return e
			}
		}
	case *parse.ActionNode:
		return allowed(v.Pipe)
	case *parse.PipeNode:
		for _, x := range v.Cmds {
			if e := allowed(x); e != nil {
				return e
			}
		}
	case *parse.CommandNode:
		for _, x := range v.Args {
			if e := allowed(x); e != nil {
				return e
			}
		}
	case *parse.IfNode:
		for _, x := range []parse.Node{v.Pipe, v.List, v.ElseList} {
			if e := allowed(x); e != nil {
				return e
			}
		}
	case *parse.RangeNode:
		for _, x := range []parse.Node{v.Pipe, v.List, v.ElseList} {
			if e := allowed(x); e != nil {
				return e
			}
		}
	case *parse.IdentifierNode:
		switch v.Ident {
		case "yaml", "json", "markdownLink", "include":
		default:
			return fmt.Errorf("function %q is not permitted", v.Ident)
		}
	case *parse.TextNode, *parse.FieldNode, *parse.VariableNode, *parse.DotNode, *parse.StringNode, *parse.NumberNode, *parse.BoolNode, *parse.NilNode:
	default:
		return fmt.Errorf("template construct %v is not permitted", n.Type())
	}
	return nil
}
func validateInputs(schema map[string]any, inputs map[string]any) (map[string]any, error) {
	out := map[string]any{}
	for k := range inputs {
		if _, ok := schema[k]; !ok {
			return nil, fmt.Errorf("unknown input %q", k)
		}
	}
	for k, v := range schema {
		field := knowledge.Mapping(v)
		if field == nil {
			return nil, fmt.Errorf("invalid input definition %q", k)
		}
		value, exists := inputs[k]
		if !exists {
			value, exists = field["default"]
		}
		required, _ := field["required"].(bool)
		if !exists || value == nil {
			if required {
				return nil, fmt.Errorf("missing required input %q", k)
			}
			continue
		}
		if e := validateValue(value, field, required); e != nil {
			return nil, fmt.Errorf("input %s: %w", k, e)
		}
		out[k] = value
	}
	return out, nil
}
func validateValue(v any, f map[string]any, required bool) error {
	kind := knowledge.String(f, "kind")
	ok := false
	switch kind {
	case "string", "text", "reference", "enum", "date", "datetime":
		s, yes := v.(string)
		ok = yes && (!required || strings.TrimSpace(s) != "")
		if yes {
			switch kind {
			case "enum":
				ok = false
				for _, x := range knowledge.List(f["values"]) {
					if x == s {
						ok = true
					}
				}
			case "date":
				_, e := time.Parse("2006-01-02", s)
				ok = e == nil
			case "datetime":
				_, e := time.Parse(time.RFC3339, s)
				ok = e == nil
			}
		}
	case "boolean":
		_, ok = v.(bool)
	case "integer":
		switch x := v.(type) {
		case int, int64, uint64:
			ok = true
		case float64:
			ok = !math.IsNaN(x) && !math.IsInf(x, 0) && math.Trunc(x) == x
		}
	case "number":
		switch x := v.(type) {
		case int, int64, uint64:
			ok = true
		case float64:
			ok = !math.IsNaN(x) && !math.IsInf(x, 0)
		}
	case "mapping":
		_, ok = v.(map[string]any)
	case "string-list", "reference-list":
		a, yes := v.([]any)
		ok = yes
		if yes {
			for _, x := range a {
				if _, yes = x.(string); !yes {
					ok = false
				}
			}
			if min, yes := f["min_items"].(int); yes && len(a) < min {
				ok = false
			}
		}
	}
	if !ok {
		return fmt.Errorf("expected %s", kind)
	}
	return nil
}

// Validate checks every declared include, including branches that may not execute
// for a particular input. Dynamic include names are deliberately unsupported so
// a saved template cannot hide a broken dependency behind runtime control flow.
func Validate(projectRoot string, docs map[string]*knowledge.Document, templatePath string) error {
	return validateDependencies(projectRoot, docs, templatePath, nil)
}

// DependencyVersions captures the exact browser-independent template assets
// inspected by validation, including inactive include branches.
func DependencyVersions(projectRoot string, docs map[string]*knowledge.Document, templatePath string) (map[string]string, error) {
	versions := map[string]string{}
	err := validateDependencies(projectRoot, docs, templatePath, versions)
	return versions, err
}

func AssetVersion(projectRoot, name string) (string, error) {
	b, e := readAsset(projectRoot, name)
	if e != nil {
		return "", e
	}
	return knowledge.Hash(b), nil
}

func validateDependencies(projectRoot string, docs map[string]*knowledge.Document, templatePath string, versions map[string]string) error {
	root := docs[templatePath]
	if root == nil {
		return fmt.Errorf("template not found: %s", templatePath)
	}
	active, done := map[string]bool{}, map[string]bool{}
	calls := 0
	var visit func(string, int) error
	visit = func(name string, depth int) error {
		calls++
		if calls > 1024 {
			return fmt.Errorf("template include budget exceeded")
		}
		if depth > MaxIncludeDepth {
			return fmt.Errorf("include depth exceeds %d", MaxIncludeDepth)
		}
		if !safePath(name) {
			return fmt.Errorf("invalid template path: %s", name)
		}
		if active[name] {
			return fmt.Errorf("include cycle at %s", name)
		}
		if done[name] {
			return nil
		}
		active[name] = true
		defer delete(active, name)
		source := ""
		if strings.HasPrefix(name, ".factory/") {
			b, e := readAsset(projectRoot, name)
			if e != nil {
				return e
			}
			source = string(b)
			if versions != nil {
				versions[name] = knowledge.Hash(b)
			}
		} else {
			d := docs[name]
			if d == nil {
				return fmt.Errorf("include target not found: %s", name)
			}
			if d.Type != root.Type {
				return fmt.Errorf("include type mismatch: %s", name)
			}
			if e := validateInputSchema(d.Frontmatter["inputs"]); e != nil {
				return fmt.Errorf("%s: %w", name, e)
			}
			p, e := Payload(d)
			if e != nil {
				return e
			}
			source = p
			if asset := knowledge.String(d.Frontmatter, "asset"); asset != "" {
				b, e := readAsset(projectRoot, asset)
				if e != nil {
					return e
				}
				source = string(b)
				if versions != nil {
					versions[asset] = knowledge.Hash(b)
				}
			}
		}
		dummy := func(...any) string { return "" }
		t, e := template.New(name).Funcs(template.FuncMap{"yaml": dummy, "json": dummy, "markdownLink": dummy, "include": dummy}).Parse(source)
		if e != nil {
			return fmt.Errorf("%s: %w", name, e)
		}
		for _, part := range t.Templates() {
			if part.Name() != name {
				return fmt.Errorf("%s: inline definitions unsupported; use include", name)
			}
			if e := allowed(part.Tree.Root); e != nil {
				return fmt.Errorf("%s: %w", name, e)
			}
			inputSchema := knowledge.Mapping(root.Frontmatter["inputs"])
			if local := docs[name]; local != nil {
				inputSchema = knowledge.Mapping(local.Frontmatter["inputs"])
			}
			includes, e := staticIncludes(part.Tree.Root, inputSchema)
			if e != nil {
				return fmt.Errorf("%s: %w", name, e)
			}
			for _, target := range includes {
				if e := visit(target, depth+1); e != nil {
					return fmt.Errorf("%s includes %s: %w", name, target, e)
				}
			}
		}
		done[name] = true
		return nil
	}
	return visit(templatePath, 0)
}
func validateInputSchema(v any) error {
	schema, ok := v.(map[string]any)
	if !ok {
		return fmt.Errorf("inputs must be a mapping")
	}
	for name, raw := range schema {
		f, ok := raw.(map[string]any)
		if !ok {
			return fmt.Errorf("input %s must have a field definition", name)
		}
		kind := knowledge.String(f, "kind")
		switch kind {
		case "string", "text", "reference", "enum", "date", "datetime", "boolean", "integer", "number", "mapping", "string-list", "reference-list":
		default:
			return fmt.Errorf("input %s has unsupported kind %q", name, kind)
		}
		required := false
		if v, exists := f["required"]; exists {
			var yes bool
			required, yes = v.(bool)
			if !yes {
				return fmt.Errorf("input %s required must be boolean", name)
			}
		}
		if v, exists := f["min_items"]; exists {
			n, ok := v.(int)
			if !ok || n < 0 || (kind != "string-list" && kind != "reference-list") {
				return fmt.Errorf("input %s min_items must be a nonnegative list bound", name)
			}
		}
		if kind == "enum" {
			a, ok := f["values"].([]any)
			if !ok || len(a) == 0 {
				return fmt.Errorf("input %s enum requires values", name)
			}
			seen := map[string]bool{}
			for _, v := range a {
				s, ok := v.(string)
				if !ok || seen[s] {
					return fmt.Errorf("input %s enum values must be unique strings", name)
				}
				seen[s] = true
			}
		}
		if def, exists := f["default"]; exists {
			if def == nil {
				if required {
					return fmt.Errorf("input %s required default cannot be null", name)
				}
			} else if e := validateValue(def, f, required); e != nil {
				return fmt.Errorf("input %s default: %w", name, e)
			}
		}
	}
	return nil
}
func staticIncludes(n parse.Node, inputs map[string]any) ([]string, error) {
	targets := []string{}
	var walk func(parse.Node) error
	walk = func(n parse.Node) error {
		if n == nil {
			return nil
		}
		switch v := n.(type) {
		case *parse.ListNode:
			if v == nil {
				return nil
			}
			for _, x := range v.Nodes {
				if e := walk(x); e != nil {
					return e
				}
			}
		case *parse.FieldNode:
			if len(v.Ident) > 1 && v.Ident[0] == "Input" {
				if _, ok := inputs[v.Ident[1]]; !ok {
					return fmt.Errorf("undeclared input %q", v.Ident[1])
				}
			}
		case *parse.ActionNode:
			return walk(v.Pipe)
		case *parse.PipeNode:
			for _, x := range v.Cmds {
				if e := walk(x); e != nil {
					return e
				}
			}
		case *parse.CommandNode:
			if len(v.Args) > 0 {
				if id, ok := v.Args[0].(*parse.IdentifierNode); ok && id.Ident == "include" {
					if len(v.Args) != 3 {
						return fmt.Errorf("include requires a literal path and context")
					}
					s, ok := v.Args[1].(*parse.StringNode)
					if !ok {
						return fmt.Errorf("dynamic include targets are unsupported")
					}
					targets = append(targets, s.Text)
				}
			}
			for _, x := range v.Args {
				if e := walk(x); e != nil {
					return e
				}
			}
		case *parse.IfNode:
			for _, x := range []parse.Node{v.Pipe, v.List, v.ElseList} {
				if e := walk(x); e != nil {
					return e
				}
			}
		case *parse.RangeNode:
			for _, x := range []parse.Node{v.Pipe, v.List, v.ElseList} {
				if e := walk(x); e != nil {
					return e
				}
			}
		}
		return nil
	}
	e := walk(n)
	return targets, e
}
