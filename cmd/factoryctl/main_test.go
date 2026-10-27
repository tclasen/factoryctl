package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type response struct {
	code           int
	stdout, stderr string
	value          map[string]any
}

func invoke(t *testing.T, stdin string, args ...string) response {
	t.Helper()
	var out, errOut bytes.Buffer
	code := run(args, strings.NewReader(stdin), &out, &errOut)
	r := response{code: code, stdout: out.String(), stderr: errOut.String()}
	if strings.HasPrefix(strings.TrimSpace(r.stdout), "{") {
		if e := json.Unmarshal(out.Bytes(), &r.value); e != nil {
			t.Fatalf("stdout is not one JSON object: %q: %v", r.stdout, e)
		}
	}
	return r
}
func wantCode(t *testing.T, r response, code int) {
	t.Helper()
	if r.code != code {
		t.Fatalf("exit=%d want=%d stdout=%s stderr=%s", r.code, code, r.stdout, r.stderr)
	}
}
func wantError(t *testing.T, r response, exit int, code string) {
	t.Helper()
	wantCode(t, r, exit)
	e, ok := r.value["error"].(map[string]any)
	if !ok || e["code"] != code {
		t.Fatalf("error code want %s: %s", code, r.stdout)
	}
	if r.stderr != "" {
		t.Fatalf("JSON error contaminated stderr: %q", r.stderr)
	}
}
func repo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, p := range []string{".git/objects", ".git/refs/heads"} {
		if e := os.MkdirAll(filepath.Join(root, p), 0755); e != nil {
			t.Fatal(e)
		}
	}
	if e := os.WriteFile(filepath.Join(root, ".git/HEAD"), []byte("ref: refs/heads/main\n"), 0644); e != nil {
		t.Fatal(e)
	}
	return root
}
func initialized(t *testing.T) string {
	t.Helper()
	root := repo(t)
	r := invoke(t, "", "--json", "init", "--path", root, "--project-id", "example", "--non-interactive")
	wantCode(t, r, 0)
	if r.stderr != "" {
		t.Fatal(r.stderr)
	}
	return root
}
func cli(t *testing.T, root, stdin string, args ...string) response {
	t.Helper()
	return invoke(t, stdin, append([]string{"--path", root, "--json", "--non-interactive"}, args...)...)
}
func createTask(t *testing.T, root, slug string) response {
	t.Helper()
	r := cli(t, root, "", "records", "create", "--type", "Task", "--slug", slug, "--input", "title="+slug, "--input", "summary=CLI integration test")
	wantCode(t, r, 0)
	if r.value["revision"] == "" || r.value["revision"] == nil {
		t.Fatalf("missing persisted revision: %s", r.stdout)
	}
	return r
}

func TestCLIInitDiscoveryAndValidation(t *testing.T) {
	root := initialized(t)
	r := cli(t, root, "", "projects", "list")
	wantCode(t, r, 0)
	p, ok := r.value["projects"].([]any)
	if !ok || len(p) != 1 {
		t.Fatalf("expected one project collection: %s", r.stdout)
	}
	r = cli(t, root, "", "validate")
	wantCode(t, r, 0)
	for _, k := range []string{"scope", "level", "complete", "errors", "warnings", "unchecked_links", "diagnostics"} {
		if _, ok := r.value[k]; !ok {
			t.Errorf("validate output missing %s: %s", k, r.stdout)
		}
	}
	r = cli(t, root, "", "init", "--project-id", "example")
	wantCode(t, r, 0)
	if len(r.value["created"].([]any)) != 0 {
		t.Fatalf("rerun rewrote defaults: %s", r.stdout)
	}
	if _, e := os.Stat(filepath.Join(root, ".factory")); !os.IsNotExist(e) {
		t.Fatal("init unexpectedly created .factory")
	}
	if _, e := os.Stat(filepath.Join(root, "factory.yaml")); !os.IsNotExist(e) {
		t.Fatal("init unexpectedly created configuration")
	}
	// Selected directory flags work after the subcommand as documented.
	r = invoke(t, "", "projects", "--json", "list", "--path", root)
	wantCode(t, r, 0)
}

