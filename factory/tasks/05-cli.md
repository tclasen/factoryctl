---
type: "Task"
title: "Deliver the human and CI command interface"
description: "Deliver the human and CI command interface."
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
owner: "developer"
depends_on:
  - "04-templates.md"
specifications:
  - "../specifications/interfaces.md"
requirements:
  - "R30"
  - "R09"
  - "R18"
  - "R22"
  - "R24"
  - "R25"
---

# Deliver the human and CI command interface

## Outcome

Expose application operations through the specified CLI with human output, JSON output, deterministic selection and explicit noninteractive failures.

## Scope and handoff

Implement the relevant modules in the [architecture](../designs/architecture.md), following the linked specifications in frontmatter. Read [operator authority](../sources/operator-brief.md) before any external effects. Historical CLI checkpoint: the document portion was in progress under the [first CLI milestone](first-cli-milestone.md), with then-uncovered criteria recorded in [CLI evidence](../sources/cli-milestone-evidence.md). The current implemented outcome is reconciled below. The [implementation plan](../designs/implementation-plan.md) establishes build order and shared contracts.

## Verification and done criteria

Exercise help/examples against a temporary Git repo, test exit codes/stdout-stderr separation, check no prompts in noninteractive mode, and inspect the Git working tree for only the requested changes.

Acceptance coverage: A06, A12 in [release acceptance](../specifications/acceptance.md). Include exact tested revision, commands, expected/observed results and remaining limitations in current task evidence when implemented. Do not equate a checked task field with proof of behavior. Run relevant negative tests as well as success journeys.

## Dependencies, risks and rollback

`depends_on` identifies prerequisites; an agent interprets readiness. Factoryctl validates references but does not schedule this work. Preserve unknown metadata and source evidence. Roll back failed implementation through ordinary Git changes by the agent, retaining useful evidence. Do not expand into deployment, purchases or external messaging.

## Collection, interaction and continuing evidence

Run S1–S3 through the same collection commands. Implement the terminal surface for [agent questions](../specifications/user-interaction.md) with free text, choices/custom answers and same-thread follow-up; test no-surface/noninteractive failure without reading MCP protocol stdin. Task 13 qualifies the broker/harness integration. Retain CLI until an explicit qualified retirement decision. Record A16/A17 evidence.

## Completion revision 2 reconciliation

Current product candidate: `a7c22ba847d58e8f4ece6a22ee9bb73302a64c02`. The explicit instruction to implement the remaining project tasks and its bounds are recorded in the [completion task](project-completion.md); the original authorization and historical checkpoints above remain distinct.

The human/JSON/noninteractive CLI and scoped collection/document/extension operations are implemented. Current native PTY evidence exercises free text, choices/custom answers, originating-tool results, same-thread follow-up and interruption; deterministic checks cover refusal and protocol-output boundaries. The CLI is explicitly retained, with no automatic retirement from first-party generated-view support (A06/A12/A16/A17).

Current command, named-test, native/browser, scene and artifact identities are consolidated in [completion milestone evidence](../sources/completion-milestone-evidence.md), with the independent [code review](../sources/completion-code-review.md) and [acceptance review](../sources/completion-acceptance-review.md). The proposed matrix accounts for 105 passing rows and blocked U045; that proposal is not a full-product acceptance or release claim. Final combined evidence and independent review remain separate gates.

Conditional closure proposal: the bounded work described here is implemented or reconciled and ready for reviewed delivery. The proposed `done` state becomes effective ONLY after this PR is independently reviewed, merged normally and verified on owning-remote main. Before those gates this is a closure proposal, not a delivered result. The final independent reviewer must assess this task's evidence basis. Coordinator retains delivery/reconciliation ownership; no earlier PR/main result is invented by this update. Task 11 and the root completion task retain the unresolved static-linkage/full-qualification boundary.
