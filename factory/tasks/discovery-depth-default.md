---
type: Task
title: Default project discovery to immediate subfolders
description: Record and qualify the operator-requested discovery depth boundary.
status: draft
factory:
  format: 1
  assertion:
    frame: mixed
    certainty: supported
    basis: authoritative-statement
task_state: done
priority: normal
owner: coordinator
depends_on: []
specifications:
  - ../specifications/project-discovery.md
  - ../specifications/acceptance.md
requirements:
  - R01
---

# Default project discovery to immediate subfolders

## Authority and scope

The operator requested a default recursion limit of 1, inspecting only folders immediately inside the current folder. The [operator source](../sources/operator-brief.md) records the request. Scope is this task, that source, R01 in the [product specification](../specifications/product.md), discovery specification, A01 acceptance and [foundation task](01-project-foundation.md). The original scope was documentation only and made no application implementation claim; current executable evidence is reconciled below.

## Plan and acceptance

Plan revision 1, authorized by the operator request: align the six named documents, check consistency and whitespace, obtain fresh independent review, and deliver through a normal PR to origin/main. The scan root is depth 0; immediate subfolders are depth 1; grandchildren are excluded. Existing current-working-tree and explicit-path selection remain. The same default applies to an explicit scan root. No new depth option is designed here. This is a bounded requirement correction, not a new architecture or costly-to-reverse interface; separate architectural design review is unnecessary.

## Run and ownership

Coordinator owns the six-document mutation claim on branch `docs/discovery-depth-one`; Reviewer has read-only access to the candidate. Baseline: `46dd295`. No blocking dependencies. Budget is this bounded documentation change and one independent review; normal diagnosed retries only. Requested review configuration: inherited model, fresh context, no implementation authorship; actual model telemetry unknown. Rollback is a scoped revert through PR. Organization files and unrelated untracked files are excluded.

## Verification and delivery

Document inspection and `git diff --check` validate the specified boundary. At the original documentation checkpoint runtime discovery tests remained planned in the foundation task; documentation checks do not establish executable behavior. Independent review and remote-main delivery are pending. Coordinator retains ownership until verified delivery; if blocked, resume review or PR delivery from the named branch. This record is design-change tracking; the foundation implementation was then planned, with current evidence reconciled below.

## Completion revision 2 reconciliation

Current product candidate: `a7c22ba847d58e8f4ece6a22ee9bb73302a64c02`. The explicit instruction to implement the remaining project tasks and its bounds are recorded in the [completion task](project-completion.md); the original authorization and historical checkpoints above remain distinct.

The depth-one correction is documented and implemented: immediate children are considered, grandchildren/deeper repositories excluded by default, and explicit selection remains supported. Current discovery and S1–S3 checks qualify executable behavior separately from the original documentation-only review. No runtime recursion option or Git subprocess is inferred (A01/A16).

Current command, named-test, native/browser, scene and artifact identities are consolidated in [completion milestone evidence](../sources/completion-milestone-evidence.md), with the independent [code review](../sources/completion-code-review.md) and [acceptance review](../sources/completion-acceptance-review.md). The proposed matrix accounts for 105 passing rows and blocked U045; that proposal is not a full-product acceptance or release claim. Final combined evidence and independent review remain separate gates.

Conditional closure proposal: the bounded work described here is implemented or reconciled and ready for reviewed delivery. The proposed `done` state becomes effective ONLY after this PR is independently reviewed, merged normally and verified on owning-remote main. Before those gates this is a closure proposal, not a delivered result. The final independent reviewer must assess this task's evidence basis. Coordinator retains delivery/reconciliation ownership; no earlier PR/main result is invented by this update. Task 11 and the root completion task retain the unresolved static-linkage/full-qualification boundary.
