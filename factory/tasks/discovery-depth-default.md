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
task_state: in_review
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

The operator requested a default recursion limit of 1, inspecting only folders immediately inside the current folder. The [operator source](../sources/operator-brief.md) records the request. Scope is this task, that source, R01 in the [product specification](../specifications/product.md), discovery specification, A01 acceptance and [foundation task](01-project-foundation.md). No application implementation exists or is claimed.

## Plan and acceptance

Plan revision 1, authorized by the operator request: align the six named documents, check consistency and whitespace, obtain fresh independent review, and deliver through a normal PR to origin/main. The scan root is depth 0; immediate subfolders are depth 1; grandchildren are excluded. Existing current-working-tree and explicit-path selection remain. The same default applies to an explicit scan root. No new depth option is designed here. This is a bounded requirement correction, not a new architecture or costly-to-reverse interface; separate architectural design review is unnecessary.

## Run and ownership

Coordinator owns the six-document mutation claim on branch `docs/discovery-depth-one`; Reviewer has read-only access to the candidate. Baseline: `46dd295`. No blocking dependencies. Budget is this bounded documentation change and one independent review; normal diagnosed retries only. Requested review configuration: inherited model, fresh context, no implementation authorship; actual model telemetry unknown. Rollback is a scoped revert through PR. Organization files and unrelated untracked files are excluded.

## Verification and delivery

Document inspection and `git diff --check` validate the specified boundary. Runtime discovery tests remain planned in the foundation task; documentation checks do not establish executable behavior. Independent review and remote-main delivery are pending. Coordinator retains ownership until verified delivery; if blocked, resume review or PR delivery from the named branch. This record is design-change tracking; the foundation implementation task remains planned.
