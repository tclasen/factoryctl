---
type: "Specification"
title: "Product scope and requirements"
description: "Product scope and requirements."
status: "draft"
factory:
  format: 1
  assertion:
    frame: "desired"
    certainty: "mixed"
    basis: "mixed"
sources:
  -
    id: "operator"
    resource: "../sources/operator-brief.md"
    title: "Operator requirements and scope"
    author: "human:operator"
---

# Product scope and requirements

The requirements below are synthesized from [operator statements](../sources/operator-brief.md). Their identifiers are stable within this specification. They describe desired behavior, not delivered functionality. Exact implementation details in linked contracts are proposed mechanisms.

## Governing principle

**R29 — Minimal interface.** Factoryctl is a small Markdown viewer/editor, standardized document-tool provider and harness launcher. All factory/domain workflow meaning is interpreted by the AI harness from Markdown. Core validation covers structure and safe storage only: no transition guards, conditional workflow requirements, readiness/approval evaluator or business-rule language. New knowledge-work workflows use different Markdown and optional UI assets without backend domain logic. See [the accepted boundary](../decisions/minimal-interface.md).

## Repository and persistence

- **R01 — Repository discovery.** Support a single project in a Git checkout, independent projects in one monorepo, and projects in sibling Git repositories. Discover immediate subfolders only by default (recursion limit 1); an explicit path overrides CWD. A project owns its bundle/state independently of repository identity. Exercise all three layouts throughout implementation and testing.
- **R02 — Storage boundary.** No root configuration file is required or read. Optional `.factory/` contains assets; `factory/` alone contains knowledge metadata. A project works without either a root config file or `.factory/`.
- **R03 — File persistence.** All durable factory knowledge is human-readable OKF Markdown. No database, project config file, durable index, conversation store or event ledger is introduced by factoryctl. Optional personal runtime preferences are human-readable Markdown outside the repo, as scoped by [personal runtime preferences](personal-runtime.md). Git-versioned records are the only initial tracking backend; external backends are deferred.
- **R04 — Git responsibility.** Git is mandatory and supplies history. Factoryctl does not commit, stage, fetch, push, merge, create branches/worktrees or resolve Git conflicts. The external agent/user performs those operations. Read-only repository inspection uses filesystem metadata or in-process libraries only. Factoryctl never spawns Git, directly or through a helper; agent harnesses may do so.

## Knowledge model

- **R05 — OKF.** Every Markdown file under `factory/` is a valid OKF concept, including the factory, definitions, profiles, prompts and templates. No index files or generated historical logs.
- **R06 — Extensibility.** Supply agent-profile, task and decision types; allow new types and disabling supplied types. Definitions describe fields, required values, relationships and views. Runtime definitions also supply type-specific tools and lifecycle actions; see [extensions](extensions.md). Infrastructure types describe the format itself.
- **R07 — Navigation.** Every concept supports readable Markdown, a generic list/detail presentation, outgoing links and always-visible backlinks. Resolve links across discovered projects and expose unresolved targets honestly.
- **R08 — Knowledge quality.** Distinguish desired state, observed state and proposed change; facts, assumptions and questions; authoritative/directly observed evidence and LLM inference. Preserve source attribution at claim level where needed.

## Interaction

- **R09 — CLI.** Support interactive use and deterministic noninteractive CI workflows, with useful errors and machine-readable output available on stdout.
- **R10 — GUI.** Use the identical project-collection UI for all three repository layouts, including a one-project collection. Launch a browser GUI with project overview/switcher, all-project lists, list/detail views for every type, configurable kanban, direct field/Markdown editing and an agent conversation panel.
- **R11 — Agent actions.** Prompt templates define actions on records and compose reusable prompts. An action creates a fresh conversation; subsequent messages continue it. Several conversations can run concurrently.
- **R12 — Codex.** Integrate current Codex CLI as the first supported adapter, selected by the contributor outside project Git, including streaming output, follow-ups, input/approval requests, errors and cancellation. Codex owns authentication and conversation persistence.
- **R13 — MCP.** Expose discovery, context, read, creation and safe editing tools to the agent. CLI, GUI and MCP share the same validation and file behavior.
- **R14 — Profiles.** Store shared profile instructions as runtime-neutral OKF nodes. Runtime, executable, model and reasoning preferences belong to the contributor outside Git. Do not generate Codex-specific instruction/config files.
- **R15 — Presentation customization.** Support user-authored HTML, CSS and JS and nested templates. Default assets are embedded; customization is optional and versioned.

## Operation and collaboration

