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
task_state: done
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

Extend the existing definition/validation services with the [declarative extension contract](../specifications/extensions.md): structural definitions, thin namespaced document tools and Markdown instruction/prompt references. No backend plugin runtime, executable hook or custom-code handler is allowed. At the original documentation-only checkpoint this was planned implementation, outside that update; the current implemented outcome is reconciled below.

## Sequence and acceptance

After foundation/storage/structural validation, pin the metadata schema and supported core-operation mappings. Exercise two different types and install a third without rebuilding the binary. Prove shared generic/specialized write enforcement and safe definition reload with in-process fixtures.

Completion covers the core prototype portion of A14 in [acceptance](../specifications/acceptance.md) plus regressions for A03/A04/A06. Test invalid/duplicate definitions, unsupported executable handlers, stale document/definition revisions and structural enum checks. Tasks 06–08 own transport/harness integration, task 09 owns arbitrary browser UI including custom drag/drop, and tasks 10–11 verify the complete A14 journey. Those downstream checks are not prerequisites for this task.

## Risks and rollback

Exact structural metadata syntax is an implementation choice to qualify; do not introduce a workflow rule language, transition guard or state-conditioned validation to replace agent interpretation. Arbitrary local computation belongs to prompt-backed harness work. No event ledger, scheduler or multi-file transaction is introduced. Roll back through agent-owned Git changes, preserving user documents and surfacing invalid definitions.

## Minimal-interface acceptance

Own the core fixtures in A15 and R29: differing workflow instructions with identical schemas must not change write acceptance. Expose the prose for agent interpretation and preserve structural rejection cases. Reuse existing definition and storage modules; no separate domain-policy engine is required.

## Completion revision 2 reconciliation

Current product candidate: `a7c22ba847d58e8f4ece6a22ee9bb73302a64c02`. The explicit instruction to implement the remaining project tasks and its bounds are recorded in the [completion task](project-completion.md); the original authorization and historical checkpoints above remain distinct.

Declarative namespaced read/create/update/delete bindings, request-time registry reload, definition identities and shared validation are implemented in the existing knowledge/application/MCP/HTTP/CLI services. Current tests reject invalid/duplicate/disabled definitions, executable handlers and stale document/definition tokens. Browser prompt/custom-board integration has current evidence. Workflow prose remains harness-interpreted, and direct harness filesystem access remains an advisory boundary rather than factoryctl-enforced confinement (A14/A15).

Current command, named-test, native/browser, scene and artifact identities are consolidated in [completion milestone evidence](../sources/completion-milestone-evidence.md), with the independent [code review](../sources/completion-code-review.md) and [acceptance review](../sources/completion-acceptance-review.md). The proposed matrix accounts for 105 passing rows and blocked U045; that proposal is not a full-product acceptance or release claim. Final combined evidence and independent review remain separate gates.

Conditional closure proposal: the bounded work described here is implemented or reconciled and ready for reviewed delivery. The proposed `done` state becomes effective ONLY after this PR is independently reviewed, merged normally and verified on owning-remote main. Before those gates this is a closure proposal, not a delivered result. The final independent reviewer must assess this task's evidence basis. Coordinator retains delivery/reconciliation ownership; no earlier PR/main result is invented by this update. Task 11 and the root completion task retain the unresolved static-linkage/full-qualification boundary.
