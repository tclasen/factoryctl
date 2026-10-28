---
type: Reference
title: Visible intent browser and CLI regression review
description: Current local demo, source provenance, executed checks and complete matrix gaps.
status: draft
factory:
  format: 1
  assertion:
    frame: observed
    certainty: supported
    basis: direct-observation
sources:
  - id: local-demo
    resource: ../../artifacts/milestone-review/20261028/demo.mp4
    title: Browser-first combined Playwright demo
  - id: current-matrix
    resource: ../../artifacts/milestone-review/20261028/requirements-coverage.json
    title: Current 100-row evidence and gap accounting
---

# Visible intent review evidence

The single stable local handoff is `/Users/Shared/factoryctl/artifacts/demo/latest.mp4`. The actual browser explorer appears first (0:00–1:25.48), followed by a fresh CLI regression replay (1:25.48–8:24.52). Requirement IDs and explanations appear on the real scenes; this is not a requirements-only film. All generated media and evidence remain ignored and local, outside Git.

## Candidate and verification

Product source revision `8c6dc7cf166214da47a7f0442285f065c301e3db`; subsequent recorder and evidence-document changes do not change Go/module/embedded browser sources. Current source manifests include Go, HTML, CSS and JavaScript. Both captures built fresh executables. Browser binary SHA-256 `2372c0534b47cb22b5caa295916ffba5fc33c96729b146006346046b179c7e22`; CLI binary SHA-256 `7576b1d7eb3dc576170d5e2f1ee0ce345d42e092d7153979185753641292a973`. Browser capture owns its server and records build information, URL and identical pre/post source hashes, rather than trusting an unrelated running server. Dirty repository metadata reflects recorder/document work, not changed product code.

Seven actual browser scenes contain ten assertions; 23 CLI replay scenes contain 64 commands and 112 assertions. Browser checks cover the current plan, all five recorded intent states, bounded accepted scope, source navigation, raw Markdown, superseded assistant attribution and absence of save/approval actions. No page errors, external requests or source mutations occurred. The CLI capture exercises prior command families and S1–S3 scenarios, including expected failures and backend regressions.

Fresh uncached Go tests contain 85 named passing events and no failures; race tests and vet pass. Fifteen self-contained evidence-checker tests pass, including stale embedded assets, missing rows, unsupported pass claims, annotation timings and stable-alias identity. These synthetic checker fixtures do not qualify product media. Repository guard tests pass. Execution evidence is macOS ARM64 only.

The combined H.264 export is 1440×1000, 25 fps, 12,613 frames, 504.52 seconds and 10,842,193 bytes. SHA-256: `43eda96b9ec6c7fb09eb68154d6cff7a55babe2952241fa9ced0dac5e9528a45`. Its archive and stable alias are identical. Combination preserves the exact sum of source frames/durations, adding bottom padding to the CLI segment without obscuring content. Full decode and actual exported frame/contact-sheet inspection passed; this is not a claim of uninterrupted human viewing.

## Complete accounting, incomplete qualification

The current local `MATRIX.md`, JSON rows, command outputs, named-test events and annotation timings account for U001–U100: nine demonstrated rows, 65 partial/unverified and 26 planned. `make review-check` passes. `make review-gate` fails as expected because full product requirements remain unmet. U096–U099 document and demonstrate the current working-method repair without claiming that a single run guarantees future agent behavior. U100 covers the bounded inspectable planning journey.

The browser is a read-only startup snapshot with search/type/intent filters, local source links and inert Markdown. Recorded labels are claims, not verified human approval. It does not implement editing, conversations, MCP Apps, canonical cross-project URL navigation/backlinks, watcher refresh or acceptance submission. The legacy CLI `gui` refusal is not evidence that this new `web` explorer is absent. Missing capabilities remain explicit gaps, never passes from an unsupported-operation error.

The live bundle retains a user-owned Finder `.DS_Store` diagnostic; the recorder neither removes nor hides it. The clean Markdown corpus with its referenced local evidence validates with zero errors/warnings and 28 unchecked external links. A clean clone without local artifacts reports missing artifact references; it must generate local evidence before those checks. The first browser attempt is retained locally as superseded evidence because its external-server provenance was weaker. CLI annotation initially detected a blank-to-title paint at 0.04s; the bounded initial-paint correction retains every later title transition and strict timing checks. An intermediate checker-test quoting error was fixed before its passing run.

## Review and delivery boundary

Independent intent-fit readiness review preceded technical acceptance. Independent final review checks this exact product source, recorder changes, current row dispositions and exported media. The [visible-intent task](../tasks/visible-intent.md) closes only after that review and verified protected PR/main delivery of code and documents. The PR records the final reviewed candidate and observed delivery. No release publication, full-matrix acceptance, CLI retirement or MCP Apps qualification is implied.
