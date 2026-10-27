---
type: "Task"
title: "Prepare the factoryctl implementation design"
description: "Prepare the factoryctl implementation design."
status: "draft"
factory:
  format: 1
  assertion:
    frame: "observed"
    certainty: "supported"
    basis: "direct-observation"
sources:
  -
    id: "operator"
    resource: "../sources/operator-brief.md"
    title: "Operator requirements and scope"
    author: "human:operator"
task_state: "in_review"
priority: "high"
owner: "design-author"
depends_on: []
specifications:
  - "../specifications/product.md"
---

# Prepare the factoryctl implementation design

## Outcome and authority

Produce an internally consistent, human-navigable OKF design package covering the confirmed [operator requirements](../sources/operator-brief.md). This task is documentation and planning only. It does not implement factoryctl or claim its runtime acceptance tests pass.

## Current ownership and plan

The authoring agent owns `factory/**` in this repository for this task. No organization files, other projects, deployed services or legacy repository migrations are in scope. Fresh independent reviewers receive read-only access to the pinned candidate. No overlapping mutation owner was present when the empty checkout was inspected.

1. Capture the authoritative interview and current environment observations.
2. Define the OKF extension, discovery, personal preference, template, GUI, CLI, MCP and runtime contracts.
3. Materialize definitions, profiles, templates, prompts and linked implementation tasks using that contract.
4. Check YAML/OKF structure, local links, definitions, templates, required fields, task dependencies and requirement coverage.
5. Obtain an independent design review; resolve blockers and recheck affected content.
6. Deliver the documentation through the owning repository's Git workflow, preserving its independent history.

## Acceptance

All confirmed requirements map to implementation tasks and observable acceptance checks. A fresh implementer can identify interfaces, failure behavior, dependencies, limits and proposed decisions. Every bundle Markdown file has a type and explicit knowledge classification. No index, history duplication, application database or Codex-specific configuration file is introduced. `.factory/` remains absent as a real example of the default-assets design.

## Review and delivery

The [verification record](../sources/design-validation.md) records current checks and limitations. The [roadmap](../designs/implementation-plan.md) is the implementation handoff. Task state describes the current work only; Git retains its revisions.

## Rollback and constraints

Revert the design change through Git if needed; do not delete unrelated files or rewrite history. No runtime budget or spending is authorized by this task. No application code or cloud resources are required for design validation.

## Current review scope

The initial design and release policy were delivered through PRs 1 and 2. The author is revising the design to remove project configuration and separate personal runtime preferences, including the operator's deferred external-tracking-backend idea. Independent review and normal PR delivery are required for this revision. No application implementation is in scope.
