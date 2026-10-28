---
type: Task
title: Reconcile project-specific delivery-status design
description: Document the accepted delivery direction and distinguish implementation proposals.
status: draft
task_state: done
owner: coordinator
priority: high
requirements: [R36, R37]
factory:
  format: 1
  assertion:
    frame: desired
    certainty: supported
    basis: authoritative-statement
sources:
  - id: operator
    resource: ../sources/operator-brief.md#delivery-status-and-project-specific-criteria
    title: Documentation request and delivery-status correction
    author: human:operator
---

# Delivery-status documentation reconciliation

Request: “Please update whatever documents you need to for the discussion we just had.” Interpretation: reconcile the source record, product/GUI/acceptance requirements, proposed delivery design, current plan and traceability. Authorization covers documentation and normal reviewed PR delivery, not implementation of a resolver, network service or new storage mechanism.

Coordinator owns AGENTS.md and the scoped factory/ documentation in this PR. Preserve unrelated files and all existing local media. This is a documentation update, not a new product milestone; do not create a requirements-only film or claim the prior demo qualifies newly added requirements. Exact schema/storage/adapter design remains proposed.

Plan: capture agreed direction and local-only correction; preserve the minimal-interface boundary; correct the stale plan observation; independently review source fit and consistency; validate Markdown/link structure; deliver through this repository's protected PR/main process. No organization-policy changes.

Acceptance: web visibility required and CLI optional; no inherited PR/remote policy; explicit undefined/stale/unavailable criteria/evidence; candidate/criteria revision binding; post-integration observations avoid recursive ticket commits; local branch and working-tree mechanisms remain identified proposals. At the original design-only checkpoint the status implementation checks were planned A21, not executed passes; current implementation evidence is reconciled below.

Run delivery-status-doc-review: fresh independent Reviewer, inherited model/reasoning (requested; actual telemetry unknown), read-only access to this task, project AGENTS, source discussion, minimal-interface decision and the exact candidate diff. Claim: review only, no implementation or Git writes. Coordinator alone owns edits and delivery. Budget: bounded documentation review, one correction cycle if needed; outcome/source fit precedes consistency review. Report exact approved head/tree, findings and limitations. Closure is determined from post-merge delivery evidence, not this ticket asserting its own merge.

## Completion revision 2 reconciliation

Current product candidate: `a7c22ba847d58e8f4ece6a22ee9bb73302a64c02`. The explicit instruction to implement the remaining project tasks and its bounds are recorded in the [completion task](project-completion.md); the original authorization and historical checkpoints above remain distinct.

The project-owned criteria/assessment and delivery-observation design is reconciled with its implemented local-working-tree resolver and browser presentation. Current evidence covers shared task/plan resolution, undefined criteria, provenance and changed-since-verification states. Acceptance remains an attributed claim; branch/remote content remains externally recorded and unverified by factoryctl. No inherited hosting/PR rule, Git subprocess, backend approval policy or mandatory tracking database is introduced (A21).

Current command, named-test, native/browser, scene and artifact identities are consolidated in [completion milestone evidence](../sources/completion-milestone-evidence.md), with the independent [code review](../sources/completion-code-review.md) and [acceptance review](../sources/completion-acceptance-review.md). The proposed matrix accounts for 105 passing rows and blocked U045; that proposal is not a full-product acceptance or release claim. Final combined evidence and independent review remain separate gates.

Conditional closure proposal: the bounded work described here is implemented or reconciled and ready for reviewed delivery. The proposed `done` state becomes effective ONLY after this PR is independently reviewed, merged normally and verified on owning-remote main. Before those gates this is a closure proposal, not a delivered result. The final independent reviewer must assess this task's evidence basis. Coordinator retains delivery/reconciliation ownership; no earlier PR/main result is invented by this update. Task 11 and the root completion task retain the unresolved static-linkage/full-qualification boundary.
