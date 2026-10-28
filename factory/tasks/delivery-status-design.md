---
type: Task
title: Reconcile project-specific delivery-status design
description: Document the accepted delivery direction and distinguish implementation proposals.
status: draft
task_state: in_review
owner: coordinator
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

Acceptance: web visibility required and CLI optional; no inherited PR/remote policy; explicit undefined/stale/unavailable criteria/evidence; candidate/criteria revision binding; post-integration observations avoid recursive ticket commits; local branch and working-tree mechanisms remain identified proposals. Required status implementation checks are planned A21, not executed passes.

Run delivery-status-doc-review: fresh independent Reviewer, inherited model/reasoning (requested; actual telemetry unknown), read-only access to this task, project AGENTS, source discussion, minimal-interface decision and the exact candidate diff. Claim: review only, no implementation or Git writes. Coordinator alone owns edits and delivery. Budget: bounded documentation review, one correction cycle if needed; outcome/source fit precedes consistency review. Report exact approved head/tree, findings and limitations. Closure is determined from post-merge delivery evidence, not this ticket asserting its own merge.
