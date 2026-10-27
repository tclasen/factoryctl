---
type: "Task"
title: "Qualify performance, portability and complete release behavior"
description: "Qualify performance, portability and complete release behavior."
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
  - "10-knowledge-work.md"
specifications:
  - "../specifications/acceptance.md"
requirements:
  - "R01"
  - "R03"
  - "R04"
  - "R16"
  - "R17"
  - "R21"
  - "R22"
  - "R24"
  - "R25"
---

# Qualify performance, portability and complete release behavior

## Outcome

Run the full Mac acceptance suite and benchmark, cross-compile the six target OS/architecture combinations where supported, inspect linkage and prepare user/operator documentation.

## Scope and handoff

Implement the relevant modules in the [architecture](../designs/architecture.md), following the linked specifications in frontmatter. Read [operator authority](../sources/operator-brief.md) before any external effects. This task is planned, not started; source paths named in the architecture do not yet exist. The [implementation plan](../designs/implementation-plan.md) establishes build order and shared contracts.

## Verification and done criteria

Record exact revision, runtime/dependency versions, corpus/hardware, errors and achieved metrics. Label Linux/Windows and non-host architecture execution untested. Verify files-only persistence, no app Git mutations, graceful shutdown and external-proxy deployment instructions without deploying a service.

Acceptance coverage: A12, A13 in [release acceptance](../specifications/acceptance.md). Include exact tested revision, commands, expected/observed results and remaining limitations in current task evidence when implemented. Do not equate a checked task field with proof of behavior. Run relevant negative tests as well as success journeys.

## Dependencies, risks and rollback

`depends_on` identifies prerequisites; an agent interprets readiness. Factoryctl validates references but does not schedule this work. Preserve unknown metadata and source evidence. Roll back failed implementation through ordinary Git changes by the agent, retaining useful evidence. Do not expand into deployment, purchases or external messaging.

## Commit and version qualification

Implement and verify the [commit and release policy](../decisions/release-policy.md): validate new commit and PR titles, calculate versions from Conventional Commits, start at 0.1.0, and prevent an unauthorized transition to or release of 1.0.0. Exercise the policy examples and negative authorization case. The 0.x line carries no compatibility promise. Release checks are planned, not currently installed.
