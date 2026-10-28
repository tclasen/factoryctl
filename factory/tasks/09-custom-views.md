---
type: "Task"
title: "Support optional trusted custom presentation"
description: "Support optional trusted custom presentation."
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
task_state: done
priority: "normal"
owner: "unassigned"
depends_on:
  - "08-gui.md"
specifications:
  - "../specifications/templates.md"
  - "../specifications/gui.md"
requirements:
  - "R31"
  - "R02"
  - "R15"
  - "R23"
---

# Support optional trusted custom presentation

## Outcome

Support optional .factory/ assets with explicit trust, digest invalidation, safe generic fallback and a constrained iframe message bridge.

## Scope and handoff

Implement the relevant modules in the [architecture](../designs/architecture.md), following the linked specifications in frontmatter. Read [operator authority](../sources/operator-brief.md) before any external effects. Historical planning checkpoint: this task was planned and not started, and the architecture source paths did not yet exist. The current implemented outcome is reconciled below. The [implementation plan](../designs/implementation-plan.md) establishes build order and shared contracts.

## Verification and done criteria

Test arbitrary HTML/CSS/JS rendering, blocked direct parent access, channel/source validation, missing assets, escaped record data, deliberate action/edit requests and a project with no assets directory. Inspect packaged binary asset completeness.

Acceptance coverage: A10 in [release acceptance](../specifications/acceptance.md). Include exact tested revision, commands, expected/observed results and remaining limitations in current task evidence when implemented. Do not equate a checked task field with proof of behavior. Run relevant negative tests as well as success journeys.

## Dependencies, risks and rollback

`depends_on` identifies prerequisites; an agent interprets readiness. Factoryctl validates references but does not schedule this work. Preserve unknown metadata and source evidence. Roll back failed implementation through ordinary Git changes by the agent, retaining useful evidence. Do not expand into deployment, purchases or external messaging.

## Runtime-defined board acceptance

Own the browser portion of A14: implement a user-authored HTML/CSS/JavaScript kanban board over an explicitly scoped record selection, drag/drop via `requestEdit`, rejected/stale move feedback and prompt-backed `requestAction`. Preserve iframe/trust/message constraints. Backend custom-code execution is excluded; local computation belongs to the selected harness.

## Collection, interaction and continuing evidence

Integrate the qualified MCP Apps/generative-view prototype from task 13, subject to observed support. Exercise at least two materially different runtime-generated views, no predefined widget catalog requirement, browser isolation/trust, real answer routing and plain fallback across S1–S3. Record A16/A18; unsupported combinations remain explicit.

## Completion revision 2 reconciliation

Current product candidate: `a7c22ba847d58e8f4ece6a22ee9bb73302a64c02`. The explicit instruction to implement the remaining project tasks and its bounds are recorded in the [completion task](project-completion.md); the original authorization and historical checkpoints above remain distinct.

Trusted repository HTML/CSS/JavaScript presentation is implemented with digest-bound in-memory grants, sandboxed frames, exact source/channel/schema checks, scoped edit/action messages, containment and generic fallback. Current browser evidence exercises trust, custom edits/actions, invalidation and negative isolation cases. Separate first-party generated form/chart views and plain fallback are qualified; external MCP Apps interoperability remains unqualified, and no backend custom-code engine is introduced (A10/A14/A18).

Current command, named-test, native/browser, scene and artifact identities are consolidated in [completion milestone evidence](../sources/completion-milestone-evidence.md), with the independent [code review](../sources/completion-code-review.md) and [acceptance review](../sources/completion-acceptance-review.md). The proposed matrix accounts for 105 passing rows and blocked U045; that proposal is not a full-product acceptance or release claim. Final combined evidence and independent review remain separate gates.

Conditional closure proposal: the bounded work described here is implemented or reconciled and ready for reviewed delivery. The proposed `done` state becomes effective ONLY after this PR is independently reviewed, merged normally and verified on owning-remote main. Before those gates this is a closure proposal, not a delivered result. The final independent reviewer must assess this task's evidence basis. Coordinator retains delivery/reconciliation ownership; no earlier PR/main result is invented by this update. Task 11 and the root completion task retain the unresolved static-linkage/full-qualification boundary.
