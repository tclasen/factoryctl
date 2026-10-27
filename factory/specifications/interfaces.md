---
type: "Specification"
title: "CLI, MCP and local HTTP interfaces"
description: "CLI, MCP and local HTTP interfaces."
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

# CLI, MCP and local HTTP interfaces

These are proposed version-1 contracts. The [product requirements](product.md) establish the interfaces; no commands below are implemented yet. All adapters call one application service layer and [validator](validation.md).

## CLI grammar

Global flags: `--path DIR` (default CWD), `--project ID` (required for ambiguous writes), `--json`, `--non-interactive`. Flags may appear before or after subcommands. Human output is default; JSON is stdout-only structured output. Diagnostics/progress go to stderr. Do not persist CLI output automatically.

| Command | Behavior |
| --- | --- |
| `projects list` | Discover and report projects, initialization and diagnostics |
| `init --project-id ID [--starter software]` | Create absent defaults in one Git working tree |
| `records list [--type TYPE] [--query TEXT] [--limit N] [--cursor C]` | Search/filter a snapshot across selected projects |
| `records show PATH` | Record, revision, metadata, outgoing links and backlinks |
| `records create --type TYPE --slug SLUG [--input KEY=VALUE]` | Collect/validate template inputs, preview or create exclusively |
| `records update PATH --expected-revision HASH --set KEY=VALUE` | Typed frontmatter edit; repeated set flags allowed |
| `records replace PATH --expected-revision HASH --file FILE` | Replace complete document, validated before writing; `-` means stdin |
| `records delete PATH --expected-revision HASH` | Explicit deletion, reject known incoming links |
| `validate [--level okf\|factory] [--check-external] [--require-complete-links]` | Structural validation with precise scope/completeness |
| `templates diff` / `templates export --kind view` | Inspect defaults/materialize requested optional assets |
| `actions list --record PATH` | Discover applicable prompts without running them |
| `actions run PROMPT --record PATH [--input KEY=VALUE] [--runtime ADAPTER]` | Start runtime conversation; terminal streams events, stdin supplies follow-ups |
| `gui [--listen HOST:PORT] [--no-open] [--trust-assets]` | Serve embedded UI; default loopback and OS-assigned port |
| `mcp` | Serve MCP over stdio; stdout contains only protocol messages |
| `doctor` | Report versions/capabilities, effective personal runtime preferences and paths (without secrets) without mutating Git/runtime state |

Runtime actions, `gui` and `doctor` additionally accept `--runtime`, `--runtime-executable`, `--model`, `--reasoning-effort` and `--preferences` as defined by [personal runtime preferences](personal-runtime.md). These are contributor choices, never saved to project records.

`KEY=VALUE` accepts JSON values, falling back to a literal string; no shell evaluation. Inputs come from flags or stdin; noninteractive mode fails on missing required input. Slugs are explicit to avoid surprise collisions; default interactive suggestions can include a random suffix, never a shared numeric counter. IDs remain paths. No mutation retries without rereading state.

Exit codes: `0` requested operation succeeded (read commands may include diagnostics), `1` validation failed, `2` invalid arguments/personal preferences/project selection, `3` conflict/stale revision, `4` I/O or runtime unavailable/failure, `5` incomplete validation when completeness was required. Signal interruption exits conventionally. `validate --json` always includes `scope`, `level`, `complete`, `errors`, `warnings`, `unchecked_links` and diagnostics; zero errors is not a claim of complete network validation.

## Shared DTOs and mutation semantics

`NodeRef = {project_id, checkout?, path}` where path includes `.md` relative to the bundle. `checkout` is required when project ID resolves to multiple open copies, and is a process-local opaque selection token. `Node = {ref, revision, type, title, frontmatter, body, diagnostics}`. Revision is SHA-256 of exact bytes. Response data never presents a cached revision as the result of a new write before disk persistence succeeds.

List inputs: `project_id?`, `type?`, `query?` (literal case-insensitive text), `filters?` (exact scalar field matches), `limit` (1–200, default 50), `cursor?`. Sort by project name, type, title, path, with identity as tie breaker. Cursors bind to the in-memory snapshot; a changed snapshot returns `cursor_expired` and a restart instruction. No durable cursor store. Backlink pagination uses the same rules.

Create inputs: `project_id`, `checkout?`, `type`, `slug`, `inputs`, optional `template`. Update inputs: `ref`, `expected_revision`, and either `set` mapping of known top-level fields or full `document` text, never both. Creation returns node/revision; update returns node/revision; delete returns removed ref. Creation rejects existing paths; replacement requires current revision. Agent tools do not automatically commit.

