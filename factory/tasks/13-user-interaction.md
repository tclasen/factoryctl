---
type: Task
title: "Qualify user questions and generative web views"
description: "Qualify user questions and generative web views."
status: draft
factory:
  format: 1
  assertion:
    frame: desired
    certainty: mixed
    basis: mixed
task_state: in_progress
priority: high
owner: coordinator
depends_on: [00-runtime-spike.md, 06-transports.md]
specifications: [../specifications/user-interaction.md]
requirements: [R30, R31]
sources:
  - id: operator
    resource: ../sources/operator-brief.md
    title: Operator requirements and scope
    author: human:operator
intent_state: proposed
intent_owner: assistant
intent_source: ../sources/operator-brief.md#project-collections-and-user-interaction
intent_authorization: MCP Apps exploration is requested; this exact protocol and implementation sequence remain proposals requiring qualification.
---

# Qualify user questions and generative web views

## Outcome and authority

Implement a bounded prototype of [user interaction](../specifications/user-interaction.md) after transports and the runtime spike. The operator requires question/answer tools and asks for an MCP Apps evaluation, potentially replacing the interactive CLI after qualification. The operator now requested completion of this next milestone; implementation and qualification are in progress, not yet delivered.

## Plan and verification

Qualify broker attachment and answer provenance with fake runtimes, then a bounded live harness exercise. Demonstrate free text, choices with custom answers and follow-up conversation through both current interfaces. Verify unanswered requests never become consent, response routing across concurrent projects/threads, noninteractive failure, restart/disconnect and repeat submissions. Then qualify generated browser views, Apps handshake/resources/bridge or first-party host rendering, sandbox limits and fallback. Exercise all [project scenarios](../specifications/project-scenarios.md).

Produce A17/A18 evidence, an explicit capability table and an interactive-CLI retain/retire recommendation. A failed or unsupported Apps experiment remains evidence, not successful interoperability. No CLI is removed without the explicit decision described in the contract. Tasks 07–09 integrate the qualified prototype; task 11 requires current release regression evidence and video coverage.


## Current admission and implementation brief

Latest request: “Please complete the next milestone.” Its preceding scope is agent questions and rich web views, runtime/transport qualification, actual responses and same-thread follow-up, lifecycle/isolation checks, capability report and local demo. Delivery-status implementation remains separate. Scope includes existing Go/embedded browser sources, new interaction/runtime/MCP packages, tests, recorder tooling, project documentation and local artifacts only. Existing PR/review/main delivery applies; no global harness configuration, paid providers, artifact uploads or automatic CLI retirement.

Approach: reuse the document explorer; add an opt-in conversation surface and an in-memory scoped question/view broker. A selected installed Codex app-server runs through argv/stdio, with invocation-only MCP configuration. A conversation owns its runtime and MCP credential to avoid cross-thread answer routing. Runtime child IPC permits asking/waiting/presenting, never answering. Trusted web host and interactive terminal submit answers through separate host-only paths. Generated code runs in a trusted-once sandboxed browser frame with no network/filesystem access; a plain trusted answer surface remains available. Actual SDK/host support must be observed before claiming MCP Apps compatibility. No backend workflow policy engine.

Readiness trigger: new runtime/interaction interface requires fresh independent design review before implementation. Required sources: this brief, operator-brief, user-interaction, agent-runtime, personal-runtime, minimal-interface decision, A08/A17/A18 and project AGENTS. Blocking open capability questions: installed runtime invocation/MCP attachment and external Apps support; qualify first. Exact API/types/UI layout are delegated engineering choices; unsupported capabilities stay explicit rather than fabricating success. Rollback via scoped revert PR, preserving local evidence.

Claims/runs: Coordinator owns internal/web/, cmd/factoryctl/, project docs, integration, local demos and Git; interaction implementer owns internal/interaction/ and internal/mcp/; runtime implementer owns internal/runtime/ and a scoped local runtime qualification directory; fresh Reviewer owns read-only review, no implementation. Each run uses inherited model/reasoning unless separately recorded; actual model/cost telemetry unknown. Shared contracts must be agreed before cross-package integration. No concurrent edits to shared files. Commands120seconds, tests10minutes, checkpoint every15active minutes; at most2diagnosed retries per operation. Preserve user Finder files and prior videos.

Acceptance evidence separates deterministic fake-runtime tests, live runtime observations and actual browser interaction. Test free text, choices/custom text, follow-up, concurrent scope, duplicate/spoofed/stale answers, timeout, cancellation, disconnect/reconnect and restart. Demonstrate distinct generated form and chart/table views, sandbox failures and plain fallback; record capability/CLI disposition honestly. Refresh current feature regression, complete102-row accounting and one local scene-annotated demo; full-product gaps remain explicit. Independent final approval binds exact source and media before protected PR/main verification.

Qualification follow-up allocation: runtime implementer now owns README.md and a new `factory/sources/interaction-capabilities.md` capability Reference only, using this task's existing run and inherited model. Reconcile observed runtime/PTY evidence, exact first-party view limits, invocation examples and CLI retain recommendation; no edits to implementation, current matrix or other documents. Coordinator retains those and the final evidence package. Recorder implementer owns only `scripts/review-video/record.mjs` and `annotate.mjs`; final capture waits for source freeze. Reviewer remains independent and read-only.
