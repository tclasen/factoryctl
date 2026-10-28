package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tclasen/factoryctl/internal/project"
)

func fixture(t *testing.T) (*Handler, string) {
	t.Helper()
	root := t.TempDir()
	if e := os.Mkdir(filepath.Join(root, "factory"), 0755); e != nil {
		t.Fatal(e)
	}
	docs := map[string]string{
		"plan.md":    "---\ntype: Design\ntitle: Current plan\nintent_state: proposed\nintent_owner: assistant\nintent_source: request.md#request\nintent_authorization: A small read-only preview only\n---\n# Current plan\n[Request](request.md#request)\n![Remote image](https://example.invalid/track.png)\n<script>alert('bad')</script>\n\n<https://example.invalid/>\n",
		"request.md": "---\ntype: Reference\ntitle: Request\nintent_state: accepted\nintent_owner: user\nintent_authorization: Implement the bounded preview\n---\n# Request\nAccepted scope only.\n",
		"unknown.md": "---\ntype: Task\ntitle: Not an approval\nstatus: stable\nverified: true\nsources:\n - author: human:user\n   resource: request.md\n---\nHuman author and stable status do not grant acceptance.\n",
	}
	for name, body := range docs {
		if e := os.WriteFile(filepath.Join(root, "factory", name), []byte(body), 0644); e != nil {
			t.Fatal(e)
		}
	}
	secret := filepath.Join(root, "secret.md")
	os.WriteFile(secret, []byte("outside-secret"), 0644)
	if e := os.Symlink(secret, filepath.Join(root, "factory", "escape.md")); e != nil {
		t.Fatal(e)
	}
	h, e := New([]project.Project{{ID: "example", Name: "Example", Root: root, Checkout: "test-checkout"}})
	if e != nil {
		t.Fatal(e)
	}
	return h, root
}
func request(h http.Handler, method, target string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, "http://127.0.0.1:8080"+target, nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func TestIntentDoesNotInferAcceptance(t *testing.T) {
	tests := []struct {
		metadata map[string]any
		want     string
	}{
		{map[string]any{"status": "stable", "verified": true, "author": "human:user"}, "unresolved"},
		{map[string]any{"intent_state": "accepted"}, "accepted"},
		{map[string]any{"intent_state": "made-up", "factory": map[string]any{"assertion": map[string]any{"certainty": "proposal"}}}, "unresolved"},
		{map[string]any{"factory": map[string]any{"assertion": map[string]any{"certainty": "proposal"}}}, "proposed"},
		{map[string]any{"factory": map[string]any{"assertion": map[string]any{"basis": "authoritative-statement", "certainty": "supported"}}}, "unresolved"},
	}
	for _, tt := range tests {
		got, _ := IntentState(tt.metadata)
		if got != tt.want {
			t.Errorf("%v: got%s want%s", tt.metadata, got, tt.want)
		}
	}
}
func TestReadOnlySnapshotAndInertMarkdown(t *testing.T) {
	h, root := fixture(t)
	w := request(h, "GET", "/api/documents")
	if w.Code != 200 {
		t.Fatal(w.Code, w.Body.String())
	}
	var listing struct {
		Documents []Document `json:"documents"`
	}
	if e := json.Unmarshal(w.Body.Bytes(), &listing); e != nil {
		t.Fatal(e)
	}
	if len(listing.Documents) != 3 {
		t.Fatalf("unsafe/missing document inclusion: %d", len(listing.Documents))
	}
	for _, d := range listing.Documents {
		if d.Raw != "" || d.HTML != "" {
			t.Fatal("list leaked full bodies")
		}
		if d.Path == "unknown.md" && d.IntentState != "unresolved" {
			t.Fatal("inferred approval")
		}
	}
	id := documentID("test-checkout", "plan.md")
	w = request(h, "GET", "/api/document?id="+id)
	var d Document
	json.Unmarshal(w.Body.Bytes(), &d)
	if d.IntentState != "proposed" || d.IntentOwner != "assistant" || d.IntentAuthorization != "A small read-only preview only" {
		t.Fatalf("lost provenance: %+v", d)
	}
	if strings.Contains(d.HTML, "<script") || strings.Contains(d.HTML, "<img") || strings.Contains(d.HTML, "href=\"https:") {
		t.Fatalf("active/untrusted rendering: %s", d.HTML)
	}
	if !strings.Contains(d.HTML, "#doc=") || len(d.Sources) != 1 || d.Sources[0].ID != documentID("test-checkout", "request.md") {
		t.Fatalf("source link lost: %+v", d)
	}
	if !strings.Contains(d.Raw, "<script>") {
		t.Fatal("raw source altered")
	}
	original, _ := os.ReadFile(filepath.Join(root, "factory", "plan.md"))
	for _, method := range []string{"POST", "PUT", "PATCH", "DELETE"} {
		w = request(h, method, "/api/document?id="+id)
		if w.Code != http.StatusMethodNotAllowed {
			t.Fatalf("%s allowed: %d", method, w.Code)
		}
	}
	after, _ := os.ReadFile(filepath.Join(root, "factory", "plan.md"))
	if string(original) != string(after) {
		t.Fatal("read-only operation changed source")
	}
	for _, target := range []string{"/api/document?id=../../secret.md", "/api/document?id=" + documentID("test-checkout", "escape.md"), "/secret.md", "/assets/../../secret.md"} {
		w = request(h, "GET", target)
		if w.Code != 404 || strings.Contains(w.Body.String(), "outside-secret") {
			t.Fatalf("escaped scope: %s %d", target, w.Code)
		}
	}
}
func TestHostOriginAndResourceBoundaries(t *testing.T) {
	h, _ := fixture(t)
	for _, tt := range []struct{ host, origin, site string }{
		{"attacker.example:8080", "", ""}, {"127.0.0.1:8080", "https://attacker.example", ""}, {"127.0.0.1:8080", "http://127.0.0.1:9999", ""}, {"127.0.0.1:8080", "", "cross-site"},
	} {
		r := httptest.NewRequest("GET", "http://127.0.0.1:8080/api/documents", nil)
		r.Host = tt.host
		r.Header.Set("Origin", tt.origin)
		r.Header.Set("Sec-Fetch-Site", tt.site)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != 403 {
			t.Fatalf("untrusted host/origin allowed: %+v %d", tt, w.Code)
		}
	}
	w := request(h, "GET", "/")
	if w.Code != 200 || !strings.Contains(w.Header().Get("Content-Security-Policy"), "img-src 'none'") || !strings.Contains(w.Body.String(), "Recorded") && !strings.Contains(w.Body.String(), "RECORDED CLAIM") {
		t.Fatal("missing inert page policy")
	}
	if request(h, "GET", "/app.js").Code != 200 {
		t.Fatal("embedded script missing")
	}
}
func TestListenRefusesNonLoopbackBeforeBinding(t *testing.T) {
	for _, address := range []string{"0.0.0.0:8080", ":8080", "192.168.1.1:8080", "example.com:8080", "[::]:8080", "127.0.0.1"} {
		listener, e := Listen(address)
		if e == nil {
			listener.Close()
			t.Fatalf("accepted %s", address)
		}
	}
}
