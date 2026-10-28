package app

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tclasen/factoryctl/internal/knowledge"
	"github.com/tclasen/factoryctl/internal/project"
)

func extensionDefinition(typ, ns, dir string) string {
	return "---\ntype: Record Definition\ntitle: " + typ + "\ndefines: " + typ + "\ndirectory: " + dir + "\nfields:\n  state:\n    kind: enum\n    required: true\n    values: [new, done]\n  evidence:\n    kind: string\n    required: true\nextension:\n  version: 1\n  namespace: " + ns + "\n  tools:\n    read:\n      operation: read\n    create:\n      operation: create\n      inputs:\n        title: {kind: string, required: true}\n        evidence: {kind: string, required: true}\n      fields:\n        title: {input: title}\n        evidence: {input: evidence}\n        state: {literal: new}\n    move:\n      operation: update\n      inputs:\n        state: {kind: enum, required: true, values: [new, done]}\n      fields:\n        state: {input: state}\n    delete:\n      operation: delete\n---\n# Instructions\nNever move to done without agent review. This prose is not executable.\n"
}
func extensionPut(t *testing.T, p project.Project, name, source string) {
	t.Helper()
	file := filepath.Join(p.Root, "factory", name)
	if e := os.MkdirAll(filepath.Dir(file), 0755); e != nil {
		t.Fatal(e)
	}
	if e := os.WriteFile(file, []byte(source), 0644); e != nil {
		t.Fatal(e)
	}
}
func extensionCode(t *testing.T, err error, want string) {
	t.Helper()
	var e *Error
	if !errors.As(err, &e) || e.Code != want {
		t.Fatalf("expected %s, got %v", want, err)
	}
}
func extensionTool(t *testing.T, ps []project.Project, p project.Project, name string) ExtensionTool {
	t.Helper()
	r, e := ExtensionCatalog(ps, p)
	if e != nil {
		t.Fatal(e)
	}
	for _, v := range r.Tools {
		if v.Name == name {
			return v
		}
	}
	t.Fatalf("missing tool %s: %#v", name, r)
	return ExtensionTool{}
}
func TestRuntimeExtensionsReloadAndSharedStructuralBoundary(t *testing.T) {
	p, ps := storeFixture(t)
	extensionPut(t, p, "definitions/research.md", extensionDefinition("Research", "research", "research"))
	opportunity := strings.Replace(extensionDefinition("Opportunity", "bids", "opportunities"), "Never move to done without agent review.", "The agent chooses award readiness from solicitation evidence.", 1)
	extensionPut(t, p, "definitions/opportunity.md", opportunity)
	for _, name := range []string{"research_create", "bids_create"} {
		extensionTool(t, ps, p, name)
	}
	bidCreate := extensionTool(t, ps, p, "bids_create")
	bidResult, e := CallExtension(ps, p, bidCreate.Name, map[string]any{"slug": "solicitation", "expected_definition": bidCreate.DefinitionHash, "inputs": map[string]any{"title": "Solicitation", "evidence": "Public notice"}})
	if e != nil {
		t.Fatal(e)
	}
	bid := bidResult.(*knowledge.Document)
	_, e = CallExtension(ps, p, "bids_delete", map[string]any{"path": bid.Ref.Path, "expected_revision": bid.Revision, "expected_definition": bidCreate.DefinitionHash})
	if e != nil {
		t.Fatal(e)
	}
	if _, e = os.Stat(filepath.Join(p.Root, "factory", bid.Ref.Path)); !os.IsNotExist(e) {
		t.Fatal("specialized delete did not remove record")
	}
	create := extensionTool(t, ps, p, "research_create")
	result, e := CallExtension(ps, p, create.Name, map[string]any{"slug": "study", "expected_definition": create.DefinitionHash, "inputs": map[string]any{"title": "A study", "evidence": "Source attribution"}})
	if e != nil {
		t.Fatal(e)
	}
	d := result.(*knowledge.Document)
	read, e := CallExtension(ps, p, "research_read", map[string]any{"path": d.Ref.Path})
	if e != nil {
		t.Fatal(e)
	}
	if !strings.Contains(read.(map[string]any)["instructions"].(string), "agent review") {
		t.Fatal("instruction prose missing")
	}
	move := extensionTool(t, ps, p, "research_move")
	moved, e := CallExtension(ps, p, move.Name, map[string]any{"path": d.Ref.Path, "expected_revision": d.Revision, "expected_definition": move.DefinitionHash, "inputs": map[string]any{"state": "done"}})
	if e != nil {
		t.Fatal("workflow prose became a gate", e)
	}
	d = moved.(*knowledge.Document)
	before := append([]byte(nil), d.Raw...)
	_, e = CallExtension(ps, p, move.Name, map[string]any{"path": d.Ref.Path, "expected_revision": d.Revision, "expected_definition": move.DefinitionHash, "inputs": map[string]any{"state": "approved"}})
	extensionCode(t, e, "invalid_arguments")
	_, e = Mutate(ps, p, Mutation{Kind: "replace", Path: d.Ref.Path, ExpectedRevision: d.Revision, ExpectedDefinition: move.DefinitionHash, Document: []byte(strings.Replace(string(d.Raw), "evidence:", "missing:", 1))})
	extensionCode(t, e, "invalid_document")
	after, _ := os.ReadFile(filepath.Join(p.Root, "factory", d.Ref.Path))
	if string(after) != string(before) {
		t.Fatal("rejected operations changed bytes")
	}
	_, e = CallExtension(ps, p, "bids_read", map[string]any{"path": d.Ref.Path})
	extensionCode(t, e, "invalid_path")
	extensionPut(t, p, "definitions/contract.md", extensionDefinition("Contract", "contracts", "contracts"))
	extensionTool(t, ps, p, "contracts_create")
	// A body-only external definition edit changes identity without changing document revision.
	extensionPut(t, p, "definitions/research.md", extensionDefinition("Research", "research", "research")+"Changed instructions.\n")
	_, e = Mutate(ps, p, Mutation{Kind: "update", Path: d.Ref.Path, ExpectedRevision: d.Revision, ExpectedDefinition: move.DefinitionHash, Set: map[string]any{"state": "new"}})
	extensionCode(t, e, "stale_definition")
	_, e = CreateChecked(ps, p, "Research", "stale", map[string]any{"title": "Stale", "state": "new", "evidence": "x"}, create.DefinitionHash)
	extensionCode(t, e, "stale_definition")
	fresh := extensionTool(t, ps, p, "research_move")
	_, e = CallExtension(ps, p, fresh.Name, map[string]any{"path": d.Ref.Path, "expected_revision": "stale", "expected_definition": fresh.DefinitionHash, "inputs": map[string]any{"state": "new"}})
	extensionCode(t, e, "stale_revision")
	extensionPut(t, p, "definitions/research.md", strings.Replace(extensionDefinition("Research", "research", "research"), "fields:\n", "enabled: false\nfields:\n", 1))
	r, _ := ExtensionCatalog(ps, p)
	for _, v := range r.Tools {
		if strings.HasPrefix(v.Name, "research_") {
			t.Fatal("disabled tools advertised")
		}
	}
	snap, _ := Load(ps)
	if snap.Bundle(p.Checkout).Docs[d.Ref.Path] == nil {
		t.Fatal("disabled type became unreadable")
	}
	generic, e := Mutate(ps, p, Mutation{Kind: "update", Path: d.Ref.Path, ExpectedRevision: d.Revision, Set: map[string]any{"state": "new"}})
	if e != nil {
		t.Fatal("disabled type generic schema write failed", e)
	}
	d = generic
	// Duplicate namespaces quarantine affected types, preserving unrelated tools.
	extensionPut(t, p, "definitions/research.md", extensionDefinition("Research", "bids", "research"))
	r, _ = ExtensionCatalog(ps, p)
	if len(r.Tools) != 4 {
		t.Fatalf("expected only unrelated contract tools, got %#v", r.Tools)
	}
	_, e = Mutate(ps, p, Mutation{Kind: "update", Path: d.Ref.Path, ExpectedRevision: d.Revision, Set: map[string]any{"state": "done"}})
	if e != nil {
		t.Fatal("namespace conflict blocked generic schema write", e)
	}
	if _, e = os.Stat(filepath.Join(p.Root, ".factory")); !os.IsNotExist(e) {
		t.Fatal("definition-only extensions required assets")
	}
}
func TestExtensionRejectsExecutableBindingsAndReservedMappings(t *testing.T) {
	for _, replacement := range []string{"handler: 'sh -c echo bad'", "expression: '.Input.state'", "transitions: [new, done]", "operation: exec", "operation: read\n      predicate: ready"} {
		src := strings.Replace(extensionDefinition("Research", "research", "research"), "operation: read", replacement, 1)
		d := knowledge.Parse("definitions/research.md", []byte(src))
		def, ds := knowledge.ParseDefinition(d)
		if !def.ExtensionInvalid || len(ds) == 0 {
			t.Fatalf("accepted executable/policy binding %s", replacement)
		}
	}
	for _, change := range []struct{ old, new string }{{"namespace: research", "namespace: factory"}, {"title: {kind: string, required: true}", "path: {kind: string, required: true}"}, {"title: {input: title}", "type: {literal: Task}"}, {"state: {literal: new}", "state: {expression: new}"}} {
		d := knowledge.Parse("definitions/research.md", []byte(strings.Replace(extensionDefinition("Research", "research", "research"), change.old, change.new, 1)))
		def, _ := knowledge.ParseDefinition(d)
		if !def.ExtensionInvalid {
			t.Fatalf("accepted forbidden binding %s", change.new)
		}
	}
}
func TestPromptPreviewNeverLaunchesAndExposesContext(t *testing.T) {
	p, ps := storeFixture(t)
	d, e := Create(ps, p, "Task", "one", map[string]any{"title": "One", "summary": "Agent-owned workflow"})
	if e != nil {
		t.Fatal(e)
	}
	extensionPut(t, p, "prompts/review.md", "---\ntype: Prompt\ntitle: Review\napplies_to: [Task]\ninputs:\n  focus: {kind: string, required: true}\nconversation: fresh\n---\n# Template\n\n```factory-prompt\nReview {{.Record.Title}} and decide next state using {{.Input.focus}}.\n```\n")
	if _, e := RenderPrompt(ps, p, "prompts/discuss.md", "", map[string]any{}); e != nil {
		t.Fatal("recordless embedded prompt preview", e)
	}
	actions, e := Actions(ps, p, d.Ref.Path)
	if e != nil || len(actions) == 0 {
		t.Fatal("no prompt actions", e)
	}
	preview, e := RenderPrompt(ps, p, "prompts/review.md", d.Ref.Path, map[string]any{"focus": "evidence"})
	if e != nil {
		t.Fatal(e)
	}
	if preview["launched"] != false || !strings.Contains(preview["prompt"].(string), "Review One") {
		t.Fatalf("bad preview %#v", preview)
	}
	_, e = RenderPrompt(ps, p, "prompts/review.md", d.Ref.Path, map[string]any{"focus": "x", "undeclared": "y"})
	extensionCode(t, e, "invalid_document")
}
