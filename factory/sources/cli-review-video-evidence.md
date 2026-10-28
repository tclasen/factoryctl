---
type: Reference
title: CLI demo and complete matrix disposition
description: Actual recording, executed evidence, unmet requirements and milestone gate verification.
status: draft
factory:
  format: 1
  assertion:
    frame: observed
    certainty: supported
    basis: direct-observation
sources:
  - id: cli-video
    resource: ../../artifacts/cli-review/848b4ec/demo.mp4
    title: Playwright CLI walkthrough
  - id: matrix-evidence
    resource: ../../artifacts/cli-review/848b4ec/requirements-coverage.json
    title: Current row-by-row evidence and gaps
---

# CLI review video evidence

[Watch the milestone demo](https://raw.githubusercontent.com/tclasen/factoryctl/main/artifacts/demo/latest.mp4). This is the single canonical link for milestone handoffs, updated with each milestone. Requirement IDs and concrete outcome explanations appear atop the actual CLI scenes. The complete row/gap matrix remains supporting documentation, not a separate requirements presentation.

This package corrects the omitted video after the first CLI delivery and implements the [every-milestone gate](../specifications/release-evidence.md#milestone-review-delivery). It is a partial product review package, **not full end-to-end qualification**. All 95 atomic rows are accounted for: 8 mapped end-to-end passes, 61 partial/unverified, 26 planned. Unsupported GUI/MCP/runtime refusals do not pass those capabilities.

## Artifact and provenance

Product revision: `848b4ec072755aafa8d6ba9a93a878170b1bf315`. Go sources and module files are unchanged; `toolchain.json` records their hashes, checked against the current checkout.

The versioned `artifacts/cli-review/848b4ec/demo.mp4` and canonical `artifacts/demo/latest.mp4` are byte-identical: SHA-256 `23fde6a71da2c4451c0ef061e9d73d88fda399e6259013c5c56ea0291b4af2fe`. The H.264 video is 418.84 seconds, 1440×900, 10,471 frames, 7,958,457 bytes. Its 23 scenes contain 64 actual command results and 112 exit/semantic checks. The annotation pass preserves the original continuous capture duration and frames; a 100-pixel top band adds scene-specific requirement IDs and explanations without covering command output. `annotations.json` records source/output hashes, exact exported scene timings and text. Superseded unannotated and requirements-only videos remain local/history, not delivery links.

`MATRIX.md` and `requirements-coverage.json` retain every row's evidence and exact remaining gap. `coverage.json` and `chapters.md` supply capture-relative timings; annotations record actual exported transitions. `results.json` retains complete stdout/stderr and expected/actual exits. `test-events.jsonl` and `test-run.json` retain a fresh uncached Go run: 80 named pass events, no failed events. `SHA256SUMS` covers the versioned package and canonical alias.

## Recording and inspection

Playwright 1.64.1 used Chromium 153.0.8010.12 on Darwin ARM64, Node 24.17.0 and FFmpeg 9.0.2. Actual commands ran against temporary Git fixtures. The browser is visibly a CLI evaluation replay, not a product GUI, and labels output excerpts; full results remain in JSON. The silent captioned recording scrolls commands into focus at readable pacing without cutting or acceleration. Post-capture annotations explicitly label partial coverage and unsupported boundaries.

Developers decoded the exported video and inspected native frames at multiple scenes. The final poster/contact sheet are extracted from the annotated demo. Independent exact-candidate review checks visual/semantic overlays, timing, hashes, source freshness, accounting and regression tests. Frame sampling is not a claim of uninterrupted human viewing.

## Verification and retained failures

The collector rebuilt the CLI and ran S1 single-project, S2 monorepo and S3 sibling-repository journeys. Fresh normal/race tests and vet passed. Actual negative outcomes cover stale/invalid writes, referenced deletion and heading protection, duplicate selection, unsupported operations and scope boundaries. Backend test mappings remain distinct from visual E2E claims.

Preflight corrected harness expectations for idempotent init, input exits, scalar quoting and DTO naming; no product code changed. The first recording did not export because cached FFmpeg lacked VP9. Installed system FFmpeg was smoke-tested before repeating the full capture. That failed attempt was never reported as successful. The operator then corrected the presentation: the separate matrix film was superseded by annotations on the actual demo.

The independent matrix assessment downgraded an overstated non-Git claim, removed an unrelated disabled-type test mapping and distinguished metadata preservation from epistemic proof. The map stays conservative. Complete accounting can pass while full qualification fails.

## Mandatory gate and proposed closure

Project `AGENTS.md`, release policy, R32/R33 and atomic U094/U095 require this format at every milestone, including unreleased development milestones. `make review-check` validates row/source/text freshness, observed named-test passes, scene annotations, canonical link and archive/latest identity. `make review-gate` additionally requires every row to pass end to end; it intentionally fails for this partial product. Fourteen checker regression cases exercise missing, stale and overstated evidence plus scene/link/alias integrity. Tools do not replace semantic/visual review or guarantee future compliance automatically.

Closure of [the review-video task](../tasks/cli-review-video.md) is effective only after independent approval and verified PR/main delivery. This closes the omitted artifact and delivery-rule work, not unimplemented product requirements. The delivery PR records exact candidate, review, merge and remote-main evidence.
