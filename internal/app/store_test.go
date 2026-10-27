package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tclasen/factoryctl/internal/knowledge"
	"github.com/tclasen/factoryctl/internal/project"
)

func storeFixture(t *testing.T) (project.Project, []project.Project) {
	t.Helper()
	root := t.TempDir()
	for _, dir := range []string{".git/objects", ".git/refs"} {
		if e := os.MkdirAll(filepath.Join(root, dir), 0755); e != nil {
			t.Fatal(e)
		}
	}
	if e := os.WriteFile(filepath.Join(root, ".git/HEAD"), []byte("ref: refs/heads/main\n"), 0644); e != nil {
		t.Fatal(e)
	}
	if _, e := Init(root, "demo", "software"); e != nil {
		t.Fatal(e)
	}
	c, e := project.Discover(root, true)
	if e != nil || len(c.Projects) != 1 {
		t.Fatalf("discover: %v %#v", e, c)
	}
	return c.Projects[0], c.Projects
}
func TestProposedGraphAndUnrelatedRepairs(t *testing.T) {
	p, ps := storeFixture(t)
	create := func(name, body string) *knowledge.Document {
		d, e := Mutate(ps, p, Mutation{Kind: "create", Path: name, Document: []byte("---\ntype: Custom\ntitle: X\n---\n" + body)})
		if e != nil {
			t.Fatal(e)
		}
		return d
	}
	d := create("records/self.md", "# Present\n[self](#present)\n")
	old := append([]byte(nil), d.Raw...)
	_, e := Mutate(ps, p, Mutation{Kind: "replace", Path: d.Ref.Path, ExpectedRevision: d.Revision, Document: []byte(strings.Replace(string(d.Raw), "# Present", "# Gone", 1))})
	if e == nil {
		t.Fatal("removed own referenced heading")
	}
	raw, _ := os.ReadFile(filepath.Join(p.Root, "factory", d.Ref.Path))
	if string(raw) != string(old) {
		t.Fatal("rejected write changed bytes")
	}
	// Existing invalid external edits remain readable but do not prohibit healthy writes.
	if e = os.WriteFile(filepath.Join(p.Root, "factory/records/bad.md"), []byte("broken document"), 0644); e != nil {
		t.Fatal(e)
	}
	create("records/healthy.md", "healthy\n")
	snap, e := Load(ps)
	if e != nil {
		t.Fatal(e)
	}
	if snap.Validate("factory").Errors == 0 {
		t.Fatal("external error disappeared")
	}
}
func TestDefinitionChangesValidateIncomingInstances(t *testing.T) {
	p, ps := storeFixture(t)
	d, e := Create(ps, p, "Task", "one", map[string]any{"title": "One", "summary": "Body"})
	if e != nil {
		t.Fatal(e)
	}
	snap, _ := Load(ps)
	def := snap.Bundle(p.Checkout).Docs["definitions/task.md"]
	next := strings.Replace(string(def.Raw), "- planned", "- waiting", 1)
	if next == string(def.Raw) {
		t.Fatal("fixture enum replacement failed")
	}
	_, e = Mutate(ps, p, Mutation{Kind: "replace", Path: def.Ref.Path, ExpectedRevision: def.Revision, Document: []byte(next)})
	if e == nil {
		t.Fatal("definition invalidating instances accepted")
	}
	raw, _ := os.ReadFile(filepath.Join(p.Root, "factory", def.Ref.Path))
	if string(raw) != string(def.Raw) {
		t.Fatal("definition changed on rejection")
	}
	if d.Type != "Task" {
		t.Fatal("wrong type")
	}
}
func TestMinimalWorkflowBoundary(t *testing.T) {
	p, ps := storeFixture(t)
	d, e := Create(ps, p, "Task", "one", map[string]any{"title": "One", "summary": "Never move directly to done without five approvals."})
	if e != nil {
		t.Fatal(e)
	}
	updated, e := Mutate(ps, p, Mutation{Kind: "update", Path: d.Ref.Path, ExpectedRevision: d.Revision, Set: map[string]any{"task_state": "done"}})
	if e != nil {
		t.Fatal(e)
	}
	if knowledge.String(updated.Frontmatter, "task_state") != "done" {
		t.Fatal("structural enum change not applied")
	}
}
func TestInitConflictDoesNotPartiallyWrite(t *testing.T) {
	p, _ := storeFixture(t)
	if e := os.Remove(filepath.Join(p.Root, "factory/agents/builder.md")); e != nil {
		t.Fatal(e)
	}
	if e := os.WriteFile(filepath.Join(p.Root, "factory/factory.md"), []byte("changed"), 0644); e != nil {
		t.Fatal(e)
	}
	r, e := Init(p.Root, "demo", "software")
	if e == nil || len(r.Created) != 0 {
		t.Fatalf("conflict result: %#v %v", r, e)
	}
	if _, e = os.Stat(filepath.Join(p.Root, "factory/agents/builder.md")); !os.IsNotExist(e) {
		t.Fatal("init wrote before resolving conflicts")
	}
}
