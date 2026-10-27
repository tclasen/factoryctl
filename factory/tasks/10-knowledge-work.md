---
type: "Task"
title: "Qualify dogfooding and non-software workflows"
description: "Qualify dogfooding and non-software workflows."
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
  - "09-custom-views.md"
specifications:
  - "../specifications/knowledge-work.md"
requirements:
  - "R06"
  - "R08"
  - "R19"
  - "R20"
  - "R25"
---

# Qualify dogfooding and non-software workflows

## Outcome

Use the shipped software starter on factoryctl records and exercise synthetic research and contract-pursuit definitions/templates with no core code changes.

## Scope and handoff

Implement the relevant modules in the [architecture](../designs/architecture.md), following the linked specifications in frontmatter. Read [operator authority](../sources/operator-brief.md) before any external effects. This task is planned, not started; source paths named in the architecture do not yet exist. The [implementation plan](../designs/implementation-plan.md) establishes build order and shared contracts.

## Verification and done criteria

Show new type list/detail/board/actions, source-vs-inference labels, cross-project links/backlinks, assumptions, inaccessible sources and agent-edited progress. Keep synthetic bids clearly fictional; perform no real submissions or contacts.

Acceptance coverage: A11 in [release acceptance](../specifications/acceptance.md). Include exact tested revision, commands, expected/observed results and remaining limitations in current task evidence when implemented. Do not equate a checked task field with proof of behavior. Run relevant negative tests as well as success journeys.

## Dependencies, risks and rollback

`depends_on` identifies prerequisites; an agent interprets readiness. Factoryctl validates references but does not schedule this work. Preserve unknown metadata and source evidence. Roll back failed implementation through ordinary Git changes by the agent, retaining useful evidence. Do not expand into deployment, purchases or external messaging.
