---
type: "Specification"
title: "Facts, assumptions, authority and inference"
description: "Facts, assumptions, authority and inference."
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
---

# Facts, assumptions, authority and inference

**Authoritative requirement:** the [operator](../sources/operator-brief.md) requires readable distinctions between facts and assumptions, authoritative information and LLM inference, desired state and current state. The following encoding is a proposed implementation of that requirement.

## Independent dimensions

`factory.assertion` supplies document defaults:

| Field | Allowed values | Meaning |
| --- | --- | --- |
| `frame` | `desired`, `observed`, `proposed`, `mixed` | Required future behavior, observed state, proposed change, or explicitly mixed document |
| `certainty` | `requirement`, `supported`, `assumption`, `proposal`, `question`, `mixed` | How the claim is being presented; supported requires cited support, not a global truth guarantee |
| `basis` | `authoritative-statement`, `direct-observation`, `llm-inference`, `mixed`, `unattributed` | Origin of the assertion, distinct from the author who typed the file |

No numeric confidence score or mechanically inferred truth label. Missing metadata displays **Unclassified**; it never defaults to fact or verified. Requirements are commitments about desired behavior, not observations about implementation. A citation proves provenance, not the validity of an inference.

## Claim-level readability

Homogeneous documents may use defaults, but their opening prose must tell a reader what kind of information follows. Mixed documents use clearly titled sections and inline labels such as **Observed**, **Operator requirement**, **Assumption**, **LLM inference**, **Proposal**, and **Open question**. Each materially different claim cites its supporting source or states the gap. Defaults do not override a local label.

Use OKF `sources` with stable `id` and `resource`. Claims cite Markdown footnotes keyed to those IDs. Include ordinary clickable links in footnote prose as well, so Git viewers work without OKF-aware rendering. Local sources are portable relative links; external primary documents use exact URLs and, when available, a revision/date locator. Source records explain the authority's scope: an operator is authoritative about desired product behavior, a solicitation issuer about its published solicitation, and a local test about the tested revision/environment. None supplies universal authority.

```markdown
**Requirement — operator:** Support Windows.[^operator]

**Observed:** The CLI test passed on macOS ARM64 at revision `<sha>`.[^test]

**Assumption — untested:** The same path handling will work on Windows.

**LLM inference:** A structured approval channel is needed for GUI interaction.[^runtime]

[^operator]: [Operator brief](../sources/operator-brief.md).
[^test]: A specific test report for the named revision, not a placeholder in a finished record.
[^runtime]: [Runtime documentation](https://openai.com/index/unlocking-the-codex-harness/).
```

The example is explanatory, not recorded test evidence. Factoryctl displays these labels and checks source references structurally. It does not decide whether prose is true, infer authority from a hostname or block work because an assumption exists. Agent curation instructions require reviewing classifications and source conflicts.

## Verification and freshness

Preserve OKF `generated`, `verified`, `stale_after` and lifecycle `status` semantics. Authorship does not constitute verification. A human typing in the GUI does not automatically create `verified`. Syntax validation does not verify knowledge. Show human-reviewed/machine-confirmed signals only from the actual metadata, with source and time; show when content is newer than its recorded verification. A stale-after signal and missing verification remain visible, without overwriting the content or inventing a review.

Git supplies edit history. Do not add status-change narratives, previous-value fields or append-only update histories. Evidence of an actual test or research observation is substantive knowledge and may be linked; it is not a duplicated history of edits. Update current-state descriptions when observations change and retain older versions through Git.

## User experiences

Detail pages show desired/observed/proposed labels, sources, assumptions and backlinks. Generic list views can filter by classification. Record previews and agent context include the classification and source references, not just the body. A status-update prompt must distinguish work observed in code/tests from plans described in tickets. A research prompt must label its derived conclusions separately from issuer/author statements.

Acceptance includes mixed-source research, contradictory sources, missing sources, inferred deadline interpretations and requirements lacking implementation. See [knowledge-work examples](knowledge-work.md) and [acceptance](acceptance.md).
