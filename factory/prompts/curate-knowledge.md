---
type: "Prompt"
title: "Curate knowledge"
description: "Curate knowledge."
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
  - "Reference"
  - "Specification"
inputs: {}
profile: "../agents/researcher.md"
conversation: "fresh"
label: "Curate knowledge"
---

# Curate knowledge

This action starts a fresh conversation with the selected record. Follow-up messages continue that native conversation. The instructions are a proposed default which the repository owner can edit.

# Template

```factory-prompt
Review this knowledge against its sources. Identify stale, missing and conflicting evidence. Label source-stated claims, assumptions and your derived conclusions. Propose or perform scoped updates requested by the user; do not fabricate verification or execute external submissions.

Project: {{ .Project.Name }}
Record: {{ .Record.Path }}
Read revision: {{ .Record.Revision }}
Title: {{ .Record.Title }}
Use factory tools to retrieve the current document, sources and backlinks before relying on or editing it.
```

See [conversation semantics](../specifications/agent-runtime.md) and [prompt rendering](../specifications/templates.md).
