---
type: "Specification"
title: "Project discovery and initialization"
description: "Project discovery and initialization."
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

# Project discovery and initialization

**Operator requirement:** projects do not require a configuration file. Git-versioned OKF Markdown is the default and only initial tracking backend. Optional `.factory/` contains versioned assets; `factory/` alone contains the knowledge graph. Do not read, generate or require `factory.yaml`. Reintroducing a project configuration file requires a demonstrated use and a separately reviewed design.

## Project metadata and defaults — proposed mechanism

The existing `factory/factory.md` Factory node holds its required `project_id` and human-readable `title`. Init suggests a stable slug from the directory and asks or accepts `--project-id`; clones share the committed identity. An optional `bundle_uri` on that node identifies the canonical absolute HTTPS URL prefix for cross-project links, with a trailing slash. It is knowledge identity metadata, not a network sync target or credential. Use the node title as the display name, falling back to the checkout directory for generic bundles.

Keep settings on the entities they describe: `enabled: false` on a Record Definition disables its specialized behavior; omission means enabled. Do not introduce a generic settings mapping in the Factory node. Profiles and prompts describe shared work intent, never choose a contributor's runtime, executable, model or reasoning level. [Personal runtime preferences](personal-runtime.md) remain outside the repository.

The bundle and asset roots are fixed to `factory/` and `.factory/`. GUI defaults to `127.0.0.1:0` and opens the browser; `--listen` and `--no-open` override this for the invocation only. Parser limits are documented built-in constants in the [knowledge format](knowledge-format.md). No project config version, disabled-type registry, stored machine path or runtime default is needed. Discovery does not read home preferences; only an explicit runtime action or `doctor` diagnostic resolves them.

## Discovery

Resolve CWD or `--path` to a canonical directory. If it is inside a Git working tree, select that working tree only. Otherwise recursively scan for Git working-tree roots. Recognize `.git` directories and Git worktree pointer files through read-only filesystem metadata inspection or in-process libraries, never a Git subprocess; exclude bare repositories, `.git` internals, ignored directories, symlink/reparse-point directories and inaccessible paths with summarized diagnostics. Stop descending at each discovered root; users can explicitly select submodules or nested worktrees. Never initialize or fetch a repository during inspection.

A Git repository containing `factory/` is available without any config file. A valid Factory node supplies portable identity. A generic OKF bundle without that node gets a process-local identity and remains readable and editable when the checkout is explicitly selected; do not require initialization just to browse it. Do not persist temporary identities or claim they identify the same project across machines. Missing or malformed identity metadata disables portable cross-project addressing, not generic access. Duplicate IDs or canonical bundle URLs require explicit checkout selection; never select a writable checkout by accident. A Git repository without a bundle appears as Not initialized. A non-Git folder is not an eligible project.

Malformed metadata remains visible with diagnostics; actions requiring that metadata fail explicitly. Discovery returns partial successes and failures. Refresh rebuilds the in-memory project list without writing a registry.

## Initialization

`factoryctl init --project-id <id> [--starter software]` requires a selected Git working tree. It writes `factory/factory.md`, the [storage coordination node](../control/write-coordination.md), and selected default definitions/templates/prompts/profiles. It creates no root configuration, user preference file or runtime-specific file. The initial set must be dependency-complete; validate the proposed graph before writing. Software starter links must not leak into base defaults. All output is human-readable Markdown; factoryctl never stages or commits it. No `.factory/` directory is necessary because default UI assets are embedded.

Preview conflicts and never overwrite existing files. Interrupted init reports created paths; rerunning fills only absent defaults and reports differing-file conflicts. No cross-file atomicity is claimed. `--non-interactive` fails before writing if choices or conflicts remain. Type creation/deactivation is explicit, never a side effect of discovery. To disable an embedded type, materialize its definition and set `enabled: false`; removing that local definition restores its embedded default.

`factoryctl templates export --kind view` writes requested embedded view assets under `.factory/views/`. Optional `.factory/templates/` contains raw payloads referenced by OKF nodes, not another graph or configuration store. See [templates](templates.md).

## Evolution and deferred tracking backends

Binary upgrades never rewrite local records, definitions or prompts. `templates diff` compares embedded defaults with local copies; edits remain explicit. Unknown future factory formats remain generically readable; structured writes that could corrupt them fail. No compatibility or migration guarantee applies before 1.0.0; see [release policy](../decisions/release-policy.md).

The operator identified GitHub Issues/Projects, Linear, Trello and Jira as possible future tracking backends. Backend routing is a plausible future use for project configuration, not a supported feature or reason to create a placeholder file now. Before implementation, decide authoritative storage versus mirroring, record/remote ID mapping, link/backlink semantics, offline behavior, conflict ownership and credential boundaries. No sync engine, remote mutations or backend credential storage is admitted by this design. Git remains the only initial source of project tracking state.
