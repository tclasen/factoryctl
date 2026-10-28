package mcp

func object(properties map[string]any, required ...string) map[string]any {
	v := map[string]any{"type": "object", "properties": properties, "additionalProperties": false}
	if len(required) > 0 {
		v["required"] = required
	}
	return v
}
func str() map[string]any { return map[string]any{"type": "string"} }
func tool(name, description string, schema map[string]any, readOnly, destructive bool) map[string]any {
	return map[string]any{"name": name, "description": description, "inputSchema": schema, "annotations": map[string]any{"readOnlyHint": readOnly, "destructiveHint": destructive, "openWorldHint": false}}
}

// Tools intentionally never advertises an answer setter, host credential or arbitrary destination.
func Tools() []map[string]any {
	ref := object(map[string]any{"path": str(), "project_id": str(), "checkout": str()}, "path")
	choice := object(map[string]any{"id": str(), "label": str()}, "id", "label")
	return []map[string]any{
		tool("factory_user_ask", "Ask the user on the attached trusted surface. Returns a pending request; no answer is implied. Choices permit custom text.", object(map[string]any{"prompt": str(), "kind": map[string]any{"type": "string", "enum": []string{"text", "single_choice", "multiple_choice"}}, "choices": map[string]any{"type": "array", "items": choice, "maxItems": 32}, "previous_request_id": str()}, "prompt", "kind"), false, false),
		tool("factory_user_wait", "Wait up to 30 seconds for a real user submission. Pending/declined/cancelled/expired/unavailable are not answers. Do not invent an answer.", object(map[string]any{"request_id": str(), "timeout_seconds": map[string]any{"type": "integer", "minimum": 0, "maximum": 30}}, "request_id"), true, false),
		tool("factory_ui_present", "Present generated HTML/CSS/JavaScript in the attached first-party sandbox after user trust. Include a plain fallback; attach request_id for an existing question. View interaction cannot submit an answer; the trusted host confirms it. External Apps hosts are unqualified.", object(map[string]any{"title": str(), "fallback": str(), "html": str(), "data": map[string]any{}, "request_id": str()}, "title", "fallback", "html"), false, false),
		tool("factory_projects_list", "List projects in this invocation's fixed scope.", object(map[string]any{}), true, false),
		tool("factory_nodes_list", "List up to 200 records in the selected project; reports truncation explicitly.", object(map[string]any{"type": str(), "query": str(), "limit": map[string]any{"type": "integer", "minimum": 1, "maximum": 200}}), true, false),
		tool("factory_definitions_list", "Read effective structural definitions in the selected project.", object(map[string]any{}), true, false),
		tool("factory_node_read", "Read exact Markdown and current revision from the selected project.", object(map[string]any{"ref": ref}, "ref"), true, false),
		tool("factory_node_links", "Read loaded incoming/outgoing structural links for a record.", object(map[string]any{"ref": ref}, "ref"), true, false),
		tool("factory_node_create", "Create a record from the selected project's template; never overwrites an existing path.", object(map[string]any{"type": str(), "slug": str(), "inputs": map[string]any{"type": "object"}}, "type", "slug"), false, false),
		tool("factory_node_update", "Update a record with its current revision, either a scalar set map or a complete document. Structural validity does not establish approval.", object(map[string]any{"ref": ref, "expected_revision": str(), "set": map[string]any{"type": "object"}, "document": str()}, "ref", "expected_revision"), false, true),
		tool("factory_node_delete", "Delete a selected-project record only with its current revision and safe incoming references.", object(map[string]any{"ref": ref, "expected_revision": str()}, "ref", "expected_revision"), false, true),
		tool("factory_validate", "Validate selected project documents structurally; external links are not fetched.", object(map[string]any{"level": map[string]any{"type": "string", "enum": []string{"okf", "factory"}}}), true, false),
	}
}
