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

## Proposed local mechanism

All commands work without a preference file except conversation launch requires a selected supported adapter. Resolve explicit invocation values before optional user preferences, then native runtime defaults for model/effort. No repository file selects the adapter, executable, model or effort. Shared profile instructions may be sent to any supported adapter; unknown runtime-specific metadata on imported profiles is preserved as data but never interpreted as launch settings.

An optional `~/.factoryctl/preferences.md` is human-readable Markdown with YAML frontmatter. Resolve the home directory via the operating system's current-user API. `--preferences FILE` explicitly selects another user-owned file, useful for a service/container; never search the checkout or load `.env`. The factoryctl process user owns these preferences, including when a server is used. Read them only for an explicit runtime action. Do not create or update them automatically, copy them into Git, expose their raw contents to templates/MCP, or treat them as an OKF bundle. No secrets, credentials, sessions or project status belong in this file.

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

Invocation flags for `actions run` and `gui`: `--runtime ADAPTER`, `--runtime-executable PATH`, `--model NAME`, `--reasoning-effort VALUE`, `--preferences FILE`. GUI can select/change these in memory for that process without persisting. Reject an explicit missing file, invalid preferences or unsupported selection before launch. Invalid preferences never block init, record editing, validation of project records or GUI browsing. If runtime is unspecified, GUI asks for a supported adapter on the first conversation; noninteractive launch fails with guidance. No automatic executable discovery chooses a runtime on the user's behalf.

If the invocation adapter differs from the file's adapter, ignore the file's adapter-bound executable/model/effort; use invocation values and the selected adapter's native defaults. File values for executable/model/effort require a `runtime` in that file. This prevents passing Codex options to a future Claude adapter. An explicit executable path is launch data from the user, never from a repository or rendered prompt. Resolve bare defaults on the process PATH only after selecting an implemented adapter. Native adapter capability checks qualify model/effort values.

## Separation and verification

Profiles/prompts retain shared roles, instructions and record references. GUI and CLI resolve contributor choices; adapters consume an in-memory selection snapshot. Credentials and transcript persistence stay with the harness. The user preference file is a narrow exception to the no-external-project-state rule, not permission to persist factory status outside Git.

The [runtime contract](agent-runtime.md) specifies Codex's first adapter implementation. [Acceptance A08](acceptance.md#a08-native-codex-conversation) must exercise missing/invalid preferences, flag precedence, different contributors, unsupported adapters and unchanged repository bytes. Keep runtime selection/capability errors distinct from invalid knowledge. No local preference file or application implementation is created by this design.
