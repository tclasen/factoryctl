package app_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/tclasen/factoryctl/internal/app"
	"github.com/tclasen/factoryctl/internal/knowledge"
	"github.com/tclasen/factoryctl/internal/project"
)

func collectionPut(t *testing.T, p string, b []byte) {
	t.Helper()
	if e := os.MkdirAll(filepath.Dir(p), 0755); e != nil {
		t.Fatal(e)
	}
	if e := os.WriteFile(p, b, 0644); e != nil {
		t.Fatal(e)
	}
}
func collectionRepo(t *testing.T, p string) {
	t.Helper()
	collectionPut(t, filepath.Join(p, ".git", "HEAD"), []byte("ref: refs/heads/main\n"))
	if e := os.MkdirAll(filepath.Join(p, ".git", "objects"), 0755); e != nil {
		t.Fatal(e)
	}
}
func collectionRead(t *testing.T, p string) []byte {
	t.Helper()
	b, e := os.ReadFile(p)
	if e != nil {
		t.Fatal(e)
	}
	return b
}
func collectionCode(t *testing.T, e error, code string) {
	t.Helper()
	var ae *app.Error
	if !errors.As(e, &ae) || ae.Code != code {
		t.Fatalf("want %s; got %v", code, e)
	}
}
func collectionNode(t *testing.T, p project.Project, state, peer string) {
	t.Helper()
	front := map[string]any{"type": "Widget", "title": p.ID + " record", "state": state, "peer": peer, "factory": map[string]any{"format": 1}}
	b, e := knowledge.Encode(front, "\n# Target\n\nContent owned by "+p.ID+".\n")
	if e != nil {
		t.Fatal(e)
	}
	collectionPut(t, filepath.Join(p.Root, "factory", "widgets", "example.md"), b)
}
func collectionFixture(t *testing.T, scenario string) []project.Project {
	t.Helper()
	root, e := filepath.EvalSymlinks(t.TempDir())
	if e != nil {
		t.Fatal(e)
	}
	names := []string{"alpha", "beta"}
	if scenario == "S1" {
		names = []string{"alpha"}
	}
	if scenario != "S3" {
		collectionRepo(t, root)
	}
	for _, name := range names {
		p := root
		if scenario != "S1" {
			p = filepath.Join(root, name)
		}
		if scenario == "S3" {
			collectionRepo(t, p)
		}
		if e := os.MkdirAll(p, 0755); e != nil {
			t.Fatal(e)
		}
		collectionPut(t, filepath.Join(p, "factory", "factory.md"), []byte(fmt.Sprintf("---\ntype: Factory\ntitle: %s\nproject_id: %s\nbundle_uri: https://example.org/%s/factory/\n---\n# Factory\n", name, name, name)))
		collectionPut(t, filepath.Join(p, "factory", "control", "write-coordination.md"), []byte("---\ntype: Factory Control\ntitle: Writer lock\n---\n# Writer lock\n"))
		collectionPut(t, filepath.Join(p, "factory", "definitions", "widget.md"), []byte(fmt.Sprintf("---\ntype: Record Definition\ntitle: Widget\ndefines: Widget\ndirectory: widgets\nfields:\n  state:\n    kind: enum\n    required: true\n    values: [%s, %s-ready]\n  peer:\n    kind: reference\n    required: true\n    target_types: [Widget]\n---\n# Widget\n\nWorkflow prose is interpreted by the harness.\n", name, name)))
	}
	c, e := project.Discover(root, true)
	if e != nil {
		t.Fatal(e)
	}
	if len(c.Projects) != len(names) {
		t.Fatalf("collection: %#v", c)
	}
	for i, p := range c.Projects {
		peer := c.Projects[(i+1)%len(c.Projects)]
		collectionNode(t, p, p.ID, peer.BundleURI+"widgets/example.md#target")
	}
	return c.Projects
}

