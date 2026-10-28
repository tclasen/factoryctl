---
type: "Design"
title: "Current implementation plan"
description: "Implementation sequence and requirements coverage."
status: "draft"
factory:
  format: 1
  assertion:
    frame: "desired"
    certainty: "mixed"
    basis: "llm-inference"
sources:
  -
    id: "operator"
    resource: "../sources/operator-brief.md"
    title: "Operator requirements and scope"
    author: "human:operator"
---

# Current implementation plan

This is the single current desired-outcome and implementation plan. A plan entry is not a build grant. It separates accepted working-method changes, observed implementation and assistant-selected next-step proposals. The [operator source](../sources/operator-brief.md#working-method-correction-and-tangible-review) records the actual instruction and its context. A source attribution does not make every sentence in this plan operator-authored.

## Delivered browser foundation

[Visible intent](../tasks/visible-intent.md) applies the accepted working-method correction and delivers a small read-only browser explorer of the real planning bundle. You can find a document, inspect its recorded intent state, open source/authorization details and linked documents, and inspect the underlying Markdown. The demo uses actual product interactions and remains at one stable local video path.

The operator authorized implementing the working-method proposals and asked for something tangible. The assistant chose this bounded browser slice as the implementation approach. Layout, launch syntax and data representation are engineering choices. The slice does not silently accept future MCP/GUI designs or retire all terminal launch/automation commands.

## Previously admitted milestone and observed delivery

The operator requested completion of the described agent-questions and rich-web milestone. [Task13](../tasks/13-user-interaction.md#current-admission-and-implementation-brief) owns the bounded build and runtime/transport qualification: question or generated view → actual response → correlated result → same-conversation follow-up. Native attachment, browser answers/generated views and same-thread follow-up have current bounded evidence; see [the review package](../sources/interaction-milestone-evidence.md). This is not authorization to implement every remaining full-GUI or delivery-status feature, publish media, change harness configuration or retire the CLI.

## Current work and proposed continuation

| Work | Status and authority | User-visible result / boundary |
| --- | --- | --- |
| [Intent and tangible review](../decisions/intent-and-review.md) | Accepted method; current task applies it | Requests, agent choices and execution authorization remain separate; demos precede long completion prose |
| [Read-only browser explorer](../tasks/visible-intent.md) | Implemented in [PR #11](https://github.com/tclasen/factoryctl/pull/11); assistant-selected reversible means (historical delivery observation, not a live status) | Navigate real documents and recorded provenance without editing or submitting acceptance |
| [Existing document services](../tasks/first-cli-milestone.md) | Implemented bounded CLI work; sequence was assistant-selected | Reuse discovery, parsing, validation and safe storage; no further CLI feature milestone selected |
| [Native runtime qualification](../tasks/00-runtime-spike.md) and [MCP Apps/questions](../tasks/13-user-interaction.md) | Implemented bounded slice in review under the operator request; full-product qualification remains incomplete | Generated view/question → actual response → correlated tool result → same-conversation follow-up; native Codex attachment and first-party views have live evidence; external Apps interoperability remains unqualified |
| [Shared transports](../tasks/06-transports.md) and [full GUI](../tasks/08-gui.md) | Admitted under completion revision 2; implementation in qualification | Revision 2 adds editing, boards, document refresh and custom views; full evidence remains pending |
| [Extensions](../tasks/12-extension-contract.md), [custom views](../tasks/09-custom-views.md), [non-software use](../tasks/10-knowledge-work.md), [release qualification](../tasks/11-release-qualification.md) | Planned and unqualified | Existing requirements remain; no full product/release acceptance is claimed |

## Delivery-status design follow-up

[Project-specific delivery status](delivery-status.md) records the accepted separation of ticket claims and delivery evidence, required web presentation and local-only support. Exact labels, observation storage and resolver/adapters remain proposals. No managed project inherits this repository's PR requirement. Implementation is admitted under completion revision 2 without replacing the optional MCP Apps interoperability experiment. The current candidate shares a read-model between ticket and plan views; recorded assessments and objective content freshness remain distinct.

## Superseded ordering

The original plan said “The first useful slice is a CLI” and made transports depend on completing that command interface. That was an LLM proposal, not an operator requirement. This ordering is superseded by user-visible vertical slices. Historical rationale and implementation remain in Git and [the prior sequence record](../decisions/cli-first-sequence.md); completed CLI work is not retroactively unauthorized or undone.

The later paragraph describing a “revised first vertical slice” did not reconcile the older table. It is replaced by this current plan. Shared storage and document contracts can be reused without making completion of every CLI or conversation feature a prerequisite for a read-only browser. Full runtime questions and MCP Apps still require real transport/host qualification.

## Checkpoint and review

For a new or changed feature, first update the desired outcome, source/requirement/matrix mapping, acceptance criteria and affected task/dependencies here. Present that planned outcome in one to three sentences. The feature request does not authorize a prototype or implementation. Build only after an explicit instruction binds the named plan revision or clear slice; an instruction explicitly asking for both plan and build authorizes those stages in order. Plan approval alone is not build permission. Continue routine engineering and required delivery inside valid existing build scope without repeat confirmation. Resolve only material ambiguity; keep detailed planning available behind the short checkpoint.

After authorized product implementation, lead its review with a working demo and local annotated video, then the decision that experience informs. A documentation-only planning correction is reviewed as documentation and does not authorize building a demo. Keep raw Markdown, source references, accessible text and full matrix gaps available. A prototype approval applies to the demonstrated scope; it does not accept unread options or future behavior. Each reviewer checks outcome fit and authority from source instructions before implementation correctness.

## Acceptance and remaining uncertainty

The [requirements matrix](../specifications/requirements-traceability.md) is the detailed evidence map. A readable document or recorded “accepted” label is not proof that all of its claims were accepted. Unknown authority remains unresolved. The application displays source and authorization claims; it does not implement a policy or approval engine.

Current risks: actual Codex/MCP Apps capability, generated-view isolation, runtime answer routing, full safe-edit integration, portability/performance and complete matrix qualification. Preserve S1–S3 project collection requirements. Full-product completion still requires [acceptance](../specifications/acceptance.md); a partial browser slice and successful tests do not establish it. Changes to the active plan must identify superseded choices and affected dependencies instead of appending a competing roadmap.


## Observed recent changes and authorization correction

The requested page title, keyboard submission and model/default selectors are implemented. [The scoped audit](../sources/feature-intent-audit.md) records their source, matrix coverage before this correction, code/delivery evidence and absence of a distinct build instruction in the supplied exchange. U103/U104 already mapped model/default outcomes; U105/U106 now repair missing title/keyboard rows retrospectively. This is not permission to remove existing behavior or to reinterpret tests/merges as build authorization. Future chat surfaces remain planned; the shared model-settings contract does not implement them.

## Previous authorized work: planning-first documentation correction

[This documentation task](../tasks/plan-before-build.md) implements the operator's explicit instruction to update the operating method and reconcile recent feature intent. Its scope is operating instructions, source excerpts, requirements/matrix, this plan and historical authorization annotations. No product build or removal is authorized by this correction. The current task stops at reviewed document delivery; future product changes require their own plan and explicit build instruction.

For fresh context, start with the [operator source](../sources/operator-brief.md#planning-first-correction), then the [audit](../sources/feature-intent-audit.md), requirement/matrix, owning task and its separate execution-authorization evidence. Missing original-message links, mappings or build grants remain explicit gaps. The application displays those records; it does not enforce this process. Prior milestone videos remain pinned to their original source and matrix snapshots, not fresh qualification of the expanded matrix.


## Completion plan revision 2 — authorized implementation

Recorded operator instruction in the current session: “Please implement the remaining tasks for this project until it is complete.” This follows the inventory of 28 task records (8 done, 8 planned, 5 in progress, 7 in review). This is explicit build authority for the remaining existing requirements and tasks, with planning reconciliation first. It does not authenticate missing historical authorizations. The source is this agent-maintained excerpt; no original-message URL is available.

The [completion task](../tasks/project-completion.md) owns this revision, runs and claims. This revision supersedes the previous statement that no remaining implementation is admitted. Existing task metadata is an inventory, not proof of missing or completed behavior. Preserve completed features and qualify outstanding portions against U001–U106 and acceptance A01–A25.

### Outcome and work packages

1. Shared document transport and GUI: reuse app.Load/Mutate/Create and the existing browser/conversation host; provide fresh collection lists, generic details/backlinks, raw and field editing, create/delete, definition-derived keyboard/drag boards, validation and stale-edit recovery. Refresh external changes without restart. Preserve provenance, page title, keyboard chat and personal preferences.
2. Declarative extensions and presentation: parse namespaced read/create/update/delete bindings in definitions, reject backend executable handlers, reload definitions, pin definition identities on writes, expose equivalent CLI/MCP/HTTP discovery and calls, render prompt actions without implicit launch. Add explicitly trusted digest-bound repository HTML/CSS/JS frames with scoped edit/action messages, safe asset containment, fallback and no backend execution.
3. Delivery observations: present project-owned criteria/assessments and revision-bound evidence in document and plan views, support local-only targets with no remote or hosting dependency, report undefined/missing/stale evidence honestly. Content hashes are objective facts; authority and acceptance stay recorded claims. No policy evaluator or Git subprocess is added.
4. Qualification: reconcile tasks 00–13 and design follow-ups by actual evidence, add synthetic research/opportunity fixtures, run regression and negative tests, benchmark 100 repositories/10,000 nodes, cross-build six OS/architecture targets, verify static packaging, conventional-commit/version guards and repository protections. No unsupported-platform execution or release publication is authorized. Initial release candidate remains 0.1.0.
5. Full review and delivery: update current feature/matrix evidence, record one current complete Playwright demo with annotated scenes for all implemented features, retain honest gaps and failed attempts, independently review candidate/media, run accounting and qualification gates, deliver via reviewed PR/main, verify canonical local video hashes. Completion requires all mandatory criteria, not mere row accounting.

### Interfaces, dependencies and choices

The existing standard-library HTTP/embedded JS implementation remains; no frontend framework or backend plugin runtime is introduced. Generic writes use the existing locked revision-checked app service. UI and extension workers coordinate shared Handler hooks before integration. Coordinator owns cross-package integration, delivery observations, current plan/task/matrix, Git and final evidence. Backend extension and browser UI packages have disjoint source ownership documented in the task. Fresh design review precedes material interfaces; independent implementation review follows source freeze.

Delivery format is an agent-selected bounded mechanism: project Markdown holds attributable criteria and assessment claims; observations outside factory/ bind those claims to scoped content hashes and explicit targets. No observation can infer approval from task_state, a merge or a test. Missing evidence remains unverified. Inspect the concrete schema in the implementation brief before adopting it.

### Verification, risks and rollback

Run focused meaningful tests during development then make check/build/repository-check, current browser journeys and review gates. Protect Host/Origin/nonces, project scope, symlinks, stale document/definition revisions, custom-frame source/nonce/digest and answer provenance. Maintain S1–S3 collection parity. Measure performance before optimization. Existing external MCP Apps remains an optional interoperability experiment with explicit capability disposition; do not claim first-party rendering proves external Apps support or remove the CLI.

Resources: use the currently supported collaboration capacity with independent review; ordinary commands 120 seconds, tests 10 minutes, checkpoint at transitions and every 15 active minutes; two diagnosed retries per operation. Actual model/cost telemetry is unknown. Local artifacts stay ignored and never uploaded. Rollback is scoped revert PRs, preserving user files, prior artifacts and failed evidence. No deployments, real procurement/contact, new paid providers, harness-wide configuration, protection bypass, release publication or v1.0 transition.
