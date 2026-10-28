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

One repo from CWD/explicit path, multiple repos in immediate subfolders with default recursion limit 1, exclusion of repositories in grandchildren and deeper folders, discovery of those repositories when their parent is explicitly selected, worktree pointer files, explicit duplicates, invalid/uninitialized projects and partial failure diagnostics. Assert no factoryctl Git subprocesses or mutations and no required/generated/read factory.yaml. Identity/title come from the Factory node; generic bundles remain usable by explicit checkout selection. Check duplicate canonical URLs and metadata-only cross-project link mapping.

## A02 OKF and readable knowledge

Load the live bundle and a generic minimal OKF corpus; preserve unknown types/keys and bare/list verified forms. Read navigation and epistemic labels in plain Markdown. No index/log or non-Markdown runtime state in factory/.

## A03 Structural validation

Reject missing fields, wrong kinds/enums, bad placement, invalid YAML, missing local targets and anchors, source-ID errors, include cycles, case collisions and stale saves. Reject removal of a referenced heading, incompatible target-type changes and broken self-anchors before mutation; accept a valid self-link on creation against the proposed graph. Keep original bytes on failed mutations and allow unrelated-error repairs. Watcher exposes identical violations from external edits; unchecked external/cross-project links are never reported as verified.

## A04 Definitions and extensibility

Create a new type by files alone; list/detail work, board follows an enum, relationships appear as backlinks, setting `enabled: false` on a local definition suppresses embedded fallback while retaining generic access and structural diagnostics. Duplicate definitions and invalid schema edits report affected records.

## A05 Template rendering

Initialize with no .factory/, render real record/prompt templates, compose nested partials, reject undeclared/missing inputs and cycles, escape hostile values, validate rendered results and preserve local customization on upgrade.

## A06 Safe editing

Create/read/update/delete through CLI, demonstrate revision conflicts, competing cooperative writers, crash-safe replacement, blocked referenced deletion, path-escape rejection and lossless-edit refusal. Explain remaining races with external writers.

## A07 MCP and HTTP

Protocol negotiation and schemas work; stdin/stdout stay clean; pagination/errors agree with CLI; tool scope forbids arbitrary roots; HTTP origin/nonce/containment checks hold; no server authentication database appears.

## A08 Native Codex conversation

On the recorded current Mac/runtime, start two fresh actions, stream output, continue one thread, answer native input/approval, interrupt, and complete a revision-checked MCP edit. Resolve adapter/model/effort from invocation or optional personal preferences without writing Git; two contributor preference files selecting different adapters must not alter identical repository bytes. An unimplemented adapter returns unsupported capability rather than falling back to Codex. Reject malformed personal preferences only for runtime actions. Verify missing runtime/capabilities and ambiguous disconnects do not produce fabricated success or automatic duplicate turns.

## A09 Default GUI

Overview, project switch, lists/details for all types, backlinks even when empty, direct field/raw editing, kanban moves and conversation panel work. Verify keyboard journeys, validation feedback, stale-edit reconciliation, readable source/assumption labels and no persistent browser database.

## A10 Custom presentation

Use optional arbitrary HTML/CSS/JS assets after explicit trust, test missing assets/generic fallback, parent isolation, message-source/nonce checks and digest-change invalidation. Start another project without .factory/ and obtain full default functionality.

## A11 Knowledge-work breadth

Use custom research and opportunity definitions in synthetic scratch repositories without changing core code. Show cited facts, assumptions, LLM inference, backlinks, pipeline fields and agent interactions. No real bids, external contacts or inferred award facts.

## A12 Portable operation and persistence

Run desktop and CI modes on the current Mac; exercise local server shutdown/restart and an invocation-configured server example with optional personal preferences mounted outside the checkout without deployment. Cross-compile supported targets, inspect linkage/embedded assets, label non-host execution untested, and verify no app Git subprocesses or mutations, runtime-specific config files, database, transcript copy or persistent index.

## A13 Scale and release completeness

Benchmark 100 repositories/10,000 nodes with a declared 100 MiB corpus; measure cold scan, p95 queries, incremental refresh and memory against proposed architecture thresholds. Report real hardware, exact versions and results. Revalidate the dogfood bundle and requirement/task coverage at the release revision.

## Release evidence and platform distinction

Run unit, integration and GUI tests appropriate to each completed slice. Use fixtures with deliberately invalid documents and synthetic opportunities. Do not run real procurement actions. macOS ARM64 is the observed authoring host; all other OS/architecture combinations remain unexecuted until separately tested. A build artifact may be labeled built but untested. No CI configuration should silently introduce Linux/Windows execution contrary to the initial testing scope.

Record document validation independently from product acceptance. The current [design-validation record](../sources/design-validation.md) covers this planning package only. Requirements and tasks are linked by the [implementation plan](../designs/implementation-plan.md).

## A14 Runtime document extensions

Task 12 qualifies core behavior using in-process fixtures: load two materially different declarative types, add a third while the binary runs, verify namespaced tool mappings, enum membership, structurally required fields, definition reload, disablement, stale revisions and invalid/duplicate definitions. Reject executable backend handlers and unsupported code/expression declarations; no plugin engine is included. Shared generic/specialized write tests prove original bytes survive rejection and raw replace cannot bypass structural checks. Definition-only projects work without `.factory/`.

Downstream tasks 06–08 qualify CLI/HTTP/MCP parity, observed tool refresh or explicit reconnect with the supported harness, and prompt-backed lifecycle actions. Agent cancellation/failure is reported separately from persisted edits; no save automatically launches a harness. Verify direct writes can bypass advisory mode and are detected through revalidation; enforced mode requires observed permission denial or remains unavailable.

