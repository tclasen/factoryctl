---
type: Design
title: Project-specific delivery status
description: Separate recorded plans from revision-bound delivery observations for local and remote projects.
status: draft
factory:
  format: 1
  assertion:
    frame: desired
    certainty: mixed
    basis: mixed
intent_state: proposed
intent_owner: assistant
intent_source: ../sources/operator-brief.md#delivery-status-and-project-specific-criteria
intent_authorization: Document the accepted delivery-status direction; exact schema, storage and resolver mechanics remain proposals, not shipped behavior.
sources:
  - id: operator
    resource: ../sources/operator-brief.md#delivery-status-and-project-specific-criteria
    title: Delivery-status discussion and correction
    author: human:operator
---

# Project-specific delivery status

## Accepted direction and boundary

The web UI must display delivery status separately from a ticket's recorded planning/work state. CLI presentation is optional; if retained, both surfaces use the same observations and semantics. A task field, merge, test pass or agent statement alone cannot establish delivery. No status resolver is implemented in the current read-only explorer.

PRs are required for developing this repository. They are not required of projects managed by factoryctl. A project can be local-only, with no remote, hosting provider, issue tracker or PR. Delivery criteria belong to that project and its authority; never inherit this repository's development rules. With no explicit criteria, show **Delivery criteria not defined**, not an invented workflow or a delivered badge. This does not remove the existing Git-working-tree project-discovery requirement.

Keep the [minimal-interface boundary](../decisions/minimal-interface.md): the harness/project authority interprets acceptance policy and records attributable assessments; factoryctl presents records and objective revision/freshness facts. Do not introduce a backend policy language, workflow engine, approval inference or Git subprocesses. A shared read-model/resolver is a proposed presentation mechanism, not permission to evaluate arbitrary policy prose. Remote adapters, authorization, refresh and local Git inspection mechanisms need qualification before implementation.

## Proposed record separation

| Record | Contents | Authority |
| --- | --- | --- |
| Ticket/plan | Outcome, acceptance criteria, implementation evidence references and recorded work state | Project authors and applicable decision authority |
| Acceptance assessment | Assessor/source, criteria revision, exact candidate/content manifest, results and scope | Project-defined authority; never inferred from structural validation |
| Delivery observation | Target identity, observed revision/snapshot, verified content mapping, time, freshness and errors | Named observer with inspectable evidence |

A PR may include a ticket saying ready for delivery and its evidence references. After integration, observe the delivered content separately. No follow-up commit is required merely to certify the merge containing that ticket, avoiding recursive closure PRs. Historical narrative may be updated later, but is not the delivery signal. Prior conditional-done records are historical claims, not a new status contract.

Proposed storage: local ignored evidence outside tracked tickets, keyed by project/checkout, target, criteria revision and candidate/content hashes. No tracked append-only status ledger, credentials store or mandatory central service. Exact format, location, retention and sharing remain unresolved; reconcile existing storage rules before adopting it. Missing or lost evidence means unverified, not failed or implicitly delivered. Read-only status inspection must not publish evidence or upload demos.

## Proposed targets, without a universal PR workflow

| Explicit project target | Evidence to inspect |
| --- | --- |
| Remote branch, optionally through a PR | Project-required acceptance assessment and verified inclusion/equivalent content on the named remote/branch; PR evidence only when the project requires it |
| Local branch | Project-required acceptance assessment and verified inclusion/equivalent content on the designated local branch |
| Local working tree | Project-required acceptance assessment bound to a manifest of the exact relevant files, including additions/deletions and scoped untracked files |

These are proposed options, not automatically selected policy. A configured remote requirement does not disappear when a remote is removed or unavailable: report unavailable verification. A deliberately local project needs no remote to establish **Delivered locally**. Bare merge or branch inclusion is insufficient when project criteria require additional acceptance.

## Web presentation

Show recorded work state alongside delivery state, target, criteria/assessment provenance, candidate and observed revision, last verification time, evidence/demo references, remaining gaps and unavailable/stale reasons. The current-plan view and ticket detail must consume the same delivery read-model rather than duplicate status labels in prose. Intent labels such as Proposed/Accepted remain distinct from work and delivery states.

Proposed labels include Delivery criteria not defined, Not delivered, In review (only when applicable), Integrated—verification pending, Delivered, Delivered locally, Delivery unverified, and Changes since verification. A PR link is optional evidence, not a required control. Never turn absent evidence or a network failure into Not submitted or failed acceptance. Retain the last verified delivery and timestamp while making stale current verification explicit.

## Concurrency and revision binding

Bind acceptance and delivery evidence to the criteria revision and exact relevant content. A changed candidate, criteria, target or relevant file invalidates affected current verification; retain historical observations. Use expected-revision checks for tracked updates and a compare-and-swap identity for observation publication, so late results cannot overwrite newer candidate observations. Read a stable snapshot or compare manifests before/after checking; if content changes during verification, retry against a new snapshot or report unverified.

For a working-tree target, exclude the observation store itself from the verified scope to avoid another self-reference loop. Explicitly define relevant paths and deletions; do not treat unrelated changes as automatic acceptance invalidation. Squash/rebase integration needs inspected content equivalence, not only original commit ancestry. These algorithms remain proposed and must be tested before claiming delivery verification.

## Qualification scenarios

1. A PR-controlled project integrates its ticket and code together; delivery becomes observable after integration without a second status PR.
2. A local-only project verifies a local branch or working-tree snapshot with no hosting calls.
3. A project has no delivery criteria: the UI says so and creates no PR requirement.
4. A candidate or acceptance-criteria revision changes during checking: stale observations cannot promote the new revision.
5. A relevant file changes after local verification: retain the historical delivery and show Changes since verification.
6. Remote access or the local evidence store is unavailable: distinguish stale/unknown evidence from failed acceptance; never silently change the target.
7. Ticket and current-plan web views agree; optional CLI output uses the same semantics. A merged PR with unmet project acceptance never displays Delivered.

These are planned checks, not executed evidence. [The current plan](implementation-plan.md) keeps implementation scope proposed; this discussion does not silently replace MCP Apps exploration as the next admitted build.
