package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLIExtensionParityAndPromptPreview(t *testing.T) {
	root := initialized(t)
	def := `---
type: Record Definition
title: Research definition
defines: Research
directory: research
fields:
  title: {kind: string, required: true}
  phase: {kind: enum, required: true, values: [open, reviewed]}
extension:
  version: 1
  namespace: research
  tools:
    create:
      operation: create
      description: Create a research note
      inputs:
        title: {kind: string, required: true}
      fields:
        title: {input: title}
        phase: {literal: open}
    read:
      operation: read
      description: Read a research note
---
# Research
The harness interprets review requirements.
`
	if e := os.WriteFile(filepath.Join(root, "factory/definitions/research.md"), []byte(def), 0644); e != nil {
		t.Fatal(e)
	}
	catalog := cli(t, root, "", "extensions", "list")
	wantCode(t, catalog, 0)
	tools := catalog.value["tools"].([]any)
	if len(tools) != 2 {
		t.Fatal(catalog.stdout)
	}
	hash := tools[0].(map[string]any)["definition_hash"].(string)
	created := cli(t, root, "", "extensions", "call", "research_create", "--input", "slug=note", "--input", "expected_definition="+hash, "--input", `inputs={"title":"A sourced note"}`)
	wantCode(t, created, 0)
	read := cli(t, root, "", "extensions", "call", "research_read", "--input", "path=research/note.md")
	wantCode(t, read, 0)
	if !strings.Contains(read.stdout, "A sourced note") {
		t.Fatal(read.stdout)
	}
	wantError(t, cli(t, root, "", "extensions", "call", "research_create", "--input", "slug=other", "--input", "expected_definition=stale", "--input", `inputs={"title":"Other"}`), 3, "stale_definition")
	preview := cli(t, root, "", "actions", "render", "prompts/discuss.md", "--record", "research/note.md")
	wantCode(t, preview, 0)
	if preview.value["launched"] != false || !strings.Contains(preview.value["prompt"].(string), "A sourced note") {
		t.Fatal(preview.stdout)
	}
	show := cli(t, root, "", "records", "show", "research/note.md")
	wantCode(t, show, 0)
	if show.value["definition_hash"] != hash {
		t.Fatal(show.stdout)
	}
	// An independent definition edit invalidates the caller's earlier draft.
	if e := os.WriteFile(filepath.Join(root, "factory/definitions/research.md"), []byte(def+"\nChanged instructions.\n"), 0644); e != nil {
		t.Fatal(e)
	}
	wantError(t, cli(t, root, "", "records", "update", "research/note.md", "--expected-revision", created.value["revision"].(string), "--expected-definition", hash, "--set", "phase=reviewed"), 3, "stale_definition")
	if _, e := os.Stat(filepath.Join(root, ".factory")); !os.IsNotExist(e) {
		t.Fatal("preview wrote runtime state")
	}
	var parsed any
	if json.Unmarshal([]byte(catalog.stdout), &parsed) != nil {
		t.Fatal("invalid catalog JSON")
	}
}

func TestCLITemplateExportPreservesCustomization(t *testing.T) {
	root := initialized(t)
	custom := []byte("customized template, never overwrite\n")
	file := filepath.Join(root, "factory/templates/task.md")
	if e := os.WriteFile(file, custom, 0644); e != nil {
		t.Fatal(e)
	}
	wantCode(t, cli(t, root, "", "templates", "export", "--kind", "record"), 0)
	got, e := os.ReadFile(file)
	if e != nil || string(got) != string(custom) {
		t.Fatal("overwrote customization", e)
	}
	r := cli(t, root, "", "templates", "export", "--kind", "view")
	wantCode(t, r, 0)
	if len(r.value["created"].([]any)) != 1 {
		t.Fatal(r.stdout)
	}
	wantCode(t, cli(t, root, "", "templates", "export", "--kind", "view"), 0)
	root2 := initialized(t)
	if e := os.Symlink(t.TempDir(), filepath.Join(root2, ".factory")); e != nil {
		t.Fatal(e)
	}
	wantError(t, cli(t, root2, "", "templates", "export", "--kind", "view"), 2, "invalid_path")
}
