---
type: "Task"
title: "Qualify dogfooding and non-software workflows"
description: "Qualify dogfooding and non-software workflows."
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
  - "09-custom-views.md"
specifications:
  - "../specifications/knowledge-work.md"
requirements:
  - "R06"
  - "R08"
  - "R19"
  - "R20"
  - "R25"
---

# Qualify dogfooding and non-software workflows

## Outcome

Use the shipped software starter on factoryctl records and exercise synthetic research and contract-pursuit definitions/templates with no core code changes.

## Scope and handoff

Implement the relevant modules in the [architecture](../designs/architecture.md), following the linked specifications in frontmatter. Read [operator authority](../sources/operator-brief.md) before any external effects. Historical planning checkpoint: this task was planned and not started, and the architecture source paths did not yet exist. The current implemented outcome is reconciled below. The [implementation plan](../designs/implementation-plan.md) establishes build order and shared contracts.

## Verification and done criteria

Show new type list/detail/board/actions, source-vs-inference labels, cross-project links/backlinks, assumptions, inaccessible sources and agent-edited progress. Keep synthetic bids clearly fictional; perform no real submissions or contacts.

Acceptance coverage: A11 in [release acceptance](../specifications/acceptance.md). Include exact tested revision, commands, expected/observed results and remaining limitations in current task evidence when implemented. Do not equate a checked task field with proof of behavior. Run relevant negative tests as well as success journeys.

## Dependencies, risks and rollback

`depends_on` identifies prerequisites; an agent interprets readiness. Factoryctl validates references but does not schedule this work. Preserve unknown metadata and source evidence. Roll back failed implementation through ordinary Git changes by the agent, retaining useful evidence. Do not expand into deployment, purchases or external messaging.

## Completion revision 2 reconciliation

Current product candidate: `a7c22ba847d58e8f4ece6a22ee9bb73302a64c02`. The explicit instruction to implement the remaining project tasks and its bounds are recorded in the [completion task](project-completion.md); the original authorization and historical checkpoints above remain distinct.

Synthetic Research Fact and Opportunity definitions are exercised without domain-specific core changes: list/detail/boards, citations, assumptions/inference, backlinks, prompt previews, structural pipeline edits and a real native MCP Research Fact read/update/read. No real bid, contact, submission or award is asserted. Current dogfood validation finds no Markdown structural errors, but retains the preserved factory/.DS_Store invalid-location diagnostic and explicitly unchecked links; it is not a clean all-files/network-validation claim (A11/A13).

Current command, named-test, native/browser, scene and artifact identities are consolidated in [completion milestone evidence](../sources/completion-milestone-evidence.md), with the independent [code review](../sources/completion-code-review.md) and [acceptance review](../sources/completion-acceptance-review.md). The proposed matrix accounts for 105 passing rows and blocked U045; that proposal is not a full-product acceptance or release claim. Final combined evidence and independent review remain separate gates.

Conditional closure proposal: the bounded work described here is implemented or reconciled and ready for reviewed delivery. The proposed `done` state becomes effective ONLY after this PR is independently reviewed, merged normally and verified on owning-remote main. Before those gates this is a closure proposal, not a delivered result. The final independent reviewer must assess this task's evidence basis. Coordinator retains delivery/reconciliation ownership; no earlier PR/main result is invented by this update. Task 11 and the root completion task retain the unresolved static-linkage/full-qualification boundary.
