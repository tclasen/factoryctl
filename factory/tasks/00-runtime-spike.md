---
type: "Task"
title: "Qualify the native Codex integration"
description: "Qualify the native Codex integration."
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
  - "../specifications/agent-runtime.md"
requirements:
  - "R11"
  - "R12"
  - "R13"
  - "R14"
  - "R24"
---

# Qualify the native Codex integration

## Outcome

Create a bounded integration experiment using the installed app-server schema. Establish thread/turn creation, streaming, question/approval response, interruption and invocation-only MCP configuration before finalizing the adapter.

## Scope and handoff

Implement the relevant modules in the [architecture](../designs/architecture.md), following the linked specifications in frontmatter. Read [operator authority](../sources/operator-brief.md) before any external effects. Installed Codex 0.164.0 schema, native turns, MCP attachment, same-thread follow-up, native input and interruption have bounded qualification in [capability evidence](../sources/interaction-capabilities.md). Native approval UI is covered by deterministic tests; full A08 and platform qualification are not asserted. The [implementation plan](../designs/implementation-plan.md) establishes build order and shared contracts.

## Verification and done criteria

Inspect current protocol; demonstrate a factory-tool stub round trip without writing Codex config files; record tested version and any missing capabilities. If a required capability is absent, revise the adapter contract explicitly before downstream implementation.

Acceptance coverage: A08 in [release acceptance](../specifications/acceptance.md). Include exact tested revision, commands, expected/observed results and remaining limitations in current task evidence when implemented. Do not equate a checked task field with proof of behavior. Run relevant negative tests as well as success journeys.

## Dependencies, risks and rollback

`depends_on` identifies prerequisites; an agent interprets readiness. Factoryctl validates references but does not schedule this work. Preserve unknown metadata and source evidence. Roll back failed implementation through ordinary Git changes by the agent, retaining useful evidence. Do not expand into deployment, purchases or external messaging.
