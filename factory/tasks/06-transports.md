---
type: "Task"
title: "Expose shared MCP and local HTTP services"
description: "Expose shared MCP and local HTTP services."
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
  - "12-extension-contract.md"
  - "05-cli.md"
specifications:
  - "../specifications/interfaces.md"
requirements:
  - "R13"
  - "R17"
  - "R22"
  - "R23"
  - "R24"
  - "R25"
---

# Expose shared MCP and local HTTP services

## Outcome

Implement MCP stdio tools and HTTP/SSE adapters around the same application services. Bind project scope, diagnostics, revision tokens and in-memory event routing.

## Scope and handoff

Implement the relevant modules in the [architecture](../designs/architecture.md), following the linked specifications in frontmatter. Read [operator authority](../sources/operator-brief.md) before any external effects. This task is planned, not started; source paths named in the architecture do not yet exist. The [implementation plan](../designs/implementation-plan.md) establishes build order and shared contracts.

## Verification and done criteria

Run MCP protocol tests and HTTP integration tests for CRUD, pagination expiry, schema failures, revision conflicts, arbitrary-root rejection, Host/Origin/nonce checks and no stdout contamination.

Acceptance coverage: A07 in [release acceptance](../specifications/acceptance.md). Include exact tested revision, commands, expected/observed results and remaining limitations in current task evidence when implemented. Do not equate a checked task field with proof of behavior. Run relevant negative tests as well as success journeys.

## Dependencies, risks and rollback

`depends_on` identifies prerequisites; an agent interprets readiness. Factoryctl validates references but does not schedule this work. Preserve unknown metadata and source evidence. Roll back failed implementation through ordinary Git changes by the agent, retaining useful evidence. Do not expand into deployment, purchases or external messaging.

## Runtime extension integration

Register namespaced tools and shared CLI/HTTP operation dispatch from the [extension contract](../specifications/extensions.md). Expose structural definitions and agent-interpreted Markdown instructions on reads; preserve shared structural validation, registry conflict handling and document persistence outcomes distinct from harness action status. Qualify tool refresh/reconnect with task 07; include R26–R28 and A14 integration evidence.
