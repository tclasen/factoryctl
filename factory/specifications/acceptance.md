---
type: "Specification"
title: "Release acceptance and verification plan"
description: "Release acceptance and verification plan."
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

# Release acceptance and verification plan

This is the proposed acceptance plan. It contains no passing product-test claims. The operator confirmed the product scope, workload and Mac-only execution testing; the precise test cases and performance thresholds are LLM-designed verification criteria.

Each result must name expected behavior, observed behavior, tested revision, command/environment, verifier and limits. Passing syntax checks does not establish runtime correctness or truth of research claims. A test marked blocked/unchecked is not passed.

## A01 Repository discovery

One repo from CWD/explicit path, multiple recursively discovered repos, worktree pointer files, explicit duplicates, invalid/uninitialized projects and partial failure diagnostics. Assert no Git mutations.

## A02 OKF and readable knowledge

Load the live bundle and a generic minimal OKF corpus; preserve unknown types/keys and bare/list verified forms. Read navigation and epistemic labels in plain Markdown. No index/log or non-Markdown runtime state in factory/.

## A03 Structural validation

Reject missing fields, wrong kinds/enums, bad placement, invalid YAML, missing local targets and anchors, source-ID errors, include cycles, case collisions and stale saves. Reject removal of a referenced heading, incompatible target-type changes and broken self-anchors before mutation; accept a valid self-link on creation against the proposed graph. Keep original bytes on failed mutations and allow unrelated-error repairs. Watcher exposes identical violations from external edits; unchecked external/cross-project links are never reported as verified.

## A04 Definitions and extensibility

Create a new type by files alone; list/detail work, board follows an enum, relationships appear as backlinks, disabling a type retains generic access. Duplicate definitions and invalid schema edits report affected records.

## A05 Template rendering

Initialize with no .factory/, render real record/prompt templates, compose nested partials, reject undeclared/missing inputs and cycles, escape hostile values, validate rendered results and preserve local customization on upgrade.

## A06 Safe editing

Create/read/update/delete through CLI, demonstrate revision conflicts, competing cooperative writers, crash-safe replacement, blocked referenced deletion, path-escape rejection and lossless-edit refusal. Explain remaining races with external writers.

## A07 MCP and HTTP

Protocol negotiation and schemas work; stdin/stdout stay clean; pagination/errors agree with CLI; tool scope forbids arbitrary roots; HTTP origin/nonce/containment checks hold; no server authentication database appears.

## A08 Native Codex conversation

On the recorded current Mac/runtime, start two fresh actions, stream output, continue one thread, answer native input/approval, interrupt, and complete a revision-checked MCP edit. Verify missing runtime/capabilities and ambiguous disconnects do not produce fabricated success or automatic duplicate turns.

## A09 Default GUI

Overview, project switch, lists/details for all types, backlinks even when empty, direct field/raw editing, kanban moves and conversation panel work. Verify keyboard journeys, validation feedback, stale-edit reconciliation, readable source/assumption labels and no persistent browser database.

## A10 Custom presentation

Use optional arbitrary HTML/CSS/JS assets after explicit trust, test missing assets/generic fallback, parent isolation, message-source/nonce checks and digest-change invalidation. Start another project without .factory/ and obtain full default functionality.

## A11 Knowledge-work breadth

Use custom research and opportunity definitions in synthetic scratch repositories without changing core code. Show cited facts, assumptions, LLM inference, backlinks, pipeline fields and agent interactions. No real bids, external contacts or inferred award facts.

## A12 Portable operation and persistence

Run desktop and CI modes on the current Mac; exercise local server shutdown/restart and an externally managed configuration example without deployment. Cross-compile supported targets, inspect linkage/embedded assets, label non-host execution untested, and verify no app Git mutations, runtime-specific config files, database, transcript copy or persistent index.

## A13 Scale and release completeness

Benchmark 100 repositories/10,000 nodes with a declared 100 MiB corpus; measure cold scan, p95 queries, incremental refresh and memory against proposed architecture thresholds. Report real hardware, exact versions and results. Revalidate the dogfood bundle and requirement/task coverage at the release revision.

## Release evidence and platform distinction

Run unit, integration and GUI tests appropriate to each completed slice. Use fixtures with deliberately invalid documents and synthetic opportunities. Do not run real procurement actions. macOS ARM64 is the observed authoring host; all other OS/architecture combinations remain unexecuted until separately tested. A build artifact may be labeled built but untested. No CI configuration should silently introduce Linux/Windows execution contrary to the initial testing scope.

Record document validation independently from product acceptance. The current [design-validation record](../sources/design-validation.md) covers this planning package only. Requirements and tasks are linked by the [implementation plan](../designs/implementation-plan.md).
