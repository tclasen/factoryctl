package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/tclasen/factoryctl/internal/app"
	"github.com/tclasen/factoryctl/internal/defaults"
	"github.com/tclasen/factoryctl/internal/knowledge"
	"github.com/tclasen/factoryctl/internal/mcp"
	"github.com/tclasen/factoryctl/internal/project"
	rt "github.com/tclasen/factoryctl/internal/runtime"
	"github.com/tclasen/factoryctl/internal/web"
)

const version = "0.1.0-dev"
const help = `factoryctl — inspect and safely edit Markdown knowledge bundles

Global flags (before or after commands):
  --path DIR --project ID --checkout TOKEN --json --non-interactive

Commands:
  web [--listen 127.0.0.1:8080] [--runtime codex]  (documents and opt-in conversations)
  chat --runtime codex --prompt TEXT  (interactive terminal conversation)
  mcp  (scoped stdio document and question tools)
  projects list
  init --project-id ID [--starter software]
  records list [--type TYPE] [--query TEXT] [--limit N] [--cursor C]
  records show PATH
  records create --type TYPE --slug SLUG [--input KEY=VALUE ...]
  records update PATH --expected-revision HASH --set KEY=VALUE ...
  records replace PATH --expected-revision HASH --file FILE  (- reads stdin)
  records delete PATH --expected-revision HASH
  validate [--level okf|factory] [--require-complete-links]
  templates diff
  actions list --record PATH
  doctor

Document commands never prompt or start an agent. Conversations require explicit runtime selection.
Document browsing remains read-only; conversation agents may use their native tools.
Network link checks and asset export remain unsupported.
Example: factoryctl --path . records list --type Task --json
`

type options struct {
	flags map[string][]string
	args  []string
}

