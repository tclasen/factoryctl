---
type: "Task"
title: "Qualify performance, portability and complete release behavior"
description: "Qualify performance, portability and complete release behavior."
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
task_state: in_review
priority: "normal"
owner: "unassigned"
depends_on:
  - "10-knowledge-work.md"
specifications:
  - "../specifications/acceptance.md"
requirements:
  - "R30"
  - "R32"
  - "R33"
  - "R31"
  - "R01"
  - "R03"
  - "R04"
  - "R16"
  - "R17"
  - "R21"
  - "R22"
  - "R24"
  - "R25"
---

# Qualify performance, portability and complete release behavior

## Outcome

Run the full Mac acceptance suite and benchmark, cross-compile the six target OS/architecture combinations where supported, inspect linkage and prepare user/operator documentation.

## Scope and handoff

Implement the relevant modules in the [architecture](../designs/architecture.md), following the linked specifications in frontmatter. Read [operator authority](../sources/operator-brief.md) before any external effects. Historical planning checkpoint: this task was planned and not started, and architecture source paths did not yet exist. Qualification tooling and current results now exist, with blockers retained below. The [implementation plan](../designs/implementation-plan.md) establishes build order and shared contracts.

## Verification and done criteria

Record exact revision, runtime/dependency versions, corpus/hardware, errors and achieved metrics. Label Linux/Windows and non-host architecture execution untested. Verify files-only persistence, no app Git mutations, graceful shutdown and external-proxy deployment instructions without deploying a service.

Acceptance coverage: A12, A13 in [release acceptance](../specifications/acceptance.md). Include exact tested revision, commands, expected/observed results and remaining limitations in current task evidence when implemented. Do not equate a checked task field with proof of behavior. Run relevant negative tests as well as success journeys.

## Dependencies, risks and rollback

`depends_on` identifies prerequisites; an agent interprets readiness. Factoryctl validates references but does not schedule this work. Preserve unknown metadata and source evidence. Roll back failed implementation through ordinary Git changes by the agent, retaining useful evidence. Do not expand into deployment, purchases or external messaging.

## Commit and version qualification

Implement and verify the [commit and release policy](../decisions/release-policy.md): validate new commit and PR titles, calculate versions from Conventional Commits, start at 0.1.0, and prevent an unauthorized transition to or release of 1.0.0. Exercise the policy examples and negative authorization case. The 0.x line carries no compatibility promise. Conventional Commit/PR-title CI enforcement and pre-1.0 version calculation/guards are now installed and tested; the initial release target remains 0.1.0, while the current development binary is not a published release.

## Runtime definition integration

Verify complete A14 evidence from task 12 core fixtures, tasks 06–08 transport/harness behavior and task 09 custom browser board interactions. No backend user-code execution engine is part of the release. Record observed no-Git-subprocess and static-linkage results separately from agent Git behavior.

## Minimal-interface release gate

Verify R29/A15 across core fixtures and interface/harness/UI evidence. Workflow changes must be possible through Markdown and views without changing backend domain logic. Structural validation must never be reported as policy approval or workflow completion.

## Collection, interaction and continuing evidence

Run A01–A20 for each release candidate, including all S1–S3 journeys. Deliver the [full-feature Playwright video](../specifications/release-evidence.md), with readable pacing, smooth scrolling, explanatory overlays and interaction highlights. Reconcile the [atomic requirements matrix](../specifications/requirements-traceability.md) against implementation, prior releases and current test/video evidence. Old supported features remain in scope; deprecated features are omitted only after removal. Missing video, traceability or current required results blocks release. Record Apps disposition and any explicitly approved CLI retirement separately from actual supported features.


## Completion revision 2 qualification checkpoint — unresolved

Product candidate: `a7c22ba847d58e8f4ece6a22ee9bb73302a64c02`. Current [completion milestone evidence](../sources/completion-milestone-evidence.md) consolidates the full named race/vet/build results, native/browser/S1–S3 journeys, current dogfood diagnostics, six-target packaging, conventional-commit/version guards and local combined-demo/matrix identities. Independent [code review](../sources/completion-code-review.md) and [acceptance review](../sources/completion-acceptance-review.md) remain distinct from operator acceptance and remote-main delivery.

The source-stable attempt-009 benchmark exercises 100 repositories, 10,000 Markdown nodes and exactly 104857600 bytes. Cold load is 3270.927 ms, cached-query p95 18.089 ms, production loader refresh 2544.947 ms and sampled after-refresh RSS 523943936 bytes. The proposed one-second loader-refresh target is missed. OS cache was not evicted; RSS is sampled rather than peak, and HTTP/debounce/browser edit visibility plus browser/harness memory are excluded. These are reported engineering targets and measurements, not invented operator-approved numerical latency bounds. Earlier failed attempts remain retained.

All six CGO-disabled binaries build with embedded assets. Linux has no dynamic interpreter/dependencies; Darwin retains system-library interfaces and Windows imports a platform DLL. Foreign targets are built and inspected, never executed. U045 remains blocked: these results do not establish the stricter all-target static-binary wording, and the operator has not clarified or accepted an exception. Current dogfood validation also retains the preserved factory/.DS_Store diagnostic and explicitly unchecked links, rather than manufacturing an all-files/network pass.

`in_review` is deliberate and this task has no closure proposal while U045/static interpretation and full qualification remain unresolved. The proposed current accounting is 105 passing rows and one blocked row, not full-product completion. Final independent combined-media/matrix review, normal PR merge and owning-remote main verification remain required for milestone delivery. No tag, release publication, deployment or v1.0 transition is authorized or claimed. Coordinator owns the remaining disposition and delivery; the root completion task remains in progress.
