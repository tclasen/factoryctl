---
type: "Specification"
title: "Configuration, initialization and project discovery"
description: "Configuration, initialization and project discovery."
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

# Configuration, initialization and project discovery

This is the proposed configuration version 1. The [operator requirement](../sources/operator-brief.md) fixes the three storage locations and makes Git mandatory. The exact keys and defaults below are design proposals.

## Configuration shape

```yaml
version: 1
project:
  id: factoryctl
  name: Factoryctl
  bundle_uri: https://github.com/tclasen/factoryctl/blob/main/factory/
types:
  disabled: []
runtime:
  adapter: codex
  executable: codex
  profile: agents/builder.md
gui:
  listen: 127.0.0.1:0
  open_browser: true
limits:
  max_record_bytes: 1048576
  max_frontmatter_bytes: 262144
```

`version` is required. `project.id` is a stable slug or UUID, unique among unrelated projects, committed and shared across clones. Init derives a suggested slug from the directory and asks or accepts `--project-id`; collisions require an explicit selection/change, not a regenerated identity on every machine. `name` defaults to the directory name. `bundle_uri` is optional and identifies a canonical URL prefix for cross-project links; trailing slash required. It is not a network sync target or credential.

`runtime` is optional for browsing/editing. `adapter` defaults to `codex`, `executable` to PATH lookup `codex`; profile is a bundle-relative path. Model and effort are optional profile fields, not duplicated global defaults. Secrets are never stored in this file; Codex manages its own credentials and environment. All unspecified fields use documented embedded defaults. Unknown config keys are errors with locations and suggestions (unlike unknown OKF metadata, which is preserved).

Only this repo-root file supplies persistent factoryctl configuration. Precedence: explicit invocation flags override this file, then embedded defaults. Invocation-only flags include discovery root, selected project, listener, browser opening and asset trust. No home-directory config, implicit `.env` loading, external database, global registry or persisted machine path mapping. The bundle root and asset root are fixed to `factory/` and `.factory/` in v1; cannot redirect them into another repository.

## Discovery

Resolve CWD or `--path` to a canonical directory. If it is inside a Git working tree, select that working tree only, even when other repositories lie below it. If it is outside any working tree, recursively scan for Git working-tree roots. Recognize `.git` directories and Git worktree pointer files through read-only Git inspection; exclude bare repositories, `.git` internals, ignored directories, symlink/reparse-point directories and inaccessible paths with a summarized diagnostic. At each discovered root stop descending into it, avoiding accidental submodule/vendor enrollment. Users can explicitly point at a submodule/worktree when desired. Nested repositories are not implicitly scanned after their enclosing repo is selected; this is the proposed boundary for “recursive discovery.”

A repository containing `factory/` is an available factory even without config; assign only a temporary checkout identity until initialized and disable ambiguous cross-project writes. A Git repository without a bundle appears as **Not initialized** with an explicit init action. A non-Git folder containing `factory/` is not eligible as a project. Inspection never initializes a repository or pulls missing repositories from a remote.

Config parse errors keep the project visible in diagnostic/read-only mode; do not hide it or silently launch a default runtime. Discovery returns its partial successes and failures. Refresh rebuilds the in-memory project list; no registry file is written in the parent directory.

## Initialization

`factoryctl init --project-id <id> [--starter software]` requires a selected Git working tree. It writes `factory.yaml`, `factory/factory.md`, the [storage coordination node](../control/write-coordination.md), and selected default definitions/templates/prompts/profiles. The initialized set must be dependency-complete: materialize all local knowledge dependencies referenced by its selected defaults, or use defaults whose links remain valid without an optional starter set. Validate the planned output graph before writing. The design-repository profiles/templates are a dogfood instance, not a license to copy software-specific links into every base initialization. These are versioned human-readable files but factoryctl does not stage or commit them. No `.factory/` directory is necessary; built-in UI assets are embedded.

Init previews conflicting paths and never overwrites an existing file. Existing plain OKF content remains intact. An interrupted init reports which paths were created; rerunning fills only missing defaults and reports conflicts for differing files. No claim of a cross-file atomic transaction. `--non-interactive` requires all unresolved choices as flags and fails before writing on conflicts. Type creation/deactivation is an explicit operation, not a side effect of discovery.

`factoryctl templates export --path <repo> --kind view` materializes embedded view assets under `.factory/views/` on request. Record and prompt templates remain OKF nodes in the bundle. An optional `.factory/templates/` may hold raw template payloads referenced by those nodes; it is an asset source, not another graph or configuration store. The full resolution rules are in [templates](templates.md).

## Evolution

Upgrading the binary does not rewrite any local definition, prompt, record or config. `factoryctl templates diff` compares embedded defaults to local copies; applying an update is a user/agent edit. Version-1 readers display unknown future factory formats generically and reject structured writes that could corrupt them. Migration is an explicit agent-guided file change, tested against a preserved Git revision.
