---
type: "Decision"
title: "Separate configuration, assets and knowledge"
description: "Separate configuration, assets and knowledge."
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
decision_state: "accepted"
---

# Separate configuration, assets and knowledge

## Context and basis

[Operator requirements](../sources/operator-brief.md) define the constraints. `accepted` here means the documented scope choice was explicitly made in the interview; it does not mean a human verified every sentence or the implementation. `proposed` marks an LLM-selected technical approach for review.

## Decision

The operator explicitly selected root factory.yaml, optional .factory/ assets and factory/ as the only OKF graph. All default operation must work without .factory/. Git supplies history and the agent handles Git mutations.

## Alternatives considered — LLM analysis

A single hidden metadata directory would mix assets and knowledge; a Markdown-only config would contradict the explicitly selected YAML file.

## Consequences — expected, not measured

The graph stays readable and portable; configuration/assets are versioned but are not knowledge nodes. No generated index or transition log is required.

See [architecture](../designs/architecture.md) and [release acceptance](../specifications/acceptance.md).
