---
type: Specification
title: "Release videos and continuing requirements evidence"
description: "Release videos and continuing requirements evidence."
status: draft
factory:
  format: 1
  assertion:
    frame: desired
    certainty: mixed
    basis: mixed
sources:
  - id: operator
    resource: ../sources/operator-brief.md
    title: Operator requirements and scope
    author: human:operator
---

# Release videos and continuing requirements evidence

The operator requires a Playwright video for every release covering all implemented features, plus a full requirements traceability matrix and fresh evidence that requirements continue to work. This is a release gate, including 0.x releases; it is not a request to fabricate a demo before implementation exists.

## Complete feature coverage

For each release candidate, enumerate implemented features from the [requirements matrix](requirements-traceability.md), implementation and prior release feature inventory. Every supported feature, whether new or old, must map to a recorded Playwright scene and an assertion/result. A deprecated feature still present must still be exercised. Omit only deprecated **and removed** features, with an explicit removal decision/revision and corresponding requirement disposition. Verify no previously supported feature silently disappears from the inventory. Pure backend behavior needs executable evidence plus a browser-visible result/diagnostic or a clearly labeled evaluation panel; a visual caption alone does not prove it passed. If interactive CLI remains supported, accompany its executed assertions with a visible, labeled terminal/result capture in the browser demo; do not pretend that browser clicks tested a terminal.

Cover all three project scenarios and all supported capabilities: discovery, navigation, link/backlink behavior, independent state, edits/validation/conflicts, definitions/templates, runtime conversations/questions, custom views and qualified generative UI. Include meaningful error/recovery journeys. Unsupported/unimplemented capabilities are labeled as such and cannot be presented as demonstrated functionality.

## Recording and presentation

Use Playwright's browser video recording for the complete demonstration. Deliver a watchable, captioned full walkthrough, optionally divided into chapters, with a complete combined coverage index. Use deterministic synthetic fixtures and real UI interactions. Scroll smoothly, move attention coherently and pause after navigation or important state changes so viewers can orient and read. Avoid jump cuts or speedups that conceal waits, errors or failed outcomes. Timing values are tuned by watching the exported recording, not by assuming fixed sleeps prove readability.

Overlay the feature/requirement and what is being demonstrated or evaluated for each scene. Add visible focus rings, outlines, pointer cues or highlights around the exact controls/content needing attention. Overlays must not obscure the interaction or capture clicks; clearly distinguish explanatory annotations from application UI and observed results. Pause long enough to read the annotation and observe the resulting state. Include the release version and tested revision. If multiple videos are needed, identify their order and include timestamps per artifact.

Record assertions, screenshots/trace when needed and actual outcomes independently from the video. Assertions wait for real state; pauses exist for presentation only. Review the finished video for readability, smooth scrolling, truthful overlays/highlights and coverage before release. Report any failed/blocked scene rather than splice it into a passing story.

## Traceability and release gate

Each atomic operator requirement has a stable ID, source, specification, implementation task/code, acceptance criteria, introduction evidence, current release regression result and history/removal disposition. Every release must refresh results for every applicable row against the exact release candidate; an old pass is historical only. States are planned, implemented-unverified, pass, fail, blocked, or deprecated-and-removed. A required feature that is still planned/failed/blocked prevents release acceptance unless the operator explicitly approves a scoped exception; an optional experiment records its outcome without inventing delivered capability.

Store each release's evidence as a versioned Markdown Reference record under `factory/sources/`, linked from the live matrix. This is a release evaluation artifact, not an event ledger or duplicate task-state log. Keep video/binary/trace files outside the OKF Markdown bundle (for example release artifacts) and link their immutable locations plus hashes. Do not overwrite prior release evidence; Git retains matrix history. Each row includes release tag, candidate/merge revision, environment/tool/runtime versions, test ID/command, expected/observed result, verifier, artifacts and video scene/timestamps. Re-run affected checks if the candidate changes; approval of an earlier candidate does not validate later code.

Before publication, verify artifact links/checksums, complete feature inventory, per-requirement current results, video coverage and the human readability review. Missing video, stale evidence or an unexplained coverage hole blocks the release. Task 11 owns this gate; each implementation task supplies its own earlier evidence. No recurring automation, release publication or unsupported-platform test is authorized by writing this contract.
