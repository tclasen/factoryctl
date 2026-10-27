---
type: "Specification"
title: "OKF knowledge format and record definitions"
description: "OKF knowledge format and record definitions."
status: "draft"
factory:
  format: 1
  assertion:
    frame: "desired"
    certainty: "proposal"
    basis: "llm-inference"
sources:
  -
    id: "okf"
    resource: "https://github.com/GoogleCloudPlatform/open-knowledge-format/blob/main/SPEC.md"
    title: "Open Knowledge Format v0.2"
    author: "GoogleCloudPlatform"
  -
    id: "operator"
    resource: "../sources/operator-brief.md"
    title: "Operator requirements"
---

# OKF knowledge format and record definitions

This is the proposed factoryctl format version 1, layered on [OKF v0.2](https://github.com/GoogleCloudPlatform/open-knowledge-format/blob/main/SPEC.md). The upstream specification is authoritative for OKF; factoryctl-specific constraints are identified below. The inspected upstream file blob was `c06e3eede0c910d0ecf12524c34204156f8795ac`. The format and templates in this repository are the initial conformance corpus, not proof of an implemented reader.

## Bundle and identity

`factory/` is the bundle root. Only UTF-8 `.md` concept files belong in it. Each has one YAML frontmatter mapping and Markdown body. `type` is the only required OKF field. No `index.md`, `INDEX.md`, or `log.md` is generated. Case-insensitive variants of these names are reserved in the factoryctl producer profile. Asset files, source code, build products, personal preferences and native harness state stay outside the bundle; no project config file is required.

OKF concept identity is the exact POSIX path relative to `factory/`, without `.md`. Display titles do not determine identity. Standard relative Markdown links are preferred for Git-host readability; `/...` links are bundle-relative as OKF specifies. Filenames created by factoryctl use lowercase kebab-case. Existing legal filenames are read without renaming. Two paths differing only by case are an ambiguity diagnostic, with writes to those paths disabled for cross-platform portability.

A project has a portable `project_id` on its Factory node as specified by [project discovery](project-discovery.md). In-memory identities are `(project_id, checkout-key, concept-path)`. Checkout-key is the canonical local working-tree path, never persisted. Clones share project identity; the GUI still distinguishes checkout locations. Two loaded checkouts with the same project ID require explicit selection for cross-project mutations; never silently choose a worktree.

Cross-project references use ordinary absolute URLs to another project's optional `bundle_uri` on its Factory node. Factoryctl maps that exact prefix to a loaded bundle. Remaining paths/anchors are decoded once, normalized and checked for containment. Other URLs remain external links. No remote fetching occurs merely to render a node. Missing projects/targets are unresolved references, not invalid OKF. Backlinks cover the currently loaded snapshot only and show that scope. Moves require updating incoming paths; v1 does not promise permanent IDs through renames or infer a migration from a Git rename heuristic.

## Frontmatter extension

Preserve unknown keys, key order, comments and body text when editing known fields where the parser permits; if a representation cannot be edited losslessly, offer a visible whole-document preview before a human saves, and return `lossless_edit_unavailable` to a field-only API request. Never quietly reserialize unrelated content. Reject duplicate YAML keys, multiple YAML documents, custom tags, non-string mapping keys, alias expansion and over-limit inputs as parse/producer errors. Arbitrary unknown *field names and type values* remain permitted.

OKF's `status` is reserved for `draft | stable | deprecated`. It is not ticket progress. `generated`, `verified`, `sources`, `stale_after`, `resource`, `title`, `description` and `tags` retain upstream semantics. Unknown lifecycle values produce a warning rather than erasing a document. A bare `verified` mapping is normalized to a list for display and preserved on disk unless edited. Never manufacture `verified` from commits, authorship, file age or a successful syntax check.

The optional `factory` mapping contains:

| Field | Meaning |
| --- | --- |
| `format` | Factory extension version, currently integer `1` |
| `assertion` | Record-level epistemic defaults, specified in [epistemics](epistemics.md) |

Custom domain fields such as `task_state`, `priority`, `owner`, `decision_state`, `implements` and `depends_on` are top-level producer extensions. A raw OKF document without `factory` remains readable and editable as a generic node. The application must not reject an entire bundle because of one unknown type or factory schema error.

## Record definitions: the T-box

Definitions are ordinary OKF concepts of `type: Record Definition` under `definitions/`. They are themselves listed, inspected, linked and edited. A small built-in meta-schema bootstraps parsing definitions and infrastructure nodes, even if their display definitions are disabled. Definitions extend factoryctl's typed editor, not OKF validity itself.

```yaml
type: Record Definition
title: Task record definition
defines: Task
directory: tasks
fields:
  task_state:
    kind: enum
    required: true
    values: [planned, ready, in_progress, blocked, in_review, done, cancelled]
    default: planned
  depends_on:
    kind: reference-list
    target_types: [Task]
    required: false
views:
  list: [title, task_state, priority, owner]
  board: {field: task_state}
template: ../templates/task.md
actions: [../prompts/discuss.md, ../prompts/request-status.md]
```

`defines` is an exact, case-sensitive type string. `directory` is bundle-relative, contains no traversal and owns the default creation location; readers still show imported records elsewhere, but factory validation reports incorrect placement and structured writes must use the declared directory. An optional `filename` fixes the basename for singleton infrastructure records. One active definition per type. Duplicate definitions produce a conflict diagnostic for that type, never last-file-wins. Two different types may share a directory. Names do not control authority.

Supported field kinds in v1: `string`, `text`, `boolean`, `integer`, `number`, `date`, `datetime`, `enum`, `string-list`, `reference`, `reference-list`, `mapping`. Definitions may also specify nonnegative `min_items` for list fields. Dates are ISO calendar dates; datetimes require offsets; retain original textual YAML representation. Kinds validate shape, not truth. Properties: `required` (default false), `default`, `description`, `values` (enum only), `target_types` (references only). Additional schema features such as expressions, inheritance, custom validators and arbitrary JSON Schema are deferred. `mapping` accepts a plain mapping without executable validation. Nested field paths are not supported by definitions v1; `factory` has a built-in structural schema.

Reference fields contain normal relative paths or absolute URLs, with optional anchors. Their outgoing edges are combined with body links, `resource`, and `sources[].resource` when resolvable. Display relationship names from field names; body links are untyped relationships. Deduplicate backlinks by source node and show individual occurrences/labels on demand. Unresolved references stay visible. Dependency cycles are reported; the application does not enforce execution readiness or approval policy.

`views.list` chooses fields; every type receives default list/detail views without it. `views.board.field` must identify an enum; columns follow `values` order. No stored board index or duplicated ticket status. Invalid/missing values appear in a visible diagnostic column. A combined board requires selected types/definitions to agree on field and enum values; otherwise show separate boards instead of guessing an ordering.

`template` and `actions` reference OKF template/prompt nodes. Actions may also be discovered through their applicability lists. Every action is available in the detail view when applicable; declarations are deduplicated by node identity.

## Resolution and disabling

Initialization materializes chosen built-in definitions/templates/prompts as versioned OKF nodes. At runtime a valid local definition replaces a same-named embedded default as a whole, without deep merging. If no local definition exists, the enabled embedded default remains available as a labeled virtual entry; it is not falsely presented as a persisted OKF node. Missing explicit custom references are diagnostics, not silent fallback. A local definition with `enabled: false` disables type-specific creation, actions and specialized views, but existing nodes remain visible generically. Omitted `enabled` defaults to true. A disabled local definition suppresses the same-type embedded fallback; it is not ignored during resolution. Two local definitions for the same type are a conflict even if one is disabled. The metadata must be a boolean. Disabling does not erase its schema or skip structural checks on existing records; only specialized interaction is disabled. The same rule applies when removing a starter set. Upgrades never overwrite local files on startup.

## Safe storage

The storage contract is single-record atomic mutation, not a multi-record transaction. APIs require a SHA-256 revision token of exact file bytes for update/delete. Under a per-checkout advisory writer lock, re-read bytes, compare the token, construct and validate the proposed post-write graph including structurally affected incoming dependents as specified in [continuous validation](validation.md), write a temporary sibling, flush, and atomically replace using the platform's supported operation. Creation is exclusive and never replaces an existing path. Deletion requires an expected revision and explicit user/agent action; deletion is blocked when known incoming links would become dangling, returning the affected references.

The stable lock anchor is `control/write-coordination.md`, itself a human-readable OKF `Factory Control` node explaining its purpose. Factoryctl holds a native file lock on it, never stores lock ownership or transitions in its text and never replaces it during normal writes. Verify the locked file's identity still matches the path before mutation; fail/retry if a Git operation replaced it. Lock handles release on process exit. No persistent lease, database or event log is created. The kernel lock coordinates factoryctl processes, not arbitrary editors or Git. Agents must avoid Git checkout/reset operations while a mutation is active. CAS plus atomic replacement cannot provide a universal transaction with a non-cooperating writer; this limitation is explicit and tested where reproducible.

Symlinks/reparse points inside the bundle and asset trees are never followed for writes; all path components and final targets must remain under the selected root. Reject NULs, absolute filesystem paths, traversal, ambiguous encodings and unsupported special files in mutation paths. Read-only `resource` references may identify repo-contained project artifacts outside the bundle; those are resources, not concept nodes, and never authorize writes or automatic file serving. A record with Git conflict markers is shown as conflicted; typed edits and agent actions on it are disabled until the agent/user resolves it. Never attempt semantic Git merges.

Proposed limits: 1 MiB per concept, 256 KiB frontmatter, 32 nesting levels and 10,000 nodes in the benchmark corpus. Over-limit files are listed with diagnostics and not parsed for actions. Limits are documented built-in constants, not persisted configuration values or silent truncation. File-watch events trigger rescans; portable polling is the correctness fallback. Snapshot data, search indexes and reverse edges live only in memory.

The [continuous validation contract](validation.md) adds strict producer checks for links, anchors, fields and placement while retaining permissive OKF reading. See [interface errors](interfaces.md), [concurrency acceptance](acceptance.md) and [format implementation task](../tasks/02-knowledge-store.md).
