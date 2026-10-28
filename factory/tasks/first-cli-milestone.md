---
type: Task
title: Build the first usable CLI milestone
description: Implement and independently verify the initial document CLI.
status: draft
task_state: done
priority: high
owner: coordinator
depends_on: []
requirements: [R01, R02, R03, R04, R05, R06, R07, R08, R09, R15, R16, R17, R22, R24, R25]
factory:
  format: 1
---

# First CLI milestone

- Target System: project
- Change Classification: project
- Authorized Scope: factoryctl source, tests, build files and project-local documentation; normal branch/PR delivery to tclasen/factoryctl main
- Affected Paths or Resources: cmd/, internal/, testdata/, go.mod, go.sum, README.md, .gitignore, factory/ task and evidence records
- Owning Role: Coordinator
- Related Requirements: R01–R09, R15–R17, R22, R24, R25
- Related Tasks: 01-project-foundation through 05-cli
- Approval Requirements: operator request “please build it now”; independent exact-candidate review and ordinary protected PR delivery

## Admission and plan revision 1

Source request: “What is the next milestone on factoryctl?” followed by “please build it now.” The assistant selected and presented the CLI milestone; the operator authorized that bounded implementation on 2026-10-27. CLI-first sequencing was an assistant proposal, not an original operator instruction. The [current plan](../designs/implementation-plan.md) supersedes that ordering. Baseline: 1d6d0d4. Preserve untracked Finder files. No deployment, published release, organization changes, harness configuration changes or external messaging is included.

Required sources: [operator brief](../sources/operator-brief.md), [architecture](../designs/architecture.md), [implementation plan](../designs/implementation-plan.md), [minimal interface](../decisions/minimal-interface.md), specifications for project discovery, knowledge format, validation, templates, interfaces, project scenarios and acceptance, and tasks 01–05. Existing reviewed design supplies the architecture; fresh readiness review checks the bounded implementation plan before coding. A CLI milestone does not claim MCP, GUI, conversation integration or full release acceptance.

Steps: (1) independently check design readiness; (2) implement discovery and document/definition services; (3) implement safe storage, deterministic templates and CLI atop those services; (4) test success and negative cases, validate this bundle and build a CGO-disabled executable; (5) update requirement mappings and actual task status, independently review pinned candidate, repair findings; (6) commit, PR, checks, merge and verify owning remote main. Rollback uses an ordinary revert; never overwrite unrelated changes.

Acceptance: meaningful tests for depth-one discovery and S1–S3 collections, parsing and structural diagnostics, unknown metadata preservation, path/symlink safety and revision conflicts, templates and noninteractive JSON/exit behavior. The binary must read and validate its own bundle without network, installed agent or Git subprocesses. Document all unimplemented contracts explicitly; do not promote whole-release acceptance from a partial slice.

## Runs and exclusive claims

Coordinator owns admission, integration, documentation and Git index/remote operations. Developer run CLI-implementation owns application source/tests/build files and may coordinate explicit disjoint child scopes before mutation. Reviewer run CLI-readiness owns read-only design assessment; reviewer run CLI-final owns read-only pinned implementation assessment. Native worker identities remain session context. Shared checkout claims are procedural, not filesystem isolation.

Model selection: inherit the available parent model for implementation and fresh review; actual model/reasoning and token/cost telemetry are unknown. Tools: shell, file editing, native collaboration and connected GitHub. No new providers or paid runners. Checkpoint at transitions and every 15 active minutes; tests bounded to ten minutes and ordinary operations to two minutes; at most two diagnosed retries. Implementation dependencies follow tasks 01–05. Review cannot approve its own authored code.

## Current evidence and next action

Source checkout initially contains no implementation. Go 1.27.3 is available on Darwin ARM64. Shell GitHub credentials are invalid; connected GitHub tools will be assessed for normal delivery. This is a delivery risk, not an implementation blocker. No release is being published. Next: independent readiness review, then implementation and verification.

### Readiness checkpoint

Fresh read-only Reviewer assessed baseline commit 1d6d0d40a80dadb1b17a5320efa7ee0a4fc061e3 and this admission (initial SHA-256 d8dbafbf53ecc550cf14bbfe1d2edbf63a18ee3d39d573f79b649a9df568579e). Architectural fit, comprehension and execution readiness passed for the bounded CLI. Finding F01 requires keeping uncovered portions of tasks 01–05 pending, especially watchers and runtime question integration. This is one independent perspective, not implementation acceptance. Required source reading was reviewer-reported; full host trace/model telemetry remains unavailable.

Implementation now proceeds under that readiness assessment. The developer delegated a disjoint discovery claim covering internal/project/ and its tests; all other source remains with the developer lead. Coordinator retains factory/ documentation and Git operations. Connected GitHub access confirms the owning remote, no pre-existing open PR, and a main ruleset requiring PR delivery and thread resolution without bypass. Network-enabled read-only Git access also succeeds. Next: integration tests and pinned independent implementation review.

### Implementation and review checkpoint

The initial CLI implementation was committed as `51e95bb493915ab639175f89109c4cfe2a42e378` and published in [PR 8](https://github.com/tclasen/factoryctl/pull/8). Normal/race tests, vet, host builds and independent command journeys passed; clean-corpus validation found no Markdown errors. Fresh implementation review nevertheless reproduced IR01 (FIFO read hang) and IR02 (silent JSON output failure), retained in [CLI evidence](../sources/cli-milestone-evidence.md). Developer repair ownership covers source/tests only; Coordinator owns evidence and delivery. Exact repaired-candidate review remains the next acceptance gate. No release acceptance or downstream task completion is claimed.

### Proposed bounded milestone closure

The done state is proposed for PR 8 and becomes effective only after independent approval of the repaired exact candidate, normal merge and verified owning-remote main delivery. Before those gates, this task remains in review operationally. The documented first usable document CLI is the completed scope; tasks 01–05 retain uncovered criteria and all downstream runtime/UI/release work remains pending. Developer and child implementation claims are released; Coordinator retains delivery ownership until the gate passes. Acceptance evidence and exact final identities are recorded in PR 8. The retained project checkout needs no temporary task-worktree cleanup.
