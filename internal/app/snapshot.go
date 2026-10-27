// Package app is the transport-independent document service. It interprets
// structure only: workflow prose never controls acceptance of a mutation.
package app

import (
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/tclasen/factoryctl/internal/defaults"
	"github.com/tclasen/factoryctl/internal/knowledge"
	"github.com/tclasen/factoryctl/internal/project"
	"github.com/tclasen/factoryctl/internal/templates"
)

type Error struct {
	Code            string                 `json:"code"`
	Message         string                 `json:"message"`
	Retryable       bool                   `json:"retryable"`
	Diagnostics     []knowledge.Diagnostic `json:"diagnostics,omitempty"`
	CurrentRevision string                 `json:"current_revision,omitempty"`
}

func (e *Error) Error() string    { return e.Code + ": " + e.Message }
func Fail(code, msg string) error { return &Error{Code: code, Message: msg} }

type Bundle struct {
	Project     project.Project
	Docs        map[string]*knowledge.Document
	Definitions map[string]knowledge.Definition
	Diagnostics []knowledge.Diagnostic
}
type Link struct {
	Source   knowledge.Ref  `json:"source"`
	Target   string         `json:"target"`
	Resolved *knowledge.Ref `json:"resolved,omitempty"`
	Status   string         `json:"status"`
	Field    string         `json:"field,omitempty"`
}
type Snapshot struct {
	Bundles     []*Bundle
	Diagnostics []knowledge.Diagnostic
	Links       []Link
	Revision    string
}
type Validation struct {
	Scope          []project.Project      `json:"scope"`
	Level          string                 `json:"level"`
	Complete       bool                   `json:"complete"`
	Errors         int                    `json:"errors"`
	Warnings       int                    `json:"warnings"`
	UncheckedLinks int                    `json:"unchecked_links"`
	Diagnostics    []knowledge.Diagnostic `json:"diagnostics"`
}

