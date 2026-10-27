---
type: "Decision"
title: "Keep factoryctl a minimal Markdown interface"
description: "Keep factoryctl a minimal Markdown interface."
status: "draft"
factory:
  format: 1
  assertion:
    frame: "desired"
    certainty: "supported"
    basis: "authoritative-statement"
sources:
  - id: operator
    resource: "../sources/operator-brief.md"
    title: "Operator requirements and scope"
    author: "human:operator"
decision_state: "accepted"
---

# Keep factoryctl a minimal Markdown interface

## Decision and authority

The [operator's minimal-interface directive](../sources/operator-brief.md#minimal-interface-directive--controlling-operator-requirement) governs the product boundary. Factoryctl is a thin interface to Markdown, not an implementation of a software factory or any other knowledge-work factory. This is an accepted scope requirement, not a claim of shipped behavior.

Factoryctl discovers, renders and navigates documents; provides standardized read/write tools; validates document structure and safe writes; displays configurable browser views; and connects explicit prompts/context to the selected agent harness. The harness interprets the Markdown to understand roles, workflows, policies and the meaning of the work, and executes that work with its own tools. Git belongs to the agent/user.

## Standardization versus workflow meaning

| Factoryctl checks or presents | Agent interprets from Markdown |
| --- | --- |
| Required fields, types, enum membership and document placement | Whether a task is ready or should advance |
| Parseable content and valid structural references | Whether evidence is sufficient or a review is approved |
| Revisions, path scope and safe atomic updates | Dependencies, role responsibilities and approval policy |
| Lists, links, boards and prompt actions | Scheduling, sequencing and factory loops |

A `task_state` enum describes valid labels. The core may reject an unknown label, but it does not compare old/new labels against allowed transitions or require extra evidence solely because the state changed. A field can be structurally mandatory for every record of a type; conditional workflow requirements are prose for the agent. Structural validity never proves work completion, truth or approval.

A custom kanban view may render columns and request a field update through the shared safe writer. It can also request a harness prompt, such as asking the agent to review and update a ticket. Browser interactions do not turn the backend into a policy engine. Agent prompts, tool descriptions and document bodies are not executable server-side rules.

## Implementation restraint

Reuse the existing reader, renderer, structural validator, safe writer and harness adapter. New domain types are Markdown definitions, templates, instructions/prompts and optional browser assets. Type-specific tools are thin descriptions/input schemas over the same generic document operations; they contain no domain procedure. No package manager, workflow scheduler, policy interpreter, general rule language or executable backend plugin is needed.

For each proposed feature, ask: can a different knowledge-work factory use the same binary with different Markdown and views? If supporting its workflow requires new backend domain logic, move that responsibility into agent-interpreted Markdown. Generic improvements to document usability or integrity remain possible under ordinary scope review; this test does not ban maintenance of the interface itself.

## Supersession and verification

This decision removes the earlier transition-guard and state-dependent-validation proposal in PR #4. It preserves structural validation, revision checks, link integrity, custom browser UI and prompt-backed harness actions. Core workflow enforcement is excluded, not postponed.

[A15](../specifications/acceptance.md#a15-minimal-interface-boundary) verifies the boundary using different workflows against identical document schemas and binary code. No application implementation exists yet; implementation tasks must demonstrate these properties before release.
