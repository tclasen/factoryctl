---
type: Reference
title: "Personal chat model selection: targeted evidence"
description: "Native selection, personal settings and browser regression evidence for the requested follow-up."
status: draft
---

# Personal chat model selection: targeted evidence

This is the requested [model-selector follow-up](../tasks/chat-model-settings.md), not a new milestone or full-product qualification. U103 covers current browser selectors and the shared contract for future clients; U104 covers personal global/project persistence. Future chat surfaces are not implemented by this change. The earlier milestone video remains historical for its pinned source and is not evidence for these new controls.

Executed environment: macOS arm64, Go 1.27.3, Node 26.12.1, installed Codex 0.164.0. The tested CLI binary (`CGO_ENABLED=0`) SHA-256 is `819ac3d82d00dc6c04f7d9528ce2fe85e1516892eaac4b268f81447c474b313a`; the browser binary (`CGO_ENABLED=1`) is `da4fde71111849e1d860667645be5f3e61c6d0ba6fa662852ef15eb76442a0bc`. A fresh build from current source is byte-identical to the browser binary; all files match the CLI source inventory. Source inventories and sanitized native metadata remain in ignored `artifacts/chat-model-settings/20261028/`. Test configuration is isolated from the user's actual home defaults. No videos or generated artifacts are committed or uploaded.

## Executed checks

- `go test -race ./cmd/... ./internal/...` and `go vet ./cmd/... ./internal/...` pass. Scoped packages exclude ignored local probe sources. Settings tests exercise exact pair decoding, inheritance and explicit native defaults, precedence, private atomic persistence, reset, canonical project isolation, cross-process revision conflicts, cancellation, immediate directory/file/lock symlink refusal and read-only errors. Catalog tests cover pagination, malformed responses and native-default resolution.
- Host API tests pass for authentication, loaded-checkout identity, incompatible model/effort rejection, stale saves, omitted-selection refusal, restart persistence and malformed-file behavior without blocking document browsing or explicit current-chat choices.
- The deterministic selector browser fixture passes 13 checks, including duplicate human project IDs with independent checkout identities, explicit saves, busy-state controls, follow-up selection, polling/focus, conflicts and unavailable stored choices. These are mock API checks, not native execution claims.
- The existing browser interaction fixture passes 21 regression checks, including Enter/Shift+Enter, IME/repeat guards, question/view trust and host isolation.
- A real PTY chat reads isolated saved project defaults, overriding a different global effort. Native `turn/start` receives `gpt-6-astra` / `low`; the agent returns `CLI_SETTINGS_OK`, the turn completes and the PTY exits 0. Settings/source/binary stay unchanged. Result: `artifacts/chat-model-settings/20261028/runtime-cli/result.json`, SHA-256 `f80a934b3ef721b26984227ffe94a5bc8849ef8b2847a73bb402f72c21b13077`.
- Six builds succeed for macOS/Linux/Windows on amd64/arm64. These are compile checks; only macOS arm64 has execution evidence. Commands and hashes are in `artifacts/chat-model-settings/20261028/builds/results.json`.
- Repository-content checks pass. Knowledge validation retains exactly the pre-existing `factory/.DS_Store` invalid-location diagnostic; the user file is preserved.

## Browser native qualification

The real browser run passes three successive turns on native thread `01a18093-61d1-7032-bbe8-9ffe85b0d4c2`: Astra/low, Sol/low, then explicit empty selection resolving to native Astra/medium. Distinct assistant replies, desired/resolved selections, actual native `turn/start` model/effort and completion notifications agree. The same recorded journey saves global and project defaults, verifies another project inherits global, restarts its own server and observes persistence, then removes the project override and verifies inheritance in a new chat.

Local targeted video: `artifacts/chat-model-settings/20261028/browser-live/attempt-4/model-settings-demo.mp4`, SHA-256 `9db1388cf8743e0de698486f0a2672e35cca9dbdc4642cb05c4522f4566af398`. It preserves the complete recording without cuts or acceleration.

Results and scene timestamps are in `artifacts/chat-model-settings/20261028/browser-live/attempt-4/result.json`. U103 scenes occupy approximately 00:01–00:52; U104 scenes occupy 00:53–01:27. Captions explain the observed outcome alongside the real controls, with at least 6.5-second holds. This is a targeted settings demonstration, not a replacement for the historical full-matrix milestone video. Visual limitation: the later assistant reply cards are below the sticky composer in completion scenes; the selections and captions are readable, while those exact replies are established by API assertions/native trace rather than fully readable response-card pixels.

Attempts 1–2 failed during recorder setup before native turns (synthetic Git initialization order and missing browser video dependency). Attempt 3 completed all native turns, then failed a recorder assertion because it matched stale Save feedback before the second save finished. All are preserved; only the corrected complete attempt 4 passes the browser journey. No native capacity failure or product repair was needed.

## Delivery boundary

Independent review checks outcome/authority separately from implementation correctness, and final approval binds the exact candidate and evidence. Completion takes effect only after ordinary PR merge and verified owning-remote main containment. This follow-up does not close the earlier milestone's remaining requirements or establish full 104-row end-to-end qualification.
