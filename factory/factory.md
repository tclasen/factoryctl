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

**Current scope correction:** product implementation is paused. The [operator correction](sources/operator-brief.md#compatibility-scope-correction--2026-10-28) removes the custom runtime/assembly/ABI/loader and legacy-browser compatibility expansion from active scope. Ordinary OS libraries are permitted; the requested 2016 OS coverage remains unqualified where standard tools cannot provide it. Earlier static-interpretation and implementation-admission statements below are historical where superseded. No product completion or baseline compliance is claimed. Only the documentation scope correction is active.

Factoryctl is a minimal CLI, local GUI and MCP interface to Markdown. The AI harness interprets that Markdown to understand and run a software factory or any other kind of knowledge-work factory. Factoryctl makes the files navigable and supplies standardized document tools; it does not implement the factory's workflow or policy.

The [minimal-interface decision](decisions/minimal-interface.md) governs implementation: keep structural validation and safe file access in factoryctl, and keep roles, readiness, approvals, lifecycle decisions and execution in the agent's interpretation of Markdown.

## Desired state — operator requirements

The [product specification](specifications/product.md) defines the accepted scope from the [operator interview](sources/operator-brief.md). The [knowledge format](specifications/knowledge-format.md) describes the OKF bundle and its extension. [Project discovery](specifications/project-discovery.md) works without a root config file; [personal runtime preferences](specifications/personal-runtime.md) stay outside project Git. The [template contract](specifications/templates.md) explains shared customization.

## Current state — observed

The current implementation candidate is `a7c22ba847d58e8f4ece6a22ee9bb73302a64c02`. CLI, HTTP and MCP share structural document services; the browser supplies project collections, overview/list/detail views, backlinks, raw/field/create/delete editing, enum boards and external-change refresh. Declarative namespaced tools, prompt preview/explicit launch, digest-bound repository HTML/CSS/JavaScript views, and local delivery-assessment/content-freshness display are implemented. Workflow interpretation and authority remain with the harness and project Markdown.

Native Codex conversations, trusted free-text/choice/custom answers, generated views with plain fallback, same-thread follow-up and personal model/default settings have current execution evidence. Single-project, monorepo and sibling-repository journeys exercise the same controls, including cross-project links, target-safe edits and concurrent native question isolation. The [completion task](tasks/project-completion.md) owns the build and qualification scope; the [current plan](designs/implementation-plan.md) distinguishes implementation from acceptance.

The local package `artifacts/completion-review/a7c22ba-final-v2/` contains the combined film and source-bound evidence: 169 named race-test passes, six inspected target builds and a prepared matrix proposing 105 passes with U045 blocked. Final independent review and delivery are pending. Linux static linkage was observed; Darwin system libraries and Windows platform DLLs leave the operator's strict static-binary requirement unanswered. No exception, full qualification, project completion or release is claimed.

The declared 100-repository/10,000-node/100 MiB benchmark measured loader refresh at about 2.545 seconds, missing the proposed one-second engineering target. This is not an operator-approved latency limit or an end-to-end browser refresh measurement. Product execution is qualified only on the current Mac; foreign binaries were built and inspected, not run. External MCP Apps interoperability remains unqualified and the CLI is retained.

The [first CLI evidence](sources/cli-milestone-evidence.md), [interaction milestone evidence](sources/interaction-milestone-evidence.md), [visible-intent slice](tasks/visible-intent.md) and [design validation](sources/design-validation.md) remain historical, bounded records. They are not substituted for the current completion candidate's evidence.

## Understanding the design

Read the [architecture](designs/architecture.md), [interaction specification](specifications/gui.md), [CLI/MCP/HTTP contracts](specifications/interfaces.md), and [Codex integration](specifications/agent-runtime.md). The [implementation plan](designs/implementation-plan.md) links buildable slices, requirements and checks. [Release acceptance](specifications/acceptance.md) defines completion without claiming those checks have run.

## Knowledge beyond software

[Research curation and contract pursuit](specifications/knowledge-work.md) demonstrate custom types, source attribution and configurable boards. The [researcher profile](agents/researcher.md) illustrates how instructions remain knowledge nodes. Software specifications and designs are an optional starter set, exercised here to build the product itself.

## Knowledge quality

[Epistemic conventions](specifications/epistemics.md) separate requirements, current observations, assumptions and LLM proposals. Markdown headings and links carry the same distinctions as the GUI. Source claims are not promoted to verified facts by being copied into this graph. Git supplies historical revisions; files describe current knowledge.

## Development and releases

Follow the [commit and release policy](decisions/release-policy.md). The initial release target is 0.1.0; there is no released executable yet.

## Runtime-extensible document behavior

The [extension contract](specifications/extensions.md) describes user-defined types, structural definitions and agent-interpreted lifecycle instructions and namespaced tools loaded at runtime. Custom code runs in the browser UI; arbitrary local work is requested through prompts to the selected harness. Factoryctl remains a thin document interface and harness launcher, with no Git subprocesses or built-in agent scheduler. The [extension task](tasks/12-extension-contract.md) and [implementation brief](sources/completion-extensions-brief.md) describe the runtime-loaded fixed document bindings and rejection of executable backend handlers. Current browser and native-runtime checks support the implemented boundary; final whole-product acceptance remains separate.

## Required collections, interaction and release evidence

The product must exercise [one project, monorepo projects and sibling repositories](specifications/project-scenarios.md) with the same project-collection UI and depth-1 discovery. [Agent questions](specifications/user-interaction.md) add free text, choices and follow-up; MCP Apps and generated browser UI require a prototype before deciding whether to retire the interactive CLI. Current CLI and browser evidence exercises the three project layouts, including actual native questions and first-party generated views. External MCP Apps interoperability remains unqualified; the CLI is retained.

The [requirements traceability matrix](specifications/requirements-traceability.md) maps each operator requirement to implementation tasks and acceptance. Historical milestone evidence remains distinct from the current GUI, collection and runtime qualification package; final review and the unresolved U045 row still prevent full-product qualification. Every release requires current regression evidence and a [full-feature Playwright video](specifications/release-evidence.md), covering older features as well as new ones with readable pacing, smooth scrolling, overlays and highlights.

## Milestone demo

The stable local demo handoff is `/Users/Shared/factoryctl/artifacts/demo/latest.mp4`; the prepared combined film and 106-row accounting are versioned in `artifacts/completion-review/a7c22ba-final-v2/`. The [earlier CLI review package](sources/cli-review-video-evidence.md) remains historical evidence. Independent film/matrix review supports the partial package and U045 remains blocked, so full-matrix qualification is incomplete. Every milestone delivery must follow the [mandatory demonstration gate](specifications/release-evidence.md#milestone-review-delivery), including current evidence for old and new implemented features.

## Current planning authority

The [accepted working-method correction](decisions/intent-and-review.md) separates requests, agent proposals and execution authorization. The [CLI-first order](decisions/cli-first-sequence.md) is historical and superseded; it was never an original operator sequencing instruction. The browser foundation and bounded conversation slice make the knowledge inspectable. The [completion revision](designs/implementation-plan.md#completion-plan-revision-2--authorized-implementation) now authorizes remaining implementation and qualification; evidence distinguishes implemented behavior from full-product acceptance.


## Completion implementation in qualification

The explicit instruction to implement the remaining tasks is recorded in [the completion task](tasks/project-completion.md). Current code adds generic web editing and boards, refresh/backlinks, trusted digest-bound repository views, declarative namespaced CLI/HTTP/MCP tools, prompt preview/launch, and delivery-assessment/content-freshness display. Independent review found and prompted repairs for truncated lists, stale template dependencies, custom asset containment, cross-project links and cache invalidation. The current combined evidence now covers these features, S1–S3 interactions and synthetic knowledge-work examples. Independent final review, resolution of U045, qualification gates and verified remote-main delivery remain outstanding. No full-product completion or release is claimed by this paragraph.
