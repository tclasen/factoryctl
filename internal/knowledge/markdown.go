package knowledge

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type BodyLink struct {
	Target string
	Line   int
}

var anchorPattern = regexp.MustCompile(`(?i)\b(?:id|name)\s*=\s*["']([^"']+)["']`)

// MarkdownInfo uses CommonMark ASTs so examples in code are never graph edges.
func MarkdownInfo(body string) ([]BodyLink, map[string]bool) {
	source := []byte(body)
	root := goldmark.DefaultParser().Parse(text.NewReader(source))
	links := []BodyLink{}
	anchors := map[string]bool{}
	counts := map[string]int{}
	_ = ast.Walk(root, func(n ast.Node, enter bool) (ast.WalkStatus, error) {
		if !enter {
			return ast.WalkContinue, nil
		}
		switch v := n.(type) {
		case *ast.Link:
			links = append(links, BodyLink{string(v.Destination), 1})
		case *ast.Image:
			links = append(links, BodyLink{string(v.Destination), 1})
		case *ast.AutoLink:
			links = append(links, BodyLink{string(v.URL(source)), 1})
		case *ast.Heading:
			var textBody strings.Builder
			_ = ast.Walk(v, func(c ast.Node, en bool) (ast.WalkStatus, error) {
				if en {
					if t, ok := c.(*ast.Text); ok {
						textBody.Write(t.Segment.Value(source))
					}
					if t, ok := c.(*ast.String); ok {
						textBody.Write(t.Value)
					}
				}
				return ast.WalkContinue, nil
			})
			base := Anchor(textBody.String())
			idx := counts[base]
			counts[base]++
			name := base
			if idx > 0 {
				name += fmtInt(idx)
			}
			anchors[name] = true
		case *ast.RawHTML:
			for i := 0; i < v.Segments.Len(); i++ {
				seg := v.Segments.At(i)
				for _, m := range anchorPattern.FindAllSubmatch(seg.Value(source), -1) {
					anchors[string(m[1])] = true
				}
			}
		case *ast.HTMLBlock:
			for i := 0; i < v.Lines().Len(); i++ {
				seg := v.Lines().At(i)
				for _, m := range anchorPattern.FindAllSubmatch(seg.Value(source), -1) {
					anchors[string(m[1])] = true
				}
			}
		}
		return ast.WalkContinue, nil
	})
	return links, anchors
}
func fmtInt(n int) string {
	const digits = "0123456789"
	s := ""
	for n > 0 {
		s = string(digits[n%10]) + s
		n /= 10
	}
	return "-" + s
}
func Anchor(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '_' || r == '-' {
			b.WriteRune(r)
		} else if unicode.IsSpace(r) {
			b.WriteByte('-')
		}
	}
	return b.String()
}
