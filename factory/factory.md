---
type: "Factory"
title: "Factoryctl: purpose, desired state and current state"
description: "Factoryctl: purpose, desired state and current state."
status: "draft"
factory:
  format: 1
  assertion:
    frame: "mixed"
    certainty: "mixed"
    basis: "mixed"
sources:
  -
    id: "operator"
    resource: "sources/operator-brief.md"
    title: "Operator requirements and scope"
    author: "human:operator"
project_id: "factoryctl"
bundle_uri: "https://github.com/tclasen/factoryctl/blob/main/factory/"
---

# Factoryctl

Factoryctl is a minimal CLI, local GUI and MCP interface to Markdown. The AI harness interprets that Markdown to understand and run a software factory or any other kind of knowledge-work factory. Factoryctl makes the files navigable and supplies standardized document tools; it does not implement the factory's workflow or policy.

The [minimal-interface decision](decisions/minimal-interface.md) governs implementation: keep structural validation and safe file access in factoryctl, and keep roles, readiness, approvals, lifecycle decisions and execution in the agent's interpretation of Markdown.

## Desired state — operator requirements

The [product specification](specifications/product.md) defines the accepted scope from the [operator interview](sources/operator-brief.md). The [knowledge format](specifications/knowledge-format.md) describes the OKF bundle and its extension. [Project discovery](specifications/project-discovery.md) works without a root config file; [personal runtime preferences](specifications/personal-runtime.md) stay outside project Git. The [template contract](specifications/templates.md) explains shared customization.

## Current state — observed

This repository contains Go document services and a CLI alongside the planning bundle. The current [visible-intent slice](tasks/visible-intent.md) adds read-only browser exploration; see [one current plan](designs/implementation-plan.md) for observed work versus proposed continuation. The [first CLI milestone](tasks/first-cli-milestone.md) implements project discovery, record reading/search, structural validation, initialization, templates and revision-checked document editing. Build and usage instructions are in the repository README. [CLI evidence](sources/cli-milestone-evidence.md) distinguishes executed checks from pending acceptance. The [agent-interaction slice](tasks/13-user-interaction.md) now adds opt-in native Codex conversations, scoped MCP document/question/view tools, trusted answers and isolated generated browser views. [Current evidence](sources/interaction-milestone-evidence.md) distinguishes live results from remaining gaps. No release or full GUI/full-matrix qualification is claimed. [Design validation](sources/design-validation.md) remains historical document evidence.

## Understanding the design

Read the [architecture](designs/architecture.md), [interaction specification](specifications/gui.md), [CLI/MCP/HTTP contracts](specifications/interfaces.md), and [Codex integration](specifications/agent-runtime.md). The [implementation plan](designs/implementation-plan.md) links buildable slices, requirements and checks. [Release acceptance](specifications/acceptance.md) defines completion without claiming those checks have run.

## Knowledge beyond software

[Research curation and contract pursuit](specifications/knowledge-work.md) demonstrate custom types, source attribution and configurable boards. The [researcher profile](agents/researcher.md) illustrates how instructions remain knowledge nodes. Software specifications and designs are an optional starter set, exercised here to build the product itself.

## Knowledge quality

[Epistemic conventions](specifications/epistemics.md) separate requirements, current observations, assumptions and LLM proposals. Markdown headings and links carry the same distinctions as the GUI. Source claims are not promoted to verified facts by being copied into this graph. Git supplies historical revisions; files describe current knowledge.

## Development and releases

Follow the [commit and release policy](decisions/release-policy.md). The initial release target is 0.1.0; there is no released executable yet.

## Runtime-extensible document behavior

The [extension contract](specifications/extensions.md) describes user-defined types, structural definitions and agent-interpreted lifecycle instructions and namespaced tools loaded at runtime. Custom code runs in the browser UI; arbitrary local work is requested through prompts to the selected harness. Factoryctl remains a thin document interface and harness launcher, with no Git subprocesses or built-in agent scheduler. The [extension prototype](tasks/12-extension-contract.md) must qualify this boundary before implementation claims. Type-derived extension registration remains pending; first-party generated browser presentation and native Codex launch now have bounded implementation evidence.

## Required collections, interaction and release evidence

The product must exercise [one project, monorepo projects and sibling repositories](specifications/project-scenarios.md) with the same project-collection UI and depth-1 discovery. [Agent questions](specifications/user-interaction.md) add free text, choices and follow-up; MCP Apps and generated browser UI require a prototype before deciding whether to retire the interactive CLI. The document CLI exercises project collections; actual harness questions and first-party generated browser UI now have bounded evidence. External MCP Apps interoperability remains unqualified; the CLI is retained.

The [requirements traceability matrix](specifications/requirements-traceability.md) maps each operator requirement to implementation tasks and acceptance. Initial CLI and current interaction evidence remain separate from unimplemented full-GUI and release requirements. Every release requires current regression evidence and a [full-feature Playwright video](specifications/release-evidence.md), covering older features as well as new ones with readable pacing, smooth scrolling, overlays and highlights.

## Milestone demo

The [CLI review package](sources/cli-review-video-evidence.md) contains one scene-annotated CLI demo at a stable link, with complete matrix gap accounting in supporting documents. Full-matrix end-to-end qualification remains incomplete. Every future milestone delivery must follow the [mandatory demonstration gate](specifications/release-evidence.md#milestone-review-delivery), including current evidence for old and new implemented features.

## Current planning authority

The [accepted working-method correction](decisions/intent-and-review.md) separates requests, agent proposals and execution authorization. The [CLI-first order](decisions/cli-first-sequence.md) is historical and superseded; it was never an original operator sequencing instruction. A read-only browser slice makes the planning documents inspectable while full GUI/editing, conversations and MCP Apps remain unqualified.
