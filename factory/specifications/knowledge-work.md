---
type: "Specification"
title: "Research curation and contract-pursuit examples"
description: "Research curation and contract-pursuit examples."
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

# Research curation and contract-pursuit examples

**Operator-selected acceptance domains:** research/maintaining a knowledge graph, and discovering government-contract opportunities through submitting bids and tracking awards. The concrete records below are synthetic design examples, not actual opportunities, procurement advice or verified research. Factoryctl remains a metadata interface; the agent performs research, application work and external actions under the user's authority.

## Research and knowledge curation

A project-defined `Research Note` definition uses `question` (required text), `curation_state` (enum: open, researching, synthesized, reviewed), `sources` (OKF provenance), and ordinary Markdown links. A separate `Reference` node describes a source, its authority/scope and locator. The researcher's prompt distinguishes source-stated facts from inferred synthesis and unresolved assumptions. It updates current knowledge in place; Git retains earlier versions.

Journey: create a note from a template, discuss the question with the researcher profile, add two source nodes, cite their claims, link a conflicting interpretation, and inspect backlinks from the note and both sources. An unavailable source is marked unchecked. A new authoritative document can change the conclusion without rewriting a chronological log. A reviewer can inspect what was observed versus inferred entirely from Markdown.

Success: add the type/definition/template and prompt using files alone; see list/detail and an optional curation board without changing Go/JS source. A copied URL is not a verified source. The graph displays both conflicting statements and the assumption that would reconcile them until evidence settles it.

## Government-contract pursuit

A custom `Opportunity` definition can declare these fields:

| Field | Kind | Purpose |
| --- | --- | --- |
| `solicitation_id` | string, required | Issuer's published identifier |
| `issuer` | string, required | Issuing organization |
| `opportunity_state` | enum, required | discovered, qualifying, pursuing, preparing, submitted, awarded, not_awarded, withdrawn |
| `deadline` | datetime, optional | Published timestamp with explicit offset, when actually known |
| `owner` | string, optional | Accountable human/agent role |
| `tasks` | reference-list to Task | Work needed to assess/prepare/submit |
| `sources` | OKF sources | Solicitation, amendments, submission receipt and award notice |

Deadline interpretations with missing time zones remain assumptions in prose; do not invent a datetime to satisfy a form. Set board field to `opportunity_state`. Prompts include **Assess fit**, **Check amendments**, **Draft bid plan**, and **Request pursuit update**. Each opens a fresh scoped conversation; the agent researches and modifies files using its available tools. Templates do not grant permission to send a bid or contact anyone.

Journey: discover a synthetic opportunity, record the original source, qualify fit with explicitly labeled assumptions, create linked preparation tasks, record submission evidence when supplied, and update to awarded only with the relevant evidence/authorized statement. Factoryctl structurally validates records but does not enforce procurement eligibility, approve bids or decide whether an award claim is true. Separate checkouts collaborate through normal Git review and synchronization by their users/agents.

Proposal documents and bid attachments are project deliverables outside `factory/`; the graph stores their metadata and resource links. Project-source links outside the bundle may be ordinary resource references; they are never automatically imported into the OKF graph. File references outside the bundle are read-only links with explicit repo containment checks when opened, not arbitrary file-serving permissions.

## Cross-project exercise

A research project can link to a capability specification in another loaded project by its canonical bundle URI. Local mapping exposes the target and backlinks; an unloaded project remains an external/unresolved edge with an honest scope label. Factoryctl never fetches that repository or edits both checkouts as a hidden side effect.

See [domain acceptance task](../tasks/10-knowledge-work.md) and [epistemic conventions](epistemics.md).