func TestCLIRecordLifecycleAndRejectedWrites(t *testing.T) {
	root := initialized(t)
	r := createTask(t, root, "example")
	revision := r.value["revision"].(string)
	file := filepath.Join(root, "factory/tasks/example.md")
	r = cli(t, root, "", "records", "list", "--type", "Task", "--query", "CLI INTEGRATION")
	wantCode(t, r, 0)
	if len(r.value["nodes"].([]any)) != 1 {
		t.Fatalf("query/type filtering failed: %s", r.stdout)
	}
	r = cli(t, root, "", "records", "show", "tasks/example.md")
	wantCode(t, r, 0)
	original := r.value["document"].(string)
	if !strings.Contains(original, "CLI integration test") {
		t.Fatal("show lost body")
	}
	if r.value["schema"] != "available" {
		t.Fatal("missing schema")
	}
	r = cli(t, root, "", "records", "update", "tasks/example.md", "--expected-revision", revision, "--set", "task_state=done")
	wantCode(t, r, 0)
	updatedRevision := r.value["revision"].(string)
	if updatedRevision == revision {
		t.Fatal("update did not change revision")
	}
	wantError(t, cli(t, root, "", "records", "update", "tasks/example.md", "--expected-revision", revision, "--set", "task_state=planned"), 3, "stale_revision")
	before, e := os.ReadFile(file)
	if e != nil {
		t.Fatal(e)
	}
	wantError(t, cli(t, root, "", "records", "update", "tasks/example.md", "--expected-revision", updatedRevision, "--set", "task_state=bogus"), 1, "invalid_document")
	after, e := os.ReadFile(file)
	if e != nil {
		t.Fatal(e)
	}
	if !bytes.Equal(before, after) {
		t.Fatal("rejected edit changed bytes")
	}
	replacement := strings.Replace(string(after), "CLI integration test", "Replacement retained body", 1)
	r = cli(t, root, replacement, "records", "replace", "tasks/example.md", "--expected-revision", updatedRevision, "--file", "-")
	wantCode(t, r, 0)
	newRevision := r.value["revision"].(string)
	wantError(t, cli(t, root, "", "records", "create", "--type", "Task", "--slug", "example", "--input", "title=Duplicate", "--input", "summary=Duplicate"), 3, "already_exists")
	r = cli(t, root, "", "records", "delete", "tasks/example.md", "--expected-revision", newRevision)
	wantCode(t, r, 0)
	if _, e := os.Stat(file); !os.IsNotExist(e) {
		t.Fatal("delete did not remove selected file")
	}
	wantError(t, cli(t, root, "", "records", "show", "tasks/example.md"), 2, "not_found")
}

func TestCLIPaginationExpiresAfterExternalEdit(t *testing.T) {
	root := initialized(t)
	createTask(t, root, "one")
	createTask(t, root, "two")
	r := cli(t, root, "", "records", "list", "--type", "Task", "--limit", "1")
	wantCode(t, r, 0)
	cursor, _ := r.value["next_cursor"].(string)
	if cursor == "" {
		t.Fatal("missing next cursor")
	}
	if len(r.value["nodes"].([]any)) != 1 {
		t.Fatal("limit ignored")
	}
	r = cli(t, root, "", "records", "list", "--type", "Task", "--limit", "1", "--cursor", cursor)
	wantCode(t, r, 0)
	if len(r.value["nodes"].([]any)) != 1 {
		t.Fatal("second page absent")
	}
	f, e := os.OpenFile(filepath.Join(root, "factory/tasks/one.md"), os.O_APPEND|os.O_WRONLY, 0)
	if e != nil {
		t.Fatal(e)
	}
	_, e = f.WriteString("\nExternal editor change.\n")
	ce := f.Close()
	if e != nil || ce != nil {
		t.Fatalf("external edit: %v %v", e, ce)
	}
	wantError(t, cli(t, root, "", "records", "list", "--type", "Task", "--limit", "1", "--cursor", cursor), 2, "cursor_expired")
}

