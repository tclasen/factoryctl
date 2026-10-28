---
type: Task
title: "Submit agent chat from the keyboard"
description: "Enter sends the chat message; Shift+Enter inserts a newline."
status: draft
task_state: done
priority: normal
owner: coordinator
intent_state: requested
intent_owner: human:operator
---

# Submit agent chat from the keyboard

User request: make Enter, Shift+Enter or Command+Enter submit agent chat without tabbing to the button, choosing a familiar convention. Assistant approach: Enter sends; Shift+Enter inserts a newline. Preserve IME composition, suppress held-key repeats and respect disabled/busy state. Show the shortcut beside the composer.

The assistant historically inferred immediate authorization for a small keyboard interaction fix, targeted browser verification and normal reviewed PR/main delivery. That authorization inference is corrected below; it was not a separate user build instruction. The delivered change was not a new milestone or a replacement full-matrix qualification claim. Prior milestone recordings remain historical evidence of their pinned source. No videos or artifacts are uploaded.

Coordinator owns the embedded chat assets, targeted browser fixture check and this record. Independent reviewer has read-only scope over this request and exact candidate; inherited model choice, actual telemetry unknown. Verify initial and follow-up submission, newline, composition, repeat, blank and unavailable states. Existing source authority and project PR/content safeguards apply. Rollback is a scoped revert PR.


## Verification and delivery condition

Targeted Playwright checks pass for initial and follow-up Enter submission, Shift+Enter newline, IME composition/229, held-key repeat, blank and busy states, and the accessible hint. The existing browser interaction fixture also passes its host-isolation/native-request checks. `go test ./internal/web` passes. Project validation retains only the pre-existing Finder `.DS_Store` diagnostic. Local results are under ignored `artifacts/chat-keyboard/20261028/`.

Completion takes effect only after independent exact-candidate approval, normal PR merge and owning-remote main verification. This fix does not requalify the historical milestone video or close Task13's remaining matrix gaps.

## Authorization and traceability correction

The [exact user request](../sources/operator-brief.md#document-title-and-keyboard-submission--recovered-source-excerpts) expresses the desired keyboard behavior. The task originally retained a paraphrase but no dedicated matrix row. U106/R40/A24 is retrospective reconciliation, not evidence that a row existed at implementation time. No separate build instruction for this change is present in the supplied exchange; the earlier immediate-build interpretation is superseded by [planning-first instructions](../decisions/intent-and-review.md). PR #15 and its verification establish observed implementation, not authorization of its timing. The done state is retained as delivery history. See [audit](../sources/feature-intent-audit.md).
