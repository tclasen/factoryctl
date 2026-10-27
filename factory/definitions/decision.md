---
type: "Record Definition"
title: "Decision record definition"
description: "Decision record definition."
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
defines: "Decision"
directory: "decisions"
fields:
  decision_state:
    kind: "enum"
    required: true
    values:
      - "proposed"
      - "accepted"
      - "superseded"
    default: "proposed"
  supersedes:
    kind: "reference-list"
    required: false
    target_types:
      - "Decision"
views:
  list:
    - "title"
    - "type"
    - "status"
---

# Decision record definition

This definition declares the factoryctl editor/validator contract for **Decision** nodes. Its structured fields are in frontmatter; [format rules](../specifications/knowledge-format.md) define their meaning. The type remains ordinary OKF and unknown metadata is preserved.

New records belong in `decisions/` within the bundle. This is a proposed version-1 definition used by the design package itself. Required/type/placement checks follow [continuous validation](../specifications/validation.md); they do not enforce workflow policy or establish factual truth.

[The factory](../factory.md) explains the desired product and current implementation state. [Epistemic conventions](../specifications/epistemics.md) apply to every instance. The software-only Specification and Design definitions are an optional starter set. Reference is used here because authoritative and observational sources need their own readable provenance nodes; infrastructure definitions describe the format rather than adding business workflow types.
