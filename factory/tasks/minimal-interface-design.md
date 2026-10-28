---
type: "Task"
title: "Document the minimal interface boundary"
description: "Document the minimal interface boundary."
status: "draft"
factory:
  format: 1
  assertion:
    frame: "desired"
    certainty: "supported"
    basis: "authoritative-statement"
sources:
  - id: operator
    resource: "../sources/operator-brief.md"
    title: "Operator requirements and scope"
    author: "human:operator"
task_state: "done"
priority: "normal"
owner: "coordinator"
depends_on: []
specifications:
  - "../specifications/product.md"
  - "../specifications/extensions.md"
requirements:
  - "R29"
---

# Document the minimal interface boundary

## Request, scope and authority

The operator requested that the minimal factoryctl boundary be documented and the repository brought up to date. This documentation-only task owns the decision, affected product/specification/design/task text and this record. No source implementation, harness configuration, organization governance or unrelated files are in scope. The operator authorizes ordinary PR/merge delivery.

## Plan, ownership and review run

Revision 1: establish the controlling minimal-interface decision; remove core transition guards and state-conditioned validation; preserve structural checks, custom browser UI and harness prompts; add R29/A15 and task coverage; check metadata, links and the baseline diff; obtain independent review; merge via PR and verify remote main. Coordinator alone writes the candidate; the fresh reviewer has a read-only claim. Required review sources are operator brief, minimal-interface decision, product, extensions, validation, interfaces, knowledge format, architecture, implementation plan, acceptance and affected tasks. Other files are read only on contradiction/link triggers.

Review checks architectural simplicity, comprehension and implementation readiness, particularly residual domain-policy enforcement and preservation of existing scope. The actual model/usage telemetry is unknown unless exposed; use an inherited model in fresh context and a bounded review with targeted rechecks. No new services/providers or application tests are required for this prose change. Rollback is a normal corrective PR, preserving published history.

## Acceptance and delivery

Expected outcome: a fresh implementer can distinguish structural standardization from agent-owned workflow interpretation; changing workflow prose does not require backend policy logic. Author checks cover YAML frontmatter/base metadata, relative Markdown links and baseline diff whitespace. The independent review must bind the exact final candidate. The done state is proposed closure conditional on independent approval, normal PR merge and verified remote-main content equality; no completion is claimed before those observations. At that documentation-only checkpoint runtime behavior remained unimplemented and untested. Current implementation and its limits are separately recorded in [completion milestone evidence](../sources/completion-milestone-evidence.md).

Author checks: all 64 Markdown frontmatters parse with required base metadata; all 267 relative Markdown link paths resolve; full baseline diff whitespace is checked against a21838a. These are documentation checks, not runtime tests.