func TestCLIArgumentsAndUnsupportedCapabilities(t *testing.T) {
	root := initialized(t)
	wantError(t, cli(t, root, "", "records", "create", "--type", "Task", "--slug", "missing-input"), 1, "invalid_document")
	wantError(t, cli(t, root, "", "records", "list", "--limit", "0"), 2, "invalid_arguments")
	wantError(t, cli(t, root, "", "records", "show", "../factory.md"), 2, "invalid_path")
	wantError(t, cli(t, root, "", "records", "list", "--unknown"), 2, "invalid_arguments")
	for _, args := range [][]string{{"gui"}, {"mcp"}, {"actions", "run", "prompts/discuss.md"}, {"templates", "export", "--kind", "view"}, {"validate", "--check-external"}} {
		wantError(t, cli(t, root, "", args...), 4, "unsupported_capability")
	}
	r := invoke(t, "", "--path", root, "records", "list", "--limit", "0")
	wantCode(t, r, 2)
	if r.stdout != "" || !strings.Contains(r.stderr, "invalid_arguments") {
		t.Fatalf("human errors must use stderr: %#v", r)
	}
	r = invoke(t, "", "--help")
	wantCode(t, r, 0)
	if !strings.Contains(r.stdout, "records replace") || r.stderr != "" {
		t.Fatal("help omitted supported commands")
	}
	r = cli(t, root, "", "doctor")
	wantCode(t, r, 0)
	if r.value["git_subprocesses"] != false {
		t.Fatal("doctor must not claim Git subprocesses")
	}
}

func TestCLIJSONMalformedArgumentsRegardlessOfFlagOrder(t *testing.T) {
	// JSON mode must remain machine-readable when the malformed flag precedes it.
	r := invoke(t, "", "--unknown", "--json")
	wantError(t, r, 2, "invalid_arguments")
}

func TestCLICollectionSelectionProtectsSameNamedRecords(t *testing.T) {
	mono := repo(t)
	alpha, beta := filepath.Join(mono, "alpha"), filepath.Join(mono, "beta")
	for _, p := range []struct{ root, id string }{{alpha, "alpha"}, {beta, "beta"}} {
		if e := os.Mkdir(p.root, 0755); e != nil {
			t.Fatal(e)
		}
		wantCode(t, cli(t, p.root, "", "init", "--project-id", p.id), 0)
		createTask(t, p.root, "example")
	}
	betaFile := filepath.Join(beta, "factory/tasks/example.md")
	before, e := os.ReadFile(betaFile)
	if e != nil {
		t.Fatal(e)
	}
	r := cli(t, mono, "", "records", "list", "--type", "Task")
	wantCode(t, r, 0)
	if len(r.value["nodes"].([]any)) != 2 {
		t.Fatalf("collection missing records: %s", r.stdout)
	}
	wantError(t, cli(t, mono, "", "records", "show", "tasks/example.md"), 2, "ambiguous_project")
	r = cli(t, mono, "", "records", "show", "tasks/example.md", "--project", "alpha")
	wantCode(t, r, 0)
	revision := r.value["node"].(map[string]any)["revision"].(string)
	wantCode(t, cli(t, mono, "", "records", "update", "tasks/example.md", "--project", "alpha", "--expected-revision", revision, "--set", "task_state=done"), 0)
	after, e := os.ReadFile(betaFile)
	if e != nil {
		t.Fatal(e)
	}
	if !bytes.Equal(before, after) {
		t.Fatal("alpha edit changed beta same-named record")
	}
}

func TestCLIIncomingReferencesBlockDestructiveEdits(t *testing.T) {
	root := initialized(t)
	target := createTask(t, root, "target")
	createTask(t, root, "source")
	show := cli(t, root, "", "records", "show", "tasks/target.md")
	wantCode(t, show, 0)
	document := show.value["document"].(string) + "\n# Shared Heading\n"
	target = cli(t, root, document, "records", "replace", "tasks/target.md", "--expected-revision", target.value["revision"].(string), "--file", "-")
	wantCode(t, target, 0)
	show = cli(t, root, "", "records", "show", "tasks/source.md")
	wantCode(t, show, 0)
	source := show.value["document"].(string) + "\n[Target](target.md#shared-heading)\n"
	sourceRevision := show.value["node"].(map[string]any)["revision"].(string)
	wantCode(t, cli(t, root, source, "records", "replace", "tasks/source.md", "--expected-revision", sourceRevision, "--file", "-"), 0)
	revision := target.value["revision"].(string)
	wantError(t, cli(t, root, "", "records", "delete", "tasks/target.md", "--expected-revision", revision), 1, "referenced_node")
	wantError(t, cli(t, root, strings.Replace(document, "# Shared Heading", "# Removed Heading", 1), "records", "replace", "tasks/target.md", "--expected-revision", revision, "--file", "-"), 1, "invalid_document")
	after, e := os.ReadFile(filepath.Join(root, "factory/tasks/target.md"))
	if e != nil {
		t.Fatal(e)
	}
	if string(after) != document {
		t.Fatal("rejected referenced-target mutation changed bytes")
	}
}
