---
type: Task
title: Define project collections and agent-driven user interaction
description: Reconcile three required project layouts, user questions and an MCP Apps evaluation.
status: draft
factory:
  format: 1
  assertion:
    frame: mixed
    certainty: mixed
    basis: mixed
task_state: done
priority: high
owner: coordinator
depends_on: []
specifications:
  - ../specifications/project-discovery.md
  - ../specifications/interfaces.md
  - ../specifications/gui.md
  - ../specifications/acceptance.md
requirements: [R01, R07, R10, R13, R30, R31]
---

# Project collections and agent-driven user interaction

## Request and authorized scope

The operator requires testing one-project Git checkouts, multiple independent factories in one monorepo, and sibling Git repositories with cross-links, using the same collection UI. A subsequent request adds agent-initiated free-text, multiple-choice and ongoing conversation, with MCP Apps exploration. Strong enough MCP Apps support may justify retiring the interactive CLI; this is a conditional evaluation, not a decision to remove it now. The depth-1 discovery requirement remains controlling.

Coordinator owns documentation changes under factory/specifications/, factory/designs/, factory/tasks/, factory/sources/operator-brief.md and factory/factory.md on branch `docs/project-collections-and-interaction`. Canonical delivery is origin/main of this project through PR. No application implementation, deployment or organization mutation is included. Preserve unrelated files.

## Plan revision 1 and design brief

1. Separate project roots/bundles from Git working-tree roots; define bounded discovery and common UI. Keep project definitions, assets and state independent and cross-links explicit.
2. Add a shared three-scenario acceptance matrix and wire it through incremental build tasks.
3. Specify a harness-neutral question/answer channel, routing and failure behavior; research official MCP Apps and elicitation contracts.
4. Define a bounded rich-web prototype and evidence-based interactive-CLI decision, without claiming Codex support or implementing backend custom code.
5. Check document consistency, links and whitespace; obtain fresh independent design/content review; deliver by PR and verify remote main.

This materially changes the project model and user interaction contract, so independent design review applies before application implementation. Alternatives rejected: treating repository as project identity, unrestricted recursive scanning, assuming native Apps support from MCP connectivity, or dropping the CLI before qualification. Runtime-selected browser components remain within the existing custom-browser-code boundary. No persistence database or workflow engine is introduced.

## Acceptance, run and review manifest

Base: b2bf9e189d7a10a4961b71b89c64bb86f05dfdaf. Dependencies: none for this design change. At this original design-only checkpoint implementations remained planned. Coordinator is the sole author/mutation owner; fresh Reviewer receives the committed candidate and this brief, plus the governing organization/reviewer instructions and changed source set. Review checks depth boundary, all three layouts, identity/storage isolation, identical collection UI, user-answer provenance and correlation, Apps capability/fallback honesty, conditional CLI decision and task/acceptance coverage. Requested model: inherited with fresh context; actual telemetry unknown. Scope is one design/review/delivery cycle with diagnosed retries only. Rollback: scoped corrective PR.

Runtime risks include host bridge support, turn interruption and answer delivery, browser code isolation and cross-project concurrent validation. The responsible prototype/storage tasks must qualify these before capability claims. At that historical checkpoint no runtime test had run. Review and remote delivery remain pending until recorded in the PR; its merged state and reviewed head provide the delivery handoff. If blocked, retain this branch and resume the pending review or delivery step.

## Plan revision 2 — additional operator steering

The operator additionally requires a full-feature Playwright video for every release, with human pacing, smooth scrolling, explanatory overlays and interaction highlights; all old features remain covered unless deprecated and removed. They also require full atomic requirements traceability and ongoing functionality evidence. This authorizes documentation under the existing scope: add R32/R33, A19/A20, the release-evidence contract and atomic matrix; connect every build stage to evidence and task 11 to current-release verification. Acceptance adds complete source-to-requirement coverage, honest absent runtime evidence and no historical-pass substitution. Independent review covers this complete revised scope. At this original design checkpoint no actual release/video existed or was claimed; current local milestone recordings remain distinct from release publication.

## Completion revision 2 reconciliation

Current product candidate: `a7c22ba847d58e8f4ece6a22ee9bb73302a64c02`. The explicit instruction to implement the remaining project tasks and its bounds are recorded in the [completion task](project-completion.md); the original authorization and historical checkpoints above remain distinct.

The bounded collection/question/Apps-evaluation design is documented and has implementation evidence at the current candidate. Current S1–S3 recordings exercise the same collection controls, scope-safe edits and native questions/follow-up; generated views establish first-party rendering, not external Apps interoperability. The documented decision is to retain the CLI. This reconciles the original design-only checkpoint without pretending that its initial request authorized implementation.

Current command, named-test, native/browser, scene and artifact identities are consolidated in [completion milestone evidence](../sources/completion-milestone-evidence.md), with the independent [code review](../sources/completion-code-review.md) and [acceptance review](../sources/completion-acceptance-review.md). The proposed matrix accounts for 105 passing rows and blocked U045; that proposal is not a full-product acceptance or release claim. Final combined evidence and independent review remain separate gates.

Conditional closure proposal: the bounded work described here is implemented or reconciled and ready for reviewed delivery. The proposed `done` state becomes effective ONLY after this PR is independently reviewed, merged normally and verified on owning-remote main. Before those gates this is a closure proposal, not a delivered result. The final independent reviewer must assess this task's evidence basis. Coordinator retains delivery/reconciliation ownership; no earlier PR/main result is invented by this update. Task 11 and the root completion task retain the unresolved static-linkage/full-qualification boundary.
