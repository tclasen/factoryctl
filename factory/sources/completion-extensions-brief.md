---
type: Reference
title: Completion extensions implementation contract and evidence
description: Declarative binding syntax, shared APIs, focused evidence and limits for completion revision 2.
status: draft
---

# Completion extensions

This is implementation evidence for completion revision 2, not independent acceptance or a complete milestone qualification. Run completion-extensions owns extension-focused knowledge/app and MCP changes. The Coordinator supplied independent design readiness before implementation. Actual model and cost telemetry are unknown.

## Declarative contract

Record Definition nodes retain their ordinary structural fields. Optional `extension` is a mapping with exactly `version: 1`, `namespace` and `tools`. Namespace is lowercase alphanumeric, starts with a letter, and cannot equal `factory`. Tool-local names start with a lowercase letter and contain lowercase letters, digits and underscores. Full tool names concatenate namespace, underscore and local name. This prevents core shadowing and ambiguous namespace splits.

Each tool has `operation: read | create | update | delete`, optional description, optional inputs using existing structural field specifications, and optional fields mappings. Each field source is exactly `{input: declared_name}` or `{literal: value}`. Read/delete accept no inputs/field mappings. Create controls type and destination through its owning definition and explicit slug; update maps declared fields; read/delete require an owning-type record path. The fixed operation envelope uses path or slug, expected_revision for update/delete, expected_definition for mutations, and optional inputs. Envelope names cannot be mapped or declared as user inputs. Unknown handlers, expressions, predicates, transitions and other binding keys are errors; no executable binding API exists. Literal text and mapping data are never interpreted as code.

Duplicate namespaces disable the affected specialized registry/actions with explicit diagnostics while unrelated tools and structurally valid generic writes remain usable. Invalid structural schemas still block affected generic writes. Disabled types remove specialized tools/actions and creation while preserving generic reads and structurally checked updates/deletion. Definition-only projects need no `.factory/` assets.

## Shared application APIs

`ExtensionCatalog(projects, selected)` returns tools, input schemas, operation/type, definition_hash, snapshot revision, diagnostics and reconnect_required. `CallExtension(projects, selected, name, args)` invokes only the fixed document operations. `DefinitionHash(definition)` includes complete definition source revision, effective metadata and conflict flags, so instruction-only edits also invalidate drafts. `CreateChecked(projects, selected, type, slug, inputs, expectedDefinition)` checks caller identity then sends the same identity to Mutate; Mutate rechecks it under the bundle writer lock. `Mutation` adds ExpectedDefinition and DefinitionType, the latter enforcing specialized scope under the same lock. Existing legacy Create pins the definition it used to render before writing, and legacy Mutate callers still receive fresh structural graph validation even without an identity token.

`Actions(projects, selected, recordPath)` discovers applicability and explicit definition action references with disabled reasons. `RenderPrompt(projects, selected, promptPath, recordPath, inputs)` returns rendered prompt, profile instructions, bounded context and revisions with launched:false. It never starts a runtime. Prompt context includes Record:nil when none is selected, plus Project, Selection and Links; the existing template engine supplies declared Input values.

MCP core create/update/delete require expected_definition; definitions-list supplies hashes plus unknown_definition_hash for generic unknown-type creation. Record reads include raw source, definition, hash, instructions, tools and actions. `factory_extensions_list` reloads bindings; `factory_prompt_render` previews without launch. The same Server dynamically reloads tools/list from disk; initialize and catalog explicitly tell caching clients to reconnect. No tool-list notification capability or actual harness automatic refresh is claimed. Direct filesystem enforcement is reported advisory, not enforced.

## Focused verification

Executed `go test ./internal/app ./internal/knowledge ./internal/mcp` successfully after allowing local listeners for the existing MCP httptest suite. Writable GOCACHE/GOTMPDIR were under `/private/tmp`. Initial sandbox attempts failed at tcp6 loopback listen with operation-not-permitted; no tests were weakened to bypass that restriction. Two command attempts also used the workspace parent instead of the Go module and failed module discovery; corrected invocations used the project directory.

