---
type: "Design"
title: "Research basis for the software-development starter set"
description: "Research basis for the software-development starter set."
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

# Research basis for the software-development starter set

This report summarizes primary documentation inspected during the design interview. Descriptions of the tools are source-backed observations; recommendations for factoryctl are LLM-derived proposals. The operator accepted an optional software starter set and living specifications maintained alongside Git history.

## Compared approaches

| Approach | Observed convention | Proposed adaptation |
| --- | --- | --- |
| GitHub Spec Kit | Feature specification, implementation plan, task checklist and optional supporting research/contracts | Separate behavior from design and preserve explicit acceptance criteria |
| OpenSpec | Persistent capability specs plus proposals/deltas/design/tasks for changes; customizable workflow schemas | Keep current capability specs; use configurable definitions/templates without requiring a delta workflow |
| Kiro | Requirements/design/tasks with structured EARS-style behaviors | Use testable requirements and concrete error scenarios inside each specification |
| BMad | Briefs/PRDs/design and linked initiatives/epics/stories with requirement references | Trace requirements into tasks while avoiding extra mandatory record families |
| Superpowers | Plans carry interfaces, files, task boundaries and verification for fresh implementers | Make each implementation task a checkable slice with explicit handoff |
| SpecWeave | Increment spec plus CLI-owned metadata/ledger | Borrow acceptance/evidence links; retain current state in Markdown and Git history instead of adopting the ledger |
| Tessl framework description | Component capabilities linked to tests and APIs | Borrow traceability, not its historical generator directives |

Primary sources: [Spec Kit specification template](https://github.com/github/spec-kit/blob/main/templates/spec-template.md), [plan template](https://github.com/github/spec-kit/blob/main/templates/plan-template.md), [OpenSpec schema](https://github.com/Fission-AI/OpenSpec/blob/main/schemas/spec-driven/schema.yaml), [Kiro requirements-first workflow](https://kiro.dev/docs/specs/feature-specs/requirements-first/), [BMad requirements](https://docs.bmad-method.org/plan/define-requirements-and-a-specification/), [Superpowers planning source](https://github.com/obra/superpowers/blob/main/skills/writing-plans/SKILL.md), [SpecWeave file reference](https://spec-weave.com/docs/reference/metadata-reference/), and [Tessl framework description](https://tessl.io/blog/tessl-launches-spec-driven-framework-and-registry).

Tessl's former `reference/spec-syntax` URL returned Page Not Found during research; its older syntax is not treated as a current interoperability contract. Tool documentation evolves; this is a dated design comparison, not a claim of permanent feature parity or an empirical benchmark of code quality.

## Complementary conventions

[Cucumber/Gherkin](https://cucumber.io/docs/gherkin/reference/) provides concrete Given/When/Then scenarios, but Markdown scenarios alone are not executable tests. [Nygard's ADR convention](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions) preserves context, decision and consequences. [arc42](https://arc42.org/overview/) supplies a completeness check for architecture, interfaces, deployment, quality and risks. These are compatible document practices, not one universal file format.

OpenSpec's own [schema-resolution spec](https://github.com/Fission-AI/OpenSpec/blob/main/openspec/specs/schema-resolution/spec.md) is a concrete dogfooding example: override rules are expressed as requirements and scenarios. This shows the convention is used in that repository; it does not prove that every implementation conforms to every scenario.

## Decision for this project

Use Specification and Design as optional custom domain types alongside built-in Task and Decision. Requirements are stable identifiers within capability specifications. A design can be an implementation plan without creating another mandatory type. Each task is its own record and owns its one current progress field; plans link to it instead of keeping duplicate status checklists. Completion evidence belongs with the task or a linked substantive observation; Git supplies historical changes.

The starter set is ordinary definitions, record templates, prompts and instructions. It has no special workflow engine and does not make software development the mandatory factory model. [Knowledge-work examples](../specifications/knowledge-work.md) test that distinction.
