package app

import (
	"github.com/tclasen/factoryctl/internal/knowledge"
	"regexp"
)

// QueryMatcher provides literal Unicode case-insensitive matching without
// allocating lowercase document copies. Compile once per query.
func QueryMatcher(query string) func([]byte) bool {
	if query == "" {
		return func([]byte) bool { return true }
	}
	r, err := regexp.Compile("(?i)" + regexp.QuoteMeta(query))
	if err != nil {
		return func([]byte) bool { return false }
	}
	return r.Match
}

// DocumentMatcher uses bounded-memory case-folded Horspool search on ASCII
// documents, with the Unicode matcher for all other text. Document.Parse records
// the ASCII fact once; queries never rescan or normalize the complete corpus.
func DocumentMatcher(query string) func(*knowledge.Document) bool {
	unicodeMatch := QueryMatcher(query)
	if query == "" {
		return func(*knowledge.Document) bool { return true }
	}
	pattern := []byte(query)
	fold := func(c byte) byte {
		if c >= 'A' && c <= 'Z' {
			return c + ('a' - 'A')
		}
		return c
	}
	for i, c := range pattern {
		if c >= 128 {
			return func(d *knowledge.Document) bool { return unicodeMatch(d.Raw) }
		}
		pattern[i] = fold(c)
	}
	var shifts [256]int
	for i := range shifts {
		shifts[i] = len(pattern)
	}
	for i := 0; i < len(pattern)-1; i++ {
		shifts[pattern[i]] = len(pattern) - 1 - i
	}
	return func(d *knowledge.Document) bool {
		if !d.ASCII {
			return unicodeMatch(d.Raw)
		}
		text := d.Raw
		for end := len(pattern) - 1; end < len(text); {
			j := len(pattern) - 1
			k := end
			for j >= 0 && fold(text[k]) == pattern[j] {
				j--
				k--
			}
			if j < 0 {
				return true
			}
			end += shifts[fold(text[end])]
		}
		return false
	}
}
