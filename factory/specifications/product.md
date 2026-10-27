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

## Repository and persistence

- **R01 — Repository discovery.** Operate from a Git checkout or recursively from a directory containing checkouts; an explicit path overrides CWD. Preserve a multi-project model for one or many projects.
- **R02 — Storage boundary.** Root `factory.yaml` contains configuration, optional `.factory/` contains assets, and `factory/` alone contains knowledge metadata. A configured project works with no `.factory/` directory.
- **R03 — File persistence.** All durable factory knowledge is human-readable OKF Markdown. No database, external config file, durable index, conversation store or event ledger is introduced by factoryctl.
- **R04 — Git responsibility.** Git is mandatory and supplies history. Factoryctl does not commit, stage, fetch, push, merge, create branches/worktrees or resolve Git conflicts. The external agent/user performs those operations. Read-only repository inspection is permitted.

## Knowledge model

- **R05 — OKF.** Every Markdown file under `factory/` is a valid OKF concept, including the factory, definitions, profiles, prompts and templates. No index files or generated historical logs.
- **R06 — Extensibility.** Supply agent-profile, task and decision types; allow new types and disabling supplied types. Definitions describe fields, required values, relationships and views. Infrastructure types describe the format itself.
- **R07 — Navigation.** Every concept supports readable Markdown, a generic list/detail presentation, outgoing links and always-visible backlinks. Resolve links across discovered projects and expose unresolved targets honestly.
- **R08 — Knowledge quality.** Distinguish desired state, observed state and proposed change; facts, assumptions and questions; authoritative/directly observed evidence and LLM inference. Preserve source attribution at claim level where needed.

## Interaction

- **R09 — CLI.** Support interactive use and deterministic noninteractive CI workflows, with useful errors and machine-readable output available on stdout.
- **R10 — GUI.** Launch a browser GUI with project overview/switcher, all-project lists, list/detail views for every type, configurable kanban, direct field/Markdown editing and an agent conversation panel.
- **R11 — Agent actions.** Prompt templates define actions on records and compose reusable prompts. An action creates a fresh conversation; subsequent messages continue it. Several conversations can run concurrently.
- **R12 — Codex.** Integrate current Codex CLI as the first runtime, including streaming output, follow-ups, input/approval requests, errors and cancellation. Codex owns authentication and conversation persistence.
- **R13 — MCP.** Expose discovery, context, read, creation and safe editing tools to the agent. CLI, GUI and MCP share the same validation and file behavior.
- **R14 — Profiles.** Store profile instructions and optional model/reasoning settings as OKF nodes. Do not generate Codex-specific instruction/config files.
- **R15 — Presentation customization.** Support user-authored HTML, CSS and JS and nested templates. Default assets are embedded; customization is optional and versioned.

## Operation and collaboration

- **R16 — Portable executable.** Deliver a single compiled executable with embedded UI/defaults for macOS, Linux and Windows, ARM64/x86-64, suitable for desktop, server and CI.
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

## Release boundary

Release 1 includes every R01–R25 behavior above, with the platform-testing limitation in R21. Excluded: built-in business-rule enforcement, automation scheduler, runtime replacement, managed credentials, shared multi-user server accounts, automated Git synchronization, migration of earlier factories, full ontology reasoning, a visual graph canvas, and an asset marketplace. Backlinks and cross-project graph relationships are included even without a graph canvas.

See [acceptance criteria](acceptance.md) and the [requirement-to-task map](../designs/implementation-plan.md). No design choice silently expands external-action authority.
