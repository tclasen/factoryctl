---
type: "Design"
title: "Implementation sequence and requirements coverage"
description: "Implementation sequence and requirements coverage."
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

# Implementation sequence and requirements coverage

This is a proposed build plan for the [confirmed product scope](../specifications/product.md). All implementation tasks are currently planned. The design task prepares the package; it does not implement these tasks. Agent execution follows repository instructions and user authority. No scheduler or approval engine is encoded in this plan.

## Buildable slices

| Task | Prerequisites | Requirement coverage | Acceptance |
| --- | --- | --- | --- |
| [Qualify the native Codex integration](../tasks/00-runtime-spike.md) | None | R11, R12, R13, R14, R24 | A08 |
| [Build the executable foundation and project discovery](../tasks/01-project-foundation.md) | None | R01, R02, R04, R16, R17 | A01, A12 |
| [Implement OKF reading and safe file storage](../tasks/02-knowledge-store.md) | 01-project-foundation | R03, R05, R07, R08, R22 | A02, A06 |
| [Implement definitions and live structural validation](../tasks/03-definition-validation.md) | 02-knowledge-store | R06, R07, R08, R24, R25 | A03, A04 |
| [Qualify runtime document extensions](../tasks/12-extension-contract.md) | 03-definition-validation | R06, R13, R16, R26, R27, R28 | A14 |
| [Implement deterministic record and prompt templates](../tasks/04-templates.md) | 03-definition-validation | R06, R11, R14, R15, R19, R25 | A05 |
| [Deliver the human and CI command interface](../tasks/05-cli.md) | 04-templates | R09, R18, R22, R24, R25 | A06, A12 |
| [Expose shared MCP and local HTTP services](../tasks/06-transports.md) | 05-cli, 12-extension-contract | R13, R17, R22, R23, R24, R25 | A07 |
| [Implement native Codex conversation integration](../tasks/07-codex-adapter.md) | 00-runtime-spike, 06-transports | R11, R12, R13, R14, R18, R24 | A08 |
| [Deliver the default project, record, board and conversation GUI](../tasks/08-gui.md) | 07-codex-adapter | R07, R08, R10, R11, R17, R22, R24, R25 | A09 |
| [Support optional trusted custom presentation](../tasks/09-custom-views.md) | 08-gui | R02, R15, R23 | A10 |
| [Qualify dogfooding and non-software workflows](../tasks/10-knowledge-work.md) | 09-custom-views | R06, R08, R19, R20, R25 | A11 |
| [Qualify performance, portability and complete release behavior](../tasks/11-release-qualification.md) | 10-knowledge-work | R01, R03, R04, R16, R17, R21, R22, R24, R25 | A12, A13 |

The runtime spike and project foundation can be investigated independently. Storage, definitions and structural validation establish the common contracts before transport/UI code. The first useful slice is a CLI that reads and validates this very bundle. The next adds MCP and native conversations; the GUI follows shared service contracts. Custom views and non-software examples prove extensibility before release qualification.

## Rules for implementers

Keep one authoritative status in each task; this document deliberately does not duplicate it. Follow links for current progress. Read the applicable specification and decision before coding. Proposals are identified in frontmatter/prose; if empirical evidence invalidates a proposal, update the living spec/design and explain the decision in a decision record. Git retains prior revisions.

Tests belong to the task that introduces the behavior. A final release task verifies integration and packaging; it does not defer basic field/link validation or failure handling until the end. Use real temporary Git repositories for storage/discovery tests, fake protocol servers for failure determinism and a bounded live Codex exercise for runtime qualification. Initial executable tests run on the current Mac only.

## Open technical risks and owners

- Runtime spike: exact current app-server capability/override behavior; owned by task 00, prerequisite to task 07.
- Storage: portable lock/replacement semantics and external-editor race limits; owned by task 02.
- Extensions: arbitrary JS isolation and message-bridge behavior; owned by task 09.
- Portability: single-binary interpretation and unexecuted non-Mac targets; owned by task 11 and the [portability decision](../decisions/portable-binary.md).
- Performance thresholds are engineering proposals; task 11 measures them and records any justified revision. The confirmed scale target is 100 repositories/10,000 nodes.

These risks have explicit acceptance work; they are not hidden claims of compatibility. No unresolved business decision is delegated to a library choice. Framework/library versions are bounded implementation choices, selected and pinned when their task begins.

## Completion and handoff

Completion requires A01–A14 in [acceptance](../specifications/acceptance.md), with failures and permitted platform exclusions explicit. Document actual capability in [the factory overview](../factory.md) as it ships. Keep intended behavior in specifications and observed verification in task/evidence records. Git operations are performed by the agent/user under the repository's delivery workflow, never by factoryctl itself.

## Revised first vertical slice

Qualify declarative/WASM extension registration with two types and a third runtime installation before finalizing GUI/MCP consumers. Then demonstrate rendered Markdown → explicit GUI action → selected harness connected to MCP → type-aware read/update → validated atomic persistence → refreshed view. Existing foundation/storage/validation work remains prerequisite; implementation and release checks must cover R26–R28/A14. The extension contract and harness permissions/tool refresh are empirical risks, not capabilities implied by Markdown specifications.
