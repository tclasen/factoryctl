---
type: "Specification"
title: "Runtime document types, tools and lifecycle hooks"
description: "Runtime document types, tools and lifecycle hooks."
status: "draft"
factory:
  format: 1
  assertion:
    frame: "desired"
    certainty: "proposal"
    basis: "llm-inference"
sources:
  - id: operator
    resource: "../sources/operator-brief.md"
    title: "Operator requirements and scope"
    author: "human:operator"
---

# Runtime document types, tools and lifecycle hooks

The operator-confirmed requirement is runtime extensibility of document types and their behavior. Declarative packages with optional WebAssembly are the agreed architecture direction. The following contract is a proposed implementation design, not a working plugin API. It extends [definitions](knowledge-format.md), [interfaces](interfaces.md) and [validation](validation.md).

## Core and extension responsibilities

The Go binary owns the CLI, local HTTP/HTML GUI, MCP server, harness adapters, package loader, embedded WASM engine, generic rendering and safe storage. It embeds default UI assets/templates and needs no Python, Node, separate WASM runtime or native plugin library. External harnesses and the browser are explicit integrations. Build-time tools are not runtime dependencies.

Packages supply type definitions, schemas, templates, presentation hints, tool schemas, lifecycle transitions and optional custom validation/hooks. Domain types are not Go switch statements. Tasks might expose assign/complete, decisions propose/supersede, and research records cite/review. These are examples, not mandatory policy. Factoryctl runs only explicitly requested operations and their bounded hooks; agents retain planning, scheduling, Git and autonomous loops.

## Package shape and storage

Keep human-readable manifests as Record Definition nodes under `factory/definitions/`, extending the existing frontmatter with an optional `extension` mapping. A package may group several definitions by package ID/version; all members must agree on identity and API compatibility. A definition-only package needs no `.factory/` directory. Optional compiled modules and presentation assets live under `.factory/extensions/<package-id>/`, outside the OKF graph; paths are checkout-relative, contained and checked against traversal/symlink escape. No root project config, registry database, generated index, transition ledger or automatic package download is introduced. Installing a shared package means explicitly adding its files, normally through the user's agent and Git.

The proposed manifest contains package ID/version, host API version, namespace, type/tool declarations, input/output schemas, lifecycle guards, hook entrypoints, asset paths/digests and requested capabilities. The exact serializable schema and WASM ABI must be pinned and validated in the extension prototype before implementation consumers depend on them. Existing definitions lacking `extension` remain supported declarative types. Reuse the existing field-kind contract; tool envelopes use JSON Schema. Complex domain validation belongs to optional WASM, not arbitrary embedded shell expressions.

Keep existing local/default definition resolution: local definitions replace embedded defaults; explicit disablement suppresses defaults; duplicate types remain errors. Tool namespaces must be unique and may not shadow core tool names. Reject an incompatible API, duplicate registration, missing module, invalid digest, schema or capability before activation. Do not silently skip a required validator and allow writes. Affected records remain readable with diagnostics; mutations requiring unavailable rules fail closed. Unrelated types remain usable.

## Activation and runtime changes

Load packages at startup and support explicit reload after installation/editing, without rebuilding or restarting the binary. Validate a complete candidate registration set, then switch the in-memory registry atomically. Invalid candidates do not partially register tools. If on-disk definitions/modules no longer match the active snapshot, affected mutations stop until a valid reload; do not quietly continue enforcing stale rules. Package removal or disablement removes specialized actions/tools but retains generic inspection. Rule changes revalidate affected documents and report incompatibilities; never automatically migrate stored content.

Each request pins registry identity, type rules and module digest. Recheck these identities before persistence; a concurrent reload produces a conflict rather than applying obsolete rules. Existing conversations can continue, but stale tool calls fail explicitly. Advertise tool-list changes when the negotiated MCP capability supports it; adapters must qualify whether the harness refreshes tools. Otherwise show that reconnect/restart is required. Never claim a newly installed tool is available to a connected client until observed.

