---
type: "Task"
title: "Support optional trusted custom presentation"
description: "Support optional trusted custom presentation."
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
  - "08-gui.md"
specifications:
  - "../specifications/templates.md"
  - "../specifications/gui.md"
requirements:
  - "R02"
  - "R15"
  - "R23"
---

# Support optional trusted custom presentation

## Outcome

Support optional .factory/ assets with explicit trust, digest invalidation, safe generic fallback and a constrained iframe message bridge.

## Scope and handoff

Implement the relevant modules in the [architecture](../designs/architecture.md), following the linked specifications in frontmatter. Read [operator authority](../sources/operator-brief.md) before any external effects. This task is planned, not started; source paths named in the architecture do not yet exist. The [implementation plan](../designs/implementation-plan.md) establishes build order and shared contracts.

## Verification and done criteria

Test arbitrary HTML/CSS/JS rendering, blocked direct parent access, channel/source validation, missing assets, escaped record data, deliberate action/edit requests and a project with no assets directory. Inspect packaged binary asset completeness.

Acceptance coverage: A10 in [release acceptance](../specifications/acceptance.md). Include exact tested revision, commands, expected/observed results and remaining limitations in current task evidence when implemented. Do not equate a checked task field with proof of behavior. Run relevant negative tests as well as success journeys.

## Dependencies, risks and rollback

`depends_on` identifies prerequisites; an agent interprets readiness. Factoryctl validates references but does not schedule this work. Preserve unknown metadata and source evidence. Roll back failed implementation through ordinary Git changes by the agent, retaining useful evidence. Do not expand into deployment, purchases or external messaging.
