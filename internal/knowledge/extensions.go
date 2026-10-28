package knowledge

import (
	"fmt"
	"regexp"
)

// Extension bindings are data. They never evaluate expressions or invoke code.
type Extension struct {
	Namespace string
	Tools     map[string]Binding
}
type Binding struct {
	Operation, Description string
	Inputs, Fields         map[string]any
}

var extensionName = regexp.MustCompile(`^[a-z][a-z0-9]*$`)
var bindingName = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

func ParseExtension(raw any, fields map[string]any) (*Extension, error) {
	m := Mapping(raw)
	if m == nil {
		return nil, fmt.Errorf("extension must be a mapping")
	}
	for k := range m {
		if k != "version" && k != "namespace" && k != "tools" {
			return nil, fmt.Errorf("unsupported extension property: %s", k)
		}
	}
	if m["version"] != 1 {
		return nil, fmt.Errorf("extension version must be 1")
	}
	ns := String(m, "namespace")
	// Single-token namespaces make full names unambiguous and cannot shadow factory_*.
	if !extensionName.MatchString(ns) || ns == "factory" {
		return nil, fmt.Errorf("namespace must be lowercase alphanumeric and must not be factory")
	}
	tools := Mapping(m["tools"])
	if tools == nil {
		return nil, fmt.Errorf("tools must be a mapping")
	}
	ext := &Extension{Namespace: ns, Tools: map[string]Binding{}}
	for name, value := range tools {
		if !bindingName.MatchString(name) {
			return nil, fmt.Errorf("invalid tool name: %s", name)
		}
		t := Mapping(value)
		if t == nil {
			return nil, fmt.Errorf("tool must be a mapping")
		}
		for k := range t {
			switch k {
			case "operation", "description", "inputs", "fields":
			default:
				return nil, fmt.Errorf("unsupported tool binding property: %s", k)
			}
		}
		b := Binding{Operation: String(t, "operation"), Description: String(t, "description"), Inputs: Mapping(t["inputs"]), Fields: Mapping(t["fields"])}
		switch b.Operation {
		case "read", "create", "update", "delete":
		default:
			return nil, fmt.Errorf("unsupported core operation")
		}
		if v, ok := t["inputs"]; ok && Mapping(v) == nil {
			return nil, fmt.Errorf("inputs must be a mapping")
		}
		if v, ok := t["fields"]; ok && Mapping(v) == nil {
			return nil, fmt.Errorf("fields must be a mapping")
		}
		for input, raw := range b.Inputs {
			if !bindingName.MatchString(input) || input == "path" || input == "slug" || input == "inputs" || input == "expected_revision" || input == "expected_definition" {
				return nil, fmt.Errorf("invalid or reserved input: %s", input)
			}
			spec := Mapping(raw)
			if spec == nil || !Kinds[String(spec, "kind")] {
				return nil, fmt.Errorf("invalid input kind")
			}
			// Reuse the complete structural schema validator, including enum/default checks.
			d := Parse("definitions/input.md", []byte("---\ntype: Record Definition\ntitle: Input\ndefines: Input\ndirectory: inputs\nfields: {}\n---\n"))
			d.Frontmatter["fields"] = map[string]any{input: spec}
			_, ds := ParseDefinition(d)
			if len(ds) > 0 {
				return nil, fmt.Errorf("input %s: %s", input, ds[0].Message)
			}
		}
		if (b.Operation == "read" || b.Operation == "delete") && (len(b.Fields) > 0 || len(b.Inputs) > 0) {
			return nil, fmt.Errorf("read/delete do not accept field mappings or inputs")
		}
		for field, value := range b.Fields {
			if _, ok := fields[field]; !ok && field != "title" && (field != "summary" || b.Operation != "create") && field != "description" {
				return nil, fmt.Errorf("undeclared mapped field: %s", field)
			}
			if field == "type" || field == "path" || field == "slug" || field == "expected_revision" || field == "expected_definition" {
				return nil, fmt.Errorf("reserved mapped field")
			}
			source := Mapping(value)
			if len(source) != 1 {
				return nil, fmt.Errorf("field mapping requires exactly one input or literal")
			}
			if input, ok := source["input"]; ok {
				n, ok := input.(string)
				if !ok || b.Inputs[n] == nil {
					return nil, fmt.Errorf("field mapping references undeclared input")
				}
			} else if _, ok := source["literal"]; !ok {
				return nil, fmt.Errorf("only input or literal field mappings are supported")
			}
		}
		ext.Tools[name] = b
	}
	return ext, nil
}
