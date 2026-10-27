package knowledge

import (
	"fmt"
	"math"
	"path"
	"strings"
	"time"
)

type Definition struct {
	Type, Directory, Filename, Path string
	Enabled                         bool
	Fields                          map[string]any
	Raw                             map[string]any
	Invalid                         bool
}

var Kinds = map[string]bool{"string": true, "text": true, "boolean": true, "integer": true, "number": true, "date": true, "datetime": true, "enum": true, "string-list": true, "reference": true, "reference-list": true, "mapping": true}

func CheckField(value any, spec map[string]any) error {
	required, _ := spec["required"].(bool)
	kind := String(spec, "kind")
	if !Kinds[kind] {
		return fmt.Errorf("unsupported field kind %q", kind)
	}
	if value == nil {
		if required {
			return fmt.Errorf("required field is absent or null")
		}
		return nil
	}
	switch kind {
	case "string", "text", "reference", "enum":
		s, ok := value.(string)
		if !ok {
			return fmt.Errorf("expected string")
		}
		if required && strings.TrimSpace(s) == "" {
			return fmt.Errorf("required string is empty")
		}
		if kind == "enum" {
			found := false
			for _, v := range List(spec["values"]) {
				if v == s {
					found = true
				}
			}
			if !found {
				return fmt.Errorf("value is not in the declared enum")
			}
		}
	case "boolean":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("expected boolean")
		}
	case "integer":
		switch v := value.(type) {
		case int, int64, uint64:
		case float64:
			if math.IsInf(v, 0) || math.IsNaN(v) || math.Trunc(v) != v {
				return fmt.Errorf("expected finite integer")
			}
		default:
			return fmt.Errorf("expected integer")
		}
	case "number":
		switch v := value.(type) {
		case int, int64, uint64:
		case float64:
			if math.IsInf(v, 0) || math.IsNaN(v) {
				return fmt.Errorf("expected finite number")
			}
		default:
			return fmt.Errorf("expected number")
		}
	case "date", "datetime":
		s, ok := value.(string)
		if t, yes := value.(time.Time); yes {
			s = t.Format(time.RFC3339)
			if kind == "date" {
				s = t.Format("2006-01-02")
			}
			ok = true
		}
		layout := "2006-01-02"
		if kind == "datetime" {
			layout = time.RFC3339
		}
		if !ok {
			return fmt.Errorf("expected ISO %s", kind)
		}
		if _, err := time.Parse(layout, s); err != nil {
			return fmt.Errorf("expected ISO %s", kind)
		}
	case "string-list", "reference-list":
		a, ok := value.([]any)
		if !ok {
			return fmt.Errorf("expected list")
		}
		for _, v := range a {
			if _, ok := v.(string); !ok {
				return fmt.Errorf("expected list of strings")
			}
		}
		if min, ok := spec["min_items"].(int); ok && len(a) < min {
			return fmt.Errorf("expected at least %d items", min)
		}
	case "mapping":
		if _, ok := value.(map[string]any); !ok {
			return fmt.Errorf("expected mapping")
		}
	}
	return nil
}
func ParseDefinition(d *Document) (Definition, []Diagnostic) {
	def := Definition{Type: String(d.Frontmatter, "defines"), Directory: String(d.Frontmatter, "directory"), Filename: String(d.Frontmatter, "filename"), Path: d.Ref.Path, Enabled: true, Fields: Mapping(d.Frontmatter["fields"]), Raw: d.Frontmatter}
	ds := []Diagnostic{}
	bad := func(field, msg string) {
		def.Invalid = true
		ds = append(ds, Diagnostic{Code: "invalid_definition", Severity: "error", Path: d.Ref.Path, Field: field, Line: 1, Column: 1, Message: msg})
	}
	if def.Type == "" {
		bad("defines", "defines must be a nonempty string")
	}
	if def.Directory == "" || path.IsAbs(def.Directory) || path.Clean(def.Directory) != def.Directory || strings.Contains(def.Directory, "\\") || def.Directory == ".." || strings.HasPrefix(def.Directory, "../") {
		bad("directory", "expected a normalized bundle-relative directory")
	}
	if def.Filename != "" && (path.Base(def.Filename) != def.Filename || !strings.HasSuffix(def.Filename, ".md")) {
		bad("filename", "expected a Markdown basename")
	}
	if v, ok := d.Frontmatter["enabled"]; ok {
		b, yes := v.(bool)
		if !yes {
			bad("enabled", "enabled must be boolean")
		} else {
			def.Enabled = b
		}
	}
	if def.Fields == nil {
		bad("fields", "fields must be a mapping")
	}
	for name, raw := range def.Fields {
		s := Mapping(raw)
		if s == nil {
			bad("fields."+name, "field schema must be a mapping")
			continue
		}
		k := String(s, "kind")
		if !Kinds[k] {
			bad("fields."+name, "unsupported field kind")
		}
		if strings.Contains(name, ".") {
			bad("fields."+name, "nested field paths are unsupported")
		}
		if v, ok := s["required"]; ok {
			if _, ok := v.(bool); !ok {
				bad(name, "required must be boolean")
			}
		}
		if v, ok := s["min_items"]; ok {
			n, ok := v.(int)
			if !ok || n < 0 || (k != "string-list" && k != "reference-list") {
				bad(name, "min_items must be a nonnegative integer for a list")
			}
		}
		if k == "enum" {
			a, ok := s["values"].([]any)
			if !ok || len(a) == 0 {
				bad(name, "enum values must be a nonempty string list")
			}
			seen := map[string]bool{}
			for _, v := range a {
				str, ok := v.(string)
				if !ok || seen[str] {
					bad(name, "enum values must be unique strings")
				}
				seen[str] = true
			}
		}
		if v, ok := s["default"]; ok {
			if e := CheckField(v, s); e != nil {
				bad(name, "invalid default: "+e.Error())
			}
		}
		for key := range s {
			switch key {
			case "kind", "required", "default", "description", "values", "target_types", "min_items":
			default:
				bad(name, "unsupported structural field property: "+key)
			}
		}
	}
	if board := Mapping(Mapping(d.Frontmatter["views"])["board"]); board != nil {
		field := String(board, "field")
		if String(Mapping(def.Fields[field]), "kind") != "enum" {
			bad("views.board.field", "board field must name an enum")
		}
	}
	// Executable or conditional policy declarations are deliberately unsupported.
	if ext := Mapping(d.Frontmatter["extension"]); ext != nil {
		bad("extension", "runtime extension registration is not supported in this CLI milestone")
	}
	return def, ds
}
