package web

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tclasen/factoryctl/internal/interaction"
)

func hostRequest(h *Handler, method, path, token, body string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, "http://127.0.0.1:8080"+path, strings.NewReader(body))
	if token != "" {
		r.Header.Set("X-Factory-Host", token)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}
func TestHostAndToolAnswerCapabilitiesAreSeparate(t *testing.T) {
	h, _ := fixture(t)
	secret := h.EnableConversations(RuntimeConfig{})
	defer h.Close()
	scope, tool, _ := h.chat.broker.NewScope("example")
	h.chat.conversations["c"] = &conversation{ID: "c", scope: scope}
	q, e := h.chat.broker.Ask(scope, interaction.AskInput{Prompt: "Actual answer?", Kind: "text"})
	if e != nil {
		t.Fatal(e)
	}
	body, _ := json.Marshal(map[string]any{"request_id": q.ID, "status": "answered", "text": "yes"})
	for _, token := range []string{"", tool, "wrong"} {
		w := hostRequest(h, "POST", "/api/conversations/c/answer", token, string(body))
		if w.Code != 403 {
			t.Fatal("untrusted answer accepted", w.Code)
		}
	}
	r := httptest.NewRequest("POST", "http://127.0.0.1:8080/internal/tool", strings.NewReader(`{"name":"factory_user_answer","arguments":{}}`))
	r.Header.Set("Authorization", "Bearer "+tool)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if !strings.Contains(w.Body.String(), "error") {
		t.Fatal("tool could answer")
	}
	r = httptest.NewRequest("POST", "http://127.0.0.1:8080/api/conversations/c/answer", bytes.NewReader(body))
	r.Header.Set("X-Factory-Host", secret)
	r.Header.Set("Origin", "http://evil.invalid")
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 403 {
		t.Fatal("cross-origin answer accepted")
	}
	if w := hostRequest(h, "POST", "/api/conversations/c/answer", secret, string(body)); w.Code != 200 {
		t.Fatal(w.Code, w.Body.String())
	}
	for _, path := range []string{"/", "/app.js", "/api/documents", "/api/conversations"} {
		w := hostRequest(h, "GET", path, "", "")
		if strings.Contains(w.Body.String(), secret) || strings.Contains(w.Body.String(), tool) {
			t.Fatal("credential leaked", path)
		}
	}
}
func TestGeneratedViewHasSeparateOpaqueSandboxPolicy(t *testing.T) {
	h, _ := fixture(t)
	h.EnableConversations(RuntimeConfig{})
	defer h.Close()
	scope, _, _ := h.chat.broker.NewScope("example")
	h.chat.conversations["c"] = &conversation{ID: "c", scope: scope}
	v, e := h.chat.broker.Present(scope, interaction.ViewInput{Title: "Form", Fallback: "Plain text", HTML: "<script>parent.postMessage({text:'draft'},'*')</script>"})
	if e != nil {
		t.Fatal(e)
	}
	if request(h, "GET", "/view/"+v.ID).Code != 404 {
		t.Fatal("untrusted direct view executed")
	}
	grant := hostRequest(h, "POST", "/api/conversations/c/trust-view", h.chat.secret, `{"view_id":"`+v.ID+`"}`)
	if grant.Code != 200 {
		t.Fatal(grant.Code, grant.Body.String())
	}
	var mount map[string]string
	_ = json.Unmarshal(grant.Body.Bytes(), &mount)
	if request(h, "GET", mount["frame_url"]).Code != 404 {
		t.Fatal("top-level frame bypassed trust boundary")
	}
	r := httptest.NewRequest("GET", "http://127.0.0.1:8080"+mount["frame_url"], nil)
	r.Header.Set("Sec-Fetch-Dest", "iframe")
	r.Header.Set("Sec-Fetch-Site", "same-origin")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatal(w.Code, w.Body.String())
	}
	csp := w.Header().Get("Content-Security-Policy")
	for _, part := range []string{"sandbox allow-scripts;", "connect-src 'none'", "frame-src 'none'", "frame-ancestors 'self'", "form-action 'none'"} {
		if !strings.Contains(csp, part) {
			t.Fatal("unsafe frame policy", csp)
		}
	}
	if strings.Contains(csp, "allow-same-origin") {
		t.Fatal(csp)
	}
	parent := request(h, "GET", "/")
	if strings.Contains(parent.Header().Get("Content-Security-Policy"), "unsafe-inline") {
		t.Fatal("parent policy weakened")
	}
	if request(h, "GET", "/view/unknown").Code != 404 {
		t.Fatal("unknown frame exposed")
	}
}
func TestNativeAnswersRequireExplicitSupportedDisposition(t *testing.T) {
	if validateNativeAnswer(nativeRequest{Method: "unknown"}, json.RawMessage(`{"decision":"accept"}`)) == nil {
		t.Fatal("unknown request accepted")
	}
	q := nativeRequest{Method: "item/tool/requestUserInput", Params: json.RawMessage(`{"questions":[{"id":"a"}]}`)}
	if validateNativeAnswer(q, json.RawMessage(`{"answers":{}}`)) == nil {
		t.Fatal("empty answer accepted")
	}
	if e := validateNativeAnswer(q, json.RawMessage(`{"answers":{"a":{"answers":["actual"]}}}`)); e != nil {
		t.Fatal(e)
	}
}

