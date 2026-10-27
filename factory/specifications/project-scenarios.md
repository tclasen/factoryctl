---
type: Specification
title: "Required project collection scenarios"
description: "Required project collection scenarios."
status: draft
factory:
  format: 1
  assertion:
    frame: desired
    certainty: mixed
    basis: mixed
sources:
  - id: operator
    resource: ../sources/operator-brief.md
    title: Operator requirements and scope
    author: human:operator
---

# Required project collection scenarios

These are operator-required regression scenarios, exercised as each implementation slice becomes available and again at release. They are planned fixtures, not passing tests. Create temporary Git repositories through test setup; factoryctl itself must never spawn Git.

| Scenario | Fixture and launch location | Expected collection |
| --- | --- | --- |
| S1 Single project | Launch at `solo/` with `.git/` and `factory/` | Exactly one project; the ordinary collection overview and All projects selector remain visible |
| S2 Monorepo | Launch at `mono/` with `.git/`, `alpha/factory/` and `beta/factory/`; neither child has `.git` | Exactly alpha and beta, each with distinct identity, definitions, prompts and state, sharing one Git root |
| S3 Sibling repositories | Launch at non-Git `workspace/` containing `alpha/.git`, `alpha/factory/`, `beta/.git`, `beta/factory/` | Exactly alpha and beta, with the same records, interactions and collection UI as S2 |

Use identical alpha/beta project content in S2 and S3 except canonical bundle URLs. Give both a `tasks/example.md` with different text to expose mistaken addressing. Give their definitions/customizations different values to detect accidental sharing. Link alpha's record to beta's record and provide the reverse backlink; use each Factory node's distinct canonical `bundle_uri`. URLs may share the monorepo remote but have distinct bundle paths. No remote fetching is needed for loaded targets.

Required journeys: enumerate projects; select one or All projects; list and read; navigate a cross-project link and backlink; edit the intended target and refresh; reject stale revisions and path escapes; show a missing/ambiguous target honestly. Verify edits do not alter the other project's same-named file or definitions. CLI, HTTP and MCP must resolve the same project identities. Repeat GUI navigation, editing, conversation and agent-question flows with the same controls in all layouts. S1 has explicit empty cross-project-link state, not a special single-project UI.

Boundary fixtures: a root bundle plus two children yields three projects; root without bundle yields only its children; deeper `group/hidden/factory/` is excluded until `group/` is explicitly selected; discovering a child Git repository does not restart the depth counter. Include CWD inside project source, explicit subdirectory selection, worktree pointer files, duplicate IDs/URLs, ignored/inaccessible paths, symlinks and a collection with no projects. Absence of a child `.git` must not exclude S2. Generic bundles remain browseable by explicit selection. No `.factory/` is required.

Stage ownership: task 01 establishes discovery fixtures; tasks 02–04 verify project identity, independent definitions/templates, links and writes; tasks 05–06 verify transport parity; tasks 07–08 verify conversation routing and common UI; task 11 reruns the complete matrix. See A16 in [acceptance](acceptance.md). Record tested revisions, expected/observed outcomes and limits per scenario; do not postpone basic layout tests to release.
