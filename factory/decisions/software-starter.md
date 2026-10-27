---
type: "Decision"
title: "Use living specs and linked implementation records"
description: "Use living specs and linked implementation records."
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

# Use living specs and linked implementation records

## Context and basis

[Operator requirements](../sources/operator-brief.md) define the constraints. `accepted` here means the documented scope choice was explicitly made in the interview; it does not mean a human verified every sentence or the implementation. `proposed` marks an LLM-selected technical approach for review.

## Decision

The operator accepted an optional software-development starter set and living specs maintained through Git. Use custom Specification and Design definitions alongside Task and Decision.

## Alternatives considered — LLM analysis

Duplicated task checklists or a mandatory delta-spec workflow add state and process that the operator did not request.

## Consequences — expected, not measured

Factoryctl dogfoods extensibility; software workflow is an optional set of records/prompts. Exact template text remains an LLM design proposal.

See [architecture](../designs/architecture.md) and [release acceptance](../specifications/acceptance.md).
