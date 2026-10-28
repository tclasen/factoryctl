---
type: "Task"
title: "Build the executable foundation and project discovery"
description: "Build the executable foundation and project discovery."
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
task_state: done
priority: "high"
owner: "developer"
depends_on: []
specifications:
  - "../specifications/project-discovery.md"
requirements:
  - "R01"
  - "R02"
  - "R04"
  - "R16"
  - "R17"
---

# Build the executable foundation and project discovery

## Outcome

Create the Go executable skeleton, embedded-default mechanism, Factory-node metadata reader and read-only Git discovery. Implement explicit project selection and zero/one/many-repository results. Default discovery to recursion limit 1: inspect only immediate subfolders of the selected directory.

## Scope and handoff

Implement the relevant modules in the [architecture](../designs/architecture.md), following the linked specifications in frontmatter. Read [operator authority](../sources/operator-brief.md) before any external effects. Historical CLI checkpoint: the document portion was in progress under the [first CLI milestone](first-cli-milestone.md), with then-uncovered criteria recorded in [CLI evidence](../sources/cli-milestone-evidence.md). The current implemented outcome is reconciled below. The [implementation plan](../designs/implementation-plan.md) establishes build order and shared contracts.

## Verification and done criteria

Test CWD inside a repo, a parent directory with immediate-child repositories and excluded grandchild/deeper repositories, the same depth boundary with `--path`, explicit selection of a deeper repository or its parent, worktree .git files, inaccessible paths, symlink loops, duplicate project IDs, malformed Factory metadata, no factory.yaml and no .factory/ folder. Prove discovery performs no Git mutations.

Acceptance coverage: A01, A12 in [release acceptance](../specifications/acceptance.md). Include exact tested revision, commands, expected/observed results and remaining limitations in current task evidence when implemented. Do not equate a checked task field with proof of behavior. Run relevant negative tests as well as success journeys.

## Dependencies, risks and rollback

`depends_on` identifies prerequisites; an agent interprets readiness. Factoryctl validates references but does not schedule this work. Preserve unknown metadata and source evidence. Roll back failed implementation through ordinary Git changes by the agent, retaining useful evidence. Do not expand into deployment, purchases or external messaging.

## Collection, interaction and continuing evidence

Create and execute S1–S3 discovery fixtures from the [project matrix](../specifications/project-scenarios.md). Separate project root from Git root, include monorepo bundles without child .git markers, retain depth 1 and verify a one-project collection. Cover A16 as well as A01.

## Completion revision 2 reconciliation

Current product candidate: `a7c22ba847d58e8f4ece6a22ee9bb73302a64c02`. The explicit instruction to implement the remaining project tasks and its bounds are recorded in the [completion task](project-completion.md); the original authorization and historical checkpoints above remain distinct.

The Go executable, embedded defaults and filesystem-only project discovery are implemented. Current checks cover depth-one inclusion/exclusion, explicit selection, Git worktree metadata and S1–S3 project collections without a required project configuration file or factoryctl Git subprocess. Six-target builds are packaging evidence; the strict static-linkage boundary belongs to task 11, not a blanket foundation completion claim (A01/A16).

Current command, named-test, native/browser, scene and artifact identities are consolidated in [completion milestone evidence](../sources/completion-milestone-evidence.md), with the independent [code review](../sources/completion-code-review.md) and [acceptance review](../sources/completion-acceptance-review.md). The proposed matrix accounts for 105 passing rows and blocked U045; that proposal is not a full-product acceptance or release claim. Final combined evidence and independent review remain separate gates.

Conditional closure proposal: the bounded work described here is implemented or reconciled and ready for reviewed delivery. The proposed `done` state becomes effective ONLY after this PR is independently reviewed, merged normally and verified on owning-remote main. Before those gates this is a closure proposal, not a delivered result. The final independent reviewer must assess this task's evidence basis. Coordinator retains delivery/reconciliation ownership; no earlier PR/main result is invented by this update. Task 11 and the root completion task retain the unresolved static-linkage/full-qualification boundary.