Current named behavioral evidence includes TestRuntimeExtensionsReloadAndSharedStructuralBoundary, TestExtensionRejectsExecutableBindingsAndReservedMappings, TestPromptPreviewNeverLaunchesAndExposesContext, TestDynamicExtensionToolsAndDefinitionPinnedTransport and the pre-existing app/knowledge/MCP suites. Fixtures load materially different Research and Opportunity instructions, add Contract without rebuilding, exercise mapped create/read/update, enum acceptance despite workflow prose, required-field and unknown-enum rejection, preserved bytes on rejection, stale document/definition identities, body-only instruction edits, disablement, duplicate namespaces, generic schema writes during extension conflict, scope isolation, forbidden handler/expression/predicate metadata, and same-server MCP discovery refresh. Existing MCP mutation tests now send required identity tokens.

Focused tests are intermediate evidence only. Full CLI/HTTP/GUI parity, custom browser assets, actual harness refresh/reconnect, current complete requirements demo, broad regression, independent review and PR/main delivery remain Coordinator/integration responsibilities. No Git command, publish, release, artifact upload or external message was performed by this worker.

## Worker handoff checkpoint

Checkpoint: 2026-10-28T04:11:44.707479+00:00. Task project-completion; run completion-extensions; role Developer; native task identity /root/completion_extensions. Project-only authorized source claims are released to the Coordinator for integration and independent review after this handoff. No Git operations were performed; branch, PR, merge and observed remote-main state were not inspected by this worker. Delivery remains pending with the Coordinator. Requested outcome comes from the operator completion instruction recorded in the current implementation plan/task; declarative syntax and APIs are implementation choices approved for bounded construction after fresh design review. No further independent worker continuation is promised.

Latest focused command `GOCACHE=/private/tmp/factoryctl-go-cache GOTMPDIR=/private/tmp go test ./internal/app ./internal/knowledge ./internal/mcp` exited 0. Parent broad regression and end-to-end review remain required. Runtime/cost counters unavailable; commands and tests stayed inside the 120-second/10-minute bounds. Sandbox-local listener failures were diagnosed and an escalated run succeeded.

Tested source SHA-256 values (source freeze for this worker only):

| Path | SHA-256 |
| --- | --- |
| `internal/knowledge/schema.go` | `77db6ae3413e0b38d6658caaf9b6fcc7fd44e674feb8d830b59aeeac17e8a09f` |
| `internal/knowledge/extensions.go` | `c9e2577ba6e10bcf11416f8ba3feb423d8291811c02fd35cea9b2662a6213dbf` |
| `internal/app/snapshot.go` | `bbaa74e91d00fac030eb69cb3974a2d4314f356dbda60906810f70d4437929e2` |
| `internal/app/store.go` | `804d58dda4bbbab4ce0dc5f2ee1036cea83dc1262df0cd80a966e6792475b548` |
| `internal/app/create.go` | `fac4485eee59678d5a5061d9d3a4a5434eaca33df2268ea3855dda1c2a9e12d8` |
| `internal/app/extensions.go` | `7d3bf74a9a8334126df51fecf0b48f972dc28565e5a1f1ac9a9c1e016d3af100` |
| `internal/app/actions.go` | `a0e2d7401f2ebecb270aa581d96b60cdd851d6fe8e6883c94019f28dea422fc3` |
| `internal/app/extensions_test.go` | `dd9a5649ef4476055d486800850e69e9c6d4059a9330d611977fff76721b8a66` |
| `internal/mcp/server.go` | `7367d604bb687f5feeee1c69c403fa22a21e698f628966c52481cceca3bb0a03` |
| `internal/mcp/tools.go` | `8ac8cafe4bcdaca024bb14659538c4f7214d6ea1ada2b72eebe5a9507c5fc5be` |
| `internal/mcp/server_test.go` | `a746b3365ab32388da66878a93a6d7ef0f78cd17838ec6b22e4b91f5eb26281f` |

Exact next safe action: Coordinator runs shared CLI/HTTP/GUI integration and independent final review against the current candidate, then carries out already-authorized repository delivery and full-matrix local demonstration. No completion claim follows from this handoff.