- **R16 — Portable executable.** Deliver a single static executable with embedded UI/defaults for macOS, Linux and Windows, ARM64/x86-64, suitable for desktop, server and CI.
- **R17 — Distributed collaboration.** Different users run against their own checkouts and exchange Git changes. One trusted user per served instance; authentication and shared-user infrastructure are external.
- **R18 — Automation boundary.** Agents own scheduling, workflow policy, task dispatch and recursive factory loops. Factoryctl may start a conversation that asks an agent to run such a loop; it is not a workflow engine.
- **R19 — Dogfooding.** Build and manage factoryctl using its own record format. Ship optional software-development definitions/templates; maintain living specs and one authoritative current task state.
- **R20 — Knowledge-work breadth.** Demonstrate research/knowledge curation and government-contract pursuit without software-specific core behavior.
- **R21 — Scale and verification.** Design for 100 repositories and 10,000 total nodes. Execute initial platform tests only on the operator's current Mac. Distinguish build targets from executed platform tests.

## Proposed integrity requirements

These engineering requirements are LLM-derived from the confirmed concurrent editing and arbitrary-asset requirements, not additional operator workflow policy.

- **R22 — Safe mutation.** Reject stale edits, path escapes and malformed outputs; preserve unknown metadata. Never silently overwrite a file changed after it was read. Document coordination limits with external writers.
- **R23 — Extension trust.** Treat custom executable presentation as code and require an explicit per-process trust action. Ordinary knowledge Markdown remains inert. Keep project selection and agent permissions visible.
- **R24 — Observable failures.** Invalid records remain discoverable with diagnostics; unresolved references and untested capabilities remain visibly unresolved. Unsupported runtime capabilities produce actionable errors.

## Additional confirmed requirement

- **R25 — Continuous structural validation.** Validate documents as they are written: syntax, mandatory fields, field value types, link targets/anchors and record placement. Apply the same checks in CLI, GUI and MCP; surface errors introduced by external editors. This requirement was explicitly added by the operator. See the [validation contract](validation.md).

## Runtime document extensions — confirmed scope

- **R26 — Runtime type packages.** Users add or change document types, tools, structural definitions, lifecycle instructions and prompt-backed actions without rebuilding factoryctl. Custom HTML/CSS/JavaScript runs only in the browser UI. Arbitrary local execution is delegated to the selected AI harness through a prompt; no custom backend code is loaded.
- **R27 — Shared controlled document access.** Type-specific MCP tools and GUI/CLI actions share validation, revision checks and persistence. Generic writes apply the same structural checks. Read/context tools return type definitions, agent-interpreted instructions and revision information; returning instructions does not enforce their workflow meaning. Report whether harness-level access restrictions are verified or merely advisory.
- **R28 — Extension boundary.** Validate and namespace declarative definitions/tools, preserve browser UI isolation, define revision/conflict behavior and refresh available tools safely. No backend plugins, custom code interpreters, executable lifecycle callbacks or generic shell endpoints. Agents retain local execution, Git and scheduling.

## Release boundary

Release 1 includes every R01–R33 requirement above, including the explicit qualification/disposition rather than guaranteed Apps interoperability in R31, with the platform-testing limitation in R21. Excluded: domain business-rule enforcement (structural field checks are included), transition guards, state-conditioned requirements, policy evaluators, backend user-code execution, automation scheduler, runtime replacement, managed credentials, shared multi-user server accounts, automated Git synchronization, migration of earlier factories, full ontology reasoning, a visual graph canvas, and an asset marketplace. Backlinks and cross-project graph relationships are included even without a graph canvas.

See [acceptance criteria](acceptance.md) and the [requirement-to-task map](../designs/implementation-plan.md). No design choice silently expands external-action authority.

## User dialogue and rich web evaluation

- **R30 — Agent questions.** Expose a harness-callable tool for free-text questions and multiple choice, with user replies returned to the requesting agent and follow-up conversation in the same thread. Both current interactive CLI and web contracts support it; headless mode fails explicitly when no answer surface is available. See [user interaction](user-interaction.md).
- **R31 — Rich web qualification.** Prototype MCP Apps and runtime-generated browser UI, record actual harness/host capabilities and fallback behavior, and evaluate whether the web interface can replace the interactive CLI. CLI retirement is conditional on that evidence and an explicit product decision. The web launcher, MCP service and any retained automation commands are separate scope choices. No native Apps support is assumed from MCP connectivity alone.

- **R32 — Full-release demonstration.** Every release includes a human-paced Playwright video of every implemented feature, old and new, with explanatory overlays, interaction highlights and smooth scrolling. Deprecated features remain covered until actually removed. See [release evidence](release-evidence.md).
- **R33 — Requirements traceability.** Maintain a full operator-requirement matrix, implementation mapping and fresh regression evidence for each release. Preserve historical evidence and explicitly report gaps, failures and approved removals. See [the matrix](requirements-traceability.md).
