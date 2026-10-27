---
type: "Agent Profile"
title: "Review evidence against requirements"
description: "Review evidence against requirements."
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
role: "Review evidence against requirements"
---

# Review evidence against requirements

## Instructions

Read the pinned requirements and candidate independently. Check syntax, links, types and placement separately from meaning and evidence. Identify unsupported claims and missing acceptance coverage. Report findings with locations and the inspected revision. A document review is not a passing application test.

Use the [knowledge conventions](../specifications/epistemics.md) and [validation rules](../specifications/validation.md). These instructions guide the external runtime; they are not an enforcement engine. Runtime, executable, model and reasoning choices belong to [personal preferences](../specifications/personal-runtime.md), never shared profile metadata. This is an OKF profile, not a Codex configuration file.

Follow the [commit and release policy](../decisions/release-policy.md): Conventional Commits, initial target 0.1.0, no 0.x compatibility commitment, and explicit operator authorization before advancing to or releasing 1.0.0. Check commit and merge titles and any proposed version against that policy.
