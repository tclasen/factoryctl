---
type: "Task"
title: "Build the executable foundation and project discovery"
description: "Build the executable foundation and project discovery."
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
depends_on: []
specifications:
  - "../specifications/configuration.md"
requirements:
  - "R01"
  - "R02"
  - "R04"
  - "R16"
  - "R17"
---

# Build the executable foundation and project discovery

## Outcome

Create the Go executable skeleton, embedded-default mechanism, root config parser and read-only Git discovery. Implement explicit project selection and zero/one/many-repository results.

## Scope and handoff

Implement the relevant modules in the [architecture](../designs/architecture.md), following the linked specifications in frontmatter. Read [operator authority](../sources/operator-brief.md) before any external effects. This task is planned, not started; source paths named in the architecture do not yet exist. The [implementation plan](../designs/implementation-plan.md) establishes build order and shared contracts.

## Verification and done criteria

Test CWD inside a repo, a parent directory, worktree .git files, inaccessible paths, symlink loops, duplicate project IDs, malformed config and no .factory/ folder. Prove discovery performs no Git mutations.

Acceptance coverage: A01, A12 in [release acceptance](../specifications/acceptance.md). Include exact tested revision, commands, expected/observed results and remaining limitations in current task evidence when implemented. Do not equate a checked task field with proof of behavior. Run relevant negative tests as well as success journeys.

## Dependencies, risks and rollback

`depends_on` identifies prerequisites; an agent interprets readiness. Factoryctl validates references but does not schedule this work. Preserve unknown metadata and source evidence. Roll back failed implementation through ordinary Git changes by the agent, retaining useful evidence. Do not expand into deployment, purchases or external messaging.
