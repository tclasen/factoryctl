---
type: Task
title: Complete the remaining factoryctl requirements
description: Implement and qualify remaining project tasks under completion plan revision 2.
status: draft
task_state: in_progress
priority: high
owner: coordinator
intent_state: requested
intent_owner: human:operator
---

# Complete the remaining factoryctl requirements

Original request: “Please implement the remaining tasks for this project until it is complete.” The current session is the source; this is an agent-maintained excerpt without a message URL. Intended outcome: complete all existing mandatory product requirements and honest end-to-end evidence, independent review and owning-remote main delivery. Plan: [completion revision 2](../designs/implementation-plan.md#completion-plan-revision-2--authorized-implementation). Build authority is explicit in that request, following the discussed remaining-task inventory; planning reconciliation precedes implementation.

Target System: project. Change Classification: project. Authorized scope: this repository's product, tests, tooling and documentation plus ignored local qualification artifacts; excludes video/artifact publication, release publication, v1.0, real procurement and external messages. Canonical remote is existing origin and branch main. Preserve .DS_Store user files and all existing local evidence. No organization files are mutation scope.

## Runs and exclusive claims

- Run completion-review-design: fresh Reviewer, read-only project/source inspection and completion review evidence only; review plan revision 2 and concrete component contracts before implementation. Requested model/reasoning inherited (host instruction), actual telemetry unknown; independent execution without parent history.
- Run completion-web: Developer owns internal/web/ except delivery*.go and delivery*_test.go; GUI, document HTTP, boards/custom assets and tests. Wait for design readiness before product mutations. Coordinator integrates shared CLI hooks after agreement.
- Run completion-extensions: Developer owns internal/extensions/ (new), extension-focused internal/knowledge/ and internal/app/ changes, internal/mcp/; registry, safe bindings/prompt rendering, reload and tests. CLI changes coordinated with Coordinator. Wait for design readiness.
- Run completion-coordinator: owns current plan/task/matrix/evidence, internal/delivery/ (new), delivery-only web hooks by agreement, cmd/factoryctl/, release/performance/tooling integration and Git. No concurrent index operations; only Coordinator stages/commits/pushes. Browser capture is exclusive to Coordinator until delegated.
- Run completion-review-final: fresh independent Reviewer owns only completion review evidence after source/evidence freeze. Must verify source fit, exact candidate, tests, full matrix and actual media; no self-approval.

Dependencies: design readiness before material implementation; extension core before specialized transport/custom-board integration; source freeze before final capture and independent final review; all children and mandatory criteria before completion. Budget/retry/rollback and acceptance are in the current plan. Native IDs remain session context until run reconciliation; no made-up observed model data.

## Required context and handoff

All workers load organization core, assigned role AGENTS/CONTEXT, project AGENTS, factory/factory.md, operator-brief, current plan and this task; then only specifications governing their package. Preserve thin structure-versus-harness boundary. Return exact paths, tests/results, failures, remaining gaps and next action. Each implementation handoff remains intermediate until independent review and verified PR/main delivery. No status edit can replace proof.


## Delivery observation implementation brief and reviewed contract

Fresh review is [completion design review](../sources/completion-design-review.md). Web/extensions are ready under its recorded contracts. Delivery's exact schema is a task frontmatter `delivery` mapping with `criteria` and `assessment` as bundle-relative Markdown paths and `target: {kind, name}`. A plan uses `delivery_ref` as a same-project bundle-relative task path; chained references are unsupported. Assessment Markdown records `assessor`, boolean `accepted` and `criteria_revision` hash. These are attributed claims, not programmatically verified authority.

Optional disposable ignored evidence is `.factory/delivery/<sha256 of owning task path>.json`, version 1, containing `record`, exact `target`, `criteria_revision`, `assessment_revision`, nonempty `files` mapping repository-relative paths to SHA-256 or `deleted`, RFC3339 `observed_at`, `observer` and nonempty `evidence` references. Criteria and assessment content must be included. No state field is trusted, no store inclusion or Git metadata is allowed, and regular contained files only are read with size bounds. Two-pass manifest verification and reread of owner/plan records detect concurrent changes. There is no publication endpoint, so CAS publication is not implemented or claimed. Users/agents may generate this optional evidence; deleting it leaves delivery unverified.

For explicit local-working-tree targets the UI reports the recorded acceptance claim and scoped content match with provenance, not inferred approval. Local-branch/remote-branch targets remain externally recorded and unverified by factoryctl: current working-tree bytes do not prove another target. This is an explicit capability limitation requiring qualification disposition, not fallback to local delivery. No root config, database, index, Git process, network query or backend policy engine is introduced. Original task and plan mappings must remain unchanged throughout verification; incomplete mappings are unverified diagnostics.

Run completion-qualification: Developer owns only scripts/qualification/ (new), release/version guard tooling plus .github workflow changes if necessary and corresponding tests; performance/breadth/portability qualification fixtures and evidence. No existing product source, matrix, current task, Git or media mutations. Coordinator owns final browser capture and integration. This run follows independent design readiness; additive local tooling is a routine verification implementation, no new product interface.


Run completion-runtime: Developer owns only a new portable scripts/review-video/completion-terminal.py probe and local ignored artifacts/completion-runtime/ versioned outputs. Qualify the current frozen binary through a real PTY/native Codex (free text, custom choices, follow-up same thread, interruption, explicit CLI retain disposition), sanitizing runtime identity traces and preserving failures. No product source, project Markdown, Git or media publication. Coordinator alone owns browser captures and final video combination. This is remaining A08/A17 verification under the explicit completion build instruction, not a new runtime/service grant. Runtime conversation persistence remains native; personal test preferences are isolated scratch files. Exact source hashes and binary digest bind evidence; live runtime version and host recorded, non-host execution never claimed.
