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
task_state: "planned"
priority: "high"
owner: "unassigned"
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

Implement the relevant modules in the [architecture](../designs/architecture.md), following the linked specifications in frontmatter. Read [operator authority](../sources/operator-brief.md) before any external effects. This task is planned, not started; source paths named in the architecture do not yet exist. The [implementation plan](../designs/implementation-plan.md) establishes build order and shared contracts.

## Verification and done criteria

Test unknown keys/types, verified mapping/list, invalid UTF-8, YAML hazards, limits, path containment, two cooperating processes, stale revisions, crashes and externally edited bytes. Record non-cooperating writer limitations.

Acceptance coverage: A02, A06 in [release acceptance](../specifications/acceptance.md). Include exact tested revision, commands, expected/observed results and remaining limitations in current task evidence when implemented. Do not equate a checked task field with proof of behavior. Run relevant negative tests as well as success journeys.

## Dependencies, risks and rollback

`depends_on` identifies prerequisites; an agent interprets readiness. Factoryctl validates references but does not schedule this work. Preserve unknown metadata and source evidence. Roll back failed implementation through ordinary Git changes by the agent, retaining useful evidence. Do not expand into deployment, purchases or external messaging.

## Collection, interaction and continuing evidence

Run S1–S3 identity/link/write fixtures: distinct bundles/definitions, same-named records, cross-links/backlinks, target-safe writes and per-bundle locking. Verify cross-project races are reported honestly, even in one monorepo; no cross-project transaction is claimed. Record A16 evidence.
