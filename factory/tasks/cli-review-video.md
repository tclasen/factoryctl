---
type: Task
title: Deliver the missing CLI review video
description: Record and review a Playwright walkthrough of the implemented document CLI.
status: draft
task_state: done
priority: high
owner: coordinator
depends_on: [first-cli-milestone.md]
requirements: [R32, R33]
factory:
  format: 1
---

# CLI review video

- Target System: project
- Change Classification: project
- Authorized Scope: scripts/review-video/, review dependency manifests, artifacts/cli-review/, AGENTS.md, Makefile, README.md and project evidence/instructions; ordinary PR delivery
- Affected Paths or Resources: reproducible Playwright recorder, local review video/assertions/coverage artifacts and factory/ documentation
- Owning Role: Coordinator
- Related Requirements: R32, R33
- Related Tasks: first-cli-milestone
- Approval Requirements: operator request for the missing Playwright review video; independent artifact/code review and normal PR/main delivery

## Intent and plan

The operator asked where the Playwright review video was after CLI milestone delivery. The earlier implementation handoff omitted this review artifact. Provide it now for the implemented CLI, even though no versioned product release is being published. Future milestone handoffs require review video coverage for implemented features, not only a statement that release qualification is pending.

Baseline product: main commit 848b4ec072755aafa8d6ba9a93a878170b1bf315. Read the release-evidence specification, README, interfaces and CLI implementation before building the scene inventory. No product behavior, GUI, runtime, MCP, release tag, deployment or harness configuration changes are in scope. The review panel is test presentation under the existing release-evidence contract, not a new product interface or architecture decision; artifact/code review covers it without reopening product design.

Record actual commands/assertions in deterministic synthetic Git project fixtures. Playwright displays those executed command results in a clearly labeled browser evaluation panel; browser interactions do not claim to test a product GUI. Cover every advertised command and S1–S3 collections, meaningful safe-write/error paths, custom definitions/templates, isolation and current regression evidence. Keep captions, focused highlights and human-readable pacing. Record screenshots and per-scene outcomes independently of the video. Unsupported capabilities must be labeled honestly.

Use available local Node/Playwright/browser tooling; new paid tools or services are not authorized. Keep videos outside factory/. Produce a watchable artifact, hashes, tested revision and timestamp coverage index. Inspect the finished exported video, then obtain independent review and deliver recorder/evidence through a normal PR. No product release is implied.

## Claims and verification

Developer run VIDEO-implementation exclusively owns scripts/review-video/, its dependency manifests, and artifacts/cli-review/ during capture. Coordinator owns factory/ documentation, Git index and remote delivery. Reviewer run VIDEO-review is read-only and independent. Inherit the available model; actual telemetry is unknown. No separate task-thread tool is available, so use native collaboration agents under existing organization instructions. Tests and recording may take up to ten minutes per command; checkpoint progress and failed attempts. Rollback removes only task-owned generated outputs or uses an ordinary revert.

Acceptance: watchable Playwright-recorded video with truthful real-command results, readable captions/highlights, full implemented-feature coverage and expected failures/assertions, reproducible script, artifact hashes, and an independently checked exact candidate. Prior CLI tests remain evidence only for their pinned code; run current relevant checks as part of this review artifact.

## Plan revision 2 — complete-matrix gate

The operator explicitly requires every milestone delivery to include a video exercising the complete traceability matrix end to end. Expanded authorized documentation scope includes the operator brief, product requirements, requirements matrix and delivery policy. Add U094 without renumbering prior rows. Coordinator owns a disjoint machine-readable all-row coverage map and checker; Developer retains recording ownership. The current CLI video is review evidence, not a claim that GUI/MCP/runtime or all matrix requirements pass. Produce a matrix review addendum accounting for every row, with implemented evidence and exact unmet portions. Independent review must distinguish complete accounting from actual end-to-end satisfaction. Full-matrix qualification remains incomplete where implementation is missing; no exception is inferred.

### Recording checkpoint and revised capture budget

The first full browser capture produced screenshots and passing executed results but no video: the older cached encoder lacked VP9 required by the installed Playwright. No success is claimed for that attempt. Developer repaired the task-owned encoder route to use installed system FFmpeg and verified a short exported smoke video before repeating the capture. A full uncached evidence rerun accompanies the replacement capture.

For the expanded matrix review, individual capture/conversion commands may run up to 15 minutes with progress checkpoints; this is a bounded exception to the earlier ten-minute default. A separate Developer matrix-video run owns scripts/review-video/matrix-record.mjs and matrix-review video/coverage/poster artifacts only. Main capture retains its existing source/output claim; Coordinator exclusively owns check-matrix.mjs, requirements-coverage.json, matrix-summary.json, project documentation and Git operations. Parallel headless browser contexts use disjoint video directories and share installed executables read-only.

Coordinator also owns project-local AGENTS.md and Makefile review-check/review-gate entry points under the explicit every-milestone directive. These are project instructions and test tooling; reusable organization instructions are unchanged. The accounting check can pass while full qualification intentionally fails on unmet requirements.

## Proposed review-package closure

The single annotated demo at the canonical stable link, all-row supporting gap mapping, current named-test events, reproducible recorder and tested mandatory gate are ready in [the evidence record](../sources/cli-review-video-evidence.md). The done state becomes effective only after independent exact-candidate artifact/code approval and normal verified PR/main delivery. Until then this task remains operationally in review. This closes the missing-review-artifact work only: full product matrix qualification remains incomplete, and no operator exception is inferred. Developer capture claims are released; Coordinator retains delivery ownership until verified merge.

## Plan revision 3 — one scene-annotated demo

Before PR 9 merged, the operator rejected a requirements-only video and required requirement IDs/result explanations atop actual CLI scenes plus one stable latest-demo link. Prior candidate approval no longer covers the revised format. Retain the real continuous Playwright CLI recording and add per-scene annotations; do not fake additional functionality or rerun unchanged product solely for captions. The separate matrix video is superseded and removed from the delivered tree (history retained); all-row accounting stays in documents. Add atomic U095. Developer owns annotation tooling, versioned demo.mp4 and identical artifacts/demo/latest.mp4. Coordinator owns updated metadata/checker/docs and Git delivery. Independent review must verify semantic overlays, unchanged scene timing, archive/latest byte identity and the final canonical link. Full-product matrix qualification remains incomplete.

## Plan revision 4 — local-only recordings

The operator rejected GitHub video publication as unauthorized. Coordinator owns removal of all tracked recordings, local preservation, ignore rules, local-link evidence/checker changes and ordinary PR delivery of code/documentation only. Independent Reviewer checks the candidate and local video hash. Prior uploaded commits remain historical remote objects until separately purged by GitHub; a deletion commit does not erase history. Do not bypass protected-branch history restrictions.

## Plan revision 5 — purge generated artifact history

The operator explicitly requested removing the entire artifacts folder from all Git history and is temporarily removing repository safeguards for the rewrite. Coordinator owns a local-only recovery backup, exact content-preserving history rewrite excluding artifacts/, verification and lease-protected force push after safeguards are removed. Add ignore rules, pre-commit/pre-push guards and CI checks rejecting artifacts directories and recordings. Independent Reviewer verifies guard behavior, preserved product content and rewritten history. This explicit request authorizes the bounded history rewrite; it does not authorize publishing local backups or altering unrelated history content. GitHub-owned PR references/cached objects require separate host cleanup and must not be reported erased by a branch rewrite.
