---
type: "Task"
title: "Qualify the native Codex integration"
description: "Qualify the native Codex integration."
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
priority: "high"
owner: "unassigned"
depends_on: []
specifications:
  - "../specifications/agent-runtime.md"
requirements:
  - "R11"
  - "R12"
  - "R13"
  - "R14"
  - "R24"
---

# Qualify the native Codex integration

## Outcome

Create a bounded integration experiment using the installed app-server schema. Establish thread/turn creation, streaming, question/approval response, interruption and invocation-only MCP configuration before finalizing the adapter.

## Scope and handoff

Implement the relevant modules in the [architecture](../designs/architecture.md), following the linked specifications in frontmatter. Read [operator authority](../sources/operator-brief.md) before any external effects. Installed Codex 0.164.0 schema, native turns, MCP attachment, same-thread follow-up, native input and interruption have bounded qualification in [capability evidence](../sources/interaction-capabilities.md). Native approval UI is covered by deterministic tests; full A08 and platform qualification are not asserted. The [implementation plan](../designs/implementation-plan.md) establishes build order and shared contracts.

## Verification and done criteria

Inspect current protocol; demonstrate a factory-tool stub round trip without writing Codex config files; record tested version and any missing capabilities. If a required capability is absent, revise the adapter contract explicitly before downstream implementation.

Acceptance coverage: A08 in [release acceptance](../specifications/acceptance.md). Include exact tested revision, commands, expected/observed results and remaining limitations in current task evidence when implemented. Do not equate a checked task field with proof of behavior. Run relevant negative tests as well as success journeys.

## Dependencies, risks and rollback

`depends_on` identifies prerequisites; an agent interprets readiness. Factoryctl validates references but does not schedule this work. Preserve unknown metadata and source evidence. Roll back failed implementation through ordinary Git changes by the agent, retaining useful evidence. Do not expand into deployment, purchases or external messaging.

## Completion revision 2 reconciliation

Current product candidate: `a7c22ba847d58e8f4ece6a22ee9bb73302a64c02`. The explicit instruction to implement the remaining project tasks and its bounds are recorded in the [completion task](project-completion.md); the original authorization and historical checkpoints above remain distinct.

The native Codex adapter is implemented and exercised on the current Mac: native conversations, same-thread follow-up, streaming, actual input/tool answers, interruption and a revision/definition-pinned MCP Research Fact edit. The negative native approval route is verified separately from the retained failed whole-probe recorder. Unsupported elicitation/Apps variants remain explicitly unavailable; no foreign-runtime execution is inferred (A08/A17/A18).

Current command, named-test, native/browser, scene and artifact identities are consolidated in [completion milestone evidence](../sources/completion-milestone-evidence.md), with the independent [code review](../sources/completion-code-review.md) and [acceptance review](../sources/completion-acceptance-review.md). The proposed matrix accounts for 105 passing rows and blocked U045; that proposal is not a full-product acceptance or release claim. Final combined evidence and independent review remain separate gates.

Conditional closure proposal: the bounded work described here is implemented or reconciled and ready for reviewed delivery. The proposed `done` state becomes effective ONLY after this PR is independently reviewed, merged normally and verified on owning-remote main. Before those gates this is a closure proposal, not a delivered result. The final independent reviewer must assess this task's evidence basis. Coordinator retains delivery/reconciliation ownership; no earlier PR/main result is invented by this update. Task 11 and the root completion task retain the unresolved static-linkage/full-qualification boundary.
