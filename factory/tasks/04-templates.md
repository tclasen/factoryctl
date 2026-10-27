---
type: "Task"
title: "Implement deterministic record and prompt templates"
description: "Implement deterministic record and prompt templates."
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
priority: "normal"
owner: "unassigned"
depends_on:
  - "03-definition-validation.md"
specifications:
  - "../specifications/templates.md"
requirements:
  - "R06"
  - "R11"
  - "R14"
  - "R15"
  - "R19"
  - "R25"
---

# Implement deterministic record and prompt templates

## Outcome

Implement restricted templating, typed input collection, include resolution, full-output validation and materialization of the software starter set.

## Scope and handoff

Implement the relevant modules in the [architecture](../designs/architecture.md), following the linked specifications in frontmatter. Read [operator authority](../sources/operator-brief.md) before any external effects. This task is planned, not started; source paths named in the architecture do not yet exist. The [implementation plan](../designs/implementation-plan.md) establishes build order and shared contracts.

## Verification and done criteria

Test escaping, missing/extra inputs, nested includes, cycles, output limits, invalid rendered YAML and unknown defaults. Initialize a scratch repository from the actual default templates and validate its entire bundle without .factory/.

Acceptance coverage: A05 in [release acceptance](../specifications/acceptance.md). Include exact tested revision, commands, expected/observed results and remaining limitations in current task evidence when implemented. Do not equate a checked task field with proof of behavior. Run relevant negative tests as well as success journeys.

## Dependencies, risks and rollback

`depends_on` identifies prerequisites; an agent interprets readiness. Factoryctl validates references but does not schedule this work. Preserve unknown metadata and source evidence. Roll back failed implementation through ordinary Git changes by the agent, retaining useful evidence. Do not expand into deployment, purchases or external messaging.

## Collection, interaction and continuing evidence

Run S1–S3 with different per-project templates/assets and absent .factory/ directories. Resolve assets from the selected project root; prohibit accidental sibling/parent inheritance. Record A16 evidence.