## Shared operation path

Type-specific tools are named with a validated namespace, for example `research_note_review`; exact encoding follows the qualified MCP naming contract. Each exposes a description, explicit input/output schema, accurate annotations and required operation capabilities. Tools resolve to deterministic extension operations. Prompt-backed GUI actions are separately labeled and launch the selected harness; they do not masquerade as deterministic validation hooks.

Read/context returns exact document content, revision, effective definition, applicable read/edit rules, available operations and package identity. Declarative read projections and bounded pure read/context hooks may add context or derived views, but cannot mutate files, hide diagnostics, execute source content or falsify raw bytes/revisions. Raw inspection remains available for invalid or disabled types.

Every create/update/delete, including raw replace, field editing, board moves and generic MCP tools, passes through the same active type rules. A type-specific operation can propose a change only to its selected record in v1. Hooks cannot write files directly or recursively invoke mutations. Multi-record operations and migrations remain explicit agent work, followed by validation; no false multi-file atomicity guarantee.

## Mutation lifecycle and failure semantics

1. Resolve launch scope and operation capabilities; pin the registry/module identity and validate tool inputs.
2. Acquire the existing checkout writer lock. Refresh document/dependency bytes and compare the expected revision (exclusive existence check for create).
3. Evaluate declared transition guards and bounded pre-write hooks on the original and proposed record. Hooks return diagnostics or a candidate patch; apply hooks in manifest order once, without recursive hook execution. Declarative guards and final validation must also accept the resulting candidate.
4. Validate final bytes, links, affected dependents and type rules against the proposed graph; recheck package/path identities. Persist using the existing single-record atomic storage contract.
5. Release the writer lock and run bounded post-write hooks against the committed revision. Post hooks may produce diagnostics or follow-up suggestions, never a second hidden mutation or a harness/network/process side effect. Refresh views and return committed bytes/revision with hook outcomes.

Reject, trap, timeout or resource exhaustion before persistence leaves original bytes intact. Post-hook failure returns `committed: true`, the new revision and an explicit hook failure; never report that nothing happened or invite blind replay. A crash after persistence may prevent notification/post-hooks. There is no durable hook queue or exactly-once promise: reread state and reconcile before retrying an uncertain operation. Cancellation does not imply rollback. Dry-run validation does not run post-write hooks and does not promise a later write will still succeed.

## Capabilities, limits and enforcement

Use an embedded pure-Go WASM engine with a narrow versioned host API. No ambient WASI filesystem, process execution, network, environment variables, credentials or shell/Git imports. Grant only bounded scoped reads and returned patches/diagnostics needed for the operation. Constrain memory, execution time, input/output size and call depth; qualify interruption on the chosen engine before shipping. No native plugin loading, downloaded executable helpers or unrestricted FFI. Installing code is explicit and trust binds to its content digest; repository checkout or Markdown rendering alone never executes it. Changed code/capabilities invalidate executable trust. Separate browser custom-view trust from server-side WASM trust.

MCP validates operations routed through factoryctl; it is not an OS security boundary. A harness with direct filesystem access can bypass hooks. The adapter must expose whether protected-document access restrictions are verified, unavailable or advisory. Enforced mode requires a qualified permission/sandbox arrangement; if unavailable, refuse to claim enforced mode. External agent Git operations can also replace records, so on-disk changes are revalidated and surfaced, never assumed compliant. No request to install or alter harness-wide configuration is implicit in browsing a repository.

## Presentation and acceptance

Use core-rendered list/detail/board views and declarative actions first. Preserve the existing optional trusted HTML/CSS/JavaScript feature in [templates](templates.md); custom presentation never bypasses server validation. Ordinary Markdown stays inert.

[Extension prototype](../tasks/12-extension-contract.md) and [A14](acceptance.md#a14-runtime-document-extensions) qualify two different type packages, then a third added while the same binary runs. Validate through CLI/GUI/MCP and a real supported harness. Pin API/ABI, test failures and record observed capabilities before declaring implementation ready.

