// Package knowledge parses bounded, inert OKF documents without rewriting source.
package knowledge

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"gopkg.in/yaml.v3"
)

const MaxDocument = 1 << 20
const MaxFrontmatter = 256 << 10
const MaxDepth = 32

type Diagnostic struct {
	Checkout  string `json:"checkout,omitempty"`
	Code      string `json:"code"`
	Severity  string `json:"severity"`
	ProjectID string `json:"project_id"`
	Path      string `json:"path"`
	Line      int    `json:"line"`
	Column    int    `json:"column"`
	Field     string `json:"field,omitempty"`
	Target    string `json:"target,omitempty"`
	Message   string `json:"message"`
}
type Ref struct {
	ProjectID string `json:"project_id"`
	Checkout  string `json:"checkout"`
	Path      string `json:"path"`
}
type Document struct {
	ASCII       bool           `json:"-"`
	Schema      string         `json:"schema"`
	Ref         Ref            `json:"ref"`
	Revision    string         `json:"revision"`
	Type        string         `json:"type"`
	Title       string         `json:"title"`
	Frontmatter map[string]any `json:"frontmatter"`
	Body        string         `json:"body"`
	Diagnostics []Diagnostic   `json:"diagnostics"`
	Raw         []byte         `json:"-"`
	YAML        *yaml.Node     `json:"-"`
	HeaderStart int            `json:"-"`
	HeaderEnd   int            `json:"-"`
	bodyLinks   []BodyLink
	bodyAnchors map[string]bool
}

// BodyInfo reuses immutable parse results for link and anchor validation.
// Documents returned by Parse are immutable snapshots; edits create new ones.
func (d *Document) BodyInfo() ([]BodyLink, map[string]bool) {
	if d.bodyAnchors != nil {
		return d.bodyLinks, d.bodyAnchors
	}
	return MarkdownInfo(d.Body)
}

func Hash(b []byte) string { h := sha256.Sum256(b); return hex.EncodeToString(h[:]) }
func Parse(p string, b []byte) *Document {
	d := &Document{Schema: "unavailable", Ref: Ref{Path: p}, Raw: b, Revision: Hash(b), Diagnostics: []Diagnostic{}}
	fail := func(s string) {
		d.Diagnostics = append(d.Diagnostics, Diagnostic{Code: "invalid_document", Severity: "error", Path: p, Line: 1, Column: 1, Message: s})
	}
	if len(b) > MaxDocument {
		fail("document exceeds 1 MiB")
		return d
	}
	d.ASCII = true
	for _, c := range b {
		if c >= 128 {
			d.ASCII = false
			break
		}
	}
	if !utf8.Valid(b) {
		fail("document is not UTF-8")
		return d
	}
	for _, line := range bytes.Split(b, []byte("\n")) {
		if bytes.HasPrefix(line, []byte("<<<<<<< ")) || bytes.Equal(line, []byte("=======")) || bytes.HasPrefix(line, []byte(">>>>>>> ")) {
			fail("unresolved conflict markers")
			return d
		}
	}
	lines := bytes.SplitAfter(b, []byte("\n"))
	if len(lines) < 3 || strings.TrimRight(string(lines[0]), "\r\n") != "---" {
		fail("expected YAML frontmatter delimited by ---")
		return d
	}
	start := len(lines[0])
	end := start
	found := false
	for _, line := range lines[1:] {
		if strings.TrimRight(string(line), "\r\n") == "---" {
			d.Body = string(b[end+len(line):])
			found = true
			break
		}
		end += len(line)
	}
	if !found {
		fail("missing closing frontmatter delimiter")
		return d
	}
	if end-start > MaxFrontmatter {
		fail("frontmatter exceeds 256 KiB")
		return d
	}
	dec := yaml.NewDecoder(bytes.NewReader(b[start:end]))
	var n yaml.Node
	if err := dec.Decode(&n); err != nil {
		fail(err.Error())
		return d
	}
	var extra yaml.Node
	if err := dec.Decode(&extra); err != io.EOF {
		fail("multiple YAML documents are forbidden")
		return d
	}
	if len(n.Content) != 1 || n.Content[0].Kind != yaml.MappingNode {
		fail("frontmatter must be a mapping")
		return d
	}
	if err := checkYAML(&n, 0); err != nil {
		fail(err.Error())
		return d
	}
	d.YAML = &n
	d.HeaderStart = start
	d.HeaderEnd = end
	if err := n.Decode(&d.Frontmatter); err != nil {
		fail(err.Error())
		return d
	}
	if _, nonfinite := jsonValue(d.Frontmatter); nonfinite {
		d.Diagnostics = append(d.Diagnostics, Diagnostic{Code: "nonfinite_metadata", Severity: "warning", Path: p, Line: 1, Column: 1, Message: "nonfinite YAML numbers are represented as .nan/.inf/-.inf strings in JSON; source bytes remain unchanged"})
	}
	d.Type, _ = d.Frontmatter["type"].(string)
	d.Title, _ = d.Frontmatter["title"].(string)
	if strings.TrimSpace(d.Type) == "" {
		fail("type must be a nonempty string")
	}
	d.bodyLinks, d.bodyAnchors = MarkdownInfo(d.Body)
	return d
}
func checkYAML(n *yaml.Node, depth int) error {
	if depth > MaxDepth {
		return fmt.Errorf("YAML nesting exceeds %d", MaxDepth)
	}
	if n.Kind == yaml.AliasNode || n.Anchor != "" {
		return fmt.Errorf("YAML anchors and aliases are forbidden")
	}
	switch n.Tag {
	case "", "!!map", "!!seq", "!!str", "!!null", "!!bool", "!!int", "!!float", "!!timestamp":
	default:
		return fmt.Errorf("custom YAML tags are forbidden")
	}
	if n.Kind == yaml.MappingNode {
		seen := map[string]bool{}
		for i := 0; i < len(n.Content); i += 2 {
			k := n.Content[i]
			if k.Tag != "!!str" {
				return fmt.Errorf("mapping keys must be strings")
			}
			if seen[k.Value] {
				return fmt.Errorf("duplicate key %q", k.Value)
			}
			seen[k.Value] = true
		}
	}
	for _, c := range n.Content {
		if err := checkYAML(c, depth+1); err != nil {
			return err
		}
	}
	return nil
}
func String(m map[string]any, k string) string { s, _ := m[k].(string); return s }
func Mapping(v any) map[string]any             { m, _ := v.(map[string]any); return m }
func List(v any) []any                         { a, _ := v.([]any); return a }
func Encode(m map[string]any, body string) ([]byte, error) {
	b, e := yaml.Marshal(m)
	if e != nil {
		return nil, e
	}
	return []byte("---\n" + string(b) + "---\n" + body), nil
}

