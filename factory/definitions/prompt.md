---
type: "Record Definition"
title: "Prompt record definition"
description: "Prompt record definition."
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
defines: "Prompt"
directory: "prompts"
fields:
  applies_to:
    kind: "string-list"
    required: true
    min_items: 1
  inputs:
    kind: "mapping"
    required: true
  profile:
    kind: "reference"
    required: false
    target_types:
      - "Agent Profile"
  conversation:
    kind: "enum"
    required: true
    values:
      - "fresh"
  label:
    kind: "string"
    required: false
  asset:
    kind: "string"
    required: false
views:
  list:
    - "title"
    - "type"
    - "status"
---

# Prompt record definition

This definition declares the factoryctl editor/validator contract for **Prompt** nodes. Its structured fields are in frontmatter; [format rules](../specifications/knowledge-format.md) define their meaning. The type remains ordinary OKF and unknown metadata is preserved.

New records belong in `prompts/` within the bundle. This is a proposed version-1 definition used by the design package itself. Required/type/placement checks follow [continuous validation](../specifications/validation.md); they do not enforce workflow policy or establish factual truth.

[The factory](../factory.md) explains the desired product and current implementation state. [Epistemic conventions](../specifications/epistemics.md) apply to every instance. The software-only Specification and Design definitions are an optional starter set. Reference is used here because authoritative and observational sources need their own readable provenance nodes; infrastructure definitions describe the format rather than adding business workflow types.
