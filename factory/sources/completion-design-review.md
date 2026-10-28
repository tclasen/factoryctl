---
type: Reference
title: Completion revision 2 design readiness review
description: Independent source-fit and component-contract review before remaining implementation.
status: draft
factory:
  format: 1
  assertion:
    frame: observed
    certainty: supported
    basis: llm-inference
---

# Completion design review

Reviewer: delegated Reviewer `/root/completion_design_review`; run `completion-review-design`, owning [completion task](../tasks/project-completion.md). One fresh perspective, no product authorship, no parent drafting history supplied. Model/reasoning inherited by requested dispatch; actual host/model/cost telemetry and full native trace unavailable. No context marker probe run. Review is design readiness, not product acceptance or delivery approval.

## Target, scope and loaded sources

Baseline HEAD: `71c05d172e7a7a83628257acd2b2d58b3eb1f18c`. Initial dirty state contained the current plan, untracked completion task and preserved `.DS_Store` files. No product changes were present at inspection. Only this review path is Reviewer mutation scope. Coordinator owns Git/delivery.

Read: organization core, authority/scope, Reviewer instructions/context, design-review procedure/template, intent-provenance, acceptance and handoff rules; project AGENTS, factory overview, operator brief, completion plan/task; extension/interface/GUI specifications and delivery-status design; targeted template and knowledge-format sections. Inspected existing app snapshot/store/create, knowledge schema, MCP tools/server, browser web handler, generated-frame/conversation and browser bridge implementation. Product requirement/matrix delivery mappings were searched, not comprehensively accepted. No required supplied source was unavailable; the original user-message permalink is unavailable.

Pinned source SHA-256 identities at review:

| Source | SHA-256 |
| --- | --- |
| Current plan | `b8200238168cdd9ff8c38fdd6d4509ff0a065661dc915ffa86778a5aa32229a3` |
| Completion task | `9bd1b90657ad1595c728edaebe2b692ab1acf2322d9ee4adc7469b86baa195c4` |
| Operator brief | `5dc951fb9386c6d4a40a19387dfcf9c6cfd358fcebd761cfcdcaa1c2c57e4862` |
| Extensions | `a29229f7959e3e4bbe339618462171e2fa917d568d290471afd450582d32b90c` |
| Interfaces | `fab754d55a557a25c6fc9162a9eeedcbb47f53d68d08f91ea9ebfb162f7c368e` |
| GUI | `afd10864dfa18c1358055b3a4fb2a571bf8f7a514a5909e1b888a2809ee0b6ab` |
| Templates | `dfb2f787da88a62a54f77240c27eee1086cb0ed829876eb093f3bf7948dccb65` |
| Delivery design | `1a280ec8b404f8fd42695a1abf0fc1935955a7979311873c879d437a8ee0c0ec` |

## Source fit, separate from technical fit

The independently read operator brief requests a navigable editable knowledge interface, runtime declarative types, arbitrary browser presentation, harness-owned execution and no factoryctl Git/backend custom-code engine. Delivery status must remain separate from work/intent claims and support local-only project criteria. Those outcomes fit the remaining packages.

The completion instruction supplied directly in the delegation and recorded in revision 2 is: “Please implement the remaining tasks for this project until it is complete.” This is explicit implementation wording for the discussed existing remaining scope, not mere plan assent. The excerpt is agent-maintained; no original message URL is available, and this review does not authenticate historical grants or silently repair their missing evidence. Agent-selected schemas, sequencing and labels remain engineering choices. No release/media publication, v1.0 transition, real procurement, harness configuration or new workflow engine follows from this instruction.

| Lens | Outcome | Evidence/limits |
| --- | --- | --- |
| Architectural fit and simplicity | Ready | Reuse `app.Load`, `Mutate`, `Create`, existing browser host, native runtime and template renderer. No additional service, framework, database, backend interpreter or Git engine is necessary. |
| Comprehension | Ready | Desired journeys and exclusions can be restated from accessible sources. New definition while running → specialized tools/board → revision-checked edit; explicit prompt → harness conversation; delivery claims and file freshness displayed separately. |
| Execution readiness | Ready for web/extensions; delivery schema technically ready, F01 documentation/recheck pending | Existing store already locks, reloads, validates proposed graph, rechecks observed dependency bytes and publishes atomically. Definition pinning and safe browser writes are new work, not existing verified properties. Delivery schema was absent from initial brief; concrete Coordinator proposal and candidate inspected below. |

## Concrete component contracts

### Shared document service and web

