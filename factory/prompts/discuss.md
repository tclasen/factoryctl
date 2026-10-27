---
type: "Prompt"
title: "Discuss"
description: "Discuss."
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
applies_to:
  - "*"
inputs: {}
profile: "../agents/builder.md"
conversation: "fresh"
label: "Discuss"
---

# Discuss

This action starts a fresh conversation with the selected record. Follow-up messages continue that native conversation. The instructions are a proposed default which the repository owner can edit.

# Template

```factory-prompt
Discuss the selected record with the user. Ask for missing context where it changes the answer. Clearly separate its requirements, observed facts, assumptions and your inferences. Do not change files merely to answer a question.

Project: {{ .Project.Name }}
Record: {{ .Record.Path }}
Read revision: {{ .Record.Revision }}
Title: {{ .Record.Title }}
Use factory tools to retrieve the current document, sources and backlinks before relying on or editing it.
```

See [conversation semantics](../specifications/agent-runtime.md) and [prompt rendering](../specifications/templates.md).