func TestCollectionDocumentJourneys(t *testing.T) {
	for _, scenario := range []string{"S1", "S2", "S3"} {
		t.Run(scenario, func(t *testing.T) {
			projects := collectionFixture(t, scenario)
			snap, e := app.Load(projects)
			if e != nil {
				t.Fatal(e)
			}
			if v := snap.Validate("factory"); v.Errors != 0 || !v.Complete {
				t.Fatalf("initial validation: %#v", v)
			}
			for _, p := range projects {
				b := snap.Bundle(p.Checkout)
				if b == nil {
					t.Fatal("missing bundle")
				}
				node := b.Docs["widgets/example.md"]
				if knowledge.String(node.Frontmatter, "state") != p.ID {
					t.Fatal("same-named records conflated")
				}
				if b.Definitions["Widget"].Path != "definitions/widget.md" {
					t.Fatal("missing local definition")
				}
				outgoing, incoming := 0, 0
				for _, link := range snap.Links {
					if link.Source.Checkout == p.Checkout {
						outgoing++
						if link.Status != "valid" || link.Resolved == nil {
							t.Fatalf("unresolved cross link: %#v", link)
						}
					}
					if link.Resolved != nil && link.Resolved.Checkout == p.Checkout && link.Resolved.Path == node.Ref.Path {
						incoming++
					}
				}
				if outgoing != 1 || incoming != 1 {
					t.Fatalf("outgoing=%d incoming=%d", outgoing, incoming)
				}
			}
			selected := projects[0]
			original := snap.Bundle(selected.Checkout).Docs["widgets/example.md"]
			otherBytes := map[string][]byte{}
			for _, p := range projects[1:] {
				otherBytes[p.Root] = collectionRead(t, filepath.Join(p.Root, "factory/widgets/example.md"))
			}
			updated, e := app.Mutate(projects, selected, app.Mutation{Kind: "update", Path: original.Ref.Path, ExpectedRevision: original.Revision, Set: map[string]any{"state": selected.ID + "-ready"}})
			if e != nil {
				t.Fatal(e)
			}
			if updated.Revision == original.Revision {
				t.Fatal("mutation did not change revision")
			}
			for root, b := range otherBytes {
				if !bytes.Equal(b, collectionRead(t, filepath.Join(root, "factory/widgets/example.md"))) {
					t.Fatal("write changed another project's same-named record")
				}
			}
			_, e = app.Mutate(projects, selected, app.Mutation{Kind: "update", Path: original.Ref.Path, ExpectedRevision: original.Revision, Set: map[string]any{"title": "stale"}})
			collectionCode(t, e, "stale_revision")
			if len(projects) > 1 {
				_, e = app.Mutate(projects, selected, app.Mutation{Kind: "update", Path: original.Ref.Path, ExpectedRevision: updated.Revision, Set: map[string]any{"state": projects[1].ID}})
				collectionCode(t, e, "invalid_document")
			}
			raw := collectionRead(t, filepath.Join(selected.Root, "factory", original.Ref.Path))
			_, e = app.Mutate(projects, selected, app.Mutation{Kind: "replace", Path: original.Ref.Path, ExpectedRevision: updated.Revision, Document: bytes.Replace(raw, []byte("# Target"), []byte("# Removed"), 1)})
			collectionCode(t, e, "invalid_document")
			if !bytes.Equal(raw, collectionRead(t, filepath.Join(selected.Root, "factory", original.Ref.Path))) {
				t.Fatal("rejected heading removal changed bytes")
			}
			if len(projects) > 1 {
				_, e = app.Mutate(projects, selected, app.Mutation{Kind: "replace", Path: original.Ref.Path, ExpectedRevision: updated.Revision, Document: bytes.Replace(raw, []byte("type: Widget"), []byte("type: Other"), 1)})
				collectionCode(t, e, "invalid_document")
				if !bytes.Equal(raw, collectionRead(t, filepath.Join(selected.Root, "factory", original.Ref.Path))) {
					t.Fatal("rejected target-type change changed bytes")
				}
				_, e = app.Mutate(projects, selected, app.Mutation{Kind: "delete", Path: original.Ref.Path, ExpectedRevision: updated.Revision})
				collectionCode(t, e, "referenced_node")
			}
			_, e = app.Mutate(projects, selected, app.Mutation{Kind: "create", Path: "../escape.md", Document: raw})
			collectionCode(t, e, "invalid_path")
			reloaded, e := app.Load(projects)
			if e != nil {
				t.Fatal(e)
			}
			if v := reloaded.Validate("factory"); v.Errors != 0 {
				t.Fatalf("after rejected mutations: %#v", v)
			}
		})
	}
}

func TestCollectionAmbiguousCanonicalURL(t *testing.T) {
	projects := collectionFixture(t, "S3")
	projects[0].BundleURI = projects[1].BundleURI
	s, e := app.Load(projects)
	if e != nil {
		t.Fatal(e)
	}
	found := false
	for _, l := range s.Links {
		if l.Source.Checkout == projects[0].Checkout && l.Status == "broken" && l.Resolved == nil {
			found = true
		}
	}
	if !found {
		t.Fatalf("ambiguous canonical URL was silently resolved: %#v", s.Links)
	}
}
