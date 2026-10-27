---
type: Reference
title: CLI review videos and complete matrix disposition
description: Actual recordings, executed evidence, unmet requirements and milestone gate verification.
status: draft
factory:
  format: 1
  assertion:
    frame: observed
    certainty: supported
    basis: direct-observation
sources:
  - id: cli-video
    resource: ../../artifacts/cli-review/848b4ec/cli-review.mp4
    title: Playwright CLI walkthrough
  - id: matrix-video
    resource: ../../artifacts/cli-review/848b4ec/matrix-review.mp4
    title: Complete requirements accounting video
  - id: matrix-evidence
    resource: ../../artifacts/cli-review/848b4ec/requirements-coverage.json
    title: Current row-by-row evidence and gaps
---

# CLI review video evidence

This package corrects the omitted review video after the first CLI delivery and implements the operator's [every-milestone gate](../specifications/release-evidence.md#milestone-review-delivery). It is a partial product review package, **not full end-to-end qualification of the requirements matrix**. All 94 atomic rows are accounted for: 8 have mapped end-to-end passes, 60 have implemented but partial/unverified coverage, and 26 remain planned. Unsupported GUI/MCP/runtime operations are shown as refusals, never capability passes.

## Delivered artifacts and provenance

Product revision: `848b4ec072755aafa8d6ba9a93a878170b1bf315`. Product Go sources and module files were unchanged; `toolchain.json` records their hashes and the checker compares them against the current checkout. Recorders/evidence/gate instructions are separate changes. Artifacts live outside the Markdown bundle in `artifacts/cli-review/848b4ec/`:

| Artifact | Observed result | SHA-256 |
| --- | --- | --- |
| cli-review.mp4 | 418.84 seconds; H.264, 1440×900; 5,990,970 bytes; 23 chapters, 64 actual command results, 112 exit/semantic checks | b8666d3feddc555143590a194d1276e1ba22b6800f91fa341ab0a6b3420f077f |
| matrix-review.mp4 | 246.32 seconds; H.264, 1440×1000; 4,518,689 bytes; all 94 rows in 24 pages | 0f0845ec752ad2528eedd06d3c5f2704ab511b98165d5ae79f750cd2fa86e104 |
| matrix-input.json | Frozen, independently assessed matrix input used for the recording; later main-video timestamp enrichment is separate | 0367c8018c79b12e0088f98160bc0f5b799af3346ef018a8facb1f2da8b2ecd3 |

`chapters.md` and `coverage.json` provide CLI timestamps; `matrix-video-coverage.json` maps every row to the second video's intervals. `MATRIX.md` is a readable combined index. `requirements-coverage.json` supplies exact partial/unmet portions and named backend cases. `results.json` retains actual command output and expected/observed exits. `test-events.jsonl` and `test-run.json` retain a fresh uncached named-test run: 80 named pass events, zero failed events. `SHA256SUMS` covers the compact versioned package. Raw WebM captures/screenshots remain local intermediates; they are not required to play the committed MP4s.

## Recording and inspection

Playwright 1.64.1 used Chromium 153.0.8010.12 on Darwin ARM64 with Node 24.17.0 and installed FFmpeg 9.0.2. CLI commands ran against deterministic temporary Git fixtures. The browser is visibly labeled a replay/evaluation panel, not a product GUI; it displays summaries/excerpts of real results, while full outputs remain in JSON. The first video scrolls each command into focus and pauses; the second highlights each row and exposes all gaps. Both recordings are silent and captioned. Dense matrix pages can be paused; they do not demonstrate missing functionality.

Developers inspected exported contact sheets and native-size frames, checked page overflow and unique row coverage, used ffprobe for media properties, and decoded the matrix MP4 completely with FFmpeg. Coordinator inspected the exported CLI contact sheet and matrix poster. Visual review is sampled, not a claim that a human watched every frame continuously. Independent exact-candidate artifact review remains a delivery gate.

## Verification and failures retained

The collector rebuilt the current product and ran actual CLI journeys for S1 single project, S2 monorepo and S3 sibling repositories. Fresh normal/race tests and vet passed. Missing-input, conflicting/stale/invalid writes, referenced deletion/heading edits, unsupported operations, duplicate selection and scope behavior have explicit expected outcomes. Grouped backend evidence is mapped to named cases rather than blanket assertions that package success proves all requirements.

Preflight corrected harness expectations for idempotent init, missing-input exit code, scalar quoting and repository-root DTO naming; no product behavior was changed to obtain passing outcomes. The first long recording produced no video because cached FFmpeg lacked the VP9 encoder required by Playwright. This was not reported as a completed video. A task-owned route to installed FFmpeg was verified by exported smoke recording before the successful full rerun. Browser sandbox escalation used ordinary authorized tooling; no provider, deployment, download or release publication was introduced.

The matrix reviewer corrected an overstated non-Git end-to-end claim, removed an unrelated disabled-type test mapping, and separated metadata preservation from epistemic proof. The final map is deliberately conservative. Complete row accounting passes, while full qualification remains incomplete. Eight checker regression tests pass: incomplete qualification, missing/duplicate rows, stale requirement text, stale product source, missing video, unexecuted named test, false partial-to-pass promotion, and missing matrix-video row.

## Mandatory ongoing gate and proposed closure

Project `AGENTS.md`, release policy, R32/R33, atomic U094 and the demonstration specification now require this evidence for every milestone delivery, including unreleased development milestones. `make review-check` verifies structural accounting, source/text freshness, actual named-test passes and video row mappings; `make review-gate` additionally requires every row to pass end to end and intentionally fails for the present product. These tools do not replace independent semantic/visual review or establish future compliance automatically.

Closure of [the review-video task](../tasks/cli-review-video.md) is proposed only for delivering this review package and gate after independent approval and normal verified PR/main delivery. It does not close the remaining product requirements or waive full-matrix qualification. The delivery PR records exact reviewed commit/tree, artifact hashes, review outcome, merge commit and remote-main verification.
