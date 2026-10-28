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

This is the single current build plan. It separates accepted working-method changes, observed implementation and assistant-selected next-step proposals. The [operator source](../sources/operator-brief.md#working-method-correction-and-tangible-review) records the actual instruction and its context. A source attribution does not make every sentence in this plan operator-authored.

## Current tangible slice

[Visible intent](../tasks/visible-intent.md) applies the accepted working-method correction and delivers a small read-only browser explorer of the real planning bundle. You can find a document, inspect its recorded intent state, open source/authorization details and linked documents, and inspect the underlying Markdown. The demo uses actual product interactions and remains at one stable local video path.

The operator authorized implementing the working-method proposals and asked for something tangible. The assistant chose this bounded browser slice as the implementation approach. Layout, launch syntax and data representation are engineering choices. The slice does not silently accept future MCP/GUI designs or retire all terminal launch/automation commands.

## Current work and proposed continuation

| Work | Status and authority | User-visible result / boundary |
| --- | --- | --- |
| [Intent and tangible review](../decisions/intent-and-review.md) | Accepted method; current task applies it | Requests, agent choices and execution authorization remain separate; demos precede long completion prose |
| [Read-only browser explorer](../tasks/visible-intent.md) | In progress; assistant-selected reversible means under the current implementation request | Navigate real documents and recorded provenance without editing or submitting acceptance |
| [Existing document services](../tasks/first-cli-milestone.md) | Implemented bounded CLI work; sequence was assistant-selected | Reuse discovery, parsing, validation and safe storage; no further CLI feature milestone selected |
| [Native runtime qualification](../tasks/00-runtime-spike.md) and [MCP Apps/questions](../tasks/13-user-interaction.md) | Required direction to explore; exact next implementation scope is proposed | Generated view/question → actual response → correlated tool result → same-conversation follow-up; installed host support remains unverified |
| [Shared transports](../tasks/06-transports.md) and [full GUI](../tasks/08-gui.md) | Remaining capability scope; implement only the services needed by an admitted vertical slice | Editing, board and conversation interfaces remain outstanding; reading documents does not require completed runtime integration |
| [Extensions](../tasks/12-extension-contract.md), [custom views](../tasks/09-custom-views.md), [non-software use](../tasks/10-knowledge-work.md), [release qualification](../tasks/11-release-qualification.md) | Planned and unqualified | Existing requirements remain; no full product/release acceptance is claimed |

## Superseded ordering

The original plan said “The first useful slice is a CLI” and made transports depend on completing that command interface. That was an LLM proposal, not an operator requirement. This ordering is superseded by user-visible vertical slices. Historical rationale and implementation remain in Git and [the prior sequence record](../decisions/cli-first-sequence.md); completed CLI work is not retroactively unauthorized or undone.

The later paragraph describing a “revised first vertical slice” did not reconcile the older table. It is replaced by this current plan. Shared storage and document contracts can be reused without making completion of every CLI or conversation feature a prerequisite for a read-only browser. Full runtime questions and MCP Apps still require real transport/host qualification.

## Checkpoint and review

Before each admitted slice, state the outcome, why it is next, and any unresolved consequential choice in one to three sentences. If the request clearly authorizes a small reversible local prototype, build it without asking again about routine details. If alternatives materially change product scope, exposure, cost or irreversible effects and the user's short assent does not resolve them, ask only that blocking question; detail belongs behind the short checkpoint.

At review, lead with a working demo and local annotated video, then the decision that experience informs. Keep raw Markdown, source references, accessible text and full matrix gaps available. A prototype approval applies to the demonstrated scope; it does not accept unread options or future behavior. Each reviewer checks outcome fit and authority from source instructions before implementation correctness.

## Acceptance and remaining uncertainty

The [requirements matrix](../specifications/requirements-traceability.md) is the detailed evidence map. A readable document or recorded “accepted” label is not proof that all of its claims were accepted. Unknown authority remains unresolved. The application displays source and authorization claims; it does not implement a policy or approval engine.

Current risks: actual Codex/MCP Apps capability, generated-view isolation, runtime answer routing, full safe-edit integration, portability/performance and complete matrix qualification. Preserve S1–S3 project collection requirements. Full-product completion still requires [acceptance](../specifications/acceptance.md); a partial browser slice and successful tests do not establish it. Changes to the active plan must identify superseded choices and affected dependencies instead of appending a competing roadmap.
