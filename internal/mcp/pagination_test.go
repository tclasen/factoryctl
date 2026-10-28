package mcp

import (
	"context"
	"fmt"
	"github.com/tclasen/factoryctl/internal/knowledge"
	"os"
	"path/filepath"
	"testing"
)

func TestMCPPaginationBindsFiltersAndSnapshot(t *testing.T) {
	cfg := fixture(t)
	os.MkdirAll(filepath.Join(cfg.Selected.Root, "factory/records"), 0755)
	for i := 0; i < 205; i++ {
		name := fmt.Sprintf("factory/records/n-%03d.md", i)
		if e := os.WriteFile(filepath.Join(cfg.Selected.Root, name), []byte(fmt.Sprintf("---\ntype: Custom\ntitle: Note %03d\n---\n", i)), 0644); e != nil {
			t.Fatal(e)
		}
	}
	s, e := New(cfg)
	if e != nil {
		t.Fatal(e)
	}
	args := map[string]any{"type": "Custom", "limit": 200}
	first, e := s.Call(context.Background(), "factory_nodes_list", raw(args))
	if e != nil {
		t.Fatal(e)
	}
	a := first.(map[string]any)
	if len(a["nodes"].([]*knowledge.Document)) != 200 || a["next_cursor"] == "" {
		t.Fatal(a)
	}
	args["cursor"] = a["next_cursor"]
	second, e := s.Call(context.Background(), "factory_nodes_list", raw(args))
	if e != nil {
		t.Fatal(e)
	}
	b := second.(map[string]any)
	if len(b["nodes"].([]*knowledge.Document)) != 5 || b["next_cursor"] != "" {
		t.Fatal(b)
	}
	args["query"] = "Note"
	_, e = s.Call(context.Background(), "factory_nodes_list", raw(args))
	code(t, e, "cursor_expired")
	delete(args, "query")
	os.WriteFile(filepath.Join(cfg.Selected.Root, "factory/records/n-000.md"), []byte("---\ntype: Custom\ntitle: External change\n---\n"), 0644)
	_, e = s.Call(context.Background(), "factory_nodes_list", raw(args))
	code(t, e, "cursor_expired")
}
