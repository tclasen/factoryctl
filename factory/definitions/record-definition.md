---
type: "Record Definition"
title: "Record Definition record definition"
description: "Record Definition record definition."
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
defines: "Record Definition"
directory: "definitions"
fields:
  defines:
    kind: "string"
    required: true
  directory:
    kind: "string"
    required: true
  enabled:
    kind: "boolean"
    required: false
    default: true
  fields:
    kind: "mapping"
    required: true
  views:
    kind: "mapping"
    required: false
views:
  list:
    - "title"
    - "type"
    - "status"
---

# Record Definition record definition

This definition declares the factoryctl editor/validator contract for **Record Definition** nodes. Its structured fields are in frontmatter; [format rules](../specifications/knowledge-format.md) define their meaning. The type remains ordinary OKF and unknown metadata is preserved.

New records belong in `definitions/` within the bundle. This is a proposed version-1 definition used by the design package itself. Required/type/placement checks follow [continuous validation](../specifications/validation.md); they do not enforce workflow policy or establish factual truth.

[The factory](../factory.md) explains the desired product and current implementation state. [Epistemic conventions](../specifications/epistemics.md) apply to every instance. The software-only Specification and Design definitions are an optional starter set. Reference is used here because authoritative and observational sources need their own readable provenance nodes; infrastructure definitions describe the format rather than adding business workflow types.
