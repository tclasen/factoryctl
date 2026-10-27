---
type: "Factory"
title: "Factoryctl: purpose, desired state and current state"
description: "Factoryctl: purpose, desired state and current state."
status: "draft"
factory:
  format: 1
  assertion:
    frame: "mixed"
    certainty: "mixed"
    basis: "mixed"
sources:
  -
    id: "operator"
    resource: "sources/operator-brief.md"
    title: "Operator requirements and scope"
    author: "human:operator"
project_id: "factoryctl"
bundle_uri: "https://github.com/tclasen/factoryctl/blob/main/factory/"
---

# Factoryctl

Factoryctl is a proposed portable CLI, local GUI and MCP interface to knowledge work recorded in Git. People and agents share readable files. Agents interpret operating instructions and perform the work; factoryctl creates, validates, discovers, renders and edits records and connects people to their agent runtime.

## Desired state — operator requirements

The [product specification](specifications/product.md) defines the accepted scope from the [operator interview](sources/operator-brief.md). The [knowledge format](specifications/knowledge-format.md) describes the OKF bundle and its extension. [Project discovery](specifications/project-discovery.md) works without a root config file; [personal runtime preferences](specifications/personal-runtime.md) stay outside project Git. The [template contract](specifications/templates.md) explains shared customization.

## Current state — observed

This repository contains a design package and factory records. There is no application implementation or released executable. The planned CLI examples are specification examples, not commands available in this checkout. [Design validation](sources/design-validation.md) reports document checks; [environment observations](sources/environment.md) identify the limited host/runtime evidence. All implementation tasks are planned.

## Understanding the design

Read the [architecture](designs/architecture.md), [interaction specification](specifications/gui.md), [CLI/MCP/HTTP contracts](specifications/interfaces.md), and [Codex integration](specifications/agent-runtime.md). The [implementation plan](designs/implementation-plan.md) links buildable slices, requirements and checks. [Release acceptance](specifications/acceptance.md) defines completion without claiming those checks have run.

## Knowledge beyond software

[Research curation and contract pursuit](specifications/knowledge-work.md) demonstrate custom types, source attribution and configurable boards. The [researcher profile](agents/researcher.md) illustrates how instructions remain knowledge nodes. Software specifications and designs are an optional starter set, exercised here to build the product itself.

## Knowledge quality

[Epistemic conventions](specifications/epistemics.md) separate requirements, current observations, assumptions and LLM proposals. Markdown headings and links carry the same distinctions as the GUI. Source claims are not promoted to verified facts by being copied into this graph. Git supplies historical revisions; files describe current knowledge.

## Development and releases

Follow the [commit and release policy](decisions/release-policy.md). The initial release target is 0.1.0; there is no released executable yet.

## Runtime-extensible document behavior

The [extension contract](specifications/extensions.md) describes user-defined types, declarative lifecycle rules and namespaced tools loaded at runtime. Custom code runs in the browser UI; arbitrary local work is requested through prompts to the selected harness. Factoryctl remains a thin document interface and harness launcher, with no Git subprocesses or built-in agent scheduler. The [extension prototype](tasks/12-extension-contract.md) must qualify this boundary before implementation claims. These capabilities are design requirements; the repository still contains no application implementation.
