---
type: Decision
title: Preserve intent and review tangible outcomes
description: Accepted working-method correction with bounded demo-first checkpoints.
status: stable
decision_state: accepted
factory:
  format: 1
  assertion:
    frame: desired
    certainty: supported
    basis: authoritative-statement
intent_state: accepted
intent_owner: user
intent_source: ../sources/operator-brief.md#working-method-correction-and-tangible-review
intent_authorization: The operator explicitly approved the five working-method proposals and asked to implement them; this does not authorize publishing generated artifacts.
---

# Intent and tangible review

The operator accepted the proposed method in [the quoted reply](../sources/operator-brief.md#working-method-correction-and-tangible-review). The assistant owns the implementation details below.

Each task separates the exact user request with relevant context, agent interpretation/selected approach, and execution authorization. Preserve original authorship when a proposal is accepted. A commit, dependency, summary, passing test or technical review does not establish human acceptance. Distinguish requested, proposed, accepted, superseded and unresolved; mixed documents must not be granted blanket acceptance.

Keep [one current plan](../designs/implementation-plan.md). Reconcile changes and superseded dependencies. Before technical review, an independent reviewer checks whether the chosen work advances the intended user experience and whether the consequential choices are authorized. Repeated correction requires revisiting affected assumptions before another implementation/delegation.

For this operator, use a concise planned-outcome checkpoint instead of a long text-first handoff. A request to add, change, fix or remove a feature changes the plan first. It does not authorize building a prototype or production behavior. Record the source request, requirement/matrix mapping, criteria and plan revision before implementation; then require an explicit build instruction bound to that revision or clearly named slice. An instruction expressly requesting both planning and implementation authorizes that sequence. Approval of a desired outcome or plan alone does not. A short “go ahead” grants a build only when its immediately clear referent is that explicit build action and scope. Once build authority exists, routine engineering, review and delivery proceed without repeat permission.

The operator’s [planning-first correction](../sources/operator-brief.md#planning-first-correction) supersedes any reading of the earlier tangible-demo preference as authority to build from a feature request. After authorized implementation, favor a small working demo and local video; keep accessible text and detailed evidence available. Direct requests to update instructions or planning documents authorize those edits now and do not create a recursive approval loop. Existing implemented behavior and historical evidence remain observed facts, not retrospective proof of build authorization.

The browser displays recorded provenance and sources. Its labels are claims from documents, not approval decisions computed by factoryctl. Viewing a proposal, choosing a default, watching a video or remaining silent never submits acceptance. Demo approval records what behavior and scope were accepted; future concepts remain proposals. All videos and generated artifacts stay local.

This is a bounded working-method correction. It does not guarantee future agent behavior or claim measured reliability improvement. Verify it against the earlier sequence misattribution, requirements-film substitution and hosting-assumption failures.

Fresh-context review follows [the scoped feature audit](../sources/feature-intent-audit.md): exact source or explicitly labeled paraphrase → requirement and atomic row → current plan/task → separate build instruction → observed code/evidence. Missing links or authorization remain visible gaps. Adding a matrix row later repairs traceability; it cannot prove that the row or build permission existed earlier.
