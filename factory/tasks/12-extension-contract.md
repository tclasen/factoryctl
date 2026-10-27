---
type: "Task"
title: "Qualify runtime document extensions before interface implementation"
description: "Qualify runtime document extensions before interface implementation."
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
task_state: "planned"
priority: "high"
owner: "unassigned"
depends_on:
  - "03-definition-validation.md"
specifications:
  - "../specifications/extensions.md"
requirements:
  - "R26"
  - "R27"
  - "R28"
---

# Qualify runtime document extensions before interface implementation

## Outcome and scope

Implement and qualify the proposed [extension contract](../specifications/extensions.md), including declarative manifests and an embedded pure-Go WASM engine, before transport/UI consumers depend on the API. This is planned implementation, not part of the current documentation update.

## Sequence and acceptance

After foundation/storage/structural validation, define the exact versioned manifest/ABI, host capabilities, namespaced tool schemas, registry replacement and operation failure envelopes. Exercise two materially different packages and install a third without rebuilding the binary. Prove shared generic/specialized write enforcement and bounded hooks before extending MCP/GUI. Use core-rendered views first; trusted custom HTML/CSS/JS remains task 09.

Acceptance is A14 in [release acceptance](../specifications/acceptance.md), plus regressions for A03/A04/A06. Include traps, timeouts, memory/output limits, invalid/missing packages, stale registry/document revisions, partial post-hook failure, external-editor bypass and tool refresh. No Git subprocess or native plugin helper is permitted. Select/pin an engine only after testing static-build compatibility and interruption behavior.

## Risks and rollback

Keep exact package syntax/ABI marked proposed until the prototype passes. No background hook scheduler, durable event ledger or multi-file transaction is introduced. Failed engine qualification requires a documented design revision; never silently substitute an external interpreter. Rollback is an ordinary agent-owned Git change, preserving user documents and revealing unavailable extension behavior.

