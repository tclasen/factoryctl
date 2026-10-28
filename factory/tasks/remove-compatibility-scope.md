---
type: Task
title: Remove the low-level compatibility expansion from scope
description: Reconcile the operator's scope correction without resuming product implementation.
status: draft
task_state: in_review
priority: high
owner: coordinator
intent_state: requested
intent_owner: human:operator
---

# Remove the low-level compatibility expansion from scope

Source: [operator correction](../sources/operator-brief.md#compatibility-scope-correction--2026-10-28). Target/classification: project. Authorized change: project instructions, requirements, acceptance, decision, plan and task/evidence summaries only. Plan revision 4 is the current scope. This instruction authorizes documentation and its ordinary reviewed PR/main delivery; it does not resume product work or authorize deleting paused implementation files.

Coordinator owns this documentation run and serialized Git delivery. Independent Reviewer checks source/outcome fit before exact-candidate correctness. Inherited model/reasoning; actual telemetry unknown. Commands bounded at 120 seconds, two diagnosed retries, checkpoint at review/delivery. No runtime experiments, new build pipeline, old-OS research, generated media or artifact upload. Rollback is a scoped documentation PR preserving the operator correction.

Acceptance: explicitly forbid shellcode/raw or inline assembly and low-level compatibility substitutes; remove custom runtime, loader, ABI and legacy-browser expansion from active work; retain the ordinary-user portability outcome and identify the 2016 OS coverage gap honestly. Do not change unrelated product requirements, claim completion, weaken existing product-demo gates or resume paused workers. This documentation correction requires document review, not a new product milestone/demo.

The paused `feat/portable-qualification` checkout and ignored evidence are retained for recovery only. None of its product, tooling or artifact changes belongs to this delivery. No cleanup/deletion is implied. This task uses an isolated documentation branch from the previously delivered main, so paused experiments cannot enter the PR. Independent review and verified owning-remote main delivery remain required before this task closes; root completion and release qualification stay nonterminal.
