---
type: Task
title: Show the viewed document title in the browser tab
description: Use the requested FactoryCTL prefix and loaded document title.
status: draft
task_state: done
priority: normal
owner: coordinator
factory:
  format: 1
  assertion:
    frame: desired
    certainty: supported
    basis: authoritative-statement
---

# Document page title

User request: set the webpage title to the viewed document title, for example `FactoryCTL - Agent questions and rich web interaction`.

Historical assistant-inferred scope (authorization characterization corrected): update the read-only browser's title on successful document navigation, preserving the existing stale-response guard; focused browser verification and reviewed PR delivery. This is a small UI correction, not a new product milestone or full-matrix qualification claim. Existing demo evidence remains bound to its earlier source revision.

Coordinator owns the single asset edit, this task and delivery. Acceptance: direct document links and subsequent document navigation set the exact requested title; the heading and tab title stay aligned. No new storage, approval or workflow behavior.

Run page-title-review: independent Reviewer, inherited model/reasoning requested (actual telemetry unknown); read-only claim on candidate diff and focused browser results, no implementation authorship. Review source fit and exact candidate before normal protected merge; delivery is observed after merge rather than self-certified in this ticket.

## Authorization and traceability correction

The [exact request](../sources/operator-brief.md#document-title-and-keyboard-submission--recovered-source-excerpts) establishes the desired title behavior. The assistant previously treated it as immediate build permission; that inference is superseded by the [planning-first correction](../decisions/intent-and-review.md). No distinct build instruction for this change is present in the supplied source. Implementation is observed in commit `35874bf70cabd724106398a342c7c6c4ac175c98`; the old in-review label is historical ticket state, not a current delivery verdict. Dedicated U105/R39/A23 mapping is added retrospectively in this audit. It does not establish earlier traceability, build authorization or new video qualification. See [audit](../sources/feature-intent-audit.md).

## Completion revision 2 reconciliation

Current product candidate: `a7c22ba847d58e8f4ece6a22ee9bb73302a64c02`. The explicit instruction to implement the remaining project tasks and its bounds are recorded in the [completion task](project-completion.md); the original authorization and historical checkpoints above remain distinct.

Current browser navigation updates the tab title to the requested FactoryCTL prefix and viewed document title, with current browser assertions and review evidence (U105/A23). This is an observed implementation/regression result under completion revision 2. It does not repair the recorded absence of a distinct historical build instruction for the original title change.

Current command, named-test, native/browser, scene and artifact identities are consolidated in [completion milestone evidence](../sources/completion-milestone-evidence.md), with the independent [code review](../sources/completion-code-review.md) and [acceptance review](../sources/completion-acceptance-review.md). The proposed matrix accounts for 105 passing rows and blocked U045; that proposal is not a full-product acceptance or release claim. Final combined evidence and independent review remain separate gates.

Conditional closure proposal: the bounded work described here is implemented or reconciled and ready for reviewed delivery. The proposed `done` state becomes effective ONLY after this PR is independently reviewed, merged normally and verified on owning-remote main. Before those gates this is a closure proposal, not a delivered result. The final independent reviewer must assess this task's evidence basis. Coordinator retains delivery/reconciliation ownership; no earlier PR/main result is invented by this update. Task 11 and the root completion task retain the unresolved static-linkage/full-qualification boundary.
