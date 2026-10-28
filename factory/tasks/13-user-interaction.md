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
task_state: done
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
intent_authorization: The operator requested completion of the next milestone, authorizing this bounded build and qualification; exact external Apps protocol and future sequence remain proposals requiring qualification.
---

# Qualify user questions and generative web views

## Outcome and authority

Implement a bounded prototype of [user interaction](../specifications/user-interaction.md) after transports and the runtime spike. The operator requires question/answer tools and asks for an MCP Apps evaluation, potentially replacing the interactive CLI after qualification. The operator now requested completion of this next milestone; the bounded implementation and live qualification are ready for review. Full-product matrix qualification remains incomplete, and delivery requires the reviewed PR/main observation.

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


## Historical interaction-milestone checkpoint

[Interaction milestone evidence](../sources/interaction-milestone-evidence.md) records the frozen product source, current tests, actual browser form/chart and same-thread dialogue, real terminal dialogue, concurrent project isolation, multiple selections, refusal, restart loss and sandbox boundaries. [Capability evidence](../sources/interaction-capabilities.md) recommends retaining the CLI. Native browser approval/input have deterministic fixture coverage; external Apps interoperability is not qualified. General GUI editing/boards, type-derived extension registration and delivery-status implementation remain outside this slice.

The local stable demo is `/Users/Shared/factoryctl/artifacts/demo/latest.mp4`; complete current row accounting accompanies it. Full-matrix qualification remains explicitly incomplete. At that milestone, `in_review` preserved this distinction rather than calling all product requirements done. The implementation delivery is tracked by [PR #14](https://github.com/tclasen/factoryctl/pull/14); a proposed delivery record becomes effective only after merge and remote-main verification.

Final qualification allocations: runtime implementer owned only local live PTY/isolation probe artifacts after releasing runtime source; recorder implementer additionally owned `scripts/review-video/qualification.mjs` and its local capture. Coordinator owned browser/terminal recorders, combination and current evidence. Reviewer remained read-only and independent. Functional capture failures and capacity failures were retained; a subsequent visual-review revision added chart/form/draft dwell without changing product source. This was a readability revision, not a hidden replacement of a failing functional outcome. No artifacts were staged or uploaded.

## Completion revision 2 reconciliation

Current product candidate: `a7c22ba847d58e8f4ece6a22ee9bb73302a64c02`. The explicit instruction to implement the remaining project tasks and its bounds are recorded in the [completion task](project-completion.md); the original authorization and historical checkpoints above remain distinct.

Current native web/PTY and S1–S3 evidence exercises free text, choices/custom text, actual originating-tool answers, concurrent routing, same-thread follow-up and interruption. Distinct generated form/chart views, sandbox/refusal/fallback tests and personal settings have current coverage. The native approval cancellation criterion has bounded verified evidence while its whole-probe recorder failure is retained. External MCP Apps interoperability and unsupported elicitation variants remain unqualified; the CLI is retained (A08/A16/A17/A18).

Current command, named-test, native/browser, scene and artifact identities are consolidated in [completion milestone evidence](../sources/completion-milestone-evidence.md), with the independent [code review](../sources/completion-code-review.md) and [acceptance review](../sources/completion-acceptance-review.md). The proposed matrix accounts for 105 passing rows and blocked U045; that proposal is not a full-product acceptance or release claim. Final combined evidence and independent review remain separate gates.

Conditional closure proposal: the bounded work described here is implemented or reconciled and ready for reviewed delivery. The proposed `done` state becomes effective ONLY after this PR is independently reviewed, merged normally and verified on owning-remote main. Before those gates this is a closure proposal, not a delivered result. The final independent reviewer must assess this task's evidence basis. Coordinator retains delivery/reconciliation ownership; no earlier PR/main result is invented by this update. Task 11 and the root completion task retain the unresolved static-linkage/full-qualification boundary.
