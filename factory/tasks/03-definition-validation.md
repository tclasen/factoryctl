---
type: "Task"
title: "Implement definitions and live structural validation"
description: "Implement definitions and live structural validation."
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
owner: "developer"
depends_on:
  - "02-knowledge-store.md"
specifications:
  - "../specifications/knowledge-format.md"
  - "../specifications/validation.md"
  - "../specifications/epistemics.md"
requirements:
  - "R06"
  - "R07"
  - "R08"
  - "R24"
  - "R25"
---

# Implement definitions and live structural validation

## Outcome

Load definitions/meta-schema, build references/backlinks, validate fields and placement, and share structured diagnostics across save and watcher flows.

## Scope and handoff

Implement the relevant modules in the [architecture](../designs/architecture.md), following the linked specifications in frontmatter. Read [operator authority](../sources/operator-brief.md) before any external effects. Historical CLI checkpoint: the document portion was in progress under the [first CLI milestone](first-cli-milestone.md), with then-uncovered criteria recorded in [CLI evidence](../sources/cli-milestone-evidence.md). The current implemented outcome is reconciled below. The [implementation plan](../designs/implementation-plan.md) establishes build order and shared contracts.

## Verification and done criteria

Use positive/negative corpus tests for required/null/enum/type cases, wrong locations, duplicate definitions, missing targets/anchors, source references, case collisions, unknown custom types, Factory metadata and definition enablement changes. Validate the proposed graph and incoming dependents: reject heading removal and target-type changes that break references, validate creation/replacement self-links, and allow unrelated-error repairs. Reject invalid mutations without changing bytes; prove invalid external edits remain visible.

Acceptance coverage: A03, A04 in [release acceptance](../specifications/acceptance.md). Include exact tested revision, commands, expected/observed results and remaining limitations in current task evidence when implemented. Do not equate a checked task field with proof of behavior. Run relevant negative tests as well as success journeys.

## Dependencies, risks and rollback

`depends_on` identifies prerequisites; an agent interprets readiness. Factoryctl validates references but does not schedule this work. Preserve unknown metadata and source evidence. Roll back failed implementation through ordinary Git changes by the agent, retaining useful evidence. Do not expand into deployment, purchases or external messaging.

## Collection, interaction and continuing evidence

Run structural validation and definition-isolation checks across S1–S3. Distinct monorepo projects must not share overrides or enablement. Record A16 evidence as these features land.

## Completion revision 2 reconciliation

Current product candidate: `a7c22ba847d58e8f4ece6a22ee9bb73302a64c02`. The explicit instruction to implement the remaining project tasks and its bounds are recorded in the [completion task](project-completion.md); the original authorization and historical checkpoints above remain distinct.

Runtime structural definitions, generic access, enum boards, backlinks and current-byte refresh/revalidation are implemented. Current tests/browser journeys cover disabled/invalid definitions, required/type/enum/placement/link failures, document/definition conflicts, draft-preserving reconciliation and invalid-definition isolation. External/network citations remain explicitly unchecked; structure does not establish truth or approval (A03/A04/A16).

Current command, named-test, native/browser, scene and artifact identities are consolidated in [completion milestone evidence](../sources/completion-milestone-evidence.md), with the independent [code review](../sources/completion-code-review.md) and [acceptance review](../sources/completion-acceptance-review.md). The proposed matrix accounts for 105 passing rows and blocked U045; that proposal is not a full-product acceptance or release claim. Final combined evidence and independent review remain separate gates.

Conditional closure proposal: the bounded work described here is implemented or reconciled and ready for reviewed delivery. The proposed `done` state becomes effective ONLY after this PR is independently reviewed, merged normally and verified on owning-remote main. Before those gates this is a closure proposal, not a delivered result. The final independent reviewer must assess this task's evidence basis. Coordinator retains delivery/reconciliation ownership; no earlier PR/main result is invented by this update. Task 11 and the root completion task retain the unresolved static-linkage/full-qualification boundary.