// Edit splices only plain single-line scalar values, preserving every unrelated
// byte. Complex YAML edits explicitly require the full-document replace API.
func (d *Document) Edit(set map[string]any) ([]byte, error) {
	if d.YAML == nil {
		return nil, fmt.Errorf("invalid_document")
	}
	if d.YAML.Content[0].Style == yaml.FlowStyle {
		return nil, fmt.Errorf("lossless_edit_unavailable: flow mappings require records replace")
	}
	result := append([]byte(nil), d.Raw...)
	type change struct {
		start, end int
		value      string
	}
	changes := []change{}
	for key, value := range set {
		encoded, e := json.Marshal(value)
		if e != nil {
			return nil, e
		}
		var target *yaml.Node
		root := d.YAML.Content[0]
		for i := 0; i < len(root.Content); i += 2 {
			if root.Content[i].Value == key {
				target = root.Content[i+1]
				break
			}
		}
		if target == nil {
			return nil, fmt.Errorf("lossless_edit_unavailable: adding fields requires records replace")
		}
		if target.Kind != yaml.ScalarNode || target.Style == yaml.LiteralStyle || target.Style == yaml.FoldedStyle || strings.Contains(target.Value, "\n") {
			return nil, fmt.Errorf("lossless_edit_unavailable: complex fields require records replace")
		}
		lines := bytes.SplitAfter(d.Raw[d.HeaderStart:d.HeaderEnd], []byte("\n"))
		offset := d.HeaderStart
		for i := 0; i < target.Line-1; i++ {
			offset += len(lines[i])
		}
		lineRunes := []rune(string(lines[target.Line-1]))
		if target.Column-1 > len(lineRunes) {
			return nil, fmt.Errorf("lossless_edit_unavailable")
		}
		offset += len(string(lineRunes[:target.Column-1]))
		lineEnd := bytes.IndexByte(d.Raw[offset:], '\n')
		if lineEnd < 0 {
			return nil, fmt.Errorf("lossless_edit_unavailable")
		}
		source := string(d.Raw[offset : offset+lineEnd])
		length := 0
		if strings.HasPrefix(source, "\"") {
			for i := 1; i < len(source); i++ {
				if source[i] == '\\' {
					i++
					continue
				}
				if source[i] == '"' {
					length = i + 1
					break
				}
			}
		} else if strings.HasPrefix(source, "'") {
			for i := 1; i < len(source); i++ {
				if source[i] == '\'' {
					if i+1 < len(source) && source[i+1] == '\'' {
						i++
						continue
					}
					length = i + 1
					break
				}
			}
		} else {
			length = len(strings.TrimRight(source, " \r"))
			for i := 1; i < len(source); i++ {
				if source[i] == '#' && (source[i-1] == ' ' || source[i-1] == '\t') {
					length = len(strings.TrimRight(source[:i], " \t"))
					break
				}
			}
		}
		if length == 0 {
			return nil, fmt.Errorf("lossless_edit_unavailable")
		}
		var probe yaml.Node
		if err := yaml.Unmarshal([]byte(source[:length]), &probe); err != nil || len(probe.Content) != 1 || probe.Content[0].Value != target.Value {
			return nil, fmt.Errorf("lossless_edit_unavailable: multiline scalar requires records replace")
		}
		changes = append(changes, change{offset, offset + length, string(encoded)})
	}
	for len(changes) > 0 {
		idx := 0
		for i := range changes {
			if changes[i].start > changes[idx].start {
				idx = i
			}
		}
		c := changes[idx]
		result = append(append(append([]byte{}, result[:c.start]...), []byte(c.value)...), result[c.end:]...)
		changes = append(changes[:idx], changes[idx+1:]...)
	}
	return result, nil
}
