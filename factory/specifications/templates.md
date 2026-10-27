---
type: "Specification"
title: "Record, prompt and presentation templates"
description: "Record, prompt and presentation templates."
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

# Record, prompt and presentation templates

Templates are the principal extension mechanism. The [operator](../sources/operator-brief.md) requires new record types, composable prompts and arbitrary HTML/CSS/JS. This document proposes a concrete version-1 contract.

## Record templates

An OKF node of `type: Record Template` has `for_type` and an `inputs` mapping with the same simple field kinds as [definitions](knowledge-format.md). Its body contains one `# Template` section with exactly one fenced `factory-template` block. Other sections explain usage. The fenced payload is an unevaluated text template which must render a complete OKF document. A node may instead declare `asset: .factory/templates/name.md.tmpl` (project-root-relative); exactly one of an inline payload or asset is allowed. Raw assets do not join the graph.

Use Go `text/template` syntax, restricted to field lookup, `if`, `range`, literals and named partials. Functions: `yaml` (encode a scalar/list/mapping as valid YAML), `json` (encode JSON for output text only), `markdownLink`, and `include`. No shell, file read, network, environment expansion or arbitrary expression execution. `missingkey=error`; required inputs are gathered interactively or provided explicitly in CI/MCP. Unknown input names are rejected. Values must be serialized, not interpolated unescaped into YAML.

`include "prompts/shared-context.md" .` resolves a bundle node with an appropriate template payload or an explicitly named asset within `.factory/`. Resolve paths relative to the bundle root for node includes and selected project root for `.factory/` asset includes. Reject cycles, type mismatches, path escapes and depth >16 before execution. Read content is literal data and is not recursively interpreted as template source. Render into memory, parse/validate the complete output and only then create a file. Limit combined output to the built-in record limit.

Context is a versioned typed object: `.Project` (`ID`, `Name`), `.Record` (optional `Path`, `Title`, `Type`, `Revision`, `Frontmatter`, `Body`), `.Input` (declared user values), `.Selection` (record paths) and `.Links` (resolved references/backlinks). No implicit access to other project contents. Types without an explicit template use an embedded generic OKF template: collect title, summary and all required declared fields, serialize their values, and validate the result. This permits creating agent profiles and custom types without requiring a bespoke template. Preview shows the exact rendered document and destination. Definition defaults are applied before rendering; resulting frontmatter must satisfy the selected definition. No template can choose a destination outside its selected project.

## Prompt nodes and actions

A `Prompt` node supplies `applies_to` (type strings or `*`), `inputs`, optional `profile` (an OKF reference to an Agent Profile, relative to the prompt document unless bundle-absolute), `conversation: fresh`, and a `# Template` fenced `factory-prompt` payload or explicit asset. `label` is the GUI button label; it defaults to the node title. Prompt includes use the same engine/limits, with a maximum rendered prompt of 128 KiB. This limit is a factoryctl proposal, not a claim about a model's context capacity.

A prompt without a profile uses no shared role overlay unless the user selects one in the GUI; native runtime defaults remain available. A selected profile supplies instructions only, never runtime/model settings. A missing personal runtime selection requires contributor input before any one-click action can launch.

The user sees the selected project, record, profile and rendered prompt before sending when the action has required inputs or explicitly requests preview. A one-click action with no extra inputs sends its documented prompt immediately to a fresh session. A follow-up in that conversation does not rerender the original action or create another session. The [runtime contract](agent-runtime.md) supplies context and tools.

Prompt metadata does not grant permissions. A profile specifies shared work intent; runtime/model/effort settings come only from personal preferences or invocation options, and the harness retains its native sandbox and approval policy. Merely rendering a prompt, viewing a file or listing actions never starts an agent or shell command. Disabled types and unresolved prompt/profile references show a reason instead of a broken button.

## Views and asset ownership

Embedded generic list/detail/board renderers require no `.factory/`. A definition may set `views.detail_template: .factory/views/task.html` and `views.assets` to explicit project-root-relative CSS/JS paths. Those references are optional. An absent explicit custom asset is a diagnostic with a safe generic-view fallback; it is never fetched from the network. Runtime does not compile TypeScript, install packages or run an asset build. Custom assets must be directly browser-usable HTML/CSS/JS.

Server interpolation uses Go `html/template` contextual escaping with a JSON-serializable view model, not the command-oriented prompt function set. User source HTML/CSS/JS is executable code. Markdown record bodies are sanitized before generic rendering; they do not become trusted scripts because the repository has a custom view.

Custom views run in an iframe with `sandbox="allow-scripts"`, no same-origin privilege, and a restrictive CSP. Local asset references are resolved and packaged into the iframe without giving it directory browsing; no external resources by default. The parent passes a record snapshot, or an explicitly scoped board selection with each record revision, and a narrowly defined message bridge: `ready`, `resize`, `navigateRecord`, `requestAction`, and `requestEdit`. Each message includes a fresh channel nonce; parent checks source window, nonce, project scope and exact message schema. An opaque iframe origin alone is not identity. Parent-mediated mutations use the same revision checks and explicit interaction as native views. No direct shell/runtime APIs are exposed to the frame.

Trust is explicit for each process invocation/project; opening the generic GUI is safe by default. `--trust-assets` or a clear GUI action enables the selected project's executable assets for that process. File changes invalidate the cached asset digest and require trust again before execution. No browser localStorage/database or durable external trust file. This is browser-code isolation, not a claim that arbitrary JS is harmless; it can read the record data intentionally passed to it. Stronger capability sharing and configurable outbound networking are deferred.

## Compatibility

Local record/prompt nodes and exported assets are authoritative when explicitly selected. Defaults fill only absent optional customization. Nested template expansion is deterministic and never invokes the LLM. Source edits appear after refresh without restarting the binary. Template validation errors identify the node/asset, line, include chain and missing input. [Acceptance](acceptance.md) includes escaped values, cycles, unknown inputs, missing assets, hostile paths and zero-assets startup.

## User-defined interactive boards

A custom HTML/CSS/JavaScript view may implement columns, filtering and ticket drag/drop. A drop requests a revision-checked state edit through `requestEdit`; the parent applies shared validation and returns the resulting state or diagnostic. A stale/invalid move restores the observed state. `requestAction` selects a declared prompt for the harness, never an arbitrary shell command. These assets run only in the browser; factoryctl does not execute their JavaScript or load custom code on the backend. See [runtime types](extensions.md) and A14 for the complete acceptance journey.