Task 09 qualifies arbitrary browser HTML/CSS/JavaScript with a user-defined kanban board: scoped record snapshots, drag/drop requestEdit, stale/invalid move rejection, action buttons requesting harness prompts, iframe/message checks and trust invalidation. No browser script gains direct filesystem, shell or backend code-execution access. Tasks 10–11 integrate and verify the full journey and evidence; task 12 does not depend on downstream GUI/harness completion.

Prove factoryctl launches no Git/helper shell for Git, while the harness may invoke Git or other local code. Inspect static linkage separately from single-file packaging and report platform exceptions without claiming compliance. No runtime capability or linkage result is asserted by this design document.

## A15 Minimal interface boundary

Use two factories with identical document schemas but different Markdown instructions for roles, review, readiness and workflow. Changing that prose must not change core write acceptance or require backend code changes. Reads must expose the instructions faithfully so the harness can interpret them. A valid enum-to-enum update succeeds through generic/specialized tools and board editing even when workflow prose advises against that move; an unknown enum value, missing structurally required field, bad reference or stale revision still fails. No state-dependent required-field evaluator, transition table or approval gate is executed by the core.

Show a prompt asking the harness to review a ticket and decide its next state. Factoryctl reports the conversation and observed document updates without claiming that structural validation proved approval, completion or truth. Custom browser UI still supports drag/drop through standardized writes; local arbitrary code and Git remain harness work. Reject workflow-policy handlers as tool bindings; do not misrepresent preserved unknown Markdown metadata as supported enforcement.

Task 12 owns core boundary fixtures; tasks 06–09 cover interface/harness/UI behavior and task 11 verifies complete release evidence. These are planned tests, not implemented or passing product behavior. This acceptance criterion enforces the [minimal-interface decision](../decisions/minimal-interface.md).

## A16 Project collection regression matrix

Run every stage-appropriate journey in [S1–S3](project-scenarios.md) as features land and rerun the complete matrix at release. Prove identical collection UI, depth-1 discovery, distinct monorepo project identities/definitions/state, cross-links/backlinks and target-safe edits. Include same-named records, duplicate identities, absent targets and root-plus-child bundles. Separate documentary coverage from executed test evidence.

## A17 Agent questions and dialogue

Through an actual harness tool call, ask free-text and choice questions, submit real user answers, and continue with follow-up in the same conversation. Exercise both current CLI and web contracts until an explicit CLI retirement decision. Verify custom answers, refusal/cancel/timeout, concurrent request isolation, headless failure, stale/spoofed submissions and restart/reconnect behavior. No answer or approval may be synthesized from a default or timeout. Record separate fake-runtime and live-runtime evidence.

## A18 MCP Apps and generative UI qualification

Record the [prototype](../tasks/13-user-interaction.md) results for at least two distinct runtime-generated views, host/resource/bridge interoperability, actual user-to-agent replies and plain fallback. Verify generated code is browser-only, scoped and isolated. Pin versions and actual capability observations. Produce an explicit interactive-CLI decision with gaps, parity evidence and remaining executable/MCP/automation needs. An unsupported Apps result can satisfy the honest evaluation requirement with CLI retained; it cannot pass a rich-UI capability claim.

## A19 Full-feature release video

For every release, verify a Playwright recording covering the complete implemented feature inventory, including older and still-present deprecated features. Check the [recording contract](release-evidence.md): smooth scrolling, human-readable pauses, explanatory overlays, precise interaction highlights and finished-video review. Reconcile all scenes against requirements and executable results; missing features, missing artifacts or captions standing in for unrun tests fail the gate. Only deprecated-and-removed features may be omitted with evidence.

## A20 Continuing requirements traceability

Verify every operator requirement has an atomic row in the [matrix](requirements-traceability.md), an implementation/acceptance mapping and current release evidence or an explicit unfulfilled disposition. Check no previous result is silently reused as a current pass, no required feature disappears, and release artifacts/versions/revisions/timestamps are verifiable. Failed/blocked mandatory requirements prevent acceptance absent an explicit operator exception. Keep design-only evidence distinct from implemented behavior.


## A21 Delivery-status visibility and local-only projects — planned

For R36/R37, demonstrate all scenarios in [delivery status](../designs/delivery-status.md#qualification-scenarios): post-integration verification without recursive closure commits; local branch/working-tree delivery without a remote; no inherited PR policy; explicit undefined criteria; revision changes and late observations; missing/stale evidence; consistent web ticket/current-plan display. Verify acceptance interpretation remains project/harness-owned and a merge alone cannot pass unmet criteria. Optional CLI presentation, if implemented, must share semantics. No implementation or execution pass is claimed.


## A22 Model/reasoning selection and personal defaults

For R38, demonstrate advertised model/reasoning selectors, new and same-thread follow-up selection, and explicit reset to native defaults. Verify native requests carry the resolved pair, including model switch and return. Save distinct global/project defaults, restart/reload to recover them, reset project inheritance, and prove unrelated canonical project roots (including duplicate IDs) remain isolated. Exercise precedence, unsupported pairs, malformed/missing/null stored values, stale revisions and competing process saves, private atomic persistence, symlink rejection and storage failures. Choosing or sending alone must not save. Existing chat choices and document browsing survive malformed saved settings. Future interfaces reuse the shared contract; unimplemented future surfaces are not presented as qualified. Preserve keyboard submission and existing conversation scope boundaries.
