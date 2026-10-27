---
type: "Record Definition"
title: "Task record definition"
description: "Task record definition."
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
defines: "Task"
directory: "tasks"
fields:
  task_state:
    kind: "enum"
    required: true
    values:
      - "planned"
      - "ready"
      - "in_progress"
      - "blocked"
      - "in_review"
      - "done"
      - "cancelled"
    default: "planned"
  priority:
    kind: "enum"
    required: true
    values:
      - "high"
      - "normal"
      - "low"
    default: "normal"
  owner:
    kind: "string"
    required: false
  depends_on:
    kind: "reference-list"
    required: false
    target_types:
      - "Task"
  specifications:
    kind: "reference-list"
    required: false
    target_types:
      - "Specification"
  requirements:
    kind: "string-list"
    required: false
views:
  list:
    - "title"
    - "task_state"
    - "priority"
    - "owner"
  board:
    field: "task_state"
template: "../templates/task.md"
actions:
  - "../prompts/discuss.md"
  - "../prompts/request-status.md"
---

# Task record definition

This definition declares the factoryctl editor/validator contract for **Task** nodes. Its structured fields are in frontmatter; [format rules](../specifications/knowledge-format.md) define their meaning. The type remains ordinary OKF and unknown metadata is preserved.

New records belong in `tasks/` within the bundle. This is a proposed version-1 definition used by the design package itself. Required/type/placement checks follow [continuous validation](../specifications/validation.md); they do not enforce workflow policy or establish factual truth.

[The factory](../factory.md) explains the desired product and current implementation state. [Epistemic conventions](../specifications/epistemics.md) apply to every instance. The software-only Specification and Design definitions are an optional starter set. Reference is used here because authoritative and observational sources need their own readable provenance nodes; infrastructure definitions describe the format rather than adding business workflow types.
