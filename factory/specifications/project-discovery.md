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

The bundle and asset roots are fixed to `factory/` and `.factory/` relative to each project root, not necessarily the repository root. GUI defaults to `127.0.0.1:0` and opens the browser; `--listen` and `--no-open` override this for the invocation only. Parser limits are documented built-in constants in the [knowledge format](knowledge-format.md). No project config version, disabled-type registry, stored machine path or runtime default is needed. Discovery does not read home preferences; only an explicit runtime action or `doctor` diagnostic resolves them.

## Discovery

A **project root** owns one `factory/` bundle and optional `.factory/` assets. A **repository root** owns Git history and may contain several project roots. Projects in a monorepo have distinct Factory identities, definitions, prompts, records, custom assets and internal state; sharing Git does not merge these. No implicit parent/sibling definition or asset inheritance occurs.

Resolve CWD or `--path` to a canonical directory. Use the selected directory as the discovery root (depth 0), inspect it and its immediate subfolders (depth 1), and never enumerate grandchildren by default. Detect the enclosing Git working tree through ancestor metadata inspection when needed; this establishes Git membership, not permission to broaden the downward scan. When CWD is inside an existing project's source subtree and has no local project candidates, resolve its nearest ancestor project within the enclosing working tree as the single selected project. An explicit `--path` selects that directory's scope without automatic ancestor-project fallback.

A project candidate has its own `factory/` bundle and lies within a Git working tree; its own `.git` marker is not required. Inspect a candidate's fixed `factory/` bundle and `.git` marker to identify it; subsequently reading its bundle records is not recursive project discovery. Recognize Git directories and worktree pointer files via filesystem metadata/in-process libraries, never Git subprocesses. Exclude bare repositories, Git internals, ignored directories, symlink/reparse-point directories and inaccessible paths with diagnostics. Do not descend below depth 1 even after finding a repository or project. Never initialize or fetch during inspection.

At a monorepo root, load the root project if present and immediate-child projects. A root with no bundle is a collection container, not an extra project, when initialized child projects exist. At a non-Git parent, immediate-child repository bundles become projects. Do not automatically expand a discovered child repository into its own child projects: that would exceed the original depth limit. A grouped layout such as `mono/packages/a/factory/` is selected with `--path mono/packages`; a deeper individual project can be selected directly. This preserves the default depth-1 limit in every topology.

A valid Factory node supplies portable identity. A generic OKF bundle without it remains readable and editable when its project root is explicitly selected, using a process-local identity; no initialization is required just to browse it. Missing/malformed metadata disables portable cross-project addressing, not generic access. Duplicate IDs or canonical bundle URLs require explicit project-instance selection; never pick a writable location by accident, even when duplicates share a repository. A repository with no discovered bundle is Not initialized. Arbitrary non-Git folders are collection containers only, never eligible projects. An explicitly selected subdirectory inside a Git working tree can be initialized as a project.

Malformed metadata remains visible with diagnostics. Discovery returns partial successes/failures, deduplicates canonical project roots, and rebuilds the in-memory collection on refresh without a registry. All interfaces use the same project collection; one project is a one-element collection. See the [three-scenario test matrix](project-scenarios.md).

## Initialization

`factoryctl init --project-id <id> [--starter software]` requires one explicitly selected project root within a Git working tree, including a monorepo subdirectory. Ambiguous collection selection fails before writing. It writes `factory/factory.md`, the [storage coordination node](../control/write-coordination.md), and selected default definitions/templates/prompts/profiles. It creates no root configuration, user preference file or runtime-specific file. The initial set must be dependency-complete; validate the proposed graph before writing. Software starter links must not leak into base defaults. All output is human-readable Markdown; factoryctl never stages or commits it. No `.factory/` directory is necessary because default UI assets are embedded.

Preview conflicts and never overwrite existing files. Interrupted init reports created paths; rerunning fills only absent defaults and reports differing-file conflicts. No cross-file atomicity is claimed. `--non-interactive` fails before writing if choices or conflicts remain. Type creation/deactivation is explicit, never a side effect of discovery. To disable an embedded type, materialize its definition and set `enabled: false`; removing that local definition restores its embedded default.

`factoryctl templates export --kind view` writes requested embedded view assets under `.factory/views/`. Optional `.factory/templates/` contains raw payloads referenced by OKF nodes, not another graph or configuration store. See [templates](templates.md).

## Evolution and deferred tracking backends

Binary upgrades never rewrite local records, definitions or prompts. `templates diff` compares embedded defaults with local copies; edits remain explicit. Unknown future factory formats remain generically readable; structured writes that could corrupt them fail. No compatibility or migration guarantee applies before 1.0.0; see [release policy](../decisions/release-policy.md).

The operator identified GitHub Issues/Projects, Linear, Trello and Jira as possible future tracking backends. Backend routing is a plausible future use for project configuration, not a supported feature or reason to create a placeholder file now. Before implementation, decide authoritative storage versus mirroring, record/remote ID mapping, link/backlink semantics, offline behavior, conflict ownership and credential boundaries. No sync engine, remote mutations or backend credential storage is admitted by this design. Git remains the only initial source of project tracking state.
