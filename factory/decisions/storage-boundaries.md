---
type: "Decision"
title: "Separate knowledge, assets and personal preferences"
description: "Separate knowledge, assets and personal preferences."
status: "draft"
factory:
  format: 1
  assertion:
    frame: "desired"
    certainty: "mixed"
    basis: "mixed"
sources:
  -
    id: "operator"
    resource: "../sources/operator-brief.md"
    title: "Operator requirements and scope"
    author: "human:operator"
decision_state: "accepted"
---

# Separate knowledge, assets and personal preferences

## Operator-authorized decision

The project requires no `factory.yaml`; remove it and do not load or generate it. Keep optional versioned `.factory/` assets and `factory/` as the OKF knowledge graph. Default operation needs no asset directory. Git supplies project status/history; agents handle Git operations. Runtime choice is contributor-specific, outside project Git.

## Mechanism — LLM design proposal

Keep project identity on the existing Factory node, type enablement on each definition and shared work intent on profiles/prompts. Use built-in operational defaults and invocation flags. An optional user-home Markdown file holds only personal runtime preferences. See [project discovery](../specifications/project-discovery.md) and [personal runtime preferences](../specifications/personal-runtime.md). This is not a replacement project config hidden inside Markdown.

## Deferred possibility

The operator identified remote tracking backends as a possible future reason for project config. GitHub Issues/Projects, Linear, Trello and Jira are not implemented or selected. Prove a routing need and resolve storage/synchronization semantics before adding a config file. Git-versioned Markdown remains the sole initial backend.

## Consequences — expected, not measured

A project can be shared by contributors with different runtime preferences without editing Git state. Application source and build inputs remain outside the knowledge graph. No index, registry or duplicate status history is introduced.
