---
type: Specification
title: "Agent questions and rich web interaction"
description: "Agent questions and rich web interaction."
status: draft
factory:
  format: 1
  assertion:
    frame: desired
    certainty: mixed
    basis: mixed
sources:
  - id: operator
    resource: ../sources/operator-brief.md
    title: Operator requirements and scope
    author: human:operator
---

# Agent questions and rich web interaction

The operator requires an agent-callable question tool, free text, multiple choice and continuing dialogue in both current user interfaces. They favor a rich web interface and propose retiring the interactive CLI if MCP Apps and generative UI prove sufficient. The bounded first-party implementation now uses the question/wait/present names below. External Apps resource advertising and interoperability remain proposals, not implemented capabilities; see [capability evidence](../sources/interaction-capabilities.md).

## Shared question contract

`factory_user_ask` accepts a prompt, response kind (`text`, `single_choice`, `multiple_choice`), optional choices with stable IDs/labels, and optional previous request ID for a follow-up. Choice prompts also allow a custom text answer. Scope and conversation identity come from the launching connection, never from an agent-supplied destination. Validate duplicate/empty choice IDs and incompatible schemas. Return a process-local request ID and `pending` immediately, plus the displayed question and a text fallback; no answer is implied by delivery. This permits rendering before the tool returns, avoids blocking rich-view initialization, and lets the harness choose its next action.

`factory_user_wait` accepts an owned request ID and performs a bounded wait, returning `pending`, `answered`, `declined`, `cancelled`, `expired` or `unavailable`. An answered result contains the user's submitted text and/or selected option IDs verbatim. A wait timeout returns pending, never a default answer; expiry means the question can no longer be answered. The client may wait again or continue unrelated work. Initial proposed wait bound is 30 seconds, qualified against the harness's tool timeout. No polling loop or automatic agent turn is scheduled by factoryctl. Duplicate wait calls are read-only and return the retained result for the current process.

Only an actual user submission through the authenticated host response channel can create an answered result. No model-callable answer setter exists. Displaying, selecting an initial default, dismissing a view, timing out or receiving agent prose does not submit an answer or approval. Native harness approvals remain separate and use their native response channel. Schema validation of an answer is not workflow authorization.

## Routing, lifecycle and conversation

The factoryctl process hosting the CLI/web conversation owns an in-memory broker. A launched MCP child attaches through an invocation-scoped local IPC endpoint and an unguessable credential delivered out of model context, not committed or written into project/runtime configuration. Bind it to the project instance, native thread/turn and tool request. Qualify concrete IPC permissions and platform cleanup in the prototype. Standalone MCP may use negotiated form elicitation from its host; without a connected user surface it returns `interaction_unavailable` promptly.

The MCP server must not read its stdin as human input: stdin/stdout carry protocol traffic. Interactive CLI replies use its owning terminal surface; web replies use a nonce-protected host route. Neither route shares a response across projects, threads, requests or processes. A reply is accepted once; duplicates return the existing disposition and conflicting/stale submissions fail. New questions and user follow-ups remain in the same native conversation. Serialize turn submission when the runtime is busy, or expose supported steering behavior; never start a competing turn to deliver a tool answer. The waiting tool receives the answer through its own result. Ordinary subsequent chat uses the conversation message route.

Noninteractive mode with no attached answer surface fails explicitly. EOF, cancellation and shutdown terminate outstanding requests with a non-answer disposition. Browser disconnect leaves a request pending while the host process and conversation remain alive; reconnect retrieves its current state. Process restart loses the broker, and old IDs cannot be reused. Runtime persistence remains runtime-owned; no factoryctl transcript, answer database or durable request queue is created. A requested durable decision must be explicitly written as a normal validated document by the agent/user.

## Rich UI proposal and MCP Apps evidence

Official [MCP Apps overview](https://apps.extensions.modelcontextprotocol.io/api/documents/overview.html) describes host-rendered, sandboxed views with a host bridge and progressive text fallback. The [Apps specification](https://github.com/modelcontextprotocol/ext-apps/blob/main/specification/draft/apps.mdx) defines UI resource metadata; the [MCP elicitation specification](https://modelcontextprotocol.io/specification/2025-11-25/client/elicitation) defines client-mediated form requests and explicit response dispositions. These sources were inspected on 2026-10-27; the draft Apps contract is moving and must be pinned during implementation. None proves that the installed Codex app-server exposes Apps capabilities or forwards the necessary events.

Proposed `factory_ui_present` lets the agent supply a title, plain-text equivalent, structured data and runtime-generated HTML/CSS/JavaScript, optionally attached to an owned question. It creates a bounded, ephemeral browser resource with a unique `ui://` identifier. For an Apps-capable host, register a generic presentation tool with `_meta.ui.resourceUri` referring to an embedded view shell served as `text/html;profile=mcp-app`; that shell receives the generated content and scoped data. The fixed shell is an isolation boundary, not a fixed catalog of user-facing widgets. Capabilities, resource loading and message bridging require empirical qualification.

The factoryctl web UI may implement the host bridge itself; it must not depend on Codex displaying an iframe. The prototype must demonstrate how the broker relays tool input/results to that host even when the harness is only an MCP client. External Apps hosts are a separate interoperability case. Hosts without qualified rich support get the plain question/text surface; CLI never evaluates browser code.

Generated views can create forms, tables, charts and other interactions during conversation without predefined application components. They execute only in a sandboxed browser frame after explicit per-process trust for generated content. Default CSP denies external network and navigation; project files, credentials, shell and direct harness APIs are unavailable. Size limits, source-window/channel checks and cleanup apply. The trusted host validates every message and routes allowed reads/edits through existing project-scoped services. View submissions are bound to the owned request through host-only capabilities unavailable in model tool arguments/results; hiding a tool from the model is not sufficient authorization. Qualify that property before enabling rich answers in an external host; otherwise use the trusted plain answer surface. Ordinary Markdown remains inert. No dynamic backend handler, policy engine or automatic package installation is introduced.

## Prototype and interactive-CLI decision

[Task 13](../tasks/13-user-interaction.md) must demonstrate agent call → visible question/generated view → real user interaction → correlated tool result → follow-up in the same conversation. Cover a generated form and a materially different generated chart/table interaction, not only a fixed question widget. Test two concurrent projects/conversations, spoofed/duplicate submissions, refusal, cancellation, reconnect, lost process state, unsupported host and malformed/generated-code failures. Record runtime, MCP and Apps/SDK versions, actual negotiated capabilities and source revisions.

Web-only user-interface adoption requires equivalent project navigation, editing, conversations, questions, accessibility and error recovery, with no essential journey requiring a terminal interaction. Evaluate web-only operation when Apps support is absent too; a resilient plain web fallback can coexist with richer views. Record a decision with supported combinations, gaps and retained launch/MCP/automation commands. Until qualification and explicit approval of that decision, keep current CLI requirements. The release may report Apps unavailable with evidence and retain the CLI; it may not claim the rich feature works. Any approved removal updates requirements, acceptance, traceability and release video coverage together.