func Load(projects []project.Project) (*Snapshot, error) {
	s := &Snapshot{}
	for _, p := range projects {
		b := &Bundle{Project: p, Docs: map[string]*knowledge.Document{}}
		root := filepath.Join(p.Root, "factory")
		err := filepath.WalkDir(root, func(file string, entry fs.DirEntry, walkErr error) error {
			rel, _ := filepath.Rel(root, file)
			rel = filepath.ToSlash(rel)
			if walkErr != nil {
				b.Diagnostics = append(b.Diagnostics, diag(p.ID, rel, "io_error", walkErr.Error()))
				if entry != nil && entry.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if entry.Type()&os.ModeSymlink != 0 {
				b.Diagnostics = append(b.Diagnostics, diag(p.ID, rel, "invalid_path", "symlink is not a bundle document"))
				return nil
			}
			if entry.IsDir() {
				return nil
			}
			if !strings.HasSuffix(rel, ".md") {
				b.Diagnostics = append(b.Diagnostics, diag(p.ID, rel, "invalid_location", "only .md concept files belong in factory/"))
				return nil
			}
			data, err := ReadBounded(file)
			if err != nil {
				b.Diagnostics = append(b.Diagnostics, diag(p.ID, rel, "io_error", err.Error()))
				return nil
			}
			d := knowledge.Parse(rel, data)
			d.Ref.ProjectID = p.ID
			d.Ref.Checkout = p.Checkout
			b.Docs[rel] = d
			return nil
		})
		if err != nil {
			return nil, err
		}
		s.Bundles = append(s.Bundles, b)
	}
	s.Validate("factory")
	return s, nil
}
func ReadBounded(file string) ([]byte, error) {
	f, e := os.Open(file)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	info, e := f.Stat()
	if e != nil {
		return nil, e
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("not a regular file")
	}
	return io.ReadAll(io.LimitReader(f, knowledge.MaxDocument+1))
}
func diag(id, p, code, msg string) knowledge.Diagnostic {
	return knowledge.Diagnostic{Code: code, Severity: "error", ProjectID: id, Path: p, Line: 1, Column: 1, Message: msg}
}
func (s *Snapshot) Validate(level string) Validation {
	s.Diagnostics = nil
	s.Links = nil
	fingerprints := []string{}
	for _, b := range s.Bundles {
		startDiagnostics := len(s.Diagnostics)
		s.Diagnostics = append(s.Diagnostics, b.Diagnostics...)
		b.Definitions = map[string]knowledge.Definition{}
		local := map[string]bool{}
		for p, data := range defaults.Definitions() {
			d := knowledge.Parse(p, data)
			def, _ := knowledge.ParseDefinition(d)
			b.Definitions[def.Type] = def
		}
		for _, p := range sortedPaths(b.Docs) {
			d := b.Docs[p]
			for _, dg := range d.Diagnostics {
				dg.ProjectID = b.Project.ID
				s.Diagnostics = append(s.Diagnostics, dg)
			}
			fingerprints = append(fingerprints, b.Project.Checkout+p+d.Revision)
			if d.Type == "Record Definition" && level == "factory" {
				def, ds := knowledge.ParseDefinition(d)
				for _, dg := range ds {
					dg.ProjectID = b.Project.ID
					s.Diagnostics = append(s.Diagnostics, dg)
				}
				if local[def.Type] {
					s.Diagnostics = append(s.Diagnostics, diag(b.Project.ID, p, "schema_conflict", "duplicate local definition for "+def.Type))
					def.Invalid = true
				}
				local[def.Type] = true
				b.Definitions[def.Type] = def
			}
		}
		cases := map[string]string{}
		for _, p := range sortedPaths(b.Docs) {
			d := b.Docs[p]
			if d.YAML == nil {
				continue
			}
			lower := strings.ToLower(p)
			if other, ok := cases[lower]; ok {
				s.Diagnostics = append(s.Diagnostics, diag(b.Project.ID, p, "invalid_location", "case collision with "+other))
			}
			cases[lower] = p
			s.validateBase(b, d, level)
			if level == "factory" {
				if def, ok := b.Definitions[d.Type]; ok {
					d.Schema = "available"
					if def.Invalid {
						s.Diagnostics = append(s.Diagnostics, diag(b.Project.ID, p, "schema_conflict", "record definition is invalid or duplicated"))
					} else {
						s.validateDefined(b, d, def)
					}
				}
			}
		}
		for i := startDiagnostics; i < len(s.Diagnostics); i++ {
			s.Diagnostics[i].Checkout = b.Project.Checkout
		}
	}
	if level == "factory" {
		for _, b := range s.Bundles {
			startDiagnostics := len(s.Diagnostics)
			for _, p := range sortedPaths(b.Docs) {
				d := b.Docs[p]
				if d.YAML != nil {
					s.validateLinks(b, d)
					if d.Type == "Record Template" || d.Type == "Prompt" {
						if e := templates.Validate(b.Project.Root, b.Docs, p); e != nil {
							s.Diagnostics = append(s.Diagnostics, diag(b.Project.ID, p, "invalid_template", e.Error()))
						}
					}
				}
			}
			for i := startDiagnostics; i < len(s.Diagnostics); i++ {
				s.Diagnostics[i].Checkout = b.Project.Checkout
			}
		}
	}
	sort.Strings(fingerprints)
	s.Revision = knowledge.Hash([]byte(strings.Join(fingerprints, "\n")))
	sort.Slice(s.Diagnostics, func(i, j int) bool {
		a, c := s.Diagnostics[i], s.Diagnostics[j]
		return a.ProjectID+a.Path+a.Field+a.Code+a.Message < c.ProjectID+c.Path+c.Field+c.Code+c.Message
	})
	v := Validation{Level: level, Complete: true, Scope: []project.Project{}, Diagnostics: s.Diagnostics}
	if v.Diagnostics == nil {
		v.Diagnostics = []knowledge.Diagnostic{}
	}
	for _, b := range s.Bundles {
		v.Scope = append(v.Scope, b.Project)
	}
	for _, d := range s.Diagnostics {
		if d.Severity == "error" {
			v.Errors++
		} else {
			v.Warnings++
		}
	}
	for _, l := range s.Links {
		if l.Status == "unchecked" {
			v.UncheckedLinks++
			v.Complete = false
		}
	}
	return v
}
func (s *Snapshot) validateBase(b *Bundle, d *knowledge.Document, level string) {
	add := func(code, field, msg string) {
		dg := diag(b.Project.ID, d.Ref.Path, code, msg)
		dg.Field = field
		s.Diagnostics = append(s.Diagnostics, dg)
	}
	for _, key := range []string{"title", "description", "resource", "status"} {
		if v, exists := d.Frontmatter[key]; exists && v != nil {
			if _, ok := v.(string); !ok {
				add("invalid_document", key, "expected string")
			}
		}
	}
	if st := knowledge.String(d.Frontmatter, "status"); st != "" && st != "draft" && st != "stable" && st != "deprecated" {
		dg := diag(b.Project.ID, d.Ref.Path, "unknown_status", "unknown OKF status")
		dg.Severity = "warning"
		s.Diagnostics = append(s.Diagnostics, dg)
	}
	if v, ok := d.Frontmatter["sources"]; ok {
		a, ok := v.([]any)
		if !ok {
			add("invalid_document", "sources", "expected list of source mappings")
		} else {
			for _, v := range a {
				m := knowledge.Mapping(v)
				if strings.TrimSpace(knowledge.String(m, "resource")) == "" {
					add("invalid_document", "sources.resource", "source requires a resource string")
				}
			}
		}
	}
	if v, ok := d.Frontmatter["generated"]; ok {
		m := knowledge.Mapping(v)
		if knowledge.String(m, "by") == "" {
			add("invalid_document", "generated.by", "generated requires by")
		}
	}
	if v, ok := d.Frontmatter["verified"]; ok {
		a, yes := v.([]any)
		if !yes {
			if m := knowledge.Mapping(v); m != nil {
				a = []any{m}
			} else {
				add("invalid_document", "verified", "expected mapping or list")
			}
		}
		for _, item := range a {
			if knowledge.Mapping(item) == nil {
				add("invalid_document", "verified", "expected verification mappings")
			}
		}
	}
	if level != "factory" {
		return
	}
	base := strings.ToLower(path.Base(d.Ref.Path))
	if base == "index.md" || base == "log.md" {
		add("invalid_location", "", "reserved concept filename")
	}
	if v, ok := d.Frontmatter["factory"]; ok {
		m := knowledge.Mapping(v)
		if m == nil {
			add("invalid_document", "factory", "expected mapping")
		} else if format, ok := m["format"].(int); !ok || format != 1 {
			add("unsupported_format", "factory.format", "structured writes support factory format 1")
		}
	}
	if d.Type == "Factory" {
		id := knowledge.String(d.Frontmatter, "project_id")
		if strings.TrimSpace(id) == "" {
			add("invalid_document", "project_id", "Factory requires project_id")
		}
		if uri := knowledge.String(d.Frontmatter, "bundle_uri"); uri != "" {
			u, e := url.Parse(uri)
			if e != nil || u.Scheme != "https" || u.Host == "" || u.RawQuery != "" || u.Fragment != "" || !strings.HasSuffix(uri, "/") {
				add("invalid_document", "bundle_uri", "expected canonical HTTPS prefix ending in /")
			}
		}
	}
}
func (s *Snapshot) validateDefined(b *Bundle, d *knowledge.Document, def knowledge.Definition) {
	add := func(code, field, msg string) {
		dg := diag(b.Project.ID, d.Ref.Path, code, msg)
		dg.Field = field
		s.Diagnostics = append(s.Diagnostics, dg)
	}
	if path.Dir(d.Ref.Path) != def.Directory || (def.Filename != "" && path.Base(d.Ref.Path) != def.Filename) {
		add("invalid_location", "", "expected directory "+def.Directory+" and filename "+def.Filename)
	}
	for name, raw := range def.Fields {
		if err := knowledge.CheckField(d.Frontmatter[name], knowledge.Mapping(raw)); err != nil {
			add("invalid_document", name, err.Error())
		}
	}
}
func sortedPaths(m map[string]*knowledge.Document) []string {
	p := make([]string, 0, len(m))
	for key := range m {
		p = append(p, key)
	}
	sort.Strings(p)
	return p
}
func (s *Snapshot) Bundle(checkout string) *Bundle {
	for _, b := range s.Bundles {
		if b.Project.Checkout == checkout {
			return b
		}
	}
	return nil
}
