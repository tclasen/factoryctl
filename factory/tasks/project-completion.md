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


Run completion-collections: Developer verification owns only new scripts/review-video/completion-collections.mjs and versioned ignored artifacts/completion-collections/ outputs. Complete A16's same-control GUI navigation, target-safe field/raw edits, stale/path-negative cases, actual conversations and agent questions in S1–S3, including independent same-named project records, canonical cross-links/backlinks, per-project definitions and empty cross-project state in S1. Use the already frozen candidate; no product changes, Git or Markdown mutation. Actual questions/answers are automated synthetic fixture input, never operator acceptance. Existing actual native/source tests remain distinct from these extra topology-specific journeys. Coordinator owns final combined film and matrix; worker returns exact source/binary hashes, commands and scene times, preserving failures.

## Current checkpoint

Product sources are frozen at a7c22ba847d58e8f4ece6a22ee9bb73302a64c02. Fresh full named race/vet/build and all six target packaging checks passed. Real native MCP Research edit, terminal dialogue, generated browser form/chart, concurrent answer routing, model/default persistence and isolated native approval cancellation have current evidence. The original approval-only recorder's termination failure is retained separately from the independently verified request/reply cancellation criterion. Native MCP elicitation variant remains unsupported and explicit.

The source-stable declared-workload benchmark is complete; its proposed one-second incremental-refresh target is not attained (about 2.5 seconds). This is a reported engineering target miss, not an invented operator-approved latency requirement. Static linkage is observed for Linux, while Darwin/Windows use OS libraries; the exact interpretation of the operator's single-static-binary requirement is pending an explicit clarification. No full-product completion or release is claimed while that boundary and final combined evidence/review/delivery remain pending.


Retained local qualification artifacts and synthetic corpora exceed the 1 GiB soft active-evidence threshold in aggregate. Proposed cleanup after review is checksum-verified compression of superseded synthetic corpora and raw recording attempts into named local archives, preserving recovery and the final chosen evidence. No evidence has been silently removed. The final selected review package remains separate from earlier failures and successful historical candidates.


## Reviewed delivery checkpoint

The authorized implementation is on `feat/project-completion`; [PR 18](https://github.com/tclasen/factoryctl/pull/18) is the normal delivery path. Initial current-head CI passes repository-content and release-policy checks. Full-history and index guards exclude every generated artifact/recording; the existing origin is the same public project repository with write access. A first automatic push review rejected unverified egress; after read-only origin/visibility/authority verification the normal push was approved and succeeded. No alternate route or protection bypass was used.

Independent review approves the strengthened source-bound qualification-test discovery, source/build fit and living source/plan reconciliation. The reviewer requested actual custom drag/drop with stale/invalid recovery and the shipped GUI Initialize journey; final proof must cover those before closing the custom-view work. Empty discoverable `factory/` is the supported GUI Initialize precondition; no-bundle GUI launch is unsupported and its failed probe remains preserved. Native retries are complete; no more native execution is needed for these browser-only checks.

The 19 bounded child task closures are proposals effective only after independent exact-candidate review, ordinary PR merge and verified owning-remote main. Release qualification and this parent remain nonterminal. U045 and full-product qualification remain unresolved until the operator actually answers the requested static-linkage clarification. No required criterion is waived by a task-state edit or a passing accounting check.


The review-requested custom qualification now passes in `artifacts/completion-custom-drag/attempt-002/`: actual GUI Initialize from empty factory bundle; actual iframe drag/drop with observed revision, 409 stale rejection/restoration, 422 invalid-enum rejection/restoration and valid recovery. Nineteen non-target knowledge documents and custom assets remain byte-identical; all 75 product files, frozen binary and executed recorder are unchanged. Both owned servers are stopped; no native runtime was started. The failed no-bundle GUI launch remains in attempt 001. Final combined assembly adds this chapter without cutting or speeding source clips.


Final selected local package: `artifacts/completion-review/a7c22ba-final-v2/`, 99 scenes / 1219.00 seconds, movie SHA-256 `51e225e9ff403a9bc015fea0e80c967053082afaf4bc294a8516014e4e3b40ed`. Independent reviewer approved the actual custom-drag/init proof and v2 visuals for a prominently partial package. Canonical demo and immutable archive now match; the prior canonical file was preserved by content hash. Accounting passes every current row with no evaluator failures. Full qualification correctly fails only U045 (105 pass / 1 blocked); no exception or full-completion claim is inferred. The current source-bound evidence and living records proceed through exact-candidate review and PR 18, then verified remote-main delivery. This parent remains in progress until U045 is actually resolved and remaining qualification passes.
