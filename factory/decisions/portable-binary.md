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

Use pure-Go backend dependencies, build-time browser compilation and embedded defaults/UI. An external browser and selected agent harness remain integrations. Factoryctl never spawns Git; agents/users own all Git operations. Document types are declarative data; custom browser code is served as optional assets. No backend plugin execution engine or native plugin loader is included.

## Alternatives considered — LLM analysis

Rust is also plausible but adds no demonstrated requirement advantage here. A Node/Python runtime would conflict with the distribution goal.

## Consequences — expected, not measured

Qualify dependency linkage per target. Fully static Darwin linkage is not promised; the single-binary interpretation must remain explicit. All non-Mac runtime behavior is untested initially.

See [architecture](../designs/architecture.md) and [release acceptance](../specifications/acceptance.md).

## Clarification and unresolved platform qualification

The operator explicitly requires a single static binary; a compiled executable alone is insufficient. Linux releases must demonstrate absence of dynamic-library dependencies. macOS/Windows normally retain platform system interfaces; the earlier self-contained interpretation is a documented platform limitation, not permission to silently weaken the static requirement. Release qualification must report actual linkage and resolve any unsupported target/requirement conflict with the operator before labeling it compliant. This documentation change makes no platform compliance claim. External type definitions are data; optional UI code runs in the browser. Neither is a companion backend library required for base operation.
