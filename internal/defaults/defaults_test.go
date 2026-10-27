package defaults

import (
	"github.com/tclasen/factoryctl/internal/knowledge"
	"strings"
	"testing"
)

func TestDefaultsAreParseableAndIsolated(t *testing.T) {
	for _, software := range []bool{false, true} {
		files := Files(software, "demo", "Demo")
		if _, ok := files["control/write-coordination.md"]; !ok {
			t.Fatal("missing lock anchor")
		}
		for p, b := range files {
			d := knowledge.Parse(p, b)
			if len(d.Diagnostics) > 0 {
				t.Fatalf("%s: %+v", p, d.Diagnostics)
			}
			if strings.Contains(string(b), "operator-brief") || strings.Contains(string(b), "../specifications/") {
				t.Fatalf("defaults leak source project links: %s", p)
			}
		}
		_, hasSoftware := files["definitions/specification.md"]
		if hasSoftware != software {
			t.Fatal("software default selection leaked")
		}
		if d := knowledge.Parse("factory.md", files["factory.md"]); d.Frontmatter["project_id"] != "demo" {
			t.Fatal("wrong identity")
		}
	}
}
func TestDefaultsFreshMaps(t *testing.T) {
	a := Files(false, "a", "A")
	delete(a, "factory.md")
	if _, ok := Files(false, "b", "B")["factory.md"]; !ok {
		t.Fatal("shared mutable defaults")
	}
}
