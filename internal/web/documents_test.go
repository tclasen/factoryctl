package web

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tclasen/factoryctl/internal/app"
	"github.com/tclasen/factoryctl/internal/project"
)

func editingFixture(t *testing.T) (*Handler, project.Project) {
	t.Helper()
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, ".git", "objects"), 0755)
	os.WriteFile(filepath.Join(root, ".git", "HEAD"), []byte("ref: refs/heads/main\n"), 0644)
	if _, e := app.Init(root, "example", "software"); e != nil {
		t.Fatal(e)
	}
	p := project.Project{ID: "example", Name: "Example", Root: root, Checkout: "edit-checkout"}
	h, e := New([]project.Project{p})
	if e != nil {
		t.Fatal(e)
	}
	return h, p
}
func documentRequest(h *Handler, method, target string, v any, nonce string) *httptest.ResponseRecorder {
	b, _ := json.Marshal(v)
	r := httptest.NewRequest(method, "http://127.0.0.1:8080"+target, strings.NewReader(string(b)))
	r.Header.Set("X-Factory-Nonce", nonce)
	r.Header.Set("Origin", "http://127.0.0.1:8080")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}
func readDocument(t *testing.T, h *Handler, p project.Project, name string) Document {
	t.Helper()
	w := request(h, "GET", "/api/document?id="+documentID(p.Checkout, name))
	if w.Code != 200 {
		t.Fatal(w.Code, w.Body.String())
	}
	var d Document
	if e := json.Unmarshal(w.Body.Bytes(), &d); e != nil {
		t.Fatal(e)
	}
	return d
}
func TestDocumentCRUDWithoutRuntimeAndConflicts(t *testing.T) {
	h, p := editingFixture(t)
	s, _ := app.Load(h.projects)
	hash := app.DefinitionHash(s.Bundle(p.Checkout).Definitions["Task"])
	create := map[string]any{"ref": map[string]any{"checkout": p.Checkout}, "type": "Task", "slug": "example", "inputs": map[string]any{"title": "Example task", "summary": "An example"}, "expected_definition": hash}
	if w := documentRequest(h, "POST", "/api/v1/nodes", create, ""); w.Code != 403 {
		t.Fatal("nonce bypass", w.Code)
	}
	if w := documentRequest(h, "POST", "/api/v1/nodes", create, h.nonce); w.Code != 200 {
		t.Fatal(w.Code, w.Body.String())
	}
	d := readDocument(t, h, p, "tasks/example.md")
	payload := map[string]any{"ref": d.Ref, "expected_revision": d.Revision, "expected_definition": d.DefinitionHash, "set": map[string]any{"task_state": "done"}}
	if w := documentRequest(h, "PUT", "/api/v1/node", payload, h.nonce); w.Code != 200 {
		t.Fatal(w.Code, w.Body.String())
	}
	if w := documentRequest(h, "PUT", "/api/v1/node", payload, h.nonce); w.Code != 409 {
		t.Fatal("stale save", w.Code, w.Body.String())
	}
	d = readDocument(t, h, p, "tasks/example.md")
	original, _ := os.ReadFile(filepath.Join(p.Root, "factory", d.Path))
	payload["expected_revision"] = d.Revision
	payload["set"] = map[string]any{"task_state": "invented"}
	if w := documentRequest(h, "PUT", "/api/v1/node", payload, h.nonce); w.Code != 422 {
		t.Fatal(w.Code, w.Body.String())
	}
	after, _ := os.ReadFile(filepath.Join(p.Root, "factory", d.Path))
	if string(after) != string(original) {
		t.Fatal("failed move changed disk")
	}
	payload["expected_definition"] = "stale"
	payload["set"] = map[string]any{"task_state": "ready"}
	if w := documentRequest(h, "PUT", "/api/v1/node", payload, h.nonce); w.Code != 409 {
		t.Fatal("definition conflict", w.Code, w.Body.String())
	}
	payload["expected_definition"] = d.DefinitionHash
	delete(payload, "set")
	if w := documentRequest(h, "DELETE", "/api/v1/node", payload, h.nonce); w.Code != 200 || !strings.Contains(w.Body.String(), `"removed"`) {
		t.Fatal(w.Code, w.Body.String())
	}
	if _, e := os.Stat(filepath.Join(p.Root, "factory", d.Path)); !os.IsNotExist(e) {
		t.Fatal("delete not persisted")
	}
}
func TestDocumentRefreshAndScope(t *testing.T) {
	h, p := editingFixture(t)
	name := filepath.Join(p.Root, "factory", "external.md")
	raw := "---\ntype: Reference\ntitle: External\n---\n# External\n"
	os.WriteFile(name, []byte(raw), 0644)
	d := readDocument(t, h, p, "external.md")
	if d.Title != "External" {
		t.Fatal(d)
	}
	os.WriteFile(name, []byte(strings.ReplaceAll(raw, "External", "Changed")), 0644)
	if readDocument(t, h, p, "external.md").Title != "Changed" {
		t.Fatal("external edit not refreshed")
	}
	w := request(h, "GET", "/api/v1/node?checkout=other&path=external.md")
	if w.Code != 409 {
		t.Fatal("cross-project ref", w.Code)
	}
	w = request(h, "GET", "/api/v1/node?checkout="+p.Checkout+"&path=../secret.md")
	if w.Code != 404 {
		t.Fatal("path escape", w.Code)
	}
}
func customFixture(t *testing.T) (*Handler, project.Project, Document) {
	h, p := editingFixture(t)
	dir := filepath.Join(p.Root, ".factory", "views")
	os.MkdirAll(dir, 0755)
	os.WriteFile(filepath.Join(dir, "task.html"), []byte(`<html><body>{{.Record.Title}}<script src="task.js"></script></body></html>`), 0644)
	os.WriteFile(filepath.Join(dir, "task.js"), []byte(`window.loaded=true;`), 0644)
	f := filepath.Join(p.Root, "factory", "definitions", "task.md")
	raw, e := os.ReadFile(f)
	if e != nil {
		t.Fatal(e)
	}
	raw = []byte(strings.Replace(string(raw), "\n---\n", "\nviews:\n  detail_template: .factory/views/task.html\n---\n", 1))
	os.WriteFile(f, raw, 0644)
	s, _ := app.Load(h.projects)
	_, e = app.Create(h.projects, p, "Task", "custom", map[string]any{"title": "Custom", "summary": "Custom view test"})
	if e != nil {
		t.Fatal(e, s.Diagnostics)
	}
	return h, p, readDocument(t, h, p, "tasks/custom.md")
}
func TestCustomTrustDigestScopeAndSandbox(t *testing.T) {
	h, p, d := customFixture(t)
	g, e := h.customBundle(d.Ref)
	if e != nil {
		t.Fatal(e)
	}
	payload := map[string]any{"ref": d.Ref, "digest": g.Digest}
	w := documentRequest(h, "POST", "/api/v1/custom/trust", payload, h.nonce)
	if w.Code != 200 {
		t.Fatal(w.Code, w.Body.String())
	}
	var grant map[string]any
	json.Unmarshal(w.Body.Bytes(), &grant)
	url := grant["frame_url"].(string)
	w = request(h, "GET", url)
	if w.Code != 200 || !strings.Contains(w.Header().Get("Content-Security-Policy"), "sandbox allow-scripts") || !strings.Contains(w.Header().Get("Content-Security-Policy"), "connect-src 'none'") {
		t.Fatal(w.Code, w.Header(), w.Body.String())
	}
	if w = request(h, "GET", strings.Replace(url, "task.html", "../../../secret.md", 1)); w.Code != 404 {
		t.Fatal("asset escape", w.Code)
	}
	msg := map[string]any{"grant": grant["grant"], "nonce": "spoofed", "digest": g.Digest, "type": "requestEdit", "id": d.ID, "expected_revision": d.Revision, "set": map[string]any{"task_state": "ready"}}
	if w = documentRequest(h, "POST", "/api/v1/custom/request", msg, h.nonce); w.Code != 400 {
		t.Fatal("spoof nonce", w.Code)
	}
	msg["nonce"] = grant["nonce"]
	msg["id"] = documentID("other", d.Path)
	if w = documentRequest(h, "POST", "/api/v1/custom/request", msg, h.nonce); w.Code != 400 {
		t.Fatal("selection escape", w.Code)
	}
	msg["id"] = d.ID
	if w = documentRequest(h, "POST", "/api/v1/custom/request", msg, h.nonce); w.Code != 200 {
		t.Fatal(w.Code, w.Body.String())
	}
	if w = documentRequest(h, "POST", "/api/v1/custom/request", msg, h.nonce); w.Code != 409 {
		t.Fatal("stale custom edit", w.Code, w.Body.String())
	}
	os.WriteFile(filepath.Join(p.Root, ".factory", "views", "task.js"), []byte(`window.changed=true;`), 0644)
	if w = request(h, "GET", url); w.Code != 409 {
		t.Fatal("transitive digest invalidation", w.Code)
	}
	if w = documentRequest(h, "POST", "/api/v1/custom/trust", payload, h.nonce); w.Code != 409 {
		t.Fatal("stale trust accepted", w.Code)
	}
}
func TestCustomRejectsSymlinkAndExternalAssets(t *testing.T) {
	h, p, d := customFixture(t)
	f := filepath.Join(p.Root, ".factory", "views", "task.js")
	os.Remove(f)
	os.Symlink(filepath.Join(p.Root, "factory", "factory.md"), f)
	if _, e := h.customBundle(d.Ref); e == nil {
		t.Fatal("symlink accepted")
	}
	os.Remove(f)
	os.WriteFile(f, []byte(`import "https://attacker.invalid/x.js"`), 0644)
	if _, e := h.customBundle(d.Ref); e == nil {
		t.Fatal("external import accepted")
	}
}

