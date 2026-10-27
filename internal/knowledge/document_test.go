package knowledge

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseHazardsAndLimits(t *testing.T) {
	tests := map[string]string{
		"duplicate": "---\ntype: Task\ntype: Decision\n---\n", "nested duplicate": "---\ntype: X\nx: {a: 1, a: 2}\n---\n", "alias": "---\ntype: X\nx: &a value\ny: *a\n---\n", "tag": "---\ntype: !execute foo\n---\n", "key": "---\ntype: X\n1: x\n---\n", "multiple": "---\ntype: X\n...\ntype: Y\n---\n", "utf8": "---\ntype: X\n---\n\xff", "conflict": "---\ntype: X\n---\n<<<<<<< HEAD\n", "missing": "type: X", "required": "---\ntype: ''\n---\n", "large": strings.Repeat("x", MaxDocument+1), "deep": "---\ntype: X\nx: " + strings.Repeat("[", 40) + "1" + strings.Repeat("]", 40) + "\n---\n"}
	for name, source := range tests {
		t.Run(name, func(t *testing.T) {
			if len(Parse("x.md", []byte(source)).Diagnostics) == 0 {
				t.Fatal("accepted invalid document")
			}
		})
	}
}
func TestLosslessScalarEdit(t *testing.T) {
	raw := []byte("---\n# preserved\ntype: Task\ntitle: 'héllo' # retained\nunknown:\n  nested: [1, 2] # untouched\nstatus: draft\n---\n\nExact body with {{ template-like text }}\n")
	d := Parse("x.md", raw)
	got, e := d.Edit(map[string]any{"title": "new \"title\"", "status": "stable"})
	if e != nil {
		t.Fatal(e)
	}
	want := strings.Replace(strings.Replace(string(raw), "'héllo'", "\"new \\\"title\\\"\"", 1), "status: draft", "status: \"stable\"", 1)
	if string(got) != want {
		t.Fatalf("unexpected source edit:\n%s", got)
	}
	if !bytes.Equal(d.Raw, raw) {
		t.Fatal("mutated source buffer")
	}
	if Parse("x.md", got).Title != "new \"title\"" {
		t.Fatal("incorrect value")
	}
}
func TestUnsafeScalarRepresentationsRefused(t *testing.T) {
	for name, source := range map[string]string{"flow": "---\n{type: Task, title: hello, status: draft}\n---\n", "quoted-multiline": "---\ntype: Task\ntitle: \"hello\n  world\"\n---\n", "plain-multiline": "---\ntype: Task\ntitle: hello\n  world\n---\n", "block": "---\ntype: Task\ntitle: |\n  hello\n---\n"} {
		t.Run(name, func(t *testing.T) {
			d := Parse("x.md", []byte(source))
			if _, e := d.Edit(map[string]any{"title": "new"}); e == nil {
				t.Fatal("unsafe edit accepted")
			}
		})
	}
}
func TestMarkdownASTLinksAndAnchors(t *testing.T) {
	body := "# Héllo, 世界!\n\n# Duplicate\n# Duplicate\n\n[real][r]\n\n[r]: target.md#x\n\n`[ignored](bad.md)`\n\n```md\n[also ignored](bad.md)\n```\n\n![image](image.md)\n<a id=\"explicit\"></a>\n"
	links, anchors := MarkdownInfo(body)
	if len(links) != 2 {
		t.Fatalf("links: %#v", links)
	}
	for _, a := range []string{"héllo-世界", "duplicate", "duplicate-1", "explicit"} {
		if !anchors[a] {
			t.Errorf("missing anchor %s: %v", a, anchors)
		}
	}
}
func TestFields(t *testing.T) {
	cases := []struct {
		kind  string
		value any
		valid bool
	}{{"integer", true, false}, {"integer", 1, true}, {"boolean", "false", false}, {"mapping", map[string]any{}, true}, {"string-list", []any{}, true}, {"string-list", []any{1}, false}, {"date", "2026-02-30", false}, {"datetime", "2026-10-27T12:30:00Z", true}, {"datetime", "2026-10-27T12:30:00", false}}
	for _, c := range cases {
		e := CheckField(c.value, map[string]any{"kind": c.kind, "required": true})
		if (e == nil) != c.valid {
			t.Errorf("%s %v: %v", c.kind, c.value, e)
		}
	}
}
