---
type: Task
title: Make intent visible and deliver a small browser review slice
description: Preserve request, proposal and authorization provenance and demonstrate it with the current planning documents.
status: draft
task_state: in_review
priority: high
owner: coordinator
requirements: [R07, R10, R31, R32, R33]
factory:
  format: 1
  assertion:
    frame: desired
    certainty: supported
    basis: authoritative-statement
intent_state: proposed
intent_owner: assistant
intent_source: ../sources/operator-brief.md#working-method-correction-and-tangible-review
intent_authorization: Implement the accepted working-method corrections; the read-only browser slice is the assistant-selected reversible means of making them tangible, not a user-authored architecture mandate.
---

# Visible intent and tangible review

- Target System: project
- Change Classification: project
- Authorized Scope: project instructions/current plan/provenance records; small read-only embedded browser document explorer and its loopback launch entry point; tests and local-only demo; normal reviewed code/document PR delivery
- Affected Paths or Resources: AGENTS.md, README.md, factory/ current plan/specification/task/source documents, cmd/factoryctl/, internal/web/, associated tests; scripts/review-video/ local recording tool and ignored artifacts/
- Owning Role: Coordinator
- Related Requirements: R07, R10, R31, R32, R33
- Related Tasks: first-cli-milestone, 08-gui, 13-user-interaction
- Approval Requirements: operator's explicit approval to implement the working-method proposals; independent intent-fit and technical review; normal protected PR/main delivery; videos and generated artifacts remain local

## Request, interpretation and authority

Exact latest request: “But otherwise, please implement all of that.” Context: the operator approved separating user requests from agent interpretations and authorization, retaining proposal status, reconciling one current plan, reviewing user-outcome fit before implementation correctness, and reassessing repeated corrections. They also asked to accommodate wanting a tangible result instead of reading a long plan, preferring demo videos to text-only summaries.

Interpretation/assistant choice: implement the operating changes in existing project records and demonstrate their usefulness with a small read-only browser explorer of the real Markdown bundle. This is a reversible vertical slice, not proof of full GUI/MCP Apps integration or an instruction from the operator selecting a specific web architecture. The shared operating-method repair has independent organization ownership and delivery; this task grants no organization edits.

## Plan revision 1 and claims

Baseline e52c4e0231145a74a368f542064b928dc2565eb2. Coordinator owns factory/ documentation, README, project AGENTS, local recording integration and Git operations. Developer owns cmd/factoryctl changes, new internal/web/ and relevant tests only. Independent reviewer owns no implementation. No edits to existing user files, runtime settings, external hosting, artifact publication, destructive cleanup, full web editor, conversations or MCP host integration.

1. Reconcile the active plan and distinguish exact request/context, agent proposal and scope of authorization; superseded sequence remains historical with a clear replacement.
2. Expose real project documents in a readable browser view with search/type/status filters, source/authorization details, navigable document links and raw Markdown. Surface requested/proposed/accepted/superseded/unresolved states as recorded claims; never infer human acceptance from Git authorship, tests, merged status or document stable status.
3. Serve read-only on loopback from the existing single binary. Reuse safe discovery/parsing. Reject nonloopback binding, path escape and mutation routes; render Markdown inertly without remote resource fetches or raw executable HTML. Unknown/missing authority displays unresolved, not accepted.
4. Run meaningful HTTP/security/UI tests; record a short local Playwright walkthrough of actual product interactions with requirements annotations and visible limitations. Keep latest stable local link and immutable local archive. Complete matrix gaps stay explicit; this partial slice does not claim full qualification.
5. Independent review tests outcome fit against source requests separately from implementation correctness; reconcile feedback, ordinary PR/checks/merge and verify remote main.

Budget:60active minutes then checkpoint; commands120seconds/tests10minutes; at most2diagnosed retries per failing operation. Rollback through normal revert of task-owned changes; local prior demos preserved. Actual cost/model telemetry unknown, inherit current model. No new providers/services. Working title/UI details are engineering choices.

## Acceptance

A user can inspect the current plan and related documents, distinguish a proposal from accepted work, open the source and authorization text, see superseded ideas and unresolved decisions, and inspect raw Markdown. None of these interactions edits a file or submits acceptance. The local demo shows those journeys, not a requirements presentation. Small-slice authorization never turns unread future options into accepted decisions; publication remains separately bounded. Ordinary authorized engineering proceeds without repeated confirmation.

## Intent-fit review checkpoint

Independent Reviewer assessed the source request, current plan and scope before implementation acceptance: proceed, no blocking outcome/authority finding. This is readiness only. Final review must verify mixed documents do not imply blanket acceptance, and source/authorization details remain directly inspectable.


## Current review package and conditional closure

The [browser and CLI evidence](../sources/visible-intent-evidence.md) records the current local combined demo, assertions, source provenance and remaining gaps. The operating-method changes and bounded explorer are implemented. Independent approval and verified PR/main delivery remain the closure gate; no blanket full-matrix qualification or future-behavior guarantee is asserted. The reusable organization protocol has separately completed its reviewed delivery.
