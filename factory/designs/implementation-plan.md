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
| [Shared transports](../tasks/06-transports.md) and [full GUI](../tasks/08-gui.md) | Remaining planned capability; requires explicit build authorization for a named slice | Editing and board interfaces remain outstanding; the current slice adds conversations; reading documents does not require completed runtime integration |
| [Extensions](../tasks/12-extension-contract.md), [custom views](../tasks/09-custom-views.md), [non-software use](../tasks/10-knowledge-work.md), [release qualification](../tasks/11-release-qualification.md) | Planned and unqualified | Existing requirements remain; no full product/release acceptance is claimed |

## Delivery-status design follow-up

[Project-specific delivery status](delivery-status.md) records the accepted separation of ticket claims and delivery evidence, required web presentation and local-only support. Exact labels, observation storage and resolver/adapters remain proposals. No managed project inherits this repository's PR requirement. Implementation is not yet admitted and does not silently replace the MCP Apps exploration direction. Once implemented, ticket and plan views will share observed delivery state instead of manually maintained status summaries.

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

## Current authorized work: planning-first documentation correction

[This documentation task](../tasks/plan-before-build.md) implements the operator's explicit instruction to update the operating method and reconcile recent feature intent. Its scope is operating instructions, source excerpts, requirements/matrix, this plan and historical authorization annotations. No product build or removal is authorized by this correction. The current task stops at reviewed document delivery; future product changes require their own plan and explicit build instruction.

For fresh context, start with the [operator source](../sources/operator-brief.md#planning-first-correction), then the [audit](../sources/feature-intent-audit.md), requirement/matrix, owning task and its separate execution-authorization evidence. Missing original-message links, mappings or build grants remain explicit gaps. The application displays those records; it does not enforce this process. Prior milestone videos remain pinned to their original source and matrix snapshots, not fresh qualification of the expanded matrix.
