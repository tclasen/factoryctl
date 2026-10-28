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
task_state: done
priority: "normal"
owner: "unassigned"
depends_on:
  - "12-extension-contract.md"
  - "05-cli.md"
specifications:
  - "../specifications/interfaces.md"
requirements:
  - "R30"
  - "R13"
  - "R17"
  - "R22"
  - "R23"
  - "R24"
  - "R25"
---

# Expose shared MCP and local HTTP services

## Outcome

Implement MCP stdio tools and local HTTP adapters around the same application services. The current concrete transport contract uses polling where specified; the original proposed SSE mechanism is superseded. Bind project scope, diagnostics, revision tokens and in-memory event routing.

## Scope and handoff

Implement the relevant modules in the [architecture](../designs/architecture.md), following the linked specifications in frontmatter. Read [operator authority](../sources/operator-brief.md) before any external effects. Shared MCP document tools and scoped question/view attachment are implemented, with protocol, revision-conflict and capability-separation tests. At that earlier checkpoint browser conversation routes were implemented while general document mutation, extension dispatch and transport parity remained outstanding; the current outcome is reconciled below. See [capability evidence](../sources/interaction-capabilities.md). The [implementation plan](../designs/implementation-plan.md) establishes build order and shared contracts.

## Verification and done criteria

Run MCP protocol tests and HTTP integration tests for CRUD, pagination expiry, schema failures, revision conflicts, arbitrary-root rejection, Host/Origin/nonce checks and no stdout contamination.

Acceptance coverage: A07 in [release acceptance](../specifications/acceptance.md). Include exact tested revision, commands, expected/observed results and remaining limitations in current task evidence when implemented. Do not equate a checked task field with proof of behavior. Run relevant negative tests as well as success journeys.

## Dependencies, risks and rollback

`depends_on` identifies prerequisites; an agent interprets readiness. Factoryctl validates references but does not schedule this work. Preserve unknown metadata and source evidence. Roll back failed implementation through ordinary Git changes by the agent, retaining useful evidence. Do not expand into deployment, purchases or external messaging.

## Runtime extension integration

Register namespaced tools and shared CLI/HTTP operation dispatch from the [extension contract](../specifications/extensions.md). Expose structural definitions and agent-interpreted Markdown instructions on reads; preserve shared structural validation, registry conflict handling and document persistence outcomes distinct from harness action status. Qualify tool refresh/reconnect with task 07; include R26–R28 and A14 integration evidence.

## Collection, interaction and continuing evidence

Expose the question tools, pending-result/wait semantics, scoped host answer routes and events from [user interaction](../specifications/user-interaction.md). Exercise S1–S3 transport parity and routing, including separate monorepo project identities. Use fake-broker tests here; task 13 qualifies runtime/Apps integration. Record A16/A17 evidence.

## Completion revision 2 reconciliation

Current product candidate: `a7c22ba847d58e8f4ece6a22ee9bb73302a64c02`. The explicit instruction to implement the remaining project tasks and its bounds are recorded in the [completion task](project-completion.md); the original authorization and historical checkpoints above remain distinct.

CLI, MCP and local HTTP share document/extension services with current scope, pagination, validation, revision and Host/Origin/nonce checks. General browser mutation, namespaced dispatch and request-time definition reload are implemented. Current transport contracts use polling where specified; earlier proposed SSE behavior is not a delivered capability. Native question/generated-view routing remains separate from external Apps interoperability (A07/A14/A16/A17).

Current command, named-test, native/browser, scene and artifact identities are consolidated in [completion milestone evidence](../sources/completion-milestone-evidence.md), with the independent [code review](../sources/completion-code-review.md) and [acceptance review](../sources/completion-acceptance-review.md). The proposed matrix accounts for 105 passing rows and blocked U045; that proposal is not a full-product acceptance or release claim. Final combined evidence and independent review remain separate gates.

Conditional closure proposal: the bounded work described here is implemented or reconciled and ready for reviewed delivery. The proposed `done` state becomes effective ONLY after this PR is independently reviewed, merged normally and verified on owning-remote main. Before those gates this is a closure proposal, not a delivered result. The final independent reviewer must assess this task's evidence basis. Coordinator retains delivery/reconciliation ownership; no earlier PR/main result is invented by this update. Task 11 and the root completion task retain the unresolved static-linkage/full-qualification boundary.
