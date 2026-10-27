package project

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func tempDir(t *testing.T) string {
	t.Helper()
	p, e := filepath.EvalSymlinks(t.TempDir())
	if e != nil {
		t.Fatal(e)
	}
	return p
}
func put(t *testing.T, p, s string) {
	t.Helper()
	if e := os.MkdirAll(filepath.Dir(p), 0755); e != nil {
		t.Fatal(e)
	}
	if e := os.WriteFile(p, []byte(s), 0644); e != nil {
		t.Fatal(e)
	}
}
func repo(t *testing.T, p string) {
	put(t, filepath.Join(p, ".git", "HEAD"), "ref: refs/heads/main\n")
	if e := os.MkdirAll(filepath.Join(p, ".git", "objects"), 0755); e != nil {
		t.Fatal(e)
	}
}
func bundle(t *testing.T, p, id string) {
	put(t, filepath.Join(p, "factory", "factory.md"), "---\ntype: Factory\ntitle: "+id+"\nproject_id: "+id+"\nbundle_uri: https://example.org/"+id+"/\n---\n# Factory\n")
}
func discover(t *testing.T, p string, explicit bool) Collection {
	t.Helper()
	c, e := Discover(p, explicit)
	if e != nil {
		t.Fatal(e)
	}
	return c
}
func hasCode(c Collection, code string) bool {
	for _, d := range c.Diagnostics {
		if d.Code == code {
			return true
		}
	}
	return false
}
func TestCollectionScenarios(t *testing.T) {
	for _, scenario := range []string{"S1", "S2", "S3"} {
		t.Run(scenario, func(t *testing.T) {
			root := tempDir(t)
			want := 2
			if scenario == "S1" {
				repo(t, root)
				bundle(t, root, "solo")
				want = 1
			} else {
				if scenario == "S2" {
					repo(t, root)
				}
				for _, name := range []string{"alpha", "beta"} {
					p := filepath.Join(root, name)
					if scenario == "S3" {
						repo(t, p)
					}
					bundle(t, p, name)
					put(t, filepath.Join(p, "factory", "tasks", "example.md"), name)
				}
			}
			c := discover(t, root, false)
			if len(c.Projects) != want {
				t.Fatalf("got %#v", c)
			}
			for _, p := range c.Projects {
				if p.Root == "" || p.Checkout == "" || p.RepositoryRoot == "" {
					t.Fatal(p)
				}
				if scenario == "S2" && p.RepositoryRoot != root {
					t.Fatal(p)
				}
			}
			if scenario != "S1" && c.Projects[0].Checkout == c.Projects[1].Checkout {
				t.Fatal("project instances collapsed")
			}
		})
	}
}
func TestDepthAndFallback(t *testing.T) {
	root := tempDir(t)
	repo(t, root)
	bundle(t, root, "root")
	bundle(t, filepath.Join(root, "child"), "child")
	deep := filepath.Join(root, "group", "deep")
	bundle(t, deep, "deep")
	c := discover(t, root, false)
	if len(c.Projects) != 2 {
		t.Fatal(c)
	}
	c = discover(t, filepath.Join(root, "group"), true)
	if len(c.Projects) != 1 || c.Projects[0].ID != "deep" {
		t.Fatal(c)
	}
	sub := filepath.Join(root, "src", "pkg")
	if e := os.MkdirAll(sub, 0755); e != nil {
		t.Fatal(e)
	}
	c = discover(t, sub, false)
	if len(c.Projects) != 1 || c.Projects[0].Root != root {
		t.Fatal(c)
	}
	if c = discover(t, sub, true); len(c.Projects) != 0 {
		t.Fatal(c)
	}
	outer := tempDir(t)
	child := filepath.Join(outer, "repo")
	repo(t, child)
	bundle(t, filepath.Join(child, "nested"), "hidden")
	if c = discover(t, outer, false); len(c.Projects) != 0 {
		t.Fatal("depth counter restarted", c)
	}
}
func TestWorktreeAndInvalidMarkers(t *testing.T) {
	root := tempDir(t)
	main := filepath.Join(root, "main")
	repo(t, main)
	work := filepath.Join(root, "work")
	meta := filepath.Join(main, ".git", "worktrees", "work")
	put(t, filepath.Join(meta, "HEAD"), "ref: refs/heads/feature\n")
	put(t, filepath.Join(meta, "commondir"), "../..\n")
	put(t, filepath.Join(work, ".git"), "gitdir: "+meta+"\n")
	bundle(t, work, "work")
	c := discover(t, work, true)
	if len(c.Projects) != 1 || c.Projects[0].RepositoryRoot != work {
		t.Fatal(c)
	}
	put(t, filepath.Join(work, ".git"), "not git\n")
	c = discover(t, work, true)
	if len(c.Projects) != 0 || !hasCode(c, "invalid_repository") {
		t.Fatal(c)
	}
}
func TestDuplicatesAndMalformedMetadata(t *testing.T) {
	root := tempDir(t)
	repo(t, root)
	bundle(t, filepath.Join(root, "a"), "duplicate")
	bundle(t, filepath.Join(root, "b"), "duplicate")
	c := discover(t, root, false)
	if len(c.Projects) != 2 || !hasCode(c, "duplicate_project_id") || !hasCode(c, "duplicate_bundle_uri") {
		t.Fatal(c)
	}
	bad := filepath.Join(root, "bad")
	put(t, filepath.Join(bad, "factory", "factory.md"), "---\ntype: Factory\nproject_id: a\nproject_id: b\n---\n")
	c = discover(t, bad, true)
	if len(c.Projects) != 1 || !hasCode(c, "invalid_factory_metadata") || !strings.HasPrefix(c.Projects[0].ID, "local:") {
		t.Fatal(c)
	}
	generic := filepath.Join(root, "generic")
	put(t, filepath.Join(generic, "factory", "note.md"), "---\ntype: Note\n---\n")
	if c = discover(t, generic, true); len(c.Projects) != 1 {
		t.Fatal(c)
	}
	if c = discover(t, generic, false); len(c.Projects) != 0 {
		t.Fatal(c)
	}
}
func TestIgnoredSymlinksAndNonGit(t *testing.T) {
	root := tempDir(t)
	repo(t, root)
	for _, name := range []string{"skip", "keep", "build-out"} {
		bundle(t, filepath.Join(root, name), name)
	}
	put(t, filepath.Join(root, ".gitignore"), "skip/\nbuild-*\n!keep/\n")
	if e := os.Symlink(root, filepath.Join(root, "loop")); e != nil {
		t.Fatal(e)
	}
	c := discover(t, root, false)
	if len(c.Projects) != 1 || c.Projects[0].ID != "keep" || !hasCode(c, "symlink_excluded") || !hasCode(c, "ignored_path") {
		t.Fatal(c)
	}
	outside := tempDir(t)
	bundle(t, outside, "outside")
	if c = discover(t, outside, true); len(c.Projects) != 0 || !hasCode(c, "not_git_project") {
		t.Fatal(c)
	}
	symlinkBundle := filepath.Join(root, "linked")
	if e := os.MkdirAll(symlinkBundle, 0755); e != nil {
		t.Fatal(e)
	}
	if e := os.Symlink(filepath.Join(root, "keep", "factory"), filepath.Join(symlinkBundle, "factory")); e != nil {
		t.Fatal(e)
	}
	if c = discover(t, symlinkBundle, true); len(c.Projects) != 0 || !hasCode(c, "symlink_excluded") {
		t.Fatal(c)
	}
}
func TestMetadataHazards(t *testing.T) {
	for _, front := range []string{"type: Factory\nproject_id: &id hello\ntitle: *id", "type: Factory\nproject_id: !custom hello\ntitle: Name", "type: Factory\nproject_id: hello\ntitle: Name\nbundle_uri: http://example.org/", "type: Factory\nproject_id: hello\ntitle: Name\n42: nope"} {
		t.Run(front, func(t *testing.T) {
			root := tempDir(t)
			repo(t, root)
			put(t, filepath.Join(root, "factory", "factory.md"), "---\n"+front+"\n---\n")
			c := discover(t, root, true)
			if !hasCode(c, "invalid_factory_metadata") {
				t.Fatal(c)
			}
		})
	}
}
