package web

import (
	"context"
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tclasen/factoryctl/internal/project"
	rt "github.com/tclasen/factoryctl/internal/runtime"
)

func settingsFixture(t *testing.T) (*Handler, string, string, string) {
	t.Helper()
	h, _ := fixture(t)
	dir, e := filepath.EvalSymlinks(t.TempDir())
	if e != nil {
		t.Fatal(e)
	}
	path := filepath.Join(dir, "personal", "chat-settings.json")
	secret := h.EnableConversations(RuntimeConfig{SettingsPath: path})
	t.Cleanup(h.Close)
	h.chat.catalogLookup = func(context.Context, project.Project) (catalogSnapshot, error) {
		return catalogSnapshot{Models: []rt.Model{
			{ID: "model-a", Model: "model-a", DisplayName: "A", Default: true, DefaultEffort: "low", Efforts: []string{"low", "high"}},
			{ID: "model-b", Model: "model-b", DisplayName: "B", DefaultEffort: "medium", Efforts: []string{"medium"}},
		}, NativeDefault: rt.ChatSelection{Model: "model-b", Effort: "medium"}}, nil
	}
	return h, secret, h.projects[0].Checkout, path
}
func readSettings(t *testing.T, h *Handler, secret, id string) map[string]json.RawMessage {
	t.Helper()
	w := hostRequest(h, "GET", "/api/chat-settings?project="+url.QueryEscape(id), secret, "")
	if w.Code != 200 {
		t.Fatal(w.Code, w.Body.String())
	}
	var result map[string]json.RawMessage
	if e := json.Unmarshal(w.Body.Bytes(), &result); e != nil {
		t.Fatal(e)
	}
	return result
}
func revisionOf(m map[string]json.RawMessage) string {
	var s string
	_ = json.Unmarshal(m["revision"], &s)
	return s
}
func TestChatSettingsRequireHostAndLoadedProject(t *testing.T) {
	h, secret, id, path := settingsFixture(t)
	_, tool, _ := h.chat.broker.NewScope("example")
	for _, token := range []string{"", tool, "wrong"} {
		if w := hostRequest(h, "GET", "/api/chat-settings?project="+id, token, ""); w.Code != 403 {
			t.Fatal(w.Code)
		}
	}
	if w := hostRequest(h, "GET", "/api/chat-settings?project=not-loaded", secret, ""); w.Code != 400 {
		t.Fatal(w.Code)
	}
	readSettings(t, h, secret, id)
	if _, e := os.Stat(path); !os.IsNotExist(e) {
		t.Fatal("read created settings", e)
	}
}
func TestChatSettingsSaveConflictResetAndProjectIsolation(t *testing.T) {
	h, secret, id, path := settingsFixture(t)
	dir, e := filepath.EvalSymlinks(t.TempDir())
	if e != nil {
		t.Fatal(e)
	}
	h.chat.projects = append(h.chat.projects, project.Project{ID: h.projects[0].ID, Checkout: "second", Root: dir})
	post := func(scope, selection, revision string) int {
		body := `{"project":"` + id + `","scope":"` + scope + `","selection":` + selection + `,"revision":"` + revision + `"}`
		return hostRequest(h, "POST", "/api/chat-settings", secret, body).Code
	}
	if code := post("global", `{"model":"model-a","reasoning_effort":"high"}`, "absent"); code != 200 {
		t.Fatal(code)
	}
	first := readSettings(t, h, secret, id)
	rev := revisionOf(first)
	if string(first["source"]) != `"global"` {
		t.Fatal(string(first["source"]))
	}
	if code := post("project", `{"model":"model-b","reasoning_effort":"medium"}`, rev); code != 200 {
		t.Fatal(code)
	}
	if code := post("global", "null", rev); code != 409 {
		t.Fatal("stale write accepted", code)
	}
	current := readSettings(t, h, secret, id)
	other := readSettings(t, h, secret, "second")
	if !strings.Contains(string(current["effective"]), "model-b") || !strings.Contains(string(other["effective"]), "model-a") {
		t.Fatal("project selection leaked")
	}
	if code := post("project", "null", revisionOf(current)); code != 200 {
		t.Fatal(code)
	}
	if !strings.Contains(string(readSettings(t, h, secret, id)["effective"]), "model-a") {
		t.Fatal("reset did not inherit")
	}
	bytesBefore, _ := os.ReadFile(path)
	if code := post("global", `{"model":"model-b","reasoning_effort":"high"}`, revisionOf(readSettings(t, h, secret, id))); code != 400 {
		t.Fatal("invalid pair accepted", code)
	}
	if code := post("global", `{"model":"","reasoning_effort":"high"}`, revisionOf(readSettings(t, h, secret, id))); code != 400 {
		t.Fatal("used catalog default instead of native config default", code)
	}
	bytesAfter, _ := os.ReadFile(path)
	if string(bytesBefore) != string(bytesAfter) {
		t.Fatal("rejected save changed file")
	}
	w := hostRequest(h, "POST", "/api/chat-settings", secret, `{"project":"`+id+`","scope":"global","revision":"`+revisionOf(readSettings(t, h, secret, id))+`"}`)
	if w.Code != 400 {
		t.Fatal("missing choice silently reset defaults")
	}
	// A new host/store instance sees the persisted global selection.
	cfg := h.chat.config
	h2, _ := fixture(t)
	s2 := h2.EnableConversations(cfg)
	defer h2.Close()
	h2.chat.catalogLookup = h.chat.catalogLookup
	if !strings.Contains(string(readSettings(t, h2, s2, h2.projects[0].Checkout)["effective"]), "model-a") {
		t.Fatal("restart lost global defaults")
	}
}
func TestMalformedSettingsDoNotBlockBrowsingOrExplicitChoices(t *testing.T) {
	h, secret, id, path := settingsFixture(t)
	if e := os.MkdirAll(filepath.Dir(path), 0700); e != nil {
		t.Fatal(e)
	}
	if e := os.WriteFile(path, []byte(`{"version":0}`), 0600); e != nil {
		t.Fatal(e)
	}
	data := readSettings(t, h, secret, id)
	if len(data["settings_error"]) == 0 {
		t.Fatal("missing diagnostic")
	}
	if request(h, "GET", "/api/documents").Code != 200 {
		t.Fatal("browsing blocked")
	}
	explicit := rt.ChatSelection{Model: "model-a", Effort: "low"}
	chosen, e := h.chat.selectedDefaults(h.projects[0].Root, &explicit)
	if e != nil || chosen != explicit {
		t.Fatal(chosen, e)
	}
	if _, e := h.chat.selectedDefaults(h.projects[0].Root, nil); e == nil {
		t.Fatal("silently ignored malformed defaults")
	}
	if _, e := (&selectionInput{}).selection(); e == nil {
		t.Fatal("partial selection accepted")
	}
}
