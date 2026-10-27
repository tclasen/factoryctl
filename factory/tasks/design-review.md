---
type: Task
title: Independently review the implementation design
description: Review architectural fit, comprehension and execution readiness.
status: draft
factory:
  format: 1
  assertion:
    frame: observed
    certainty: supported
    basis: direct-observation
task_state: in_review
priority: high
owner: independent-reviewer
depends_on: []
specifications:
  - ../specifications/product.md
requirements: []
sources:
  - id: operator
    resource: ../sources/operator-brief.md
    title: Operator requirements and design authority
---

# Independently review the implementation design

## Scope and run manifest

This is read-only review of the candidate pinned by Git commit in the review dispatch. The reviewer is a fresh execution without drafting history or implementation authorship. Requested configuration: gpt-6.1-sol with high reasoning, selected for a bounded technical-design review. Actual runtime metadata is unknown unless independently exposed. No paid providers, deployments, source implementation or unrelated repository changes are authorized.

Required sources: the [operator brief](../sources/operator-brief.md), [product requirements](../specifications/product.md), [architecture](../designs/architecture.md), [implementation plan](../designs/implementation-plan.md), and every specification they define as a version-1 contract. Inspect representative actual definitions, prompts, templates and tasks to establish that the bundle follows its own format. The [validation evidence](../sources/design-validation.md) is author-reported evidence to independently assess, not a review conclusion.

## Criteria and handoff

Assess architectural fit and simplicity, comprehension by a fresh implementer, and execution readiness. Prioritize contradictory contracts, lost user requirements, invented factual certainty, incomplete validation, persistence violations, unsafe integration assumptions and unsupported completeness claims. Return file/line findings with severity, which required sources were read, check results and limitations. Recheck fixes against a new pinned candidate when they affect findings. No arbitrary minimum finding count.

The author owns all candidate edits; the reviewer has no mutation claim. Report findings through the native handoff. The coordinator records current review outcome in [design validation](../sources/design-validation.md). Review completion is separate from application implementation and from remote delivery.
