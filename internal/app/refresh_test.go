package app

import (
	"github.com/tclasen/factoryctl/internal/knowledge"
	"github.com/tclasen/factoryctl/internal/templates"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
)

func TestCachedRefreshChecksBytesAndRevalidatesGraph(t *testing.T) {
	p, ps := storeFixture(t)
	d, e := Mutate(ps, p, Mutation{Kind: "create", Path: "records/source.md", Document: []byte("---\ntype: Custom\ntitle: Source\n---\n# Anchor\n[self](#anchor)\n")})
	if e != nil {
		t.Fatal(e)
	}
	before, e := Load(ps)
	if e != nil {
		t.Fatal(e)
	}
	file := filepath.Join(p.Root, "factory", d.Ref.Path)
	st, _ := os.Stat(file)
	// Same length and restored timestamp must not hide a content/anchor change.
	raw := []byte("---\ntype: Custom\ntitle: Source\n---\n# Absent\n[self](#anchor)\n")
	if e = os.WriteFile(file, raw, 0644); e != nil {
		t.Fatal(e)
	}
	os.Chtimes(file, st.ModTime(), st.ModTime())
	after, e := LoadCached(ps, before)
	if e != nil {
		t.Fatal(e)
	}
	if before.Revision == after.Revision || after.Validate("factory").Errors == 0 {
		t.Fatal("cached graph concealed changed bytes/anchor")
	}
	if before.Bundle(p.Checkout).Docs[d.Ref.Path].Revision != d.Revision {
		t.Fatal("mutated prior snapshot")
	}
	os.Remove(file)
	removed, e := LoadCached(ps, after)
	if e != nil || removed.Bundle(p.Checkout).Docs[d.Ref.Path] != nil {
		t.Fatal("retained deleted document", e)
	}
}

func TestTemplateSourceAndAssetPinsInsideWriter(t *testing.T) {
	p, ps := storeFixture(t)
	snapshot, e := Load(ps)
	if e != nil {
		t.Fatal(e)
	}
	os.MkdirAll(filepath.Join(p.Root, ".factory/templates"), 0755)
	asset := ".factory/templates/body.tmpl"
	file := filepath.Join(p.Root, asset)
	os.WriteFile(file, []byte("original"), 0644)
	pin, e := templates.AssetVersion(p.Root, asset)
	if e != nil {
		t.Fatal(e)
	}
	os.WriteFile(file, []byte("changed"), 0644)
	candidate := []byte("---\ntype: Custom\ntitle: Candidate\n---\n")
	_, e = Mutate(ps, p, Mutation{Kind: "create", Path: "records/asset-pinned.md", Document: candidate, ExpectedSnapshot: snapshot.Revision, ExpectedAssets: map[string]string{asset: pin}})
	if e == nil {
		t.Fatal("accepted changed template asset")
	}
	if _, e = os.Stat(filepath.Join(p.Root, "factory/records/asset-pinned.md")); !os.IsNotExist(e) {
		t.Fatal("rejection wrote bytes")
	}
	os.WriteFile(filepath.Join(p.Root, "factory/templates/task.md"), []byte("changed template"), 0644)
	_, e = Mutate(ps, p, Mutation{Kind: "create", Path: "records/source-pinned.md", Document: candidate, ExpectedSnapshot: snapshot.Revision})
	if e == nil {
		t.Fatal("accepted stale template source graph")
	}
	if _, e = os.Stat(filepath.Join(p.Root, "factory/records/source-pinned.md")); !os.IsNotExist(e) {
		t.Fatal("rejection wrote bytes")
	}
}

func TestCachedDefinitionRemovalMatchesColdLoad(t *testing.T) {
	p, ps := storeFixture(t)
	file := filepath.Join(p.Root, "factory/definitions/custom.md")
	os.WriteFile(file, []byte("---\ntype: Record Definition\ntitle: Custom\ndefines: Custom\ndirectory: records\nfields: {}\n---\n"), 0644)
	if _, e := Mutate(ps, p, Mutation{Kind: "create", Path: "records/custom.md", Document: []byte("---\ntype: Custom\ntitle: Example\n---\n")}); e != nil {
		t.Fatal(e)
	}
	prior, e := Load(ps)
	if e != nil {
		t.Fatal(e)
	}
	if prior.Bundle(p.Checkout).Docs["records/custom.md"].Schema != "available" {
		t.Fatal("fixture schema unavailable")
	}
	os.Remove(file)
	cached, e := LoadCached(ps, prior)
	if e != nil {
		t.Fatal(e)
	}
	cold, e := Load(ps)
	if e != nil {
		t.Fatal(e)
	}
	if cached.Bundle(p.Checkout).Docs["records/custom.md"].Schema != "unavailable" || cached.Bundle(p.Checkout).Docs["records/custom.md"].Schema != cold.Bundle(p.Checkout).Docs["records/custom.md"].Schema {
		t.Fatal("cached schema differs from cold load")
	}
	if prior.Bundle(p.Checkout).Docs["records/custom.md"].Schema != "available" {
		t.Fatal("mutated old snapshot")
	}
}

func TestLiteralCaseInsensitiveSearch(t *testing.T) {
	for _, c := range []struct {
		query, body string
		matches     bool
	}{{"SOURCE", "A source citation", true}, {"[a-z]", "literal [A-Z] field", true}, {"[a-z]", "ordinary letter c", false}, {"σ", "Greek Σ character", true}, {"", "anything", true}} {
		if got := QueryMatcher(c.query)([]byte(c.body)); got != c.matches {
			t.Fatalf("%q in %q: %v", c.query, c.body, got)
		}
	}
}

func TestDocumentMatcherMatchesUnicodeReference(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	alphabet := []byte("aABbcXYZ []*.012\n")
	for i := 0; i < 300; i++ {
		body := make([]byte, r.Intn(300))
		query := make([]byte, r.Intn(12))
		for j := range body {
			body[j] = alphabet[r.Intn(len(alphabet))]
		}
		for j := range query {
			query[j] = alphabet[r.Intn(len(alphabet))]
		}
		d := knowledge.Parse("record.md", append([]byte("---\ntype: Custom\ntitle: T\n---\n"), body...))
		if got, want := DocumentMatcher(string(query))(d), QueryMatcher(string(query))(d.Raw); got != want {
			t.Fatalf("query=%q body=%q got=%v want=%v", query, body, got, want)
		}
	}
	for _, q := range []string{"K", "σ", "é"} {
		d := knowledge.Parse("unicode.md", []byte("---\ntype: Custom\ntitle: K Σ É\n---\n"))
		if !DocumentMatcher(q)(d) {
			t.Fatalf("Unicode fold failed %q", q)
		}
	}
}
