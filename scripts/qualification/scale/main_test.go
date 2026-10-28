package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tclasen/factoryctl/internal/app"
	"github.com/tclasen/factoryctl/internal/knowledge"
	"github.com/tclasen/factoryctl/internal/project"
)

func TestSyntheticKnowledgeWorkCorpus(t *testing.T) {
	root := filepath.Join(t.TempDir(), "corpus")
	if err := corpus(root, 2, 20, 65536); err != nil {
		t.Fatal(err)
	}
	c, err := project.Discover(root, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Projects) != 2 || len(c.Diagnostics) != 0 {
		t.Fatalf("discovery: %#v", c)
	}
	s, err := app.Load(c.Projects)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Diagnostics) != 0 {
		t.Fatalf("diagnostics: %#v", s.Diagnostics)
	}
	total, count := 0, 0
	types := map[string]bool{}
	for _, b := range s.Bundles {
		for _, d := range b.Docs {
			total += len(d.Raw)
			count++
			types[d.Type] = true
		}
	}
	if total != 65536 || count != 20 || !types["Research Fact"] || !types["Opportunity"] {
		t.Fatalf("bytes %d nodes %d types %v", total, count, types)
	}
	resolved, unchecked := 0, 0
	for _, link := range s.Links {
		if link.Status == "valid" {
			resolved++
		}
		if link.Status == "unchecked" {
			unchecked++
		}
	}
	if resolved != 20 || unchecked != 10 {
		t.Fatalf("expected two resolved Markdown links and one unchecked external source per domain node: %d resolved, %d unchecked; %#v", resolved, unchecked, s.Links[0])
	}
	if err := corpus(root, 2, 20, 65536); err == nil {
		t.Fatal("must preserve existing attempts")
	}
	path := filepath.Join(root, "repository-000/factory/opportunities/bid-001.md")
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	before = append(before, []byte("\nExternal edit.\n")...)
	if err = os.WriteFile(path, before, 0644); err != nil {
		t.Fatal(err)
	}
	next, err := app.Load(c.Projects)
	if err != nil {
		t.Fatal(err)
	}
	if next.Revision == s.Revision {
		t.Fatal("external edit did not change snapshot")
	}
}

func TestResearchOpportunityJourney(t *testing.T) {
	root := filepath.Join(t.TempDir(), "breadth")
	if err := corpus(root, 1, 10, 16384); err != nil {
		t.Fatal(err)
	}
	c, err := project.Discover(root, true)
	if err != nil {
		t.Fatal(err)
	}
	p := c.Projects[0]
	actions, err := app.Actions(c.Projects, p, "research/fact-000.md")
	if err != nil {
		t.Fatal(err)
	}
	if len(actions) != 1 || actions[0]["disabled_reason"] != "" {
		t.Fatalf("applicable evidence review action: %#v", actions)
	}
	preview, err := app.RenderPrompt(c.Projects, p, "prompts/review-evidence.md", "research/fact-000.md", nil)
	if err != nil {
		t.Fatal(err)
	}
	if preview["launched"] != false || !strings.Contains(preview["prompt"].(string), "Separate cited synthetic claims") {
		t.Fatalf("prompt preview: %#v", preview)
	}
	s, err := app.Load(c.Projects)
	if err != nil {
		t.Fatal(err)
	}
	ref := "opportunities/bid-001.md"
	d := s.Bundles[0].Docs[ref]
	updated, err := app.Mutate(c.Projects, p, app.Mutation{Kind: "update", Path: ref, ExpectedRevision: d.Revision, Set: map[string]any{"stage": "submitted"}})
	if err != nil {
		t.Fatal(err)
	}
	if knowledge.String(updated.Frontmatter, "stage") != "submitted" {
		t.Fatal("pipeline field did not change")
	}
	path := filepath.Join(p.Root, "factory", ref)
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, mutation := range []app.Mutation{
		{Kind: "update", Path: ref, ExpectedRevision: d.Revision, Set: map[string]any{"stage": "tracking"}},
		{Kind: "update", Path: ref, ExpectedRevision: updated.Revision, Set: map[string]any{"stage": "invented-stage"}},
	} {
		if _, err := app.Mutate(c.Projects, p, mutation); err == nil {
			t.Fatal("accepted stale or structurally invalid update")
		}
		after, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(before, after) {
			t.Fatal("failed mutation changed bytes")
		}
	}
	s, err = app.Load(c.Projects)
	if err != nil {
		t.Fatal(err)
	}
	research := s.Bundles[0].Docs["research/fact-000.md"]
	if !bytes.Contains(research.Raw, []byte("LLM inference")) || !bytes.Contains(research.Raw, []byte("Assumption:")) || !bytes.Contains(research.Raw, []byte("sources:")) {
		t.Fatal("research attribution lost")
	}
	backlinks := 0
	for _, link := range s.Links {
		if link.Resolved != nil && link.Resolved.Path == research.Ref.Path {
			backlinks++
		}
	}
	if backlinks != 2 {
		t.Fatalf("research backlinks: %d", backlinks)
	}
	created, err := knowledge.Encode(map[string]any{"type": "Research Fact", "title": "Synthetic follow-up", "certainty": "inference"}, "\n# Evidence\n\nLLM inference from an invented source, not a verified fact.\n")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := app.Mutate(c.Projects, p, app.Mutation{Kind: "create", Path: "research/follow-up.md", Document: created}); err != nil {
		t.Fatal(err)
	}
}

func TestCorpusRejectsInvalidScale(t *testing.T) {
	for _, args := range [][3]int{{0, 10, 1000}, {2, 11, 10000}, {2, 6, 10000}, {2, 20, 10}} {
		if err := corpus(filepath.Join(t.TempDir(), "corpus"), args[0], args[1], args[2]); err == nil {
			t.Fatalf("accepted %v", args)
		}
	}
}
