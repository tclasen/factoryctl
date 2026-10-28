---
type: "Task"
title: "Deliver the default project, record, board and conversation GUI"
description: "Deliver the default project, record, board and conversation GUI."
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
task_state: "planned"
priority: "normal"
owner: "unassigned"
depends_on:
  - "07-codex-adapter.md"
specifications:
  - "../specifications/gui.md"
requirements:
  - "R30"
  - "R07"
  - "R08"
  - "R10"
  - "R11"
  - "R17"
  - "R22"
  - "R24"
  - "R25"
---

# Deliver the default project, record, board and conversation GUI

## Outcome

Build the embedded browser UI with all-project navigation, generic list/detail views, keyboard-operable boards, epistemic/source display, direct edits, live diagnostics and conversation panels.

## Scope and handoff

Implement the relevant modules in the [architecture](../designs/architecture.md), following the linked specifications in frontmatter. Read [operator authority](../sources/operator-brief.md) before any external effects. This task is planned, not started; source paths named in the architecture do not yet exist. The [implementation plan](../designs/implementation-plan.md) establishes build order and shared contracts.

## Verification and done criteria

Exercise ticket-update end to end, wrong-type validation, stale-edit reconciliation, approval/input UI, zero backlinks, unknown types, disconnected runtime and keyboard navigation. Confirm no browser persistent database/localStorage is required.

Acceptance coverage: A09 in [release acceptance](../specifications/acceptance.md). Include exact tested revision, commands, expected/observed results and remaining limitations in current task evidence when implemented. Do not equate a checked task field with proof of behavior. Run relevant negative tests as well as success journeys.

## Dependencies, risks and rollback

`depends_on` identifies prerequisites; an agent interprets readiness. Factoryctl validates references but does not schedule this work. Preserve unknown metadata and source evidence. Roll back failed implementation through ordinary Git changes by the agent, retaining useful evidence. Do not expand into deployment, purchases or external messaging.

## Collection, interaction and continuing evidence

Exercise the identical collection shell and complete common journeys in S1–S3, including the one-project collection. Integrate question cards, actual replies, follow-up and the qualified rich-view bridge/fallback. Record A16/A17 and applicable A18 evidence, with video-ready deterministic fixtures.

## Current slice boundary

The [visible-intent task](visible-intent.md) delivers a read-only browser slice independently of completed conversation integration. The dependencies above apply to this full GUI task's conversation/editing acceptance, not to every useful browser interaction. Use the [current plan](../designs/implementation-plan.md) rather than the superseded CLI-first order.
