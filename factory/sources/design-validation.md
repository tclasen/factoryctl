---
type: "Reference"
title: "Design validation and review evidence"
description: "Design validation and review evidence."
status: "draft"
factory:
  format: 1
  assertion:
    frame: "observed"
    certainty: "supported"
    basis: "direct-observation"
sources:
  -
    id: "operator"
    resource: "operator-brief.md"
    title: "Operator requirements and scope"
    author: "human:operator"
---

# Design validation and review evidence

## Scope

This record concerns the design package, not factoryctl application behavior. The current package passed the author-run document checks described below and is awaiting independent review. No application tests, runtime integration tests or non-Mac execution tests have passed by virtue of these documents existing.

## Verification plan

Parse all YAML/frontmatter; verify OKF type and knowledge classification, definition fields/placement, link targets/anchors, prompt/profile/template references, task dependency acyclicity and R01–R25 coverage. Inspect the absence of index/log files, app persistence stores and generated Codex-specific files. Independently review architectural fit, comprehension and execution readiness.

See [design task](../tasks/design-package.md), [format](../specifications/knowledge-format.md) and [release acceptance](../specifications/acceptance.md).

## Observed author-run checks

On the current Mac, a temporary Python/Ruby validator parsed all frontmatter and factory.yaml, checked the materialized definitions and required field values/locations, resolved local body/metadata links and heading fragments, checked task dependencies for cycles and mapped R01–R25 into tasks. The pre-review run reported 57 OKF nodes, 11 definitions, 14 tasks, 303 checked local references, all 25 requirements covered, and zero errors. Seventeen distinct external URLs were identified but not checked for reachability by this validator.

A temporary Go program used the installed Go text/template engine to parse and render all six record/prompt payloads with synthetic input. Ruby parsed the three rendered record frontmatters and verified a title containing a colon and quotes survived serialization. This tests the checked-in payload syntax with the proposed scalar-encoding approach; it does not validate a nonexistent production template engine.

## Coverage limits

The temporary document checker uses a bounded Markdown-link scan with fenced examples excluded, not the future production Markdown AST parser. Its passing result applies to this authored corpus, not arbitrary user documents. External URL availability, runtime conversation, MCP handshake, production writer races, custom JavaScript isolation, performance and non-Mac behavior are untested. The production negative fixtures and complete parser behavior remain planned implementation tasks.
