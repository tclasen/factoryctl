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
task_state: "planned"
priority: "high"
owner: "unassigned"
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

Implement the relevant modules in the [architecture](../designs/architecture.md), following the linked specifications in frontmatter. Read [operator authority](../sources/operator-brief.md) before any external effects. This task is planned, not started; source paths named in the architecture do not yet exist. The [implementation plan](../designs/implementation-plan.md) establishes build order and shared contracts.

## Verification and done criteria

Use positive/negative corpus tests for required/null/enum/type cases, wrong locations, duplicate definitions, missing targets/anchors, source references, case collisions, unknown custom types, Factory metadata and definition enablement changes. Validate the proposed graph and incoming dependents: reject heading removal and target-type changes that break references, validate creation/replacement self-links, and allow unrelated-error repairs. Reject invalid mutations without changing bytes; prove invalid external edits remain visible.

Acceptance coverage: A03, A04 in [release acceptance](../specifications/acceptance.md). Include exact tested revision, commands, expected/observed results and remaining limitations in current task evidence when implemented. Do not equate a checked task field with proof of behavior. Run relevant negative tests as well as success journeys.

## Dependencies, risks and rollback

`depends_on` identifies prerequisites; an agent interprets readiness. Factoryctl validates references but does not schedule this work. Preserve unknown metadata and source evidence. Roll back failed implementation through ordinary Git changes by the agent, retaining useful evidence. Do not expand into deployment, purchases or external messaging.

## Collection, interaction and continuing evidence

Run structural validation and definition-isolation checks across S1–S3. Distinct monorepo projects must not share overrides or enablement. Record A16 evidence as these features land.
