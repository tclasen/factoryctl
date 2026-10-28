package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/tclasen/factoryctl/internal/app"
	"github.com/tclasen/factoryctl/internal/interaction"
	"github.com/tclasen/factoryctl/internal/knowledge"
	"github.com/tclasen/factoryctl/internal/project"
)

func fixture(t *testing.T) Config {
	t.Helper()
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, ".git/objects"), 0755)
	os.WriteFile(filepath.Join(root, ".git/HEAD"), []byte("ref: refs/heads/main\n"), 0644)
	if _, e := app.Init(root, "demo", "software"); e != nil {
		t.Fatal(e)
	}
	c, e := project.Discover(root, true)
	if e != nil || len(c.Projects) != 1 {
		t.Fatal(e)
	}
	return Config{Projects: c.Projects, Selected: c.Projects[0]}
}
func raw(v any) json.RawMessage { b, _ := json.Marshal(v); return b }
func code(t *testing.T, e error, want string) {
	t.Helper()
	var ae *app.Error
	if !errors.As(e, &ae) || ae.Code != want {
		t.Fatalf("want %s got %v", want, e)
	}
}
func TestScopedDocumentToolsAndStaleWrite(t *testing.T) {
	cfg := fixture(t)
	s, e := New(cfg)
	if e != nil {
		t.Fatal(e)
	}
	ctx := context.Background()
	created, e := s.Call(ctx, "factory_node_create", raw(map[string]any{"type": "Task", "slug": "example", "expected_definition": taskHash(t, cfg), "inputs": map[string]any{"title": "Example", "summary": "A task"}}))
	if e != nil {
		t.Fatal(e)
	}
	doc := created.(*knowledge.Document)
	ref := map[string]any{"path": doc.Ref.Path}
	read, e := s.Call(ctx, "factory_node_read", raw(map[string]any{"ref": ref}))
	if e != nil || read.(map[string]any)["source"].(string) != string(doc.Raw) {
		t.Fatal("exact read", e)
	}
	_, e = s.Call(ctx, "factory_node_read", raw(map[string]any{"ref": map[string]any{"path": doc.Ref.Path, "project_id": "other"}}))
	code(t, e, "invalid_path")
	_, e = s.Call(ctx, "factory_node_update", raw(map[string]any{"ref": ref, "expected_definition": taskHash(t, cfg), "expected_revision": "stale", "set": map[string]any{"title": "Lost update"}}))
	code(t, e, "stale_revision")
	updated, e := s.Call(ctx, "factory_node_update", raw(map[string]any{"ref": ref, "expected_definition": taskHash(t, cfg), "expected_revision": doc.Revision, "set": map[string]any{"title": "Updated"}}))
	if e != nil || updated.(*knowledge.Document).Title != "Updated" {
		t.Fatal("update", e)
	}
	if _, e = s.Call(ctx, "factory_node_update", raw(map[string]any{"ref": ref, "expected_definition": taskHash(t, cfg), "expected_revision": doc.Revision, "set": map[string]any{"title": "Again"}})); e == nil {
		t.Fatal("old revision reused")
	}
	for _, name := range []string{"factory_projects_list", "factory_nodes_list", "factory_definitions_list", "factory_validate"} {
		if _, e = s.Call(ctx, name, raw(map[string]any{})); e != nil {
			t.Fatal(name, e)
		}
	}
}
func TestHeadlessAndNoModelAnswerTool(t *testing.T) {
	s, e := New(Config{})
	if e != nil {
		t.Fatal(e)
	}
	for _, name := range []string{"factory_user_ask", "factory_user_wait", "factory_ui_present"} {
		_, e = s.Call(context.Background(), name, raw(map[string]any{}))
		code(t, e, "interaction_unavailable")
	}
	for _, tool := range Tools() {
		if strings.Contains(tool["name"].(string), "answer") {
			t.Fatal("answer setter exposed")
		}
	}
}
func TestBridgeScopeAndRealSubmissionRoundTrip(t *testing.T) {
	b := interaction.New()
	scope, token, e := b.NewScope("one")
	if e != nil {
		t.Fatal(e)
	}
	other, otherToken, _ := b.NewScope("two")
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		id, ok := b.ScopeForToken(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
		if !ok {
			http.Error(w, "unauthorized", 403)
			return
		}
		var request Request
		if json.NewDecoder(r.Body).Decode(&request) != nil {
			http.Error(w, "bad request", 400)
			return
		}
		result, e := b.ToolCall(r.Context(), id, request.Name, request.Arguments)
		if e != nil {
			var ae *app.Error
			errors.As(e, &ae)
			json.NewEncoder(w).Encode(map[string]any{"error": ae})
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"result": result})
	}))
	defer srv.Close()
	s, e := New(Config{BridgeURL: srv.URL, BridgeToken: token})
	if e != nil {
		t.Fatal(e)
	}
	got, e := s.Call(context.Background(), "factory_user_ask", raw(map[string]any{"prompt": "Real reply?", "kind": "text"}))
	if e != nil {
		t.Fatal(e)
	}
	var q interaction.Question
	json.Unmarshal(got.(json.RawMessage), &q)
	if q.ScopeID != scope || q.Status != "pending" {
		t.Fatal(q)
	}
	second, _ := New(Config{BridgeURL: srv.URL, BridgeToken: otherToken})
	_, e = second.Call(context.Background(), "factory_user_wait", raw(map[string]any{"request_id": q.ID, "timeout_seconds": 0}))
	code(t, e, "not_found")
	if _, e = b.Answer(other, q.ID, interaction.AnswerInput{Status: "answered", Text: "spoof"}); e == nil {
		t.Fatal("spoof")
	}
	b.Answer(scope, q.ID, interaction.AnswerInput{Status: "answered", Text: "actual user text"})
	got, e = s.Call(context.Background(), "factory_user_wait", raw(map[string]any{"request_id": q.ID, "timeout_seconds": 0}))
	if e != nil {
		t.Fatal(e)
	}
	json.Unmarshal(got.(json.RawMessage), &q)
	if q.Text != "actual user text" || q.Status != "answered" {
		t.Fatal(q)
	}
	if calls.Load() != 3 {
		t.Fatal(calls.Load())
	}
}
func TestBridgeRefusesRemoteAndRedirectCredentialLeak(t *testing.T) {
	for _, u := range []string{"https://127.0.0.1/x", "http://example.org/x", "http://localhost/x", "http://user@127.0.0.1/x", "http://127.0.0.1/x?token=x"} {
		if _, e := New(Config{BridgeURL: u, BridgeToken: "secret"}); e == nil {
			t.Fatal("unsafe endpoint", u)
		}
	}
	var reached atomic.Bool
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { reached.Store(true) }))
	defer target.Close()
	redirect := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL, http.StatusTemporaryRedirect)
	}))
	defer redirect.Close()
	s, _ := New(Config{BridgeURL: redirect.URL, BridgeToken: "secret"})
	_, e := s.Call(context.Background(), "factory_user_ask", raw(map[string]any{}))
	code(t, e, "interaction_unavailable")
	if reached.Load() {
		t.Fatal("credential redirected")
	}
}
func TestMCPProtocolAndPendingWaitCancellation(t *testing.T) {
	b := interaction.New()
	scope, token, _ := b.NewScope("one")
	q, _ := b.Ask(scope, interaction.AskInput{Prompt: "Wait", Kind: "text"})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req Request
		json.NewDecoder(r.Body).Decode(&req)
		result, e := b.ToolCall(r.Context(), scope, req.Name, req.Arguments)
		if e != nil {
			json.NewEncoder(w).Encode(map[string]any{"error": &app.Error{Code: "cancelled", Message: e.Error()}})
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"result": result})
	}))
	defer srv.Close()
	inR, inW := io.Pipe()
	outR, outW := io.Pipe()
	done := make(chan error, 1)
	go func() {
		done <- Run(context.Background(), inR, outW, Config{BridgeURL: srv.URL, BridgeToken: token})
		outW.Close()
	}()
	enc := json.NewEncoder(inW)
	dec := json.NewDecoder(outR)
	read := func() map[string]any {
		var v map[string]any
		if e := dec.Decode(&v); e != nil {
			t.Fatal(e)
		}
		return v
	}
	enc.Encode(map[string]any{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": map[string]any{"protocolVersion": "2025-11-25"}})
	if read()["result"] == nil {
		t.Fatal("initialize")
	}
	enc.Encode(map[string]any{"jsonrpc": "2.0", "method": "notifications/initialized"})
	enc.Encode(map[string]any{"jsonrpc": "2.0", "id": 2, "method": "tools/list"})
	if read()["result"] == nil {
		t.Fatal("tools")
	}
	enc.Encode(map[string]any{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": Request{Name: "factory_user_wait", Arguments: raw(map[string]any{"request_id": q.ID, "timeout_seconds": 30})}})
	enc.Encode(map[string]any{"jsonrpc": "2.0", "id": 4, "method": "ping"})
	if read()["id"].(float64) != 4 {
		t.Fatal("wait blocked protocol")
	}
	enc.Encode(map[string]any{"jsonrpc": "2.0", "method": "notifications/cancelled", "params": map[string]any{"requestId": 3}})
	if response := read(); response["id"].(float64) != 3 || response["result"].(map[string]any)["isError"] != true {
		t.Fatal(response)
	}
	inW.Close()
	select {
	case e := <-done:
		if e != nil {
			t.Fatal(e)
		}
	case <-time.After(time.Second):
		t.Fatal("EOF did not terminate")
	}
	got, _ := b.Wait(context.Background(), scope, q.ID, 0)
	if got.Status != "pending" {
		t.Fatal("RPCcancel synthesized disposition")
	}
}
func TestMCPMalformedMessagesFailWithoutProtocolNoise(t *testing.T) {
	var out bytes.Buffer
	input := strings.NewReader("not json\n{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/list\"}\n{\"jsonrpc\":\"2.0\",\"id\":2,\"method\":\"initialize\"}\n{\"jsonrpc\":\"2.0\",\"id\":3,\"method\":\"factory_user_answer\"}\n")
	if e := Run(context.Background(), input, &out, Config{}); e != nil {
		t.Fatal(e)
	}
	dec := json.NewDecoder(&out)
	for i, want := range []int{-32700, -32002, 0, -32601} {
		var r map[string]any
		if e := dec.Decode(&r); e != nil {
			t.Fatal(e)
		}
		if want != 0 && int(r["error"].(map[string]any)["code"].(float64)) != want {
			t.Fatal(i, r)
		}
	}
	var tail any
	if e := dec.Decode(&tail); e != io.EOF {
		t.Fatal("unexpected stdout", e)
	}
}

type failingWriter struct{}

func (failingWriter) Write(p []byte) (int, error) { return 0, io.ErrClosedPipe }
func TestMCPOutputFailureIsReported(t *testing.T) {
	e := Run(context.Background(), strings.NewReader("{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"initialize\"}\n"), failingWriter{}, Config{})
	if !errors.Is(e, io.ErrClosedPipe) {
		t.Fatalf("output failure lost: %v", e)
	}
}

// Codex attaches standard request metadata, including progressToken. Metadata
// belongs to the protocol envelope, not tool arguments or destination authority.
func TestMCPToolCallAcceptsMetadataWithoutForwardingIt(t *testing.T) {
	seen := make(chan map[string]json.RawMessage, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]json.RawMessage
		if e := json.NewDecoder(r.Body).Decode(&body); e != nil {
			t.Error(e)
		}
		seen <- body
		json.NewEncoder(w).Encode(map[string]any{"result": map[string]any{"status": "pending", "id": "request-1"}})
	}))
	defer srv.Close()
	inR, inW := io.Pipe()
	outR, outW := io.Pipe()
	done := make(chan error, 1)
	go func() {
		done <- Run(context.Background(), inR, outW, Config{BridgeURL: srv.URL, BridgeToken: "tool-only-token"})
		outW.Close()
	}()
	enc := json.NewEncoder(inW)
	dec := json.NewDecoder(outR)
	enc.Encode(map[string]any{"jsonrpc": "2.0", "id": 1, "method": "initialize"})
	var init map[string]any
	if e := dec.Decode(&init); e != nil {
		t.Fatal(e)
	}
	enc.Encode(map[string]any{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": map[string]any{"name": "factory_user_ask", "arguments": map[string]any{"prompt": "Choice?", "kind": "text"}, "_meta": map[string]any{"progressToken": "progress-1", "scope_id": "untrusted", "client": map[string]any{"name": "Codex"}}}})
	var response map[string]any
	if e := dec.Decode(&response); e != nil {
		t.Fatal(e)
	}
	if response["error"] != nil || response["result"].(map[string]any)["isError"] == true {
		t.Fatal(response)
	}
	body := <-seen
	if len(body) != 2 || body["_meta"] != nil || body["scope_id"] != nil || body["name"] == nil || body["arguments"] == nil {
		t.Fatalf("metadata crossed tool bridge: %#v", body)
	}
	var args map[string]any
	json.Unmarshal(body["arguments"], &args)
	if len(args) != 2 || args["prompt"] != "Choice?" || args["kind"] != "text" {
		t.Fatal(args)
	}
	inW.Close()
	if e := <-done; e != nil {
		t.Fatal(e)
	}
}

func taskHash(t *testing.T, cfg Config) string {
	t.Helper()
	snap, e := app.Load(cfg.Projects)
	if e != nil {
		t.Fatal(e)
	}
	return app.DefinitionHash(snap.Bundle(cfg.Selected.Checkout).Definitions["Task"])
}

func TestDynamicExtensionToolsAndDefinitionPinnedTransport(t *testing.T) {
	cfg := fixture(t)
	def := `---
type: Record Definition
title: Sample
defines: Sample
directory: samples
fields:
  state: {kind: enum, required: true, values: [new, done]}
extension:
  version: 1
  namespace: samples
  tools:
    inspect: {operation: read}
    create:
      operation: create
      inputs:
        title: {kind: string, required: true}
      fields:
        title: {input: title}
        state: {literal: new}
---
# Instructions
The harness interprets this prose.
`
	file := filepath.Join(cfg.Selected.Root, "factory/definitions/sample.md")
	if e := os.WriteFile(file, []byte(def), 0644); e != nil {
		t.Fatal(e)
	}
	s, e := New(cfg)
	if e != nil {
		t.Fatal(e)
	}
	has := func(name string) bool {
		ts, e := s.Tools()
		if e != nil {
			t.Fatal(e)
		}
		for _, v := range ts {
			if v["name"] == name {
				return true
			}
		}
		return false
	}
	if !has("samples_create") {
		t.Fatal("runtime tool absent")
	}
	registry, e := app.ExtensionCatalog(cfg.Projects, cfg.Selected)
	if e != nil {
		t.Fatal(e)
	}
	hash := ""
	for _, v := range registry.Tools {
		if v.Name == "samples_create" {
			hash = v.DefinitionHash
		}
	}
	result, e := s.Call(context.Background(), "samples_create", raw(map[string]any{"slug": "one", "expected_definition": hash, "inputs": map[string]any{"title": "One"}}))
	if e != nil {
		t.Fatal(e)
	}
	d := result.(*knowledge.Document)
	_, e = s.Call(context.Background(), "factory_node_update", raw(map[string]any{"ref": d.Ref, "expected_revision": d.Revision, "set": map[string]any{"state": "done"}}))
	code(t, e, "invalid_arguments")
	if e = os.WriteFile(file, []byte(def+"New instructions.\n"), 0644); e != nil {
		t.Fatal(e)
	}
	_, e = s.Call(context.Background(), "factory_node_update", raw(map[string]any{"ref": d.Ref, "expected_revision": d.Revision, "expected_definition": hash, "set": map[string]any{"state": "done"}}))
	code(t, e, "stale_definition")
	updated := strings.Replace(def, "namespace: samples", "namespace: changed", 1)
	if e = os.WriteFile(file, []byte(updated), 0644); e != nil {
		t.Fatal(e)
	}
	if has("samples_create") || !has("changed_create") {
		t.Fatal("same running server did not reload tool list")
	}
	if e = os.WriteFile(file, []byte(strings.Replace(updated, "fields:\n", "enabled: false\nfields:\n", 1)), 0644); e != nil {
		t.Fatal(e)
	}
	if has("changed_create") {
		t.Fatal("disabled type still advertised")
	}
	if _, e = s.Call(context.Background(), "factory_node_read", raw(map[string]any{"ref": d.Ref})); e != nil {
		t.Fatal("generic disabled read failed", e)
	}
}
