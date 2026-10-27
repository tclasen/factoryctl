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

Use bidirectional app-server JSON-RPC for GUI interaction; expose factoryctl record tools as a separate MCP server to Codex.

## Alternatives considered — LLM analysis

Terminal-output scraping lacks structured approval/input semantics. Rebuilding an agent runtime would expand scope and duplicate Codex persistence.

## Consequences — expected, not measured

A runtime spike must qualify actual installed protocol/settings. Sessions and credentials remain in Codex; factoryctl keeps transient routing only.

See [architecture](../designs/architecture.md) and [release acceptance](../specifications/acceptance.md).
