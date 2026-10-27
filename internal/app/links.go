package app

import (
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/tclasen/factoryctl/internal/knowledge"
)

func (s *Snapshot) validateLinks(b *Bundle, d *knowledge.Document) {
	add := func(target, field string, types []any, resource bool) {
		if target == "" {
			return
		}
		link := Link{Source: d.Ref, Target: target, Field: field, Status: "valid"}
		tb, tp, fragment, status := s.resolve(b, d, target, resource)
		link.Status = status
		if tb != nil && tp != "" {
			ref := knowledge.Ref{ProjectID: tb.Project.ID, Checkout: tb.Project.Checkout, Path: tp}
			link.Resolved = &ref
			td := tb.Docs[tp]
			if td == nil {
				link.Status = "broken"
			} else {
				if fragment != "" {
					_, anchors := knowledge.MarkdownInfo(td.Body)
					if !anchors[fragment] {
						link.Status = "broken"
					}
				}
				if len(types) > 0 {
					match := false
					for _, typ := range types {
						if typ == td.Type {
							match = true
						}
					}
					if !match {
						link.Status = "broken"
					}
				}
			}
		}
		if link.Status == "broken" {
			dg := diag(b.Project.ID, d.Ref.Path, "broken_link", "missing, unsafe or incompatible reference target/anchor")
			dg.Field = field
			dg.Target = target
			s.Diagnostics = append(s.Diagnostics, dg)
		}
		s.Links = append(s.Links, link)
	}
	links, _ := knowledge.MarkdownInfo(d.Body)
	for _, l := range links {
		add(l.Target, "body", nil, false)
	}
	if target := knowledge.String(d.Frontmatter, "resource"); looksLikeReference(target) {
		add(target, "resource", nil, true)
	}
	for _, src := range knowledge.List(d.Frontmatter["sources"]) {
		target := knowledge.String(knowledge.Mapping(src), "resource")
		if looksLikeReference(target) {
			add(target, "sources.resource", nil, true)
		}
	}
	def, ok := b.Definitions[d.Type]
	if ok {
		for field, raw := range def.Fields {
			spec := knowledge.Mapping(raw)
			kind := knowledge.String(spec, "kind")
			if kind == "reference" {
				add(knowledge.String(d.Frontmatter, field), field, knowledge.List(spec["target_types"]), false)
			}
			if kind == "reference-list" {
				for _, v := range knowledge.List(d.Frontmatter[field]) {
					target, _ := v.(string)
					add(target, field, knowledge.List(spec["target_types"]), false)
				}
			}
		}
	}
	if d.Type == "Record Definition" {
		add(knowledge.String(d.Frontmatter, "template"), "template", []any{"Record Template"}, false)
		for _, v := range knowledge.List(d.Frontmatter["actions"]) {
			target, _ := v.(string)
			add(target, "actions", []any{"Prompt"}, false)
		}
	}
}
func looksLikeReference(s string) bool {
	return strings.HasPrefix(s, "./") || strings.HasPrefix(s, "../") || strings.HasPrefix(s, "/") || strings.HasPrefix(s, "#") || strings.Contains(s, "://") || (!strings.Contains(s, " ") && strings.HasSuffix(s, ".md"))
}
func (s *Snapshot) resolve(b *Bundle, d *knowledge.Document, target string, resource bool) (*Bundle, string, string, string) {
	u, e := url.Parse(target)
	if e != nil {
		return nil, "", "", "broken"
	}
	fragment := u.Fragment
	if strings.Contains(u.Path, "\\") || strings.ContainsRune(u.Path, 0) || strings.Contains(strings.ToLower(u.EscapedPath()), "%25") {
		return nil, "", "", "broken"
	}
	if u.IsAbs() {
		if u.Scheme != "http" && u.Scheme != "https" && u.Scheme != "mailto" {
			return nil, "", "", "broken"
		}
		if u.Scheme == "mailto" {
			return nil, "", "", "unchecked"
		}
		if u.Host == "" || u.User != nil {
			return nil, "", "", "broken"
		}
		var matches []*Bundle
		prefix := ""
		for _, candidate := range s.Bundles {
			base := candidate.Project.BundleURI
			if base != "" && strings.HasPrefix(target, base) {
				if len(base) > len(prefix) {
					prefix = base
					matches = nil
				}
				if len(base) == len(prefix) {
					matches = append(matches, candidate)
				}
			}
		}
		if len(matches) == 0 {
			return nil, "", "", "unchecked"
		}
		if len(matches) != 1 {
			return nil, "", "", "broken"
		}
		rest := strings.TrimPrefix(target, prefix)
		ru, e := url.Parse(rest)
		if e != nil {
			return nil, "", "", "broken"
		}
		p := path.Clean(ru.Path)
		if p == ".." || strings.HasPrefix(p, "../") || path.IsAbs(p) {
			return nil, "", "", "broken"
		}
		return matches[0], p, ru.Fragment, "valid"
	}
	if u.Host != "" {
		return nil, "", "", "broken"
	}
	p := u.Path
	if p == "" {
		return b, d.Ref.Path, fragment, "valid"
	}
	if strings.HasPrefix(p, "/") {
		p = path.Clean(strings.TrimPrefix(p, "/"))
	} else {
		p = path.Clean(path.Join(path.Dir(d.Ref.Path), p))
	}
	if p == ".." || strings.HasPrefix(p, "../") {
		if resource {
			full := filepath.Clean(filepath.Join(b.Project.Root, "factory", filepath.FromSlash(p)))
			rel, e := filepath.Rel(b.Project.RepositoryRoot, full)
			if e != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
				return nil, "", "", "broken"
			}
			real, e := filepath.EvalSymlinks(full)
			if e != nil {
				return nil, "", "", "broken"
			}
			rr, e := filepath.Rel(b.Project.RepositoryRoot, real)
			if e != nil || rr == ".." || strings.HasPrefix(rr, ".."+string(filepath.Separator)) {
				return nil, "", "", "broken"
			}
			if info, e := os.Stat(real); e == nil && info.Mode().IsRegular() {
				return nil, "", "", "valid"
			}
		}
		return nil, "", "", "broken"
	}
	if !strings.HasSuffix(p, ".md") {
		return nil, "", "", "broken"
	}
	return b, p, fragment, "valid"
}
