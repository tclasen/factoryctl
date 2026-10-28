# Local milestone review recording

`record.mjs` builds the current document CLI, creates temporary synthetic Git projects, executes real commands and assertions, runs fresh Go test/race/vet checks, and records a browser replay of those results with Playwright `recordVideo`. The browser is explicitly labeled **CLI review harness**, not a product GUI. No commands are executed by browser clicks. Displayed output uses labeled summaries/excerpts; `results.json` retains full stdout/stderr, expected/actual exit codes and assertion labels. Commands stop the run on unexpected results.

From this directory, install the pinned recorder dependency (`npm install`) and its browser/encoder (`npx playwright install chromium`). From the repository root:

```sh
REVIEW_OUTPUT_DIR=/private/tmp/factoryctl-review-new node scripts/review-video/record.mjs
```

Use `--check-only` to execute all assertions without browser capture. `REVIEW_OUTPUT_DIR` defaults to the original `artifacts/cli-review/848b4ec` evidence directory, and a completed WebM or MP4 there is protected from overwrite. Always select a fresh destination for another capture. Fixtures and the built executable use a fresh temporary directory; `results.json` identifies that directory. The harness creates Git repositories only in those fixtures; the CLI itself never invokes Git. It does not edit product sources, the project Git index, user files or remotes. Test build cache defaults to `/private/tmp/factoryctl-review-go-cache` when `GOCACHE` is not set. Temporary fixtures remain available for inspection.

Environment overrides support existing local tooling without modifying another checkout:

- `PLAYWRIGHT_MODULE`: absolute path to an installed Playwright module; otherwise normal local resolution uses the pinned dependency.
- `REVIEW_CHROMIUM`: optional explicit Chromium executable, otherwise Playwright's installed browser.
- `PLAYWRIGHT_BROWSERS_PATH`: standard Playwright browser/encoder cache override.
- `REVIEW_OUTPUT_DIR`: fresh output directory for results, screenshots, coverage and video.

The original host had Playwright 1.64.1, Chromium cache revision 1243 and system FFmpeg. It used an explicit Chromium path and a task-owned temporary encoder cache pointing to system FFmpeg in the revision directory expected by Playwright. Browser launch required normal macOS sandbox escalation; no downloads or external services were used. A normal matching `playwright install chromium` is the portable reproduction path. This capture qualifies macOS ARM64 only.

The first run exposed incorrect harness expectations (identical init is idempotent; missing input exits 1; scalar edits may quote the replacement) and local cache/launch restrictions. An older VP8-only Playwright encoder silently failed to export the first full attempt; system FFmpeg with VP9 support resolved this, and the recorder now smoke-exports a short video before the full capture. Those expectations were corrected to the existing CLI contract before capture. No product behavior was changed to obtain passing results.

## Coverage and limitations

The 23 chapters cover all supported command families, help/version, both validation levels, type/query/pagination/cursor expiry, exact source/revision, scalar preservation, stdin/file replacement, incoming reference protection and deletion, explicit input errors, stale/existing conflicts, containment/anchor protection, software initialization, templates and prompt discovery, duplicate checkout selection, CWD versus explicit scope, S1 single repository, S2 monorepo, S3 sibling repositories, custom definitions/templates, cross-project links/backlinks and isolated state. Fresh test/race/vet output supplies grouped backend regression evidence for parsing, include bounds, special files, concurrency and atomic-write behavior; these backend cases are not each individual visual scenes. No claim is made that this is a full release traceability matrix.

The legacy `gui` entry point, MCP, runtime launches, external network validation and view export are exercised as explicitly unsupported errors. No release or unsupported-platform qualification is implied. Assertions establish observed structural behavior, not task approval or workflow completion.

`coverage.json` gives recording-relative scene start/end seconds; screenshots are taken near each chapter's opening. Per-command focus scrolls smoothly and pauses 4.5 seconds; short chapters pause at least 16 seconds. The full exported video is continuous, without acceleration or removed failures.

For a compact shareable MP4, preserve the complete WebM duration:

```sh
ffmpeg -i "$REVIEW_OUTPUT_DIR/cli-review.webm" -c:v libx264 -preset slow -crf 28 -pix_fmt yuv420p -movflags +faststart -an "$REVIEW_OUTPUT_DIR/cli-review.mp4"
ffprobe -v error -show_entries format=duration,size -show_entries stream=codec_name,width,height -of json "$REVIEW_OUTPUT_DIR/cli-review.mp4"
```

The original artifact directory includes hashes, an extracted poster/contact sheet, results and timestamp coverage. Raw capture/screenshots may remain local; the compact MP4 and machine-readable evidence are the delivery artifacts.

`python3 scripts/review-video/finalize.py OUTPUT_DIRECTORY` performs that conversion, extracts actual-video inspection frames and writes `media.json`, `chapters.md` and `SHA256SUMS`. It refuses to overwrite existing media. `toolchain.json` in the original evidence includes the inspected product-source manifest; the collector generates an equivalent source manifest for fresh runs.

## Complete matrix accounting and gate

Run `annotate.mjs` after converting the continuous capture to MP4. It adds requirement IDs and concrete result explanations to every actual scene, preserving duration and command output. It writes the immutable `demo.mp4`, `annotations.json`, and an identical `artifacts/demo/latest.mp4`. It refuses to overwrite an archived demo. The only final demo link is /Users/Shared/factoryctl/artifacts/demo/latest.mp4. Set `REVIEW_SKIP_ALIAS=1` when preparing a segment for a combined demo; update the stable alias only after reviewing the combined export. Raw source captures remain local intermediates; their hashes are recorded in annotation provenance.

The all-row matrix is supporting documentation, not a separate film. The combined browser/CLI package accounts for U001–U100 but is not fully qualified. Run `make review-check` for structural reconciliation and `make review-gate` for mandatory full end-to-end qualification. The latter intentionally fails on partial/unimplemented requirements. `node --test scripts/review-video/check-matrix.test.mjs` checks missing/stale/overstated evidence, scene annotations and stable-link integrity. Independent semantic and exported-video review remain required.

The checker also requires `ffprobe` on PATH to verify the actual exported closing-hold endpoint. All scene starts/intermediate ends retain strict timing checks; the final capture-clock tail is bounded to two seconds and must match the real media endpoint.


## Actual browser product capture

`visible-intent.mjs` builds its own temporary executable, starts its own loopback server and records seven actual browser journeys against the real planning bundle. It records binary identity, Go build information and matching source manifests before/after capture, including embedded browser assets. It rejects an externally supplied server URL and refuses to overwrite completed or failed captures. No source changes, page errors or external requests are permitted during recording.

```sh
REVIEW_OUTPUT_DIR=artifacts/intent-review/FRESH_ID node scripts/review-video/visible-intent.mjs
```

The milestone demo places this browser segment first, followed by the freshly captured and annotated CLI regression segment. Preserve both full durations without cuts or acceleration; pad different frame heights without covering content. Offset CLI scene IDs/timestamps in the combined evidence, retain both binary/source manifests, and independently review all row dispositions against current results. Historical mappings may help locate scenarios but cannot establish a current pass. The archive and stable alias must match before `make review-check`; full qualification remains a separate gate.

Checker unit tests create synthetic metadata and fake media with a test-only probe. They require no private local artifacts and do not qualify a real recording. Actual delivery uses real FFprobe, full decode and exported-frame inspection.
