---
type: "Task"
title: "Qualify runtime document extensions before interface implementation"
description: "Qualify runtime document extensions before interface implementation."
status: "draft"
factory:
  format: 1
  assertion:
    frame: "desired"
    certainty: "proposal"
    basis: "llm-inference"
sources:
  - id: operator
    resource: "../sources/operator-brief.md"
    title: "Operator requirements and scope"
    author: "human:operator"
task_state: "planned"
priority: "high"
owner: "unassigned"
depends_on:
  - "03-definition-validation.md"
specifications:
  - "../specifications/extensions.md"
requirements:
  - "R26"
  - "R27"
  - "R28"
---

# Qualify runtime document extensions before interface implementation

## Outcome and scope

Extend the existing definition/validation services with the [declarative extension contract](../specifications/extensions.md): lifecycle rules, namespaced tools and prompt references. No backend plugin runtime, executable hook or custom-code handler is allowed. This is planned implementation, not part of the documentation update.

## Sequence and acceptance

After foundation/storage/structural validation, pin the metadata schema and supported core-operation mappings. Exercise two different types and install a third without rebuilding the binary. Prove shared generic/specialized write enforcement and safe definition reload with in-process fixtures.

Completion covers the core prototype portion of A14 in [acceptance](../specifications/acceptance.md) plus regressions for A03/A04/A06. Test invalid/duplicate definitions, unsupported executable handlers, stale document/definition revisions and allowed/denied transitions. Tasks 06–08 own transport/harness integration, task 09 owns arbitrary browser UI including custom drag/drop, and tasks 10–11 verify the complete A14 journey. Those downstream checks are not prerequisites for this task.

## Risks and rollback

Exact declarative syntax is an implementation choice to qualify; do not introduce a general expression language to replace rejected backend plugins. Arbitrary local computation belongs to prompt-backed harness work. No event ledger, scheduler or multi-file transaction is introduced. Roll back through agent-owned Git changes, preserving user documents and surfacing invalid definitions.
