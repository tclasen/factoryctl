---
type: Reference
title: Completion web implementation contract and checks
description: Browser document transport, editing, boards and repository assets under completion revision 2.
status: draft
factory:
  format: 1
  assertion:
    frame: observed
    certainty: supported
    basis: executable-checks
---

# Completion web implementation contract and checks

Run `completion-web` implements the browser portion of [completion revision 2](../designs/implementation-plan.md#completion-plan-revision-2--authorized-implementation). The Coordinator delivered independent design readiness before product edits. This record is implementer evidence, not independent acceptance or final delivery.

## Interfaces

Existing `/api/documents` and `/api/document?id=…` remain compatible read routes. New `/api/v1/session` returns a process-local document nonce; every mutation requires `X-Factory-Nonce` independently of conversation/runtime availability. Existing exact loopback Host, same Origin and cross-site checks remain. No CORS or login database is added.

New routes are `/api/v1/projects`, `/definitions`, `/nodes` (GET/POST), `/node` (GET/PUT/DELETE), `/links`, `/validate`, `/actions`, `/actions/render`, `/extensions`, `/extensions/call`, `/init`, `/delivery`, and `/custom` metadata plus `/custom/trust` and `/custom/request`. Node selection uses `project_id`, opaque `checkout` and bundle-relative `path`. Writers require `expected_definition`; existing records also require `expected_revision`. Update accepts exactly one of `set` or full `document`. All writes delegate to the shared app service. Delete returns `{removed: ref}`. Lists paginate against snapshot-bound cursors. Definition responses expose structural fields, original metadata, enablement/diagnostics and effective definition identity.

Handler document caches use a dedicated mutex; conversation state retains its own lock. A fresh `LoadCached` checks current bytes and structural diagnostics on requests. Document/diagnostic identity controls renderer reuse. Canonical resolved links and incoming edges are indexed once before rendering. Independent delivery polling observes new evidence outside the bundle without resetting editor drafts. No watcher daemon or persistent browser/index database is introduced. Browser polling is 500 ms without overlapping requests; end-to-end refresh includes scan cost.

## Visible behavior

The common project selector supports All projects, overview/type counts, search, type/intent filters, generic details/raw source, backlinks including an empty state, structural diagnostics and delivery observations. Existing intent/source labels, document tab title, conversation choices/defaults and keyboard submission remain. Long lists use fixed-size next/previous pages and boards cap the complete selected collection to 200 cards per page; boards derive from a declared enum view (or the default task_state enum), group incompatible project definitions separately and show invalid/unassigned values explicitly.

Raw and JSON field editing retain the observed revision and definition. A stale record or definition save preserves the draft, shows current disk bytes and offers explicit reconciliation. Create collects declared inputs through the shared template service. Delete uses the same reference checks. Board keyboard choices and actual drag/drop issue one revision-checked field mutation; failure rebuilds the authoritative board. Prompt actions show previews including profile instructions, then draft an explicitly sent conversation; rendering never launches a harness.

## Repository assets

Custom views start with generic fallback and a trust explanation. Trust binds an unguessable frame grant and channel nonce to the selected project/type, record allowlist, definition identity and all recursively referenced/declared asset bytes. Disabled/invalid definitions cannot activate executable presentation. Assets are normalized `.factory/` paths read through an `os.Root` with regular-file and symlink checks, limits and explicit MIME types. Local script/style references are packaged into HTML; no arbitrary directory route is exposed.

Frames use `sandbox="allow-scripts"` and restrictive CSP with no connection/frame/worker capability. Parent messages require the exact source window, opaque origin, nonce, digest and an allowlisted schema. Navigation checks the grant selection. Request edits require host confirmation and pass through server grant/digest/selection validation and ordinary shared writes. Actions require a declared prompt and yield a preview for explicit host launch. Changing asset or definition bytes invalidates execution, resources and requests until fresh trust. Generated conversation frames and answer provenance keep their separate existing controls.

## Checks observed

`GOCACHE=/private/tmp/factoryctl-go-cache go test ./internal/web` passes. Added tests exercise CRUD with no runtime, missing nonce, stale record/definition, invalid enum byte preservation, external edits, target boundaries, explicit delete DTO, canonical cross-project links/backlinks, diagnostic-only refresh, sandbox/trust, message nonce/selection, stale custom edits, transitive JS digest changes, stale trust, external resources, symlinks and disabled custom fallback. Existing conversation/model/default/title/keyboard tests remain.

The reusable `scripts/review-video/completion-web-check.cjs` passes against a fresh initialized scratch project on the current Mac. It performs first-of-type create, typed and raw saves, actual drag and keyboard board moves, stale raw-save reconciliation with draft retention, delete, repository-JavaScript edit/action requests and asset digest invalidation, and asserts no browser JavaScript errors. Configure `FACTORYCTL_WEB_URL`, `FACTORYCTL_WEB_ROOT`, `PLAYWRIGHT_MODULE` and optional `REVIEW_CHROMIUM` browser executable. `REVIEW_OUTPUT_DIR` enables a new timestamped local video/results directory with actual-scene annotations, acceptance IDs, five-second reading pauses and assertion/command metadata; omit it for check-only execution. Results include actual scene start/end seconds, executed atomic requirement IDs alongside acceptance IDs, versioned artifact paths and command/environment metadata. The runner writes only its explicitly selected scratch fixture; start with a fresh fixture. Failed earlier scratch attempts were retained under `/private/tmp/factoryctl-web-completion/`; successful fixtures include `/private/tmp/factoryctl-web-completion2/` the bounded-pagination regression fixture `/private/tmp/factoryctl-web-completion4/`, and the final active-definition-conflict regression fixture `/private/tmp/factoryctl-web-completion7/`. This check is not the final complete-matrix annotated recording.

A separate live browser check verified the real planning bundle shell, tab title, draft retention across polling, backlinks, enum board and overview without writing project records. No release, external binding, unsupported-platform execution, external MCP Apps interoperability, full-matrix satisfaction or final media review is claimed here. Default listening remains loopback; an externally protected server deployment requires an externally configured proxy. Coordinator owns full qualification, delivery and final acceptance.

The optional recorder was executed successfully with Playwright 1.64.1, the explicitly selected cached Chromium 1243 and matching ffmpeg helper 1013 in `/private/tmp/factoryctl-web-pw/`. The retained component package `/private/tmp/factoryctl-web-completion-recordings/web-2026-10-28T04-24-22.790Z/` contains a successful recording and actual scene/result metadata. It predates the final stale-definition reconciliation fix and is historical component evidence; the Coordinator must capture the final candidate and stable complete film. The final check-only run on fixture 7 includes that fix and passes. Worker-owned loopback servers were stopped after checks. No recordings or generated artifacts were staged or uploaded.
