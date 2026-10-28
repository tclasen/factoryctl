---
type: "Task"
title: "Implement OKF reading and safe file storage"
description: "Implement OKF reading and safe file storage."
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
  - "01-project-foundation.md"
specifications:
  - "../specifications/knowledge-format.md"
  - "../specifications/epistemics.md"
requirements:
  - "R03"
  - "R05"
  - "R07"
  - "R08"
  - "R22"
---

# Implement OKF reading and safe file storage

## Outcome

Parse OKF Markdown into loss-preserving documents, build in-memory snapshots, and implement revision-checked create/update/delete with native writer coordination and atomic replacement.

## Scope and handoff

Implement the relevant modules in the [architecture](../designs/architecture.md), following the linked specifications in frontmatter. Read [operator authority](../sources/operator-brief.md) before any external effects. Historical CLI checkpoint: the document portion was in progress under the [first CLI milestone](first-cli-milestone.md), with then-uncovered criteria recorded in [CLI evidence](../sources/cli-milestone-evidence.md). The current implemented outcome is reconciled below. The [implementation plan](../designs/implementation-plan.md) establishes build order and shared contracts.

## Verification and done criteria

Test unknown keys/types, verified mapping/list, invalid UTF-8, YAML hazards, limits, path containment, two cooperating processes, stale revisions, crashes and externally edited bytes. Record non-cooperating writer limitations.

Acceptance coverage: A02, A06 in [release acceptance](../specifications/acceptance.md). Include exact tested revision, commands, expected/observed results and remaining limitations in current task evidence when implemented. Do not equate a checked task field with proof of behavior. Run relevant negative tests as well as success journeys.

## Dependencies, risks and rollback

`depends_on` identifies prerequisites; an agent interprets readiness. Factoryctl validates references but does not schedule this work. Preserve unknown metadata and source evidence. Roll back failed implementation through ordinary Git changes by the agent, retaining useful evidence. Do not expand into deployment, purchases or external messaging.

## Collection, interaction and continuing evidence

Run S1–S3 identity/link/write fixtures: distinct bundles/definitions, same-named records, cross-links/backlinks, target-safe writes and per-bundle locking. Verify cross-project races are reported honestly, even in one monorepo; no cross-project transaction is claimed. Record A16 evidence.

## Completion revision 2 reconciliation

Current product candidate: `a7c22ba847d58e8f4ece6a22ee9bb73302a64c02`. The explicit instruction to implement the remaining project tasks and its bounds are recorded in the [completion task](project-completion.md); the original authorization and historical checkpoints above remain distinct.

The shared OKF reader and safe writer are implemented, including unknown metadata preservation, revision/definition checks, native cooperative locking, proposed-graph validation and atomic single-record replacement. Current tests and S1–S3 browser journeys exercise containment, stale edits, target-safe writes and distinct same-path records. Uncoordinated external editors and cross-project writes do not form a transaction (A02/A06/A16).

Current command, named-test, native/browser, scene and artifact identities are consolidated in [completion milestone evidence](../sources/completion-milestone-evidence.md), with the independent [code review](../sources/completion-code-review.md) and [acceptance review](../sources/completion-acceptance-review.md). The proposed matrix accounts for 105 passing rows and blocked U045; that proposal is not a full-product acceptance or release claim. Final combined evidence and independent review remain separate gates.

Conditional closure proposal: the bounded work described here is implemented or reconciled and ready for reviewed delivery. The proposed `done` state becomes effective ONLY after this PR is independently reviewed, merged normally and verified on owning-remote main. Before those gates this is a closure proposal, not a delivered result. The final independent reviewer must assess this task's evidence basis. Coordinator retains delivery/reconciliation ownership; no earlier PR/main result is invented by this update. Task 11 and the root completion task retain the unresolved static-linkage/full-qualification boundary.
