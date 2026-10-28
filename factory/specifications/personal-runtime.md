---
type: "Specification"
title: "Personal runtime preferences"
description: "Personal runtime preferences."
status: "draft"
factory:
  format: 1
  assertion:
    frame: "desired"
    certainty: "proposal"
    basis: "llm-inference"
sources:
  -
    id: "operator"
    resource: "../sources/operator-brief.md"
    title: "Operator requirements and scope"
    author: "human:operator"
---

# Personal runtime preferences

**Operator requirement:** runtime selection belongs to each contributor and is not project Git state. A teammate may prefer another harness without changing repository files. Codex remains the first supported adapter implementation, not a mandatory project choice. Claude and other adapters are future work; choosing an unimplemented adapter must report `unsupported_capability`, never silently substitute Codex.

## Local mechanism and remaining proposals

All commands work without a preference file except conversation launch requires a selected supported adapter. Resolve explicit conversation selection before invocation flags, saved project defaults, saved global defaults, legacy personal preferences, then native runtime defaults for model/effort. No repository file selects the adapter, executable, model or effort. Shared profile instructions may be sent to any supported adapter; unknown runtime-specific metadata on imported profiles is preserved as data but never interpreted as launch settings.

An optional `~/.factoryctl/preferences.md` is human-readable Markdown with YAML frontmatter. Resolve the home directory via the operating system's current-user API. `--preferences FILE` explicitly selects another user-owned file, useful for a service/container; never search the checkout or load `.env`. The factoryctl process user owns these preferences, including when a server is used. Read them only for an explicit runtime action or an explicit `doctor` diagnostic. Do not create or update them automatically, copy them into Git, expose their raw contents to templates/MCP, or treat them as an OKF bundle. No secrets, credentials, sessions or project status belong in this file.

```yaml
---
version: 1
runtime: codex
# executable: /absolute/path/to/codex
# model: contributor-selected-model
# reasoning_effort: contributor-selected-level
---
Personal runtime preferences; not shared project state.
```

When a file is present, require version 1; `runtime`, `executable`, `model` and `reasoning_effort` are optional nonempty strings. Reject unknown keys, malformed YAML, duplicate keys, custom tags and alias expansion with field/path diagnostics. Use the same YAML resource limits as records. It is a narrow preference schema, not an arbitrary command runner or general configuration store.

Implemented invocation flags for `web` and `chat`: `--runtime ADAPTER`, `--runtime-executable PATH`, `--model NAME`, `--reasoning-effort VALUE`, `--preferences FILE`. The web host supplies per-conversation model/reasoning selectors and explicit global/project saves as described below. `doctor` reports the opt-in adapter capability without reading preferences or starting a conversation. Runtime selection through `actions run`, legacy `gui` or a richer diagnostic remains proposed. Reject an explicit missing file, invalid preferences or unsupported selection before launch. Invalid preferences never block init, record editing, validation of project records or GUI browsing. If runtime is unspecified, `web` remains a document explorer; start it with `--runtime codex` or explicit `--preferences FILE` to enable conversations. Headless `chat` fails explicitly. No automatic executable discovery chooses a runtime on the user's behalf.

If the invocation adapter differs from the file's adapter, ignore the file's adapter-bound executable/model/effort; use invocation values and the selected adapter's native defaults. File values for executable/model/effort require a `runtime` in that file. This prevents passing Codex options to a future Claude adapter. An explicit executable path is launch data from the user, never from a repository or rendered prompt. Resolve bare defaults on the process PATH only after selecting an implemented adapter. Native adapter capability checks qualify model/effort values.

## Separation and verification

Profiles/prompts retain shared roles, instructions and record references. GUI and CLI resolve contributor choices; adapters consume an in-memory selection snapshot. Credentials and transcript persistence stay with the harness. The user preference file is a narrow exception to the no-external-project-state rule, not permission to persist factory status outside Git.

The [runtime contract](agent-runtime.md) specifies Codex's first adapter implementation. [Acceptance A08](acceptance.md#a08-native-codex-conversation) must exercise missing/invalid preferences, flag precedence, different contributors, unsupported adapters and unchanged repository bytes. Keep runtime selection/capability errors distinct from invalid knowledge. The bounded implementation reads preferences without creating a local preference file; see [capability evidence](../sources/interaction-capabilities.md).


## Requested model selectors and personal project defaults

The operator subsequently requested selectors for the current chat and future planned chat capabilities, with home-directory configuration at global and per-project scope. This explicitly supersedes the earlier launch-only model-selection design. It authorizes FactoryCTL to write its own personal chat configuration after an explicit Save/Reset action, while preserving the read-only legacy `preferences.md` and never modifying Codex's configuration.

The reusable store is `~/.factoryctl/chat-settings.json` (version 1), resolved from the operating-system user. An explicit `--chat-settings` path supports isolated service/test environments. Only model/reasoning pairs belong here; no credentials, conversations, runtime executable, project status or workflow policies. The structure is a nullable global pair and a map of canonical absolute project-root paths to pairs. Paths distinguish independent monorepo projects and duplicate project IDs. Moving a project or using another checkout does not silently adopt another path's override.

A pair contains both string fields `model` and `reasoning_effort`. Empty model selects the native runtime configuration; empty effort selects the configured native effort when model is empty, or the advertised model default when model is explicit. Missing project override inherits global. Global null means no global override. Project reset removes the entry; null project entries and missing/null pair fields are malformed, not implicit native defaults.

Resolution is explicit conversation pair → invocation flags → saved project pair → saved global pair → legacy preferences → native defaults. A model invocation flag replaces the lower pair and clears lower effort unless supplied explicitly; effort-only retains lower model. Changing model in the UI resets its reasoning selection to avoid carrying an incompatible level. Catalogs come from the native runtime; no model list or capability is invented. The native configured model can differ from the catalog default. Before every turn, resolve and send concrete model and effort, including when returning to native defaults in an existing thread.

Selecting or sending does not persist. Explicit Save/Reset uses authenticated host controls, validates advertised choices and revision, and writes only this personal store. New chats resolve current saved defaults; existing chats keep their own pair until changed at an idle boundary. Malformed storage remains visible and cannot be overwritten through Save, while document browsing and explicit current-chat choices remain usable. Unknown/unavailable models stay visible and require an explicit supported selection rather than a silent substitute.

Private, bounded, strict JSON storage rejects unknown/duplicate fields, malformed pairs and symlinks at the target file or its immediate settings directory (ancestor OS aliases remain permitted). Cooperative cross-process locking protects revision comparison and atomic replacement; stale clients receive a conflict and must reload deliberately. The absent-file revision is `absent`. No repository document chooses a configuration filename or executable.

All future chat capabilities must reuse the shared resolver, store, catalog validation and explicit-save semantics. This is an implementation contract for future surfaces, not a claim that those interfaces already exist. The interactive CLI consumes the same saved defaults and retains invocation overrides. See [the selector task](../tasks/chat-model-settings.md) and acceptance A22.
