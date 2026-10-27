---
type: "Decision"
title: "Use a self-contained Go executable with embedded assets"
description: "Use a self-contained Go executable with embedded assets."
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

# Use a self-contained Go executable with embedded assets

## Context and basis

[Operator requirements](../sources/operator-brief.md) define the constraints. `accepted` here means the documented scope choice was explicitly made in the interview; it does not mean a human verified every sentence or the implementation. `proposed` marks an LLM-selected technical approach for review.

## Decision

Use pure-Go backend dependencies, build-time browser compilation and embedded defaults/UI. Git, a browser and Codex remain explicit external programs.

## Alternatives considered — LLM analysis

Rust is also plausible but adds no demonstrated requirement advantage here. A Node/Python runtime would conflict with the distribution goal.

## Consequences — expected, not measured

Qualify dependency linkage per target. Fully static Darwin linkage is not promised; the single-binary interpretation must remain explicit. All non-Mac runtime behavior is untested initially.

See [architecture](../designs/architecture.md) and [release acceptance](../specifications/acceptance.md).
