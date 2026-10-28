---
type: Task
title: "Choose chat model and reasoning with personal project defaults"
description: "Reusable runtime selection and home-directory global/project defaults."
status: draft
task_state: done
priority: high
owner: coordinator
intent_state: requested
intent_owner: human:operator
---

# Choose chat model and reasoning with personal project defaults

User request: “Please give me a selector for the model and reasoning level for both this chat interface and future planned chat capabilities (configuration stored in home directory in a config file on a global and per-project basis).” Context is FactoryCTL's browser chat. The assistant interprets future capabilities as requiring a shared settings contract, not implementing all future chat surfaces or altering the hosting Codex application's UI.

Authorized slice: model/reasoning selectors for new and continuing FactoryCTL chats, discovered from the selected native runtime; explicit save/reset of global and per-project defaults in the operating-system user's home; reusable resolution for CLI/future chat clients; targeted tests, local tangible review and normal PR/main delivery. No harness configuration changes, invented model list, agent-callable settings writer, artifact upload or new provider. Preserve existing running conversations when refreshing the app.

Assistant-selected design: `~/.factoryctl/chat-settings.json`, version 1, global selection plus project selections keyed by canonical project root (independent monorepo roots and duplicate project IDs remain isolated). Each selection is a pair of model/reasoning strings; empty means native default. Missing project override inherits global. Precedence: explicit conversation choice, explicit invocation flags, saved project pair, saved global pair, legacy personal preferences, native defaults. Save is a separate explicit host action; merely choosing or sending does not persist. Saved defaults apply to new chats, while existing chats retain their last selection until changed at an idle turn boundary. Catalog and settings access require the existing host credential; no executable/path setting is accepted through these routes. Validate model/effort against advertised runtime capability, bound catalogs and errors, keep filesystem writes atomic/private with conflict protection and reject symlink targets. Runtime selection stays local; no project Git configuration.

Claims/runs: Coordinator owns web Go/API integration, CLI resolution integration, project docs, tests/evidence and Git. Runtime implementer owns `internal/runtime/chat_settings.go`, its tests, and reusable model catalog support in `internal/runtime/`; no web/assets/docs/Git edits. UI implementer owns `internal/web/assets/` plus a dedicated selector browser fixture; no web Go/runtime/docs/Git edits. Independent reviewer is read-only, first checking authority/outcome and design readiness, then exact-candidate correctness and evidence. All runs inherit model/reasoning; actual model/cost telemetry unknown. Contract changes coordinate before cross-package edits. Tests bounded to ten minutes; retain failures and checkpoints, no conflicting edits. Root task record owns this follow-up allocation.

Acceptance: actual runtime model catalog, supported reasoning choices, initial/follow-up native request selection, separate global/project saves, reload/restart persistence, reset/inheritance, unrelated-project isolation, conflict handling, invalid selection and malformed/read-only storage errors, keyboard submit retained. Tests use an isolated explicit settings path, never silently change the user's real defaults. Existing personal preferences are preserved. Future surfaces must reuse this store/resolver and capability validation; no claim they already exist. Rollback is a scoped revert preserving user configuration and local evidence.

## Clarified design review contract

Absent selection inherits; an explicitly present pair containing two empty strings deliberately selects native runtime defaults. Both fields must be present and strings in API/stored pair objects. Global null means no global override; a missing project entry inherits, and resetting removes that entry. A model invocation flag replaces the lower pair and clears inherited effort unless the effort flag is also supplied; an effort-only flag retains the lower model and overrides only effort.

Native `config/read` provides the effective configured model/effort when the desired model is empty; this may differ from `model/list`'s default. Explicit model with empty effort uses that model's advertised default effort. Every turn sends concrete validated model and effort, including a reset, instead of assuming omitted fields reset an existing native thread. Live switch-and-return qualification is required.

Revision comparison and atomic private-file replacement occur under a cooperative cross-process advisory lock. The absent-file revision is `absent`; stale clients receive conflict without overwriting unrelated project defaults. Immediate settings-directory and file symlink rejection applies; ancestor OS aliases remain permitted. Current chats carry checkout identity separately from the human-readable project ID so duplicate IDs cannot route settings to another checkout. Existing browser regression fixture support was added to UI ownership; implementation ownership is now released to Coordinator for integration fixes.

## Verification and delivery condition

[Targeted evidence](../sources/chat-model-settings-evidence.md) records passing shared-store/API tests, browser regressions, actual CLI defaults and real browser same-thread switching/reset plus restart persistence. Independent review and normal PR/main delivery remain required. The proposed done state takes effect only after exact-candidate approval, PR merge and verified owning-remote main containment; the session handoff records that observation. No release, new milestone or full-matrix qualification is claimed. Existing user conversations are preserved when activating the updated app.
