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

For this operator, favor a small tangible working slice and a local demo video over a long text-first handoff. Use a brief outcome checkpoint, then build already-authorized reversible work. A short “go ahead” grants the clear slice, not unread alternatives or publication. Ask only unresolved consequential questions; no extra permission loop for routine implementation. Keep accessible text/transcripts and detailed evidence available.

The browser displays recorded provenance and sources. Its labels are claims from documents, not approval decisions computed by factoryctl. Viewing a proposal, choosing a default, watching a video or remaining silent never submits acceptance. Demo approval records what behavior and scope were accepted; future concepts remain proposals. All videos and generated artifacts stay local.

This is a bounded working-method correction. It does not guarantee future agent behavior or claim measured reliability improvement. Verify it against the earlier sequence misattribution, requirements-film substitution and hosting-assumption failures.
