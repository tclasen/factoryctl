---
type: "Task"
title: "Implement native Codex conversation integration"
description: "Implement native Codex conversation integration."
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
  - "13-user-interaction.md"
  - "00-runtime-spike.md"
  - "06-transports.md"
specifications:
  - "../specifications/personal-runtime.md"
  - "../specifications/agent-runtime.md"
requirements:
  - "R30"
  - "R11"
  - "R12"
  - "R13"
  - "R14"
  - "R18"
  - "R24"
---

# Implement native Codex conversation integration

## Outcome

Implement the qualified app-server adapter, profile/action context, ephemeral MCP overrides and native thread lifecycle without duplicating persistence.

## Scope and handoff

Implement the relevant modules in the [architecture](../designs/architecture.md), following the linked specifications in frontmatter. Read [operator authority](../sources/operator-brief.md) before any external effects. Historical planning checkpoint: this task was planned and not started, and the architecture source paths did not yet exist. The current implemented outcome is reconciled below. The [implementation plan](../designs/implementation-plan.md) establishes build order and shared contracts.

## Verification and done criteria

Use deterministic fake-server tests for interleaved events, unsupported approvals, disconnects and ambiguous turns; then a bounded live Mac test for two conversations, follow-up, interruption and a factoryctl MCP edit. No guessed success on runtime errors.

Acceptance coverage: A08 in [release acceptance](../specifications/acceptance.md). Include exact tested revision, commands, expected/observed results and remaining limitations in current task evidence when implemented. Do not equate a checked task field with proof of behavior. Run relevant negative tests as well as success journeys.

## Dependencies, risks and rollback

`depends_on` identifies prerequisites; an agent interprets readiness. Factoryctl validates references but does not schedule this work. Preserve unknown metadata and source evidence. Roll back failed implementation through ordinary Git changes by the agent, retaining useful evidence. Do not expand into deployment, purchases or external messaging.

Implement invocation/user-home preference resolution and adapter dispatch per [personal runtime preferences](../specifications/personal-runtime.md). Verify different contributors can use the same repository bytes, unsupported adapters fail explicitly, and profiles cannot override personal runtime/model/effort choices. Do not implement extra adapters as part of this first Codex slice.

## Collection, interaction and continuing evidence

Integrate the qualified question broker from task 13. Verify actual user answer → originating tool result → same-thread follow-up in CLI/web, concurrent-project isolation and honest unsupported Apps capability. Run S1–S3 conversation fixtures and record A16/A17/A18 integration evidence.

## Completion revision 2 reconciliation

Current product candidate: `a7c22ba847d58e8f4ece6a22ee9bb73302a64c02`. The explicit instruction to implement the remaining project tasks and its bounds are recorded in the [completion task](project-completion.md); the original authorization and historical checkpoints above remain distinct.

The native Codex app-server adapter, invocation-only MCP attachment, contributor-home preferences, model/effort resolution and thread lifecycle are implemented. Current native web/PTY evidence includes concurrent conversation routing, tool answers, follow-up, interruption and a real MCP document edit. Missing adapters/capabilities fail explicitly. Native runtime storage remains external; unsupported Apps/elicitation variants are not reported as successes (A08/A16/A17/A18/A22).

Current command, named-test, native/browser, scene and artifact identities are consolidated in [completion milestone evidence](../sources/completion-milestone-evidence.md), with the independent [code review](../sources/completion-code-review.md) and [acceptance review](../sources/completion-acceptance-review.md). The proposed matrix accounts for 105 passing rows and blocked U045; that proposal is not a full-product acceptance or release claim. Final combined evidence and independent review remain separate gates.

Conditional closure proposal: the bounded work described here is implemented or reconciled and ready for reviewed delivery. The proposed `done` state becomes effective ONLY after this PR is independently reviewed, merged normally and verified on owning-remote main. Before those gates this is a closure proposal, not a delivered result. The final independent reviewer must assess this task's evidence basis. Coordinator retains delivery/reconciliation ownership; no earlier PR/main result is invented by this update. Task 11 and the root completion task retain the unresolved static-linkage/full-qualification boundary.
