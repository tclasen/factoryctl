---
type: "Task"
title: "Implement deterministic record and prompt templates"
description: "Implement deterministic record and prompt templates."
status: "draft"
factory:
  format: 1
  assertion:
    frame: "desired"
    certainty: "proposal"
    basis: "llm-inference"
sources:
  -
    id: "operator"
    resource: "../sources/operator-brief.md"
    title: "Operator requirements and scope"
    author: "human:operator"
task_state: done
priority: "normal"
owner: "developer"
depends_on:
  - "03-definition-validation.md"
specifications:
  - "../specifications/templates.md"
requirements:
  - "R06"
  - "R11"
  - "R14"
  - "R15"
  - "R19"
  - "R25"
---

# Implement deterministic record and prompt templates

## Outcome

Implement restricted templating, typed input collection, include resolution, full-output validation and materialization of the software starter set.

## Scope and handoff

Implement the relevant modules in the [architecture](../designs/architecture.md), following the linked specifications in frontmatter. Read [operator authority](../sources/operator-brief.md) before any external effects. Historical CLI checkpoint: the document portion was in progress under the [first CLI milestone](first-cli-milestone.md), with then-uncovered criteria recorded in [CLI evidence](../sources/cli-milestone-evidence.md). The current implemented outcome is reconciled below. The [implementation plan](../designs/implementation-plan.md) establishes build order and shared contracts.

## Verification and done criteria

Test escaping, missing/extra inputs, nested includes, cycles, output limits, invalid rendered YAML and unknown defaults. Initialize a scratch repository from the actual default templates and validate its entire bundle without .factory/.

Acceptance coverage: A05 in [release acceptance](../specifications/acceptance.md). Include exact tested revision, commands, expected/observed results and remaining limitations in current task evidence when implemented. Do not equate a checked task field with proof of behavior. Run relevant negative tests as well as success journeys.

## Dependencies, risks and rollback

`depends_on` identifies prerequisites; an agent interprets readiness. Factoryctl validates references but does not schedule this work. Preserve unknown metadata and source evidence. Roll back failed implementation through ordinary Git changes by the agent, retaining useful evidence. Do not expand into deployment, purchases or external messaging.

## Collection, interaction and continuing evidence

Run S1–S3 with different per-project templates/assets and absent .factory/ directories. Resolve assets from the selected project root; prohibit accidental sibling/parent inheritance. Record A16 evidence.

## Completion revision 2 reconciliation

Current product candidate: `a7c22ba847d58e8f4ece6a22ee9bb73302a64c02`. The explicit instruction to implement the remaining project tasks and its bounds are recorded in the [completion task](project-completion.md); the original authorization and historical checkpoints above remain distinct.

Restricted record/prompt rendering, declared inputs, includes, cycle/error handling, safe initialization and the software starter are implemented. Prompt previews and explicitly invoked actions use shared project context; rendering alone launches nothing. Current regression evidence covers template/dependency pins, customization preservation and projects without optional assets (A05/A16).

Current command, named-test, native/browser, scene and artifact identities are consolidated in [completion milestone evidence](../sources/completion-milestone-evidence.md), with the independent [code review](../sources/completion-code-review.md) and [acceptance review](../sources/completion-acceptance-review.md). The proposed matrix accounts for 105 passing rows and blocked U045; that proposal is not a full-product acceptance or release claim. Final combined evidence and independent review remain separate gates.

Conditional closure proposal: the bounded work described here is implemented or reconciled and ready for reviewed delivery. The proposed `done` state becomes effective ONLY after this PR is independently reviewed, merged normally and verified on owning-remote main. Before those gates this is a closure proposal, not a delivered result. The final independent reviewer must assess this task's evidence basis. Coordinator retains delivery/reconciliation ownership; no earlier PR/main result is invented by this update. Task 11 and the root completion task retain the unresolved static-linkage/full-qualification boundary.
