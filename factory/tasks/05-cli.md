---
type: "Task"
title: "Deliver the human and CI command interface"
description: "Deliver the human and CI command interface."
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
  - "04-templates.md"
specifications:
  - "../specifications/interfaces.md"
requirements:
  - "R09"
  - "R18"
  - "R22"
  - "R24"
  - "R25"
---

# Deliver the human and CI command interface

## Outcome

Expose application operations through the specified CLI with human output, JSON output, deterministic selection and explicit noninteractive failures.

## Scope and handoff

Implement the relevant modules in the [architecture](../designs/architecture.md), following the linked specifications in frontmatter. Read [operator authority](../sources/operator-brief.md) before any external effects. This task is planned, not started; source paths named in the architecture do not yet exist. The [implementation plan](../designs/implementation-plan.md) establishes build order and shared contracts.

## Verification and done criteria

Exercise help/examples against a temporary Git repo, test exit codes/stdout-stderr separation, check no prompts in noninteractive mode, and inspect the Git working tree for only the requested changes.

Acceptance coverage: A06, A12 in [release acceptance](../specifications/acceptance.md). Include exact tested revision, commands, expected/observed results and remaining limitations in current task evidence when implemented. Do not equate a checked task field with proof of behavior. Run relevant negative tests as well as success journeys.

## Dependencies, risks and rollback

`depends_on` identifies prerequisites; an agent interprets readiness. Factoryctl validates references but does not schedule this work. Preserve unknown metadata and source evidence. Roll back failed implementation through ordinary Git changes by the agent, retaining useful evidence. Do not expand into deployment, purchases or external messaging.