Keep scope selection based on opaque checkout plus project identity and normalized bundle path. All native/raw/board/custom/extension writes converge on the same locked app writer; never add independent filesystem writes. Extend mutation input with expected effective-definition identity and definition type; checked creation must compare the caller's definition identity inside the writer lock, including template dependencies. Checking once before acquiring the lock is insufficient. Generic raw replace/delete require the same protection as specialized writes. Effective identities include applicable default/local bytes and disabled/invalid resolution; removal/replacement cannot accidentally reuse an old identity.

Fresh document requests load a current consistent snapshot under the agreed `docMu`; the browser refreshes external changes without overwriting an unsaved draft. A stale save retains the draft, exposes current content and requires reread/reconcile. Do not automatically replay uncertain writes. Lists remain bounded/paginated, deterministic, and carry snapshot identity; collection identity must survive monorepo and duplicate-ID checkouts. Backlinks always have an empty state. Invalid/disabled/unknown types retain generic readable diagnostics and appropriate repair paths.

Web's agreed `/api/v1/session` capability and `X-Factory-Nonce` protect mutations independently of whether chat/runtime is enabled. Check exact allowed Host/Origin and project scope; no permissive CORS. Bound request bodies and reject unknown envelope fields. Map stable app errors to documented status/envelope. Definition enum fields drive board columns; incompatible definitions yield separate boards; keyboard and drag use identical writes, failure restores authoritative state. Prompt listing/rendering never launches the runtime.

### Declarative extension registry

Agreed candidate shape: `extension: {version: 1, namespace: lower_snake, tools: {name: {operation: read|create|update|delete, description, inputs, fields}}}`. `inputs` use existing structural field kinds; each field binding has exactly one declared-input reference or literal value. Advertise `namespace_tool`; forbid `factory` namespace/core shadowing, duplicates and unsupported executable/expression/workflow bindings. No expressions, shell, transitions or callbacks. Strict operation envelopes reserve `path`, `slug`, `expected_revision`, `expected_definition`, `inputs`; bindings cannot override reserved selectors/tokens or target unrelated types.

Agreed app surface: `ExtensionCatalog`, `CallExtension`, `DefinitionHash`, `RenderPrompt`, plus checked creation and shared `Mutation`. Tool discovery/call reloads definitions at request time and validates complete candidate registry before publication. Invalid/duplicate affected definitions expose diagnostics and reject typed writes/specialized calls; unrelated valid types remain usable. A last-good registry must never authorize writes under newly invalid current rules. Removal/disablement removes specialized interactions without hiding old records or skipping their structural validation.

CLI, HTTP and MCP call the same catalog/dispatcher; preserve equivalent errors, scope and definition tokens. MCP request-time lists may report `reconnect_required: true`; do not claim live list-change notification or harness refresh without evidence. Prompt rendering resolves declared local Prompt nodes and context without launch; explicit host action may invoke the existing selected adapter. No MCP recursive-launch endpoint.

### Trusted repository assets

Keep generated conversation frames' answer restrictions intact; repository custom frames gain a distinct scoped grant, not host credentials. A grant binds process/project, effective definition, complete packaged HTML/CSS/JS digest and allowlisted record refs/revisions. Trust is explicit and in memory. Any dependency byte change invalidates trust before execution or edit/action acceptance. Missing/bad assets fall back to generic UI.

Use `sandbox="allow-scripts"` without same-origin privilege and restrictive CSP; package only explicit contained regular local assets, reject symlinks/traversal/special files, never serve directories or arbitrary project paths. Parent bridge checks actual source window, fresh channel nonce, digest/grant and exact message schema for `ready`, `resize`, `navigateRecord`, `requestEdit`, `requestAction`. Only granted refs and declared fields/actions are accepted. Edit uses original revision/definition tokens and common validation. Action visibly invokes a declared prompt via trusted host controls; it cannot supply arbitrary executable or harness settings. No frame can submit trusted answers or approvals.

### Delivery observations

Coordinator's concrete proposal is technically suitable with these constraints: task Markdown owns `delivery` target, criteria reference/revision and attributable assessment reference/revision/assessor/accepted claim; plan `delivery_ref` selects the same task in the same project. Optional ignored `.factory/delivery/<record-key>.json` holds versioned observation evidence, explicit target, criteria/assessment identities, manifest, observer/time and evidence references. It is deletable evidence, never a mandatory project config, database, status authority or index. Removing it means unverified. Definition-only/no-assets projects remain supported.