func TestApprovalMustMatchOfferedDecision(t *testing.T) {
	q := nativeRequest{Method: "item/commandExecution/requestApproval", Params: json.RawMessage(`{"availableDecisions":["decline","cancel"]}`)}
	if validateNativeAnswer(q, json.RawMessage(`{"decision":"accept"}`)) == nil {
		t.Fatal("unoffered approval accepted")
	}
	if e := validateNativeAnswer(q, json.RawMessage(`{"decision":"decline"}`)); e != nil {
		t.Fatal(e)
	}
}

func TestLateCompletionCannotReleaseNewTurn(t *testing.T) {
	c := &conversation{TurnID: "new", Status: "running", Native: []nativeRequest{{ID: "1"}}}
	if applyTurnEvent(c, "turn/completed", "old", "completed", nil) || c.Status != "running" || c.TurnID != "new" {
		t.Fatal("late completion released another turn")
	}
	c.Status = "interrupting"
	applyTurnEvent(c, "turn/started", "new", "", nil)
	if c.Status != "interrupting" {
		t.Fatal("start notification undid interrupt wait")
	}
	if !applyTurnEvent(c, "turn/completed", "new", "interrupted", nil) || c.Status != "interrupted" || c.TurnID != "" || len(c.Native) != 0 {
		t.Fatal("matching completion did not clear pending state")
	}
}
func TestInterruptingConversationRejectsFollowup(t *testing.T) {
	h, _ := fixture(t)
	secret := h.EnableConversations(RuntimeConfig{})
	defer h.Close()
	h.chat.conversations["c"] = &conversation{ID: "c", Status: "interrupting"}
	w := hostRequest(h, "POST", "/api/conversations/c/turn", secret, `{"text":"new turn"}`)
	if w.Code != 409 {
		t.Fatal(w.Code, w.Body.String())
	}
}

func TestDelayedStartResponseCannotBindToNextTurn(t *testing.T) {
	c := &conversation{generation: 1, Status: "running"}
	applyTurnEvent(c, "turn/started", "old", "", nil)
	applyTurnEvent(c, "turn/completed", "old", "completed", nil)
	applyTurnResponse(c, 1, "old")
	if c.TurnID != "" {
		t.Fatal("completed launch response restored old ID")
	}
	c.generation++
	c.Status = "running"
	applyTurnResponse(c, 1, "old")
	if c.TurnID != "" {
		t.Fatal("old launch response bound to new turn")
	}
	applyTurnResponse(c, c.generation, "new")
	if c.TurnID != "new" {
		t.Fatal("current response not applied")
	}
}
