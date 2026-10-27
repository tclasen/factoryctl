---
type: "Record Template"
title: "Specification creation template"
description: "Specification creation template."
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
for_type: "Specification"
inputs:
  title:
    kind: "string"
    required: true
  summary:
    kind: "text"
    required: true
---

# Specification creation template

This template creates a complete OKF node. The default classification is an unattributed proposal until the author supplies evidence. It never manufactures a verification event. Required fields and links are checked before creation.

# Template

```factory-template
---
type: Specification
title: {{ yaml .Input.title }}
status: draft
factory:
  format: 1
  assertion:
    frame: proposed
    certainty: proposal
    basis: unattributed
---

# {{ .Input.title }}

{{ .Input.summary }}
```

See [template rules](../specifications/templates.md).
