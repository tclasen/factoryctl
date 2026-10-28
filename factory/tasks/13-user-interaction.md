---
type: Task
title: "Qualify user questions and generative web views"
description: "Qualify user questions and generative web views."
status: draft
factory:
  format: 1
  assertion:
    frame: desired
    certainty: mixed
    basis: mixed
task_state: planned
priority: high
owner: unassigned
depends_on: [00-runtime-spike.md, 06-transports.md]
specifications: [../specifications/user-interaction.md]
requirements: [R30, R31]
sources:
  - id: operator
    resource: ../sources/operator-brief.md
    title: Operator requirements and scope
    author: human:operator
intent_state: proposed
intent_owner: assistant
intent_source: ../sources/operator-brief.md#project-collections-and-user-interaction
intent_authorization: MCP Apps exploration is requested; this exact protocol and implementation sequence remain proposals requiring qualification.
---

# Qualify user questions and generative web views

## Outcome and authority

Implement a bounded prototype of [user interaction](../specifications/user-interaction.md) after transports and the runtime spike. The operator requires question/answer tools and asks for an MCP Apps evaluation, potentially replacing the interactive CLI after qualification. This task is planned; no prototype exists yet.

## Plan and verification

Qualify broker attachment and answer provenance with fake runtimes, then a bounded live harness exercise. Demonstrate free text, choices with custom answers and follow-up conversation through both current interfaces. Verify unanswered requests never become consent, response routing across concurrent projects/threads, noninteractive failure, restart/disconnect and repeat submissions. Then qualify generated browser views, Apps handshake/resources/bridge or first-party host rendering, sandbox limits and fallback. Exercise all [project scenarios](../specifications/project-scenarios.md).

Produce A17/A18 evidence, an explicit capability table and an interactive-CLI retain/retire recommendation. A failed or unsupported Apps experiment remains evidence, not successful interoperability. No CLI is removed without the explicit decision described in the contract. Tasks 07–09 integrate the qualified prototype; task 11 requires current release regression evidence and video coverage.
