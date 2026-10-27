---
type: "Reference"
title: "Design validation and review evidence"
description: "Design validation and review evidence."
status: "draft"
factory:
  format: 1
  assertion:
    frame: "observed"
    certainty: "supported"
    basis: "direct-observation"
sources:
  -
    id: "operator"
    resource: "operator-brief.md"
    title: "Operator requirements and scope"
    author: "human:operator"
---

# Design validation and review evidence

## Scope and current verification

This record concerns design documents, not implemented application behavior. The current revision removes required root configuration, retains Git as the initial tracking backend and separates contributor runtime preferences from shared project knowledge. The author owns the documentation changes; the independent reviewer has a read-only claim on the pinned revision supplied in dispatch.

The initial design was delivered through [PR 1](https://github.com/tclasen/factoryctl/pull/1), and the Conventional Commits/0.x release policy through [PR 2](https://github.com/tclasen/factoryctl/pull/2). Fetched remote main contained the reviewed candidates with identical trees after each merge. Those deliveries do not qualify subsequent design revisions.

## Verification method and limitations

Validate a clean archive of the tracked candidate: YAML/frontmatter, definition field kinds and locations, links/anchors, template/profile references, task dependency cycles and R01–R25 coverage. The temporary checker now asserts absence of factory.yaml and checks the existing Factory identity and optional canonical URL, rather than parsing a required config file. Untracked Finder metadata in the working checkout is preserved and excluded from the tracked-corpus check.

The checker uses a bounded Markdown-link scan with fenced examples excluded, not the planned production Markdown AST parser. External URLs remain unchecked. No production validation, runtime conversation, personal-preference resolution, backend synchronization, performance or non-Mac execution tests have run. The new preference path, precedence and metadata rules are design proposals, not implemented features.

## Review and delivery

A fresh independent review must cover the changed contracts, identity/link handling, type disabling, personal preference scope/precedence, contributor separation, init without config and deferred backends before merge. Reviewer receives the exact commit and this bounded scope through the native dispatch, reports findings through its handoff, and never mutates the candidate. Requested configuration is gpt-6.1-sol/high with fresh context; actual runtime telemetry is unknown unless exposed by the host.

Current review and merge evidence is attached to the delivery PR and session handoff; it is not proof of application acceptance. See [design task](../tasks/design-package.md) and [release acceptance](../specifications/acceptance.md).
