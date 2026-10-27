// Package defaults supplies dependency-complete, human-readable starter documents.
// They are materialized only by explicit initialization, never by discovery.
package defaults

import "github.com/tclasen/factoryctl/internal/knowledge"

func field(kind string, required bool) map[string]any {
	return map[string]any{"kind": kind, "required": required}
}
func enumeration(required bool, def string, values ...string) map[string]any {
	f := field("enum", required)
	f["default"] = def
	f["values"] = values
	return f
}
func definition(name, dir string, fields map[string]any) map[string]any {
	return map[string]any{"type": "Record Definition", "title": name + " record definition", "defines": name, "directory": dir, "fields": fields, "status": "draft", "factory": map[string]any{"format": 1}}
}
func encode(meta map[string]any, body string) []byte {
	b, e := knowledge.Encode(meta, body)
	if e != nil {
		panic(e)
	}
	return b
}
func baseDefinitions() map[string]map[string]any {
	defs := map[string]map[string]any{
		"definitions/record-definition.md": definition("Record Definition", "definitions", map[string]any{"defines": field("string", true), "directory": field("string", true), "enabled": field("boolean", false), "fields": field("mapping", true), "views": field("mapping", false), "extension": field("mapping", false)}),
		"definitions/factory.md":           definition("Factory", ".", map[string]any{"project_id": field("string", true), "bundle_uri": field("string", false)}),
		"definitions/factory-control.md":   definition("Factory Control", "control", map[string]any{}),
		"definitions/record-template.md":   definition("Record Template", "templates", map[string]any{"for_type": field("string", true), "inputs": field("mapping", true), "asset": field("string", false)}),
		"definitions/prompt.md":            definition("Prompt", "prompts", map[string]any{"applies_to": field("string-list", true), "inputs": field("mapping", true), "profile": field("reference", false), "conversation": enumeration(true, "fresh", "fresh")}),
		"definitions/agent-profile.md":     definition("Agent Profile", "agents", map[string]any{}),
		"definitions/task.md":              definition("Task", "tasks", map[string]any{"task_state": enumeration(true, "planned", "planned", "ready", "in_progress", "blocked", "in_review", "done", "cancelled"), "priority": enumeration(true, "normal", "high", "normal", "low"), "owner": field("string", false), "depends_on": map[string]any{"kind": "reference-list", "target_types": []string{"Task"}}, "requirements": field("string-list", false)}),
		"definitions/decision.md":          definition("Decision", "decisions", map[string]any{"decision_state": enumeration(true, "proposed", "proposed", "accepted", "superseded", "rejected")}),
		"definitions/reference.md":         definition("Reference", "sources", map[string]any{}),
	}
	defs["definitions/factory.md"]["filename"] = "factory.md"
	defs["definitions/factory-control.md"]["filename"] = "write-coordination.md"
	defs["definitions/task.md"]["template"] = "../templates/task.md"
	return defs
}

// Definitions returns virtual base defaults. Software-specific types are opt-in.
func Definitions() map[string][]byte {
	out := map[string][]byte{}
	for p, m := range baseDefinitions() {
		out[p] = encode(m, "\n# "+knowledge.String(m, "title")+"\n\nFields describe document structure. Workflow instructions belong in Markdown.\n")
	}
	return out
}

// Files returns a fresh initial bundle. The caller must validate the complete
// overlay, detect all conflicts first, and create each absent path exclusively.
func Files(software bool, projectID, projectName string) map[string][]byte {
	out := Definitions()
	out["factory.md"] = encode(map[string]any{"type": "Factory", "title": projectName, "project_id": projectID, "status": "draft", "factory": map[string]any{"format": 1}}, "\n# "+projectName+"\n\nProject knowledge lives in this bundle. Structural validation does not establish truth, approval or completion.\n")
	out["control/write-coordination.md"] = encode(map[string]any{"type": "Factory Control", "title": "Writer coordination", "status": "stable", "factory": map[string]any{"format": 1}}, "\n# Writer coordination\n\nFactoryctl holds a native advisory lock on this stable file during a single-record write. Do not replace it during writes. Editors and Git do not share this lock; avoid concurrent Git operations. This file contains no lease or mutable lock state.\n")
	out["templates/task.md"] = recordTemplate("Task", "task_state: planned\npriority: normal\ndepends_on: []\n")
	out["agents/builder.md"] = encode(map[string]any{"type": "Agent Profile", "title": "Builder", "status": "draft", "factory": map[string]any{"format": 1}}, "\n# Builder\n\nRead the project documents, implement the requested work, test the result, and distinguish observations from assumptions. Runtime selection and permissions belong to the contributor.\n")
	out["prompts/discuss.md"] = encode(map[string]any{"type": "Prompt", "title": "Discuss record", "applies_to": []string{"*"}, "inputs": map[string]any{}, "conversation": "fresh", "status": "draft", "factory": map[string]any{"format": 1}}, "\n# Discuss record\n\nDiscuss the selected record without interpreting structural validity as approval.\n\n# Template\n\n```factory-prompt\nRead and discuss the selected record in project {{ .Project.Name }}.\n{{ if .Record }}{{ .Record.Body }}{{ end }}\n```\n")
	if software {
		for _, typ := range []struct{ name, dir string }{{"Specification", "specifications"}, {"Design", "designs"}} {
			m := definition(typ.name, typ.dir, map[string]any{"requirements": field("string-list", false)})
			out["definitions/"+lower(typ.name)+".md"] = encode(m, "\n# "+typ.name+"\n\nDescribe the intended software behavior and distinguish implementation evidence from proposals.\n")
			out["templates/"+lower(typ.name)+".md"] = recordTemplate(typ.name, "")
		}
	}
	return out
}
func lower(s string) string {
	switch s {
	case "Task":
		return "task"
	case "Specification":
		return "specification"
	case "Design":
		return "design"
	}
	return s
}
func recordTemplate(typ, extra string) []byte {
	return encode(map[string]any{"type": "Record Template", "title": typ + " creation template", "for_type": typ, "inputs": map[string]any{"title": field("string", true), "summary": field("text", true)}, "status": "draft", "factory": map[string]any{"format": 1}}, "\n# "+typ+" creation template\n\n# Template\n\n```factory-template\n---\ntype: "+typ+"\ntitle: {{ yaml .Input.title }}\nstatus: draft\nfactory:\n  format: 1\n"+extra+"---\n\n# {{ .Input.title }}\n\n{{ .Input.summary }}\n```\n")
}
