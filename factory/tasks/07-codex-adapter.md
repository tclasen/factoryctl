---
type: "Task"
title: "Implement native Codex conversation integration"
description: "Implement native Codex conversation integration."
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
  - "00-runtime-spike.md"
  - "06-transports.md"
specifications:
  - "../specifications/agent-runtime.md"
requirements:
  - "R11"
  - "R12"
  - "R13"
  - "R14"
  - "R18"
  - "R24"
---

# Implement native Codex conversation integration

## Outcome

Implement the qualified app-server adapter, profile/action context, ephemeral MCP overrides and native thread lifecycle without duplicating persistence.

## Scope and handoff

Implement the relevant modules in the [architecture](../designs/architecture.md), following the linked specifications in frontmatter. Read [operator authority](../sources/operator-brief.md) before any external effects. This task is planned, not started; source paths named in the architecture do not yet exist. The [implementation plan](../designs/implementation-plan.md) establishes build order and shared contracts.

## Verification and done criteria

Use deterministic fake-server tests for interleaved events, unsupported approvals, disconnects and ambiguous turns; then a bounded live Mac test for two conversations, follow-up, interruption and a factoryctl MCP edit. No guessed success on runtime errors.

Acceptance coverage: A08 in [release acceptance](../specifications/acceptance.md). Include exact tested revision, commands, expected/observed results and remaining limitations in current task evidence when implemented. Do not equate a checked task field with proof of behavior. Run relevant negative tests as well as success journeys.

## Dependencies, risks and rollback

`depends_on` identifies prerequisites; an agent interprets readiness. Factoryctl validates references but does not schedule this work. Preserve unknown metadata and source evidence. Roll back failed implementation through ordinary Git changes by the agent, retaining useful evidence. Do not expand into deployment, purchases or external messaging.
