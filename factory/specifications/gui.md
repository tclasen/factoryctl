---
type: "Specification"
title: "GUI journeys and interaction design"
description: "GUI journeys and interaction design."
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

# GUI journeys and interaction design

This proposed interaction design implements [R07–R15](product.md). All screens work for one project or many. Every concept, including a definition, prompt and infrastructure node, has list and detail views.

## Shell and navigation

A persistent project selector offers **All projects** and named checkouts, with a location disambiguator for duplicates. The overview shows per-project counts by type/state and structural-health diagnostics, plus active native conversations in this process. No Git hosting login or shared task database. Not-initialized repositories have an Init action. Invalid configs remain visible with their errors.

Left navigation derives enabled types and views from definitions; it is not an INDEX.md file. A search/filter bar works across the selected scope. Lists expose type, title, epistemic labels and configurable columns. Keyboard navigation, labeled controls, visible focus, screen-reader status announcements and adequate contrast are part of acceptance. Virtualize/paginate large lists; do not render 10,000 DOM rows at once.

## Record detail

Header: project, type, title, concept path, lifecycle, knowledge classification and validation health. Body: sanitized rendered Markdown with source footnotes and links. Side panels: typed fields, outgoing relationships, **Backlinks** (always present, including an explicit empty state), applicable actions and conversation access. Show unresolved and external references differently from verified local targets. In all-project mode a link changes the selected project only after identifying its target; unresolved links never switch to an arbitrary checkout.

Desired/current/proposed headings remain visible. Assumptions and LLM-derived conclusions retain labels in the rendered body. A compact label does not replace source context. Unknown types and disabled types retain generic details and raw editing with diagnostics; invalid syntax shows the raw file and repair guidance rather than a blank page.

## Editing and boards

Structured field editing and full Markdown editing are both available. The original revision token accompanies every save. Inline validation checks fields, placement, templates and links; server validation repeats it. A stale-save dialog shows the editor's draft and newly read content and offers reread/reconcile, never silent last-write-wins. Local unsaved text stays in browser memory; closing the tab may lose it, with the normal unsaved-changes prompt. No browser persistence is required.

The board is a projection of an enum field, usually task_state. Dragging a card performs one revision-checked field update. Keyboard users can choose the destination column with the same action. A failed write returns the card to its observed position and displays the error. No additional transition history or rank file. Within a column use the shared deterministic sort; user-defined manual ranking is deferred. Cross-project boards show project badges; incompatible state definitions get separate boards.

## Conversations

From a ticket, **Request status update** starts a fresh Codex thread with that ticket's context, selected profile and rendered prompt. Stream agent text and tool progress; show runtime questions/approvals as actionable cards. User messages continue the same thread. The record refreshes when disk changes, and the UI distinguishes an agent's textual claim from an observed file edit. Several conversation tabs can run independently; clearly label the project and selected record in each.

Other records use **Discuss** by default; custom prompts add actions. The user can inspect/select the profile and prompt before starting. Cancelling invokes native interruption and reports whether it was acknowledged; it does not revert file changes. Runtime failure preserves the knowledge views and direct editing. Reopening history delegates to available Codex thread APIs.

## Extensions and failure states

Custom executable assets initially show a trust affordance and generic fallback. Explain that enabling the view executes repository code in a constrained browser frame and exposes the selected record data to it. Do not repeat warnings on every click during the same trusted digest/process. The built-in project selector, validation errors and approval controls remain outside custom frames.

Required empty/error states: no repositories, no initialized factories, zero records, missing definition, malformed file, missing link/project, unsupported runtime setting, missing Codex, cancelled approval, stale edit and lost connection. State transitions should be understandable without reading server logs.

See [GUI implementation](../tasks/08-gui.md), [extension implementation](../tasks/09-custom-views.md), and [acceptance](acceptance.md).
