package app

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/tclasen/factoryctl/internal/knowledge"
	"github.com/tclasen/factoryctl/internal/project"
)

type ExtensionTool struct {
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	Type           string            `json:"type"`
	Operation      string            `json:"operation"`
	DefinitionHash string            `json:"definition_hash"`
	InputSchema    map[string]any    `json:"inputSchema"`
	Binding        knowledge.Binding `json:"-"`
}
type ExtensionRegistry struct {
	Tools             []ExtensionTool        `json:"tools"`
	Revision          string                 `json:"revision"`
	ReconnectRequired bool                   `json:"reconnect_required"`
	Diagnostics       []knowledge.Diagnostic `json:"diagnostics"`
}

// DefinitionHash includes effective structure and instructions. Unknown types have
// a stable identity too, so adding their first definition conflicts with old drafts.
func DefinitionHash(def knowledge.Definition) string {
	raw, err := knowledge.Encode(def.Raw, fmt.Sprintf("%s\n%s\n%s\n%t\n%t\n%t", def.Type, def.Path, def.Revision, def.Enabled, def.Invalid, def.ExtensionInvalid))
	if err != nil {
		return knowledge.Hash([]byte(fmt.Sprintf("%#v", def)))
	}
	return knowledge.Hash(raw)
}
func ExtensionCatalog(projects []project.Project, p project.Project) (ExtensionRegistry, error) {
	snap, err := Load(projects)
	if err != nil {
		return ExtensionRegistry{}, err
	}
	b := snap.Bundle(p.Checkout)
	if b == nil {
		return ExtensionRegistry{}, Fail("not_found", "selected bundle unavailable")
	}
	r := ExtensionRegistry{Tools: []ExtensionTool{}, Revision: snap.Revision, ReconnectRequired: true, Diagnostics: snap.Diagnostics}
	for typ, def := range b.Definitions {
		if def.Invalid || def.ExtensionInvalid || !def.Enabled {
			continue
		}
		ext, err := knowledge.ParseExtension(def.Raw["extension"], def.Fields)
		if err != nil {
			continue
		}
		for name, binding := range ext.Tools {
			props := map[string]any{"expected_definition": map[string]any{"type": "string"}}
			required := []string{}
			if binding.Operation == "create" {
				props["slug"] = map[string]any{"type": "string"}
				required = append(required, "slug")
			} else {
				props["path"] = map[string]any{"type": "string"}
				required = append(required, "path")
			}
			if binding.Operation != "read" {
				required = append(required, "expected_definition")
			}
			if binding.Operation == "update" || binding.Operation == "delete" {
				props["expected_revision"] = map[string]any{"type": "string"}
				required = append(required, "expected_revision")
			}
			inProps := map[string]any{}
			inRequired := []string{}
			for key, value := range binding.Inputs {
				spec := knowledge.Mapping(value)
				inProps[key] = inputJSONSchema(spec)
				if spec["required"] == true {
					inRequired = append(inRequired, key)
				}
			}
			sort.Strings(inRequired)
			if len(inProps) > 0 {
				props["inputs"] = map[string]any{"type": "object", "properties": inProps, "required": inRequired, "additionalProperties": false}
				if len(inRequired) > 0 {
					required = append(required, "inputs")
				}
			}
			r.Tools = append(r.Tools, ExtensionTool{Name: ext.Namespace + "_" + name, Description: binding.Description, Type: typ, Operation: binding.Operation, DefinitionHash: DefinitionHash(def), Binding: binding, InputSchema: map[string]any{"type": "object", "properties": props, "required": required, "additionalProperties": false}})
		}
	}
	sort.Slice(r.Tools, func(i, j int) bool { return r.Tools[i].Name < r.Tools[j].Name })
	return r, nil
}
func inputJSONSchema(spec map[string]any) map[string]any {
	kind := knowledge.String(spec, "kind")
	typ := "string"
	out := map[string]any{}
	switch kind {
	case "integer", "number", "boolean":
		typ = kind
	case "mapping":
		typ = "object"
	case "string-list", "reference-list":
		typ = "array"
		out["items"] = map[string]any{"type": "string"}
		if v, ok := spec["min_items"]; ok {
			out["minItems"] = v
		}
	}
	out["type"] = typ
	if kind == "enum" {
		out["enum"] = spec["values"]
	}
	return out
}
func CallExtension(projects []project.Project, p project.Project, name string, args map[string]any) (any, error) {
	r, err := ExtensionCatalog(projects, p)
	if err != nil {
		return nil, err
	}
	var t *ExtensionTool
	for i := range r.Tools {
		if r.Tools[i].Name == name {
			t = &r.Tools[i]
			break
		}
	}
	if t == nil {
		return nil, Fail("unsupported_capability", "unknown or unavailable extension tool")
	}
	allowed := knowledge.Mapping(t.InputSchema["properties"])
	for k := range args {
		if allowed[k] == nil {
			return nil, Fail("invalid_arguments", "unknown argument: "+k)
		}
	}
	for _, k := range t.InputSchema["required"].([]string) {
		if _, ok := args[k].(string); k != "inputs" && (!ok || knowledge.String(args, k) == "") {
			return nil, Fail("invalid_arguments", k+" is required")
		}
	}
	for _, k := range []string{"path", "slug", "expected_revision", "expected_definition"} {
		if v, ok := args[k]; ok {
			if _, ok = v.(string); !ok {
				return nil, Fail("invalid_arguments", k+" must be a string")
			}
		}
	}
	if t.Operation != "read" && knowledge.String(args, "expected_definition") != t.DefinitionHash {
		return nil, Fail("stale_definition", "effective definition changed; refresh")
	}
	inputs := knowledge.Mapping(args["inputs"])
	if v, ok := args["inputs"]; ok && v != nil && inputs == nil {
		return nil, Fail("invalid_arguments", "inputs must be a mapping")
	}
	values := map[string]any{}
	for k := range inputs {
		if t.Binding.Inputs[k] == nil {
			return nil, Fail("invalid_arguments", "unknown input: "+k)
		}
	}
	for k, raw := range t.Binding.Inputs {
		spec := knowledge.Mapping(raw)
		v, ok := inputs[k]
		if !ok {
			v = spec["default"]
		}
		if err := knowledge.CheckField(v, spec); err != nil {
			return nil, Fail("invalid_arguments", k+": "+err.Error())
		}
		if v != nil {
			values[k] = v
		}
	}
	fields := map[string]any{}
	for k, raw := range t.Binding.Fields {
		source := knowledge.Mapping(raw)
		if n, ok := source["input"].(string); ok {
			if v, exists := values[n]; exists {
				fields[k] = v
			}
		} else {
			fields[k] = source["literal"]
		}
	}
	if t.Operation == "create" {
		return CreateChecked(projects, p, t.Type, knowledge.String(args, "slug"), fields, t.DefinitionHash)
	}
	path := knowledge.String(args, "path")
	if err = SafePath(p.Root+"/factory", path, false); err != nil {
		return nil, err
	}
	if t.Operation == "read" {
		snap, err := Load(projects)
		if err != nil {
			return nil, err
		}
		b := snap.Bundle(p.Checkout)
		d := b.Docs[path]
		if d == nil {
			return nil, Fail("not_found", "record not found")
		}
		if d.Type != t.Type {
			return nil, Fail("invalid_path", "record is outside the tool type scope")
		}
		return RecordContext(b, d, r), nil
	}
	return Mutate(projects, p, Mutation{Kind: t.Operation, Path: path, Set: fields, ExpectedRevision: knowledge.String(args, "expected_revision"), ExpectedDefinition: t.DefinitionHash, DefinitionType: t.Type})
}
func RecordContext(b *Bundle, d *knowledge.Document, r ExtensionRegistry) map[string]any {
	def := b.Definitions[d.Type]
	instructions := ""
	if node := b.Docs[def.Path]; node != nil {
		instructions = node.Body
	}
	tools := []ExtensionTool{}
	for _, t := range r.Tools {
		if t.Type == d.Type {
			tools = append(tools, t)
		}
	}
	return map[string]any{"node": d, "source": string(d.Raw), "definition": def, "definition_hash": DefinitionHash(def), "instructions": instructions, "tools": tools, "access_enforcement": "advisory"}
}

// DecodeExtensionArguments is shared by CLI/MCP to reject non-object envelopes.
func DecodeExtensionArguments(raw []byte) (map[string]any, error) {
	var args map[string]any
	if err := json.Unmarshal(raw, &args); err != nil || args == nil {
		return nil, Fail("invalid_arguments", "expected an argument object")
	}
	return args, nil
}