Errors share `{code, message, retryable, diagnostics?, current_revision?}`. Stable codes include `not_found`, `ambiguous_project`, `invalid_path`, `invalid_document`, `invalid_location`, `broken_link`, `validation_incomplete`, `stale_revision`, `already_exists`, `referenced_node`, `schema_conflict`, `lossless_edit_unavailable`, `cursor_expired`, `runtime_unavailable`, `unsupported_capability`, `cancelled` and `io_error`. A retryable error does not authorize automatic replay of an uncertain mutation.

## MCP tools

Use standard MCP initialization/capability negotiation and stdio transport. Qualify a supported protocol version against current Codex during the integration spike; the researched baseline is the [2025-11-25 transport specification](https://modelcontextprotocol.io/specification/2025-11-25/basic/transports). Never print banners or application logs on stdout.

| Tool | Inputs | Output |
| --- | --- | --- |
| `factory_projects_list` | none | scoped project summaries and diagnostics |
| `factory_definitions_list` | project selection | effective definitions, source/local/default status |
| `factory_nodes_list` | list inputs above | nodes, next cursor, snapshot identity |
| `factory_node_read` | `ref` | full node plus source classification |
| `factory_node_links` | `ref`, limit/cursor | incoming/outgoing edges and unresolved scope |
| `factory_node_create` | create inputs above | created node/revision |
| `factory_node_update` | update inputs above | updated node/revision |
| `factory_node_delete` | `ref`, `expected_revision` | removed ref |
| `factory_validate` | project/ref selection, level | structural diagnostics; network checks excluded from default tool |
| `factory_prompt_render` | prompt ref, record ref, inputs | prompt/context preview without runtime launch |

Tool schemas use JSON Schema object shapes with required fields, explicit optional fields and `additionalProperties: false` for tool envelopes. Frontmatter/input maps remain extensible according to their record definitions. Domain failures return MCP tool errors with readable messages and structured diagnostics. Annotate read-only/destructive/idempotent tool characteristics accurately, but do not treat annotations as access control. Tools operate only in their launch scope; MCP cannot request arbitrary new filesystem roots. Starting more agent sessions is the harness's responsibility, not an MCP recursion endpoint in v1.

## Local HTTP contract

The GUI uses `/api/v1/projects`, `/definitions`, `/nodes` (GET/POST), `/node` (GET/PUT/DELETE with NodeRef query/body), `/links`, `/validate`, `/actions`, `/actions/render`, `/conversations` (POST/GET), `/conversations/{id}/messages` (POST), `/conversations/{id}/interrupt` (POST), and `/conversations/{id}/requests/{requestId}` (POST response to native input/approval). These use the shared DTOs. Conversation IDs are process-local mappings to native thread IDs; no transcript persistence.

`/api/v1/events` is an SSE stream of `snapshot.changed`, `diagnostics.changed`, `conversation.item`, `conversation.state`, and `runtime.request` envelopes with project/conversation identity. Ordering is per conversation, not a global durable event log. On disconnect/reconnect fetch a fresh snapshot and native conversation state; missed events are not proof no action happened. Use normal POST for bidirectional user input.

HTTP status mapping: 400 malformed arguments, 404 missing resource, 409 stale revision/ambiguous target/reference conflict, 422 invalid/incomplete document, 503 unavailable runtime, 500 unexpected I/O. Errors carry the common envelope. Require exact allowed Host/Origin checks and an in-memory session nonce on mutations; reject cross-site requests even on localhost. No permissive CORS. Default binding is loopback; explicit non-loopback bind is for a trusted externally protected deployment. No login/accounts database is introduced.

Node and asset routes enforce canonical containment and MIME types; do not serve arbitrary repo paths, `.git`, native credentials or source trees. Browser state and conversation buffers are in memory. UI raw edits are user-supplied text, not executable templates until explicitly installed and trusted.

## Runtime type tools and common enforcement

In addition to the core tools above, [runtime packages](extensions.md) register namespaced type-specific tools with validated input/output schemas. CLI and HTTP expose equivalent extension operation discovery/invocation through the shared application service; exact command/route spelling is finalized in the extension prototype. GUI buttons distinguish prompt-backed harness launches from deterministic type operations.

Core reads include effective definition, editing rules, available tools and registry identity. Every generic writer, including replace/delete, applies the same enabled declarative field and lifecycle rules as specialized tools. Required invalid/unavailable extensions block affected writes; generic APIs are not an escape hatch. Tool-list change notification or explicit reconnect is part of adapter qualification. Responses carry current revision/definition identity; agent-action status is separate from document persistence. No custom backend handler or executable callback may be registered. MCP tool annotations and agent instructions do not establish filesystem access control; expose verified/advisory/unavailable enforcement capability explicitly.
