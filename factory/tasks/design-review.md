---
type: Task
title: Independently review the implementation design
description: Review architectural fit, comprehension and execution readiness.
status: draft
factory:
  format: 1
  assertion:
    frame: observed
    certainty: supported
    basis: direct-observation
task_state: done
priority: high
owner: independent-reviewer
depends_on: []
specifications:
  - ../specifications/product.md
requirements: []
sources:
  - id: operator
    resource: ../sources/operator-brief.md
    title: Operator requirements and design authority
---

# Independently review the implementation design

## Scope and run manifest

This is read-only review of the candidate pinned by Git commit in the review dispatch. The reviewer is a fresh execution without drafting history or implementation authorship. Requested configuration: gpt-6.1-sol with high reasoning, selected for a bounded technical-design review. Actual runtime metadata is unknown unless independently exposed. No paid providers, deployments, source implementation or unrelated repository changes are authorized.

Required sources: the [operator brief](../sources/operator-brief.md), [product requirements](../specifications/product.md), [architecture](../designs/architecture.md), [implementation plan](../designs/implementation-plan.md), and every specification they define as a version-1 contract. Inspect representative actual definitions, prompts, templates and tasks to establish that the bundle follows its own format. The [validation evidence](../sources/design-validation.md) is author-reported evidence to independently assess, not a review conclusion.

## Criteria and handoff

Assess architectural fit and simplicity, comprehension by a fresh implementer, and execution readiness. Prioritize contradictory contracts, lost user requirements, invented factual certainty, incomplete validation, persistence violations, unsafe integration assumptions and unsupported completeness claims. Return file/line findings with severity, which required sources were read, check results and limitations. Recheck fixes against a new pinned candidate when they affect findings. No arbitrary minimum finding count.

The author owns all candidate edits; the reviewer has no mutation claim. Report findings through the native handoff. The coordinator records current review outcome in [design validation](../sources/design-validation.md). Review completion is separate from application implementation and from remote delivery.

## Completion revision 2 reconciliation

Current product candidate: `a7c22ba847d58e8f4ece6a22ee9bb73302a64c02`. The explicit instruction to implement the remaining project tasks and its bounds are recorded in the [completion task](project-completion.md); the original authorization and historical checkpoints above remain distinct.

Independent design readiness review and later code/acceptance review records are present. This original bounded design-review task can close as review work after its reviewed delivery; it does not author or approve the pending final combined-media/full-matrix decision. Current independent conclusions and remaining limits are linked from the completion evidence, rather than inferred from a task label.

Current command, named-test, native/browser, scene and artifact identities are consolidated in [completion milestone evidence](../sources/completion-milestone-evidence.md), with the independent [code review](../sources/completion-code-review.md) and [acceptance review](../sources/completion-acceptance-review.md). The proposed matrix accounts for 105 passing rows and blocked U045; that proposal is not a full-product acceptance or release claim. Final combined evidence and independent review remain separate gates.

Conditional closure proposal: the bounded work described here is implemented or reconciled and ready for reviewed delivery. The proposed `done` state becomes effective ONLY after this PR is independently reviewed, merged normally and verified on owning-remote main. Before those gates this is a closure proposal, not a delivered result. The final independent reviewer must assess this task's evidence basis. Coordinator retains delivery/reconciliation ownership; no earlier PR/main result is invented by this update. Task 11 and the root completion task retain the unresolved static-linkage/full-qualification boundary.
