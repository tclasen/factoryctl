---
type: Decision
title: Historical CLI-first sequence
description: Superseded assistant-selected implementation order; not an original operator instruction.
status: deprecated
decision_state: superseded
factory:
  format: 1
  assertion:
    frame: observed
    certainty: supported
    basis: authoritative-statement
intent_state: superseded
intent_owner: assistant
intent_source: ../sources/operator-brief.md#working-method-correction-and-tangible-review
intent_authorization: The operator authorized building the presented bounded CLI milestone; no source establishes that they originally selected CLI before the web UI.
---

# Historical CLI-first sequence

The assistant's initial implementation plan selected CLI → transports → conversations → GUI. The original source brief explicitly treated implementation order as an LLM proposal. The later milestone admission compressed the assistant recommendation and the user's “please build it now” into an “original request.” That did not make the ordering user-authored.

The [current plan](../designs/implementation-plan.md) supersedes that ordering with user-visible slices. Existing document services remain useful; this record neither removes the completed CLI nor asserts that its bounded execution lacked authorization. The [accepted working method](intent-and-review.md) prevents repeating the attribution error.