func parse(args []string) (options, error) {
	o := options{flags: map[string][]string{}}
	bools := map[string]bool{"json": true, "non-interactive": true, "help": true, "version": true, "require-complete-links": true, "check-external": true}
	allowed := map[string]bool{}
	for _, k := range strings.Fields("path project checkout project-id starter type query limit cursor input slug expected-revision set file level kind record listen runtime runtime-executable model reasoning-effort preferences prompt") {
		allowed[k] = true
	}
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "-h" {
			a = "--help"
		}
		if !strings.HasPrefix(a, "--") {
			o.args = append(o.args, a)
			continue
		}
		pair := strings.SplitN(strings.TrimPrefix(a, "--"), "=", 2)
		k := pair[0]
		if bools[k] {
			if len(pair) > 1 {
				return o, app.Fail("invalid_arguments", "boolean flag takes no value: "+k)
			}
			o.flags[k] = []string{"true"}
			continue
		}
		if !allowed[k] {
			return o, app.Fail("invalid_arguments", "unknown flag: "+k)
		}
		value := ""
		if len(pair) == 2 {
			value = pair[1]
		} else {
			i++
			if i >= len(args) {
				return o, app.Fail("invalid_arguments", "missing flag value: "+k)
			}
			value = args[i]
		}
		if len(o.flags[k]) > 0 && k != "input" && k != "set" {
			return o, app.Fail("invalid_arguments", "duplicate flag: "+k)
		}
		o.flags[k] = append(o.flags[k], value)
	}
	return o, nil
}
func (o options) get(k string) string {
	if a := o.flags[k]; len(a) > 0 {
		return a[0]
	}
	return ""
}
func values(items []string) (map[string]any, error) {
	m := map[string]any{}
	for _, item := range items {
		p := strings.SplitN(item, "=", 2)
		if len(p) != 2 || p[0] == "" {
			return nil, app.Fail("invalid_arguments", "expected KEY=VALUE")
		}
		if _, ok := m[p[0]]; ok {
			return nil, app.Fail("invalid_arguments", "duplicate input: "+p[0])
		}
		var v any
		if json.Unmarshal([]byte(p[1]), &v) != nil {
			v = p[1]
		}
		m[p[0]] = v
	}
	return m, nil
}
func main() { os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr)) }
func run(args []string, in io.Reader, out, errOut io.Writer) (exitCode int) {
	var outputError error
	defer func() {
		if outputError != nil {
			fmt.Fprintln(errOut, "io_error: cannot write command output:", outputError)
			exitCode = 4
		}
	}()

	out = &checkedWriter{writer: out, err: &outputError}
	o, e := parse(args)
	jsonMode := o.get("json") != ""
	for _, arg := range args {
		if arg == "--json" {
			jsonMode = true
		}
	}
	emit := func(v any) {
		if jsonMode {
			data, err := json.MarshalIndent(v, "", "  ")
			if err != nil {
				outputError = err
				data, _ = json.Marshal(map[string]any{"error": &app.Error{Code: "io_error", Message: "cannot encode command result"}})
			}
			if _, err = out.Write(append(data, '\n')); err != nil {
				outputError = err
			}
		} else {
			human(out, errOut, v)
		}
	}
	fail := func(e error) int {
		var ae *app.Error
		if !errors.As(e, &ae) {
			ae = &app.Error{Code: "io_error", Message: e.Error()}
		}
		if jsonMode {
			emit(map[string]any{"error": ae})
		} else {
			fmt.Fprintln(errOut, ae.Error())
			for _, d := range ae.Diagnostics {
				fmt.Fprintf(errOut, "%s:%d: %s: %s\n", d.Path, d.Line, d.Code, d.Message)
			}
		}
		switch ae.Code {
		case "invalid_document", "invalid_location", "broken_link", "referenced_node", "schema_conflict":
			return 1
		case "invalid_arguments", "ambiguous_project", "invalid_path", "not_found", "cursor_expired", "lossless_edit_unavailable":
			return 2
		case "stale_revision", "already_exists":
			return 3
		case "validation_incomplete":
			return 5
		default:
			return 4
		}
	}
	if e != nil {
		return fail(e)
	}
	if o.get("help") != "" || (len(o.args) == 0 && o.get("version") == "") {
		fmt.Fprint(out, help)
		return 0
	}
	if o.get("version") != "" {
		emit(version)
		return 0
	}
	if len(o.args) == 0 {
		return fail(app.Fail("invalid_arguments", "missing command"))
	}
	cmd := o.args[0]
	sub := ""
	if len(o.args) > 1 {
		sub = o.args[1]
	}
	allowedFlags := map[string]bool{}
	for _, k := range strings.Fields("path project checkout json non-interactive help version") {
		allowedFlags[k] = true
	}
	extra := ""
	expected := 2
	switch cmd {
	case "web":
		extra = "listen runtime runtime-executable model reasoning-effort preferences"
		expected = 1
	case "chat":
		extra = "runtime runtime-executable model reasoning-effort preferences prompt"
		expected = 1
	case "mcp":
		expected = 1
	case "init":
		extra = "project-id starter"
		expected = 1
	case "validate":
		extra = "level require-complete-links check-external"
		expected = 1
	case "doctor":
		expected = 1
	case "projects":
	case "templates":
		extra = "kind"
	case "actions":
		extra = "record input"
		if sub == "run" {
			expected = 3
		}
	case "records":
		switch sub {
		case "list":
			extra = "type query limit cursor"
		case "create":
			extra = "type slug input"
		case "show":
			expected = 3
		case "update":
			extra = "expected-revision set"
			expected = 3
		case "replace":
			extra = "expected-revision file"
			expected = 3
		case "delete":
			extra = "expected-revision"
			expected = 3
		}
	}
	for _, k := range strings.Fields(extra) {
		allowedFlags[k] = true
	}
	if cmd != "gui" {
		for k := range o.flags {
			if !allowedFlags[k] {
				return fail(app.Fail("invalid_arguments", "flag --"+k+" is not valid for this command"))
			}
		}
	}
	if cmd != "gui" && len(o.args) != expected {
		return fail(app.Fail("invalid_arguments", "unexpected or missing positional argument"))
	}
	if cmd == "gui" || (cmd == "actions" && sub == "run") || (cmd == "templates" && sub == "export") || o.get("check-external") != "" {
		return fail(app.Fail("unsupported_capability", "this capability is not implemented in the document CLI milestone"))
	}
	if cmd == "doctor" {
		emit(map[string]any{"version": version, "runtime": "codex-opt-in", "git_subprocesses": false, "network_required": false, "limits": map[string]int{"document_bytes": knowledge.MaxDocument, "frontmatter_bytes": knowledge.MaxFrontmatter, "yaml_depth": knowledge.MaxDepth}, "writer_coordination": "advisory per bundle; external editors and other bundles can race"})
		return 0
	}
	selectedPath := o.get("path")
	explicit := selectedPath != ""
	if selectedPath == "" {
		selectedPath, e = os.Getwd()
		if e != nil {
			return fail(e)
		}
	}
	collection, e := project.Discover(selectedPath, explicit)
	if e != nil {
		return fail(e)
	}
	if cmd == "projects" && sub == "list" {
		emit(collection)
		return 0
	}
	projects := collection.Projects
	filtered := []project.Project{}
	for _, p := range projects {
		if (o.get("project") == "" || p.ID == o.get("project")) && (o.get("checkout") == "" || p.Checkout == o.get("checkout")) {
			filtered = append(filtered, p)
		}
	}
	if (o.get("project") != "" || o.get("checkout") != "") && len(filtered) == 0 {
		return fail(app.Fail("ambiguous_project", "selection matches no discovered project"))
	}
	projects = filtered
	if cmd == "init" {
		if len(projects) > 1 {
			return fail(app.Fail("ambiguous_project", "select one project before init"))
		}
		if len(projects) == 1 && (!explicit || o.get("project") != "" || o.get("checkout") != "") {
			selectedPath = projects[0].Root
		}
		result, e := app.Init(selectedPath, o.get("project-id"), o.get("starter"))
		if e != nil {
			if len(result.Created) > 0 {
				fmt.Fprintf(errOut, "init created paths before failure: %v\n", result.Created)
			}
			return fail(e)
		}
		emit(result)
		return 0
	}
	if len(projects) == 0 {
		return fail(app.Fail("ambiguous_project", "no initialized project in this scope; use init --project-id ID"))
	}
	if cmd == "mcp" {
		if jsonMode {
			return fail(app.Fail("invalid_arguments", "MCP uses its own stdio protocol, not --json"))
		}
		if len(projects) != 1 {
			return fail(app.Fail("ambiguous_project", "MCP requires one explicit project"))
		}
		if e := mcp.Run(context.Background(), in, out, mcp.Config{Projects: projects, Selected: projects[0], BridgeURL: os.Getenv("FACTORYCTL_BRIDGE_URL"), BridgeToken: os.Getenv("FACTORYCTL_BRIDGE_TOKEN")}); e != nil {
			return fail(e)
		}
		return 0
	}
	if cmd == "chat" {
		if jsonMode || o.get("non-interactive") != "" {
			return fail(app.Fail("interaction_unavailable", "chat requires an interactive answer surface"))
		}
		if len(projects) != 1 {
			return fail(app.Fail("ambiguous_project", "chat requires one project"))
		}
		cfg, e := runtimeConfig(o)
		if e != nil {
			return fail(e)
		}
		if cfg == nil {
			return fail(app.Fail("invalid_arguments", "select --runtime codex"))
		}
		if e := web.ChatTerminal(projects, *cfg, o.get("prompt"), in, out); e != nil {
			return fail(e)
		}
		return 0
	}
	if cmd == "web" {
		if jsonMode {
			return fail(app.Fail("invalid_arguments", "web is a local browser server, not a JSON command"))
		}
		address := o.get("listen")
		if address == "" {
			address = "127.0.0.1:8080"
		}
		cfg, e := runtimeConfig(o)
		if e != nil {
			return fail(e)
		}
		if err := web.ServeWithRuntime(projects, address, out, cfg); err != nil {
			return fail(err)
		}
		return 0
	}
	snap, e := app.Load(collection.Projects)
	if e != nil {
		return fail(e)
	}
	one := func() (project.Project, error) {
		if len(projects) != 1 {
			return project.Project{}, app.Fail("ambiguous_project", "select exactly one project with --project and, for duplicates, --checkout")
		}
		return projects[0], nil
	}
	if cmd == "validate" {
		level := o.get("level")
		if level == "" {
			level = "factory"
		}
		if level != "factory" && level != "okf" {
			return fail(app.Fail("invalid_arguments", "level must be okf or factory"))
		}
		v := snap.Validate(level)
		selectedIDs := map[string]bool{}
		v.Scope = nil
		for _, p := range projects {
			selectedIDs[p.Checkout] = true
			v.Scope = append(v.Scope, p)
		}
		ds := []knowledge.Diagnostic{}
		v.Errors = 0
		v.Warnings = 0
		for _, d := range v.Diagnostics {
			if selectedIDs[d.Checkout] {
				ds = append(ds, d)
				if d.Severity == "error" {
					v.Errors++
				} else {
					v.Warnings++
				}
			}
		}
		v.Diagnostics = ds
		v.UncheckedLinks = 0
		v.Complete = true
		for _, l := range snap.Links {
			if selectedIDs[l.Source.Checkout] && l.Status == "unchecked" {
				v.UncheckedLinks++
				v.Complete = false
			}
		}
		emit(v)
		if v.Errors > 0 {
			return 1
		}
		if !v.Complete && o.get("require-complete-links") != "" {
			return 5
		}
		return 0
	}
	if cmd == "templates" && sub == "diff" {
		p, e := one()
		if e != nil {
			return fail(e)
		}
		b := snap.Bundle(p.Checkout)
		rows := []map[string]string{}
		for file, data := range defaults.Files(true, p.ID, p.Name) {
			state := "absent"
			if d := b.Docs[file]; d != nil {
				state = "different"
				if string(data) == string(d.Raw) {
					state = "equal"
				}
			}
			rows = append(rows, map[string]string{"path": file, "status": state})
		}
		sort.Slice(rows, func(i, j int) bool { return rows[i]["path"] < rows[j]["path"] })
		emit(rows)
		return 0
	}
	if cmd == "records" && sub == "list" {
		limit := 50
		if o.get("limit") != "" {
			limit, e = strconv.Atoi(o.get("limit"))
			if e != nil || limit < 1 || limit > 200 {
				return fail(app.Fail("invalid_arguments", "limit must be 1–200"))
			}
		}
		nodes := []*knowledge.Document{}
		names := map[string]string{}
		for _, b := range snap.Bundles {
			selected := false
			for _, p := range projects {
				if p.Checkout == b.Project.Checkout {
					selected = true
				}
			}
			if !selected {
				continue
			}
			names[b.Project.Checkout] = b.Project.Name
			for _, d := range b.Docs {
				if o.get("type") != "" && d.Type != o.get("type") {
					continue
				}
				if q := strings.ToLower(o.get("query")); q != "" && !strings.Contains(strings.ToLower(string(d.Raw)), q) {
					continue
				}
				nodes = append(nodes, d)
			}
		}
		sort.Slice(nodes, func(i, j int) bool {
			a, b := nodes[i], nodes[j]
			return names[a.Ref.Checkout]+"\x00"+a.Type+"\x00"+a.Title+"\x00"+a.Ref.Path+a.Ref.Checkout < names[b.Ref.Checkout]+"\x00"+b.Type+"\x00"+b.Title+"\x00"+b.Ref.Path+b.Ref.Checkout
		})
		offset := 0
		signature := knowledge.Hash([]byte(snap.Revision + o.get("type") + "\x00" + o.get("query") + "\x00" + o.get("project") + "\x00" + o.get("checkout")))
		if cursor := o.get("cursor"); cursor != "" {
			raw, e := base64.RawURLEncoding.DecodeString(cursor)
			parts := strings.Split(string(raw), ":")
			if e != nil || len(parts) != 2 || parts[0] != signature {
				return fail(app.Fail("cursor_expired", "snapshot or query changed; restart listing"))
			}
			offset, e = strconv.Atoi(parts[1])
			if e != nil || offset < 0 || offset > len(nodes) {
				return fail(app.Fail("cursor_expired", "invalid cursor; restart listing"))
			}
		}
		end := offset + limit
		if end > len(nodes) {
			end = len(nodes)
		}
		next := ""
		if end < len(nodes) {
			next = base64.RawURLEncoding.EncodeToString([]byte(signature + ":" + strconv.Itoa(end)))
		}
		emit(map[string]any{"nodes": nodes[offset:end], "next_cursor": next, "snapshot": snap.Revision, "diagnostics": snap.Diagnostics})
		return 0
	}
	p, e := one()
	if e != nil {
		return fail(e)
	}
	b := snap.Bundle(p.Checkout)
	if cmd == "actions" && sub == "list" {
		d := b.Docs[o.get("record")]
		if d == nil {
			return fail(app.Fail("not_found", "record not found"))
		}
		actions := []*knowledge.Document{}
		def, known := b.Definitions[d.Type]
		explicitActions := map[string]bool{}
		for _, raw := range knowledge.List(def.Raw["actions"]) {
			if target, ok := raw.(string); ok {
				explicitActions[path.Clean(path.Join(path.Dir(def.Path), target))] = true
			}
		}
		if !known || (def.Enabled && !def.Invalid) {
			for _, a := range b.Docs {
				if a.Type == "Prompt" {
					applicable := explicitActions[a.Ref.Path]
					for _, typ := range knowledge.List(a.Frontmatter["applies_to"]) {
						if typ == "*" || typ == d.Type {
							applicable = true
						}
					}
					if applicable {
						actions = append(actions, a)
					}
				}
			}
		}
		sort.Slice(actions, func(i, j int) bool { return actions[i].Ref.Path < actions[j].Ref.Path })
		emit(actions)
		return 0
	}
	if cmd != "records" {
		return fail(app.Fail("invalid_arguments", "unknown command"))
	}
	if sub == "create" {
		if len(o.args) != 2 || o.get("type") == "" || o.get("slug") == "" {
			return fail(app.Fail("invalid_arguments", "create requires --type and --slug"))
		}
		inputs, e := values(o.flags["input"])
		if e != nil {
			return fail(e)
		}
		d, e := app.Create(collection.Projects, p, o.get("type"), o.get("slug"), inputs)
		if e != nil {
			return fail(e)
		}
		emit(d)
		return 0
	}
	if len(o.args) != 3 {
		return fail(app.Fail("invalid_arguments", "record operation requires exactly one bundle-relative PATH"))
	}
	recordPath := o.args[2]
	if e = app.SafePath(filepath.Join(p.Root, "factory"), recordPath, false); e != nil {
		if os.IsNotExist(e) {
			return fail(app.Fail("not_found", "record not found"))
		}
		return fail(e)
	}
	d := b.Docs[recordPath]
	if d == nil {
		return fail(app.Fail("not_found", "record not found"))
	}
	if sub == "show" {
		incoming := []app.Link{}
		outgoing := []app.Link{}
		for _, l := range snap.Links {
			if l.Source == d.Ref {
				outgoing = append(outgoing, l)
			}
			if l.Resolved != nil && *l.Resolved == d.Ref {
				incoming = append(incoming, l)
			}
		}
		schema := "unavailable"
		if _, ok := b.Definitions[d.Type]; ok {
			schema = "available"
		}
		emit(map[string]any{"node": d, "document": string(d.Raw), "outgoing": outgoing, "backlinks": incoming, "schema": schema, "diagnostics": snap.Diagnostics})
		return 0
	}
	m := app.Mutation{Kind: sub, Path: recordPath, ExpectedRevision: o.get("expected-revision")}
	if m.ExpectedRevision == "" {
		return fail(app.Fail("invalid_arguments", "--expected-revision is required"))
	}
	switch sub {
	case "update":
		if len(o.flags["set"]) == 0 {
			return fail(app.Fail("invalid_arguments", "update requires --set"))
		}
		m.Set, e = values(o.flags["set"])
	case "replace":
		file := o.get("file")
		if file == "" {
			return fail(app.Fail("invalid_arguments", "replace requires --file"))
		}
		if file == "-" {
			m.Document, e = io.ReadAll(io.LimitReader(in, knowledge.MaxDocument+1))
		} else {
			m.Document, e = app.ReadBounded(file)
		}
	case "delete":
	default:
		return fail(app.Fail("invalid_arguments", "unknown record operation"))
	}
	if e != nil {
		return fail(e)
	}
	result, e := app.Mutate(collection.Projects, p, m)
	if e != nil {
		return fail(e)
	}
	emit(result)
	return 0
}

type checkedWriter struct {
	writer io.Writer
	err    *error
}

func (w *checkedWriter) Write(p []byte) (int, error) {
	n, e := w.writer.Write(p)
	if e == nil && n != len(p) {
		e = io.ErrShortWrite
	}
	if e != nil {
		*w.err = e
	}
	return n, e
}

func runtimeConfig(o options) (*web.RuntimeConfig, error) {
	if o.get("runtime") == "" && o.get("preferences") == "" {
		if o.get("runtime-executable") != "" || o.get("model") != "" || o.get("reasoning-effort") != "" {
			return nil, app.Fail("invalid_arguments", "runtime options require --runtime or --preferences")
		}
		return nil, nil
	}
	selected, e := rt.ResolvePreferences(rt.Selection{Runtime: o.get("runtime"), Executable: o.get("runtime-executable"), Model: o.get("model"), Effort: o.get("reasoning-effort"), Preferences: o.get("preferences")})
	if e != nil {
		return nil, app.Fail("invalid_arguments", e.Error())
	}
	return &web.RuntimeConfig{Executable: selected.Executable, Model: selected.Model, Effort: selected.Effort}, nil
}
