---
type: "Task"
title: "Document and review the runtime extension architecture"
description: "Document and review the runtime extension architecture."
status: "draft"
factory:
  format: 1
  assertion:
    frame: "desired"
    certainty: "proposal"
    basis: "llm-inference"
sources:
  - id: operator
    resource: "../sources/operator-brief.md"
    title: "Operator requirements and scope"
    author: "human:operator"
task_state: "done"
priority: "normal"
owner: "coordinator"
depends_on: []
specifications:
  - "../specifications/extensions.md"
requirements:
  - "R04"
  - "R16"
  - "R26"
  - "R27"
  - "R28"
---

# Document and review the runtime extension architecture

## Request, authority and ownership

Original request: document the agreed factoryctl architecture and merge it into the codebase. The operator also authorized moving the checkout beside the other maintained repositories; that local relocation preserves Git/content and is not project product behavior. Product authority is the operator's explicit request. This task owns only the design/specification/task documentation changes on its scoped documentation branch; unrelated files and organization governance are excluded.

## Plan and run

Revision 1: reconcile the existing design with the operator clarification, add the extension contract and acceptance criteria, update conflicting Git/validation/runtime statements, check links/metadata/diff, obtain fresh independent review of the committed candidate, publish a PR, merge under repository protections, and verify remote main content. Preserve static-linkage/platform uncertainty rather than weaken the requirement. No application implementation or runtime test is claimed.

Coordinator is the sole document writer; reviewer has a read-only claim over the pinned candidate. Required review sources: operator brief, product, extension specification, knowledge format, validation, interfaces, runtime, architecture, implementation plan and acceptance. Other files are read on contradiction/reference triggers. Review scope includes architectural fit, preserved requirements and internal consistency. Run configuration is the existing session for authoring; reviewer inherits the available model in a fresh context. Actual model/usage telemetry is unknown unless exposed by the harness. Use one bounded review and targeted repair/review cycles, no new services/providers; rollback through an ordinary corrective PR.

## Verification and delivery

Expected: all agreed constraints are discoverable in the living design, linked acceptance and planned tasks; Markdown remains navigable; no implementation claim. Author checks cover relative links, parseable metadata and diff whitespace. Independent review and PR/main verification remain required before task completion. Review findings and delivery evidence are reported with exact candidate/remote identities. The done state in this candidate is conditional on its own independent approval and verified PR/main delivery; until then it is proposed closure, not a delivered result.

## Revision 2 — operator correction and review repair

The operator explicitly rejected custom backend code execution. Remove the earlier WASM/plugin proposal rather than defer it. Retain declarative type rules, arbitrary browser UI code (including custom kanban drag/drop) and prompt-backed harness execution for local work. This supersedes the first candidate's extension mechanism without changing the documentation-only grant.

Independent reviewer `/root/review_factoryctl_extensions` reviewed candidate c21b6b881910c49201aad1db53b67e879221cc46. F01 (medium): prototype acceptance incorrectly required downstream GUI/harness work; revised task 12 and A14 now split core and integration ownership. F02 (low): baseline-to-candidate whitespace check found three blank EOF lines; repaired and changed verification to check the complete baseline diff. The initial working-tree-only check missed untracked files. This review was documentation-only; the revised scope requires a fresh pinned-candidate recheck before merge.

Author validation for revision 2: all 62 Markdown frontmatters parse with required base metadata; all 256 relative Markdown link paths resolve; baseline-to-working-tree `git diff --check 2011d25` passes. These are document checks only. PR #4 is the delivery vehicle; fresh review and verified remote-main delivery remain pending.

## Independent review and conditional closure

Fresh reviewer `/root/review_factoryctl_extensions` approved documentation candidate `22a1551c1ea86503db4d6e646206c71b45abe75d` (tree `963a91e6adbe8203e56280ec482c2c8807791179`) against base `2011d25df79ea6f32d488d68031df09a47748fc3`. Architecture, comprehension and prototype readiness passed; F01/F02 were independently resolved. The reviewer searched the candidate for active backend execution requirements and found only exclusions/historical correction references. Review coverage included the required source manifest plus templates and affected GUI/tasks. Metadata/link checks remain author-reported; runtime, linkage and browser/harness behavior remain untested. Native model/usage telemetry is unknown.

[PR #4](https://github.com/tclasen/factoryctl/pull/4) carries the reviewed design and this proposed closure. The Coordinator must recheck approval for this final evidence-only change, merge normally, verify the remote main tree equals the final candidate, and report the observed merge identity. The PR is the authoritative delivery receipt; this pre-merge record does not assert a merge has already occurred. No application implementation is completed by this task.
