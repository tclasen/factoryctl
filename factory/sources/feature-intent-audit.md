---
type: Reference
title: Recent feature intent and build-authorization audit
description: Separate requested outcomes, traceability gaps, inferred build authority and observed delivery for three recent UI changes.
status: draft
---

# Recent feature intent and build-authorization audit

This audit responds to the operator's [planning-first correction](operator-brief.md#planning-first-correction). Its baseline is repository commit `1f1cb43755e873ae40796c856869a5a5b6283400`, before this documentation repair. Scope is the three recent UI changes below, not a complete audit of every implemented feature. It changes documentation and operating instructions only; no feature is removed or rebuilt.

## What was documented before this correction

| Change | Requested outcome and source then retained | Matrix coverage at baseline | Observed implementation/delivery | Build-authorization finding |
| --- | --- | --- | --- | --- |
| Browser document title | [Page-title task](../tasks/document-page-title.md) retained a paraphrase and the requested example. The exact message is now copied into [the source](operator-brief.md#document-title-and-keyboard-submission--recovered-source-excerpts). | No dedicated atomic row. General browser requirements did not explicitly map this request. | Candidate `2dffb4d`; delivered title change is present in [commit 35874bf](https://github.com/tclasen/factoryctl/commit/35874bf70cabd724106398a342c7c6c4ac175c98). Task's in-review field was stale relative to observed code delivery. | The task called its implementation scope authorized. That was an assistant inference from the feature request; no distinct instruction to build this change appears in the supplied exchange. |
| Keyboard chat submission | [Keyboard task](../tasks/chat-keyboard-submit.md) retained a paraphrase, the assistant-selected Enter/Shift+Enter convention and reported targeted checks. Exact request now copied into [the source](operator-brief.md#document-title-and-keyboard-submission--recovered-source-excerpts). | No dedicated atomic row. Generic keyboard-journey acceptance did not provide a request-specific mapping. | Candidate `235a90be938e0477802b35bdd9dd3b4551843163`, [PR #15](https://github.com/tclasen/factoryctl/pull/15), merge `cf92e578281299d5efa8fd2312d2ec925bf5a4b6`. | The task's “This authorizes” statement inferred immediate implementation from desired behavior. No separate build instruction for this change appears in the supplied exchange. |
| Model/reasoning selectors and personal defaults | [Operator brief](operator-brief.md#model-and-reasoning-selectors-with-personal-project-defaults) and [task](../tasks/chat-model-settings.md) retained the exact user request. Shared configuration details were labeled assistant choices. | U103/U104, R38 and A22 were added with the implementation in commit `2272aac`; they already existed at this audit's baseline. That commit alone does not prove a pre-build approval sequence. | Final candidate `13f358b00a44045d08a13b1f38e8bcb92833aa78`, [PR #16](https://github.com/tclasen/factoryctl/pull/16), merge `1f1cb43755e873ae40796c856869a5a5b6283400`; [targeted evidence](chat-model-settings-evidence.md). Future chat surfaces remain planned. | The task's “Authorized slice” treated the request as immediate build permission. No separate instruction to build that planned slice appears in the supplied exchange. Matrix coverage, review and successful native tests do not fill that gap. |

These sources support that the outcomes were requested. They do not establish permission to build immediately or blanket acceptance of every implementation choice. The operator's correction supersedes the immediate-build interpretation. Historical implementation, tests and delivery remain observed facts; the audit does not rewrite them as approved or erase them.

## Traceability repair and the hypothetical removal question

U105/R39/A23 now map the page-title outcome. U106/R40/A24 now map keyboard submission. They are explicitly retrospective additions. U103/U104 remain the existing model/default mappings. The current matrix has 106 atomic rows; earlier videos/evidence retain their original snapshots and do not acquire fresh 106-row qualification through this edit.

A strict request-specific matrix audit at the baseline would have flagged page title and keyboard submission as traceability gaps. It would not have flagged the model/default outcomes merely for lacking matrix rows, because U103/U104 existed. A traceability gap is not proof that a feature was unwanted. The source excerpts show requested intent, while the separate build-timing gap remains. The hypothetical question grants no removal authority, and this audit neither proposes automatic deletion nor claims the product has an automatic removal mechanism. Any actual removal request first updates the requirements and plan, identifies affected behavior and dependencies, and awaits explicit build authorization under the corrected method.

## Fresh-context evidence chain

A fresh agent must inspect the following, rather than inherit the previous assistant's conclusion:

1. [Operator source](operator-brief.md): exact statement or a labeled paraphrase, speaker/context, original-message link/export if available, and any source limitations. The title and keyboard quotes were recovered in this correction; the model quote predates it in the repository.
2. [Requirements matrix](../specifications/requirements-traceability.md) and [product requirements](../specifications/product.md): each outcome's atomic ID and criteria, with retrospective additions marked honestly.
3. [Current plan](../designs/implementation-plan.md) and owning task: proposed means, affected dependencies, bounded scope and current decision state.
4. **Separate explicit build instruction**, linked or quoted with its plan revision/named slice. For these three changes, no distinct build instruction is established by the supplied exchange. Do not substitute the feature request, a previous milestone instruction for other scope, task status, technical approval, merge or tests.
5. Code revision, delivery and [acceptance evidence](../specifications/acceptance.md): what was actually built and verified, separately from requested intent and build authority.

The exact excerpts here were copied by the assistant from user messages supplied in the current conversation. No original-message permalink, exported transcript or independent authenticated record is available in this checkout. Git can show when the agent-maintained record changed; it cannot independently authenticate the speaker or missing original context. This is durable, inspectable evidence with a stated limit—not conclusive proof independent of the conversation. A fresh agent must preserve that limit and must not reconstruct missing quotations or build approval.

## Corrected operating behavior

A new feature request changes the plan first. The short response describes the planned outcome and any consequential unresolved choice. Only an explicit instruction to build a named plan/slice authorizes product implementation; an express combined plan-and-build instruction permits those steps in order. Direct requests to edit instructions/planning documents authorize that document work now. Existing authorized engineering and normal reviewed delivery do not require repeated permission. The [decision](../decisions/intent-and-review.md), project AGENTS.md and [documentation task](../tasks/plan-before-build.md) carry these rules into fresh contexts.

This repair establishes clearer instructions and records. It is not a runtime approval engine, a measured reliability improvement, or a guarantee that an agent can never misinterpret a future instruction.

## Why the workflow allowed this

The evidence exposes a shared premise, not just missing paperwork: feature requests were treated as immediate execution grants, then the tasks described their own scope as authorized. The implementation and independent reviews proceeded from that interpretation. Tests and normal PR delivery verified technical behavior and delivery, while the unresolved meaning of the user's request passed through. The model/default feature even had matrix rows; traceability alone did not prevent the execution-timing mistake.

The prior demo-first preference increased pressure toward a working result. Organization instructions already separated plans and permission, but did not explicitly make feature requests planning-only; the phrase approving a “set of changes” could also collapse plan approval into build approval. Fresh contexts could then inherit that interpretation through a task's authorization label. This explains the observed three-change failure class; it is not an exhaustive diagnosis of every project mistake or proof that a particular model will always act this way.

The correction acts at admission, before implementation delegation/mutation, independent source review and handoff. Review must separately establish that the behavior was requested and that building the named plan/slice was explicitly authorized. The original source, requirement mapping and build grant remain distinct. Missing grants cannot be repaired by passing tests or adding rows afterward. Instructions and independent scenario review are procedural controls; no universal tool interception or authenticated authorization gate is implemented. Stronger mechanical enforcement would need its own separately authorized design and build.
