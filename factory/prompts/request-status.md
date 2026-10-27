---
type: "Prompt"
title: "Request status update"
description: "Request status update."
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
  - "Task"
inputs: {}
profile: "../agents/builder.md"
conversation: "fresh"
label: "Request status update"
---

# Request status update

This action starts a fresh conversation with the selected record. Follow-up messages continue that native conversation. The instructions are a proposed default which the repository owner can edit.

# Template

```factory-prompt
Read the selected task, its dependencies, relevant source files and available verification evidence. Report the current state and distinguish observed work from intended work. If warranted, update the task current fields/body with evidence using a freshly read revision. Do not add previous/new-status narratives or an event log. Git holds historical revisions.

Project: {{ .Project.Name }}
Record: {{ .Record.Path }}
Read revision: {{ .Record.Revision }}
Title: {{ .Record.Title }}
Use factory tools to retrieve the current document, sources and backlinks before relying on or editing it.
```

See [conversation semantics](../specifications/agent-runtime.md) and [prompt rendering](../specifications/templates.md).
