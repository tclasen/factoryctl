---
type: "Decision"
title: "Use Codex app-server for conversation and MCP for records"
description: "Use Codex app-server for conversation and MCP for records."
status: "draft"
factory:
  format: 1
  assertion:
    frame: "desired"
    certainty: "mixed"
    basis: "mixed"
sources:
  -
    id: "operator"
    resource: "../sources/operator-brief.md"
    title: "Operator requirements and scope"
    author: "human:operator"
decision_state: "proposed"
---

# Use Codex app-server for conversation and MCP for records

## Context and basis

[Operator requirements](../sources/operator-brief.md) define the constraints. `accepted` here means the documented scope choice was explicitly made in the interview; it does not mean a human verified every sentence or the implementation. `proposed` marks an LLM-selected technical approach for review.

## Decision

Use a contributor-selected adapter boundary with shared runtime-neutral profiles. Implement Codex first using bidirectional app-server JSON-RPC for GUI interaction and a separate factoryctl MCP server. Selection, executable, model and reasoning preferences remain outside project Git; see [personal runtime preferences](../specifications/personal-runtime.md).

## Alternatives considered — LLM analysis

Terminal-output scraping lacks structured approval/input semantics. Rebuilding an agent runtime would expand scope and duplicate Codex persistence.

## Consequences — expected, not measured

A runtime spike must qualify actual installed protocol/settings. Sessions and credentials remain in Codex; factoryctl keeps transient routing only.

See [architecture](../designs/architecture.md) and [release acceptance](../specifications/acceptance.md).

## Domain interpretation boundary

The [minimal-interface decision](minimal-interface.md) governs every adapter: the harness interprets Markdown roles, workflow, approvals and lifecycle instructions. Factoryctl supplies readable context and standardized structural document tools, not a policy engine. Native harness permission/approval requests remain transport features and must not be confused with factoryctl evaluating project approval policy.
