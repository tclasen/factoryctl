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
task_state: done
priority: "high"
owner: "design-author"
depends_on: []
specifications:
  - "../specifications/product.md"
---

# Prepare the factoryctl implementation design

## Outcome and authority

Produce an internally consistent, human-navigable OKF design package covering the confirmed [operator requirements](../sources/operator-brief.md). This task is documentation and planning only. It does not implement factoryctl or claim its runtime acceptance tests pass.

## Original ownership and plan

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

## Historical documentation-revision checkpoint

The initial design and release policy were delivered through PRs 1 and 2. At that checkpoint the author was revising the design to remove project configuration and separate personal runtime preferences, including the operator's deferred external-tracking-backend idea. Independent review and normal PR delivery are required for this revision. No application implementation is in scope.

## Completion revision 2 reconciliation

Current product candidate: `a7c22ba847d58e8f4ece6a22ee9bb73302a64c02`. The explicit instruction to implement the remaining project tasks and its bounds are recorded in the [completion task](project-completion.md); the original authorization and historical checkpoints above remain distinct.

The bounded implementation design package is present and has informed the implemented remaining slices. The current plan, concrete interfaces and completion design/code/acceptance reviews supersede initial not-yet-implemented checkpoints. Product acceptance, static-linkage disposition and current milestone delivery are separate from closure of this documentation task; initial source/authority limitations remain recorded.

Current command, named-test, native/browser, scene and artifact identities are consolidated in [completion milestone evidence](../sources/completion-milestone-evidence.md), with the independent [code review](../sources/completion-code-review.md) and [acceptance review](../sources/completion-acceptance-review.md). The proposed matrix accounts for 105 passing rows and blocked U045; that proposal is not a full-product acceptance or release claim. Final combined evidence and independent review remain separate gates.

Conditional closure proposal: the bounded work described here is implemented or reconciled and ready for reviewed delivery. The proposed `done` state becomes effective ONLY after this PR is independently reviewed, merged normally and verified on owning-remote main. Before those gates this is a closure proposal, not a delivered result. The final independent reviewer must assess this task's evidence basis. Coordinator retains delivery/reconciliation ownership; no earlier PR/main result is invented by this update. Task 11 and the root completion task retain the unresolved static-linkage/full-qualification boundary.