Resolver is read-only. Explicit nonempty manifest maps contained repository-relative paths to exact hashes or deletion markers, excludes observation store, and includes criteria/assessment resources. Validate bounds, normalized paths, symlink containment, target and record key; compare a stable snapshot or before/after manifests. Added/deleted relevant files and changed criteria/assessment invalidate freshness. Observation `state` is recorded content, not an instruction to derive Delivered. No policy prose, `task_state`, merge or assessment boolean establishes authority.

For a local-working-tree target the objective claim is content match/freshness; present assessor's delivery claim separately, e.g. “Assessed delivered locally — content verified.” Local-branch/remote-branch targets cannot be objectively verified by current working-tree hashing; label externally recorded/unverified capability honestly and never fall back to another target. Do not claim those proposed qualification scenarios pass without independent target-content evidence. No Git/network calls or Git hosting dependency. Both plan and task resolve the identical result. No publication endpoint is proposed; if added later, observation CAS identity is required so late checks cannot replace newer observations.

## Raw findings and disposition

| Label | Location/target | Severity | Finding | Disposition/owner |
| --- | --- | --- | --- | --- |
| F01 | Revision-2 delivery interfaces paragraph; initial delivery design hash above | Blocking for delivery implementation only | Plan required concrete schema inspection, while initial schema/storage/target mechanics remained unresolved. Existing `delivery-status.md` describes proposals and cannot alone pin a safe implementation. | Coordinator supplied concrete mapping/ignored observation proposal in review. Technical constraints above resolve the architecture; Coordinator must persist agreed schema and branch/remote capability disposition in the implementation brief before dependent product work. Recheck that exact brief. Independent web/extension work may proceed. |
| F02 | Revision-2 shared write interface; `internal/app/store.go` baseline | Major contract obligation | Existing writer pins document/dependency bytes but lacks caller expected-definition token; pre-lock Create derivation permits stale schema/template intent. | Resolved at design level by agreed checked Create/Mutation contracts above. Implementation tests and final review still required. |
| F03 | Asset/bridge contract; `internal/web/conversations.go` and `assets/app.js` baseline | Major contract obligation | Existing generated-view bridge deliberately cannot mutate or invoke tools; treating it as a ready repository asset bridge could leak answer/host capability or unscoped writes. | Separate repository grants and scoped messages agreed above. Preserve generated-view restrictions; implement/verify separately. |

No finding is silently waived. F02/F03 are contract requirements, not accusations that the baseline violates its then-admitted slice. F01 is initially unresolved design readiness, not a historical implementation defect. New schema/features that exceed these bounded choices return to planning.

### Narrow delivery follow-up

Coordinator subsequently supplied `internal/delivery/delivery.go` and `internal/web/delivery.go` for concrete contract inspection. Finalized shape is task `delivery: {criteria: <bundle-relative Markdown>, assessment: <bundle-relative Markdown>, target: {kind, name}}`; plan `delivery_ref` is same-project, one-hop. Assessment contains `assessor`, boolean `accepted`, `criteria_revision`. Observation version 1 is keyed by SHA-256 of record path and binds record, target, criteria/assessment revisions, nonempty files/evidence, observer/time. Branch targets explicitly remain unverified; local match label retains the recorded-claim caveat. This resolves F01's technical ambiguity; the owning implementation brief still needs to record this agreed shape and limits.

Focused candidate observations returned to Coordinator, without broad acceptance: (1) two-pass manifest checks must also reread the original plan/task and resolved delivery owner revisions before presenting a current match, because target or delivery_ref could change during checking; (2) partially populated/invalid delivery mappings must show unverified diagnostics rather than “criteria not defined” when explicit criteria exist. These are implementation follow-ups for normal tests/final review, not architecture blockers to independent work. No tests or final candidate acceptance were performed for this narrow inspection.

## Verification and handoff

Inspected sources and Git state using read-only shell commands; no product tests run because this is pre-build design review. One exploratory search referenced nonexistent `app.go`/`server.go` paths, then used discovered filenames; no mutation or retry with uncertain effects occurred. No runtime, browser, platform, performance, full-matrix, media or final-product claim is made. No external effects or Git writes. Ordinary command budget 120 seconds, tests 10 minutes; no diagnosed operation exhausted retries.

Next safe action: Coordinator persists the agreed concrete delivery brief and capability limits, obtains the narrow F01 recheck, then proceeds with delivery implementation. Web/extension workers may implement the agreed contracts now. Coordinator retains integration, independent frozen-candidate/media review and normal PR/main delivery. This review artifact is locally prepared and delivery pending; worker completion does not close the parent or establish product acceptance.