func TestCustomDisabledDefinitionUsesGenericFallback(t *testing.T) {
	h, p, d := customFixture(t)
	f := filepath.Join(p.Root, "factory", "definitions", "task.md")
	raw, _ := os.ReadFile(f)
	raw = []byte(strings.Replace(string(raw), "\n---\n", "\nenabled: false\n---\n", 1))
	os.WriteFile(f, raw, 0644)
	if _, e := h.customBundle(d.Ref); e == nil {
		t.Fatal("disabled custom view executed")
	}
	if got := readDocument(t, h, p, d.Path); got.Raw == "" {
		t.Fatal("generic read unavailable")
	}
}
func TestExternalHealthChangeInvalidatesSnapshot(t *testing.T) {
	h, p := editingFixture(t)
	before := h.snapshot
	os.WriteFile(filepath.Join(p.Root, "factory", "bad.txt"), []byte("not markdown"), 0644)
	if e := h.refresh(); e != nil {
		t.Fatal(e)
	}
	if before == h.snapshot {
		t.Fatal("diagnostic-only edit did not invalidate snapshot")
	}
}

func TestCanonicalCrossProjectLinksAndBacklinks(t *testing.T) {
	h, root := fixture(t)
	other := t.TempDir()
	os.Mkdir(filepath.Join(other, "factory"), 0755)
	os.WriteFile(filepath.Join(other, "factory", "target.md"), []byte("---\ntype: Reference\ntitle: Other target\n---\n# Other target\n"), 0644)
	raw, _ := os.ReadFile(filepath.Join(root, "factory", "plan.md"))
	raw = append(raw, []byte("\n[Other](https://example.invalid/other/target.md)\n")...)
	os.WriteFile(filepath.Join(root, "factory", "plan.md"), raw, 0644)
	projects := append(h.projects, project.Project{ID: "other", Name: "Other", Root: other, Checkout: "other-checkout", BundleURI: "https://example.invalid/other/"})
	next, e := New(projects)
	if e != nil {
		t.Fatal(e)
	}
	d := readDocument(t, next, projects[0], "plan.md")
	target := documentID("other-checkout", "target.md")
	if !strings.Contains(d.HTML, "#doc="+target) {
		t.Fatalf("canonical link inert: %s", d.HTML)
	}
	otherDoc := readDocument(t, next, projects[1], "target.md")
	found := false
	for _, link := range otherDoc.Backlinks {
		if link.ID == d.ID {
			found = true
		}
	}
	if !found {
		t.Fatal("cross-project backlink missing")
	}
}

func TestHTTPNodeCursorBindsSelectionAndRejectsMalformedOffset(t *testing.T) {
	h, _ := fixture(t)
	w := request(h, "GET", "/api/v1/nodes?limit=1")
	var page struct {
		Next string `json:"next_cursor"`
	}
	if e := json.Unmarshal(w.Body.Bytes(), &page); e != nil || page.Next == "" {
		t.Fatal(e, w.Body.String())
	}
	if w = request(h, "GET", "/api/v1/nodes?type=Task&cursor="+page.Next); w.Code != 409 {
		t.Fatal("cursor crossed filter scope", w.Code, w.Body.String())
	}
	bits := strings.Split(page.Next, ":")
	bits[2] = "not-a-number"
	if w = request(h, "GET", "/api/v1/nodes?cursor="+strings.Join(bits, ":")); w.Code != 400 {
		t.Fatal("malformed cursor accepted", w.Code, w.Body.String())
	}
}
