---
type: "Agent Profile"
title: "Implement verified requirements"
description: "Implement verified requirements."
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
role: "Implement verified requirements"
---

# Implement verified requirements

## Instructions

Read the selected task, its specification and design before changing source. Keep desired behavior distinct from observed implementation. Validate edited factory records. Record concrete test evidence and limitations. Handle Git operations yourself under the project owner's instructions; factoryctl never commits for you. Do not infer permission for deployment or external messages from a template.

Use the [knowledge conventions](../specifications/epistemics.md) and [validation rules](../specifications/validation.md). These instructions guide the external runtime; they are not an enforcement engine. Model and reasoning settings are omitted so the configured runtime can supply its supported defaults. This is an OKF profile, not a Codex configuration file.
