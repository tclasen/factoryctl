---
type: "Specification"
title: "Runtime document types, tools and lifecycle actions"
description: "Runtime document types, tools and lifecycle actions."
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

# Runtime document types, tools and lifecycle actions

The operator requires runtime-extensible document types, with structural definitions, agent-interpreted Markdown instructions, custom browser UI and prompt-backed agent actions. No user-supplied code runs inside the factoryctl backend. This contract extends [definitions](knowledge-format.md), [interfaces](interfaces.md), [templates](templates.md) and [validation](validation.md). Exact metadata and API syntax below remain implementation proposals.

## Execution boundaries

The Go binary owns document parsing, fixed declarative validation, safe storage, CLI/HTTP/MCP interfaces, rendering and harness adapters. It embeds defaults/UI and requires no backend plugin engine, interpreter, native plugin library or downloaded executable helper. Domain types are data loaded at runtime, not compiled Go cases.

Users may supply arbitrary HTML/CSS/JavaScript for browser interactions, including defining a kanban board with drag-and-drop. Browser code uses the existing scoped message bridge to request document edits or agent actions; it has no direct filesystem, shell or harness access. Server validation remains authoritative. For arbitrary local computation or automation, an action sends a prompt to the contributor-selected AI harness. That harness executes the work using its own permissions and tools. Factoryctl never invokes Git or executes a package's shell command.

## Definition shape and storage

Use ordinary Record Definition nodes under `factory/definitions/`. Extend their optional `extension` mapping with a namespace, supported definition format version, tool declarations and references to Markdown instructions/prompts. Reuse existing fields, templates, prompt references and views rather than introduce a package manager or executable ABI. A group of definitions and optional assets may be shared as files; no separate package registry, dependency installer or root config is required. A definition-only setup needs no `.factory/` directory. Browser assets retain the existing `.factory/views/` and explicit asset path conventions.

Each tool declaration binds a name/description and typed inputs to a supported core operation: scoped read, create, update or delete. Field mappings may use declared input values or literal defaults only; they are not expressions or scripts. Use the existing simple field kinds for documents and JSON Schema for tool envelopes. The prototype pins the exact metadata schema and rejects unsupported operations, expression/code handlers, duplicate namespaces, core-name shadowing and invalid prompt/view references.

Lifecycle meaning is written as ordinary Markdown instructions and prompts for the agent. Factoryctl does not parse allowed source/destination pairs, state-dependent required fields or other workflow predicates. An enum can define state labels and an evidence field can be structurally required for all records, but the core does not decide whether a task may advance or whether evidence is sufficient. Unsupported workflow-enforcement declarations cannot be advertised as active validation; preserve unknown metadata under the normal format contract and diagnose unsupported declarations when used as a tool/validation binding. A populated field or agent assessment does not prove truth or approval. The [minimal-interface decision](../decisions/minimal-interface.md) controls this boundary.

## Resolution and reload

Preserve existing local/default definition resolution and disablement semantics. Validate a complete candidate definition set before atomically replacing the in-memory registry on reload. Invalid definitions remain visible with diagnostics; affected typed writes fail rather than silently ignoring required rules. Unrelated types remain usable. No automatic content migration or overwrite of user definitions occurs.

Requests pin the effective definition identity along with document revision. Recheck them before persistence; a changed definition produces a conflict instead of applying stale rules. Reload updates navigation and namespaced MCP tools without rebuilding the binary. Notify clients of tool-list changes where supported; otherwise report that reconnection is needed. Qualify actual refresh behavior for each harness. Removal/disablement removes specialized tools/actions, while existing documents retain generic inspection and the structural-validation behavior specified in [knowledge format](knowledge-format.md).

## Shared reads and writes

Read/context returns raw content, revision, effective structural definition, Markdown instructions and available operations/prompts. Derived projections use fixed core behavior; displaying content does not execute it. Generic inspection remains available for invalid/disabled types.

All writes, including generic MCP update/replace, raw editing and board moves, follow the same operation path: resolve scope and inputs, acquire the writer lock, refresh document/dependency bytes, compare revision/definition identity, validate document structure, fields and links against the proposed graph, then persist one record atomically. Release the lock and refresh views. Rejection leaves original bytes intact. A dropped connection after persistence is an uncertain outcome; reread before retrying. No multi-record transaction, durable hook queue or exactly-once execution promise is introduced.

## Lifecycle actions through the harness

A lifecycle action is only a label for an ordinary Prompt node, with declared inputs and relevant record context. It is an explicit user/agent-requested action, not an executable save callback or background listener. It can ask the harness to review evidence, perform arbitrary local work, or submit an ordinary document update through MCP. Only the harness runs that work. No automatic recursive agent launch, arbitrary command endpoint or code handler is added to the MCP server. MCP can render the prompt for the already-running agent; CLI/GUI action launch uses the existing adapter.

Prompt failure/cancellation does not imply rollback of already completed edits or Git effects. Show conversation state separately from persisted document state. Mere browsing, action listing or an ordinary save never launches an agent. Tool descriptions and prompts guide behavior; they do not turn agent judgment into deterministic validation.

## Custom UI: kanban example

A user-defined board receives an explicitly scoped list of record snapshots and revisions, renders columns and implements drag/drop in browser JavaScript. Dropping a ticket sends `requestEdit` with the selected record, expected revision and proposed state. The parent checks message identity/scope and uses the same validated write service as native views. A rejected or stale edit restores the observed card position and shows the diagnostic. An action such as “review and complete” instead sends `requestAction` for a declared prompt, visibly launching the harness. Neither route bypasses validation or permits direct local code execution.

Use the existing sandboxed iframe, nonce/schema checks, per-process trust and asset-digest invalidation in [templates](templates.md). Extend the view model with only the explicitly selected board's records; do not expose an unscoped workspace. Ordinary Markdown remains inert. This browser extensibility is in scope; backend user-code execution is not deferred work, it is excluded.

## Enforcement and qualification

MCP controls operations routed through factoryctl. A harness with direct file access can bypass that path; report access restrictions as verified, advisory or unavailable. Do not claim prompts force tool usage. Agent Git operations may replace documents, so external changes must be revalidated. No silent harness-wide permission/configuration changes are implied.

[Task 12](../tasks/12-extension-contract.md) qualifies declarative definitions and shared core operations with fixtures. Downstream transports, harness and custom UI tasks qualify the complete [A14](acceptance.md#a14-runtime-document-extensions) journey. Add a third definition while the binary runs, exercise a custom drag/drop board and a prompt-backed agent action, and demonstrate that no backend extension code can be registered or executed.
