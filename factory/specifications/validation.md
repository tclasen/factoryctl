---
type: "Specification"
title: "Continuous structural validation"
description: "Continuous structural validation."
status: "draft"
factory:
  format: 1
  assertion:
    frame: "desired"
    certainty: "proposal"
    basis: "llm-inference"
sources:
  -
    id: "operator"
    resource: "../sources/operator-brief.md"
    title: "Operator requirements and scope"
    author: "human:operator"
---

# Continuous structural validation

**Operator requirement:** validate syntax, mandatory fields, types, links and placement while documents are written. This is product logic, independent of whether an agent follows a workflow. It supplements [OKF conformance](knowledge-format.md) without redefining OKF's permissive consumer rules.

## Two explicit validation levels

`okf` checks upstream conformance and parses generic knowledge. Missing/unknown factory definitions and broken links do not make an otherwise valid concept invalid OKF. `factory` (the default for validation and mutation) additionally checks the project's declared record definitions, location rules, template references and link integrity. An invalid factory document remains visible/readable with diagnostics; neither a type error nor a broken link makes the entire project disappear.

All factoryctl writes validate the proposed final document against the current snapshot before persistence. Known-type writes require the definition's mandatory fields, types and directory. Generic unknown-type writes require OKF syntax, safe placement and link checks but cannot claim schema validity: the result explicitly reports `schema: unavailable`. Optional metadata remains optional. Changing a definition revalidates affected instances; block a definition edit that would make existing instances invalid, returning the affected paths. An agent may perform a deliberate multi-file migration externally and then validate it; v1 has no implicit schema migration engine.

## Required checks and diagnostics

| Area | Checks |
| --- | --- |
| Parse | UTF-8, one YAML mapping, delimiters, duplicates/tags/aliases, resource limits, Markdown parse, conflict markers |
| Fields | Required presence, null/empty distinctions, kind, enum membership, supported extension version; preserve unknown metadata |
| Placement | Inside `factory/`, correct definition directory, fixed filename when declared, no reserved indexes/logs, case collisions, symlink escapes or asset files |
| References | Local target exists, safe resolution, exact path case, referenced type when declared, anchor exists, source footnote IDs resolve |
| Templates | Exactly one payload source, valid template syntax, declared inputs, resolvable partials, no include cycle, well-formed rendered output |
| Definitions | Unique active `defines`, allowed kinds, valid default values, board field is enum, valid template/profile/action references |
| Configuration | Known config version/keys, fixed storage boundaries, unique project selection, runtime/profile references when used |

A required string must be nonempty after trimming; required lists may be empty unless the field defines `min_items` (supported nonnegative integer, default 0). Required mappings may be empty. Null fails a required field and is treated as absent for an optional field. Numbers are finite; integers exclude booleans. References may include a heading fragment. Generate heading anchors with a documented GitHub-style algorithm (lowercase, strip punctuation except hyphen/underscore, spaces to hyphens, suffix duplicates `-1`, `-2`); recognize explicit HTML anchor IDs without executing HTML. Test Unicode and duplicate headings.

Parse Markdown ASTs, not regex alone: ignore links inside code fences/inline code, resolve reference-style links and images, and match footnotes to `sources[].id` when used as provenance. Ordinary internal footnotes without a source ID remain normal Markdown; do not invent provenance. Frontmatter references follow definition kinds and OKF path-valued fields; a `sources.resource` scope descriptor is not automatically a broken URL.

Link results are `valid`, `broken`, or `unchecked`. In-bundle links and mapped loaded-project links are checked against disk. A known cross-project URI whose checkout is absent is `unchecked`, and strict affected-record writes return `validation_incomplete` rather than claim success. External HTTP(S) URLs are syntax-checked by default but reachability is `unchecked`; they do not prevent ordinary offline edits. `validate --check-external` explicitly enables bounded network checking and reports redirects/timeouts/auth-required endpoints without treating them as verified content. This is a proposed offline-friendly interpretation of link validity, not a guarantee that external resources exist. `--require-complete-links` makes unchecked link results fail CI. Do not automatically crawl or execute links while browsing. External fetch checks must reject unsafe URL schemes and private/link-local destinations unless the operator explicitly scopes them in that invocation.

## Write and watch behavior

GUI validates during editing (debounced), associates results with the draft revision, and disables Save while affected-record errors remain. Server validation always reruns on Save; client checks cannot bypass it. CLI/MCP return the same diagnostic codes before writing. Disk watcher/polling revalidates external edits, invalidates dependent templates/backlinks and shows file/line errors. A partially written external file is retried after debounce; persistent errors remain visible. Do not roll back or overwrite an external editor's content.

A stale revision returns conflict before trying to repair the document. Deletion is blocked when it would break currently known incoming references; offer those references to the user/agent for explicit repair. Rename/multi-record migration is deferred to the agent's file tools and Git, followed by full validation. A valid edit elsewhere remains possible when an unrelated existing record is invalid; report both affected-record validation and project-wide health separately.

Diagnostics have `code`, `severity`, `project_id`, `path`, `line`, `column`, optional `field`, optional `target`, and `message`. Error messages explain the expected value/location and the actual problem without echoing secrets. `valid` means structurally valid within the reported scope, not true, approved, implemented or compliant with a government solicitation.

## Verification

Negative fixtures cover missing required fields, string-vs-integer errors, enums, wrong directories, malformed frontmatter, duplicate definitions, missing/incorrect-case targets, broken anchors, missing cross-project checkouts, include cycles, stale edits and conflict markers. Tests must prove rejected writes leave original bytes unchanged. A direct external edit must produce the same error through the watcher. See [task 03](../tasks/03-definition-validation.md) and [acceptance A03](acceptance.md#a03-structural-validation).
