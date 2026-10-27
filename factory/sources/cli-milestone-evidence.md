---
type: Reference
title: First CLI milestone implementation evidence
description: Verification scope and remaining work for the initial document CLI.
status: draft
factory:
  format: 1
  assertion:
    frame: observed
    certainty: supported
    basis: direct-observation
---

# First CLI milestone evidence

The [milestone task](../tasks/first-cli-milestone.md) authorizes a usable document CLI. It does not authorize treating the entire product or initial release as complete. Initial implementation and command-level checks now exist. Results below apply to the document slice only; the final candidate still requires independent review and remote-main delivery.

## Environment and provenance

Baseline commit: `1d6d0d40a80dadb1b17a5320efa7ee0a4fc061e3`. Execution platform: Darwin ARM64; available compiler: Go 1.27.3. Other platforms may be cross-compiled but remain unexecuted. Exact implementation candidate and independent reviewer outcome will be recorded before delivery. Model/token/cost telemetry is unavailable.

## Verification scope

Check discovery for single projects, monorepo bundles and sibling repositories; parse the live bundle; exercise structural definitions, references and templates; validate scratch initialization and document CRUD; prove stale/invalid/unsafe writes leave original bytes unchanged. Run unit/integration tests, static checks and a CGO-disabled build. Check actual CLI exit codes and stdout/stderr separation. No application command may launch Git or mutate its metadata.

## Initial independently executed command checks

Coordinator built the Darwin ARM64 executable with `CGO_ENABLED=0 go build -o /tmp/factoryctl ./cmd/factoryctl`, using a temporary Go build cache. Two external Python-driven CLI journeys exercised fresh temporary Git repositories. Observed successes: software starter initialization and validation; Task create/read/update/delete; stale revisions returned exit 3 without changing bytes; invalid enum edits returned exit 1 without changing bytes; unknown metadata and inline comments survived field edits exactly; linked-record deletion and referenced-heading removal were refused; symlink/traversal writes were refused; flow-style YAML editing either preserved structure or refused safely. The scripts are session-local test drivers, not product runtime dependencies. Introducing source tests provide durable regression coverage.

The first driver needed a fixture directory creation correction and an expected exit-code correction: referenced deletion is a structural validation failure (exit 1), while stale revisions use exit 3. These were test-driver corrections, not successful product results inferred from failures.

Live checkout validation reports one error for the pre-existing untracked `factory/.DS_Store`; it is preserved. The Markdown corpus itself has no errors at this checkpoint. There are 26 unchecked external links, so `complete` is false. A clean tracked/owned corpus check is required for final evidence; default offline validation never claims those remote URLs were verified.

## Product scope still pending

Native Codex integration and its runtime spike, MCP/HTTP servers, GUI, browser extensions, user-question broker, long-running filesystem watcher, release video, full scale qualification and release publication are not acceptance claims of this CLI milestone. Any other unimplemented CLI contract must be reported explicitly in the final test evidence and command behavior. Tasks 01–05 retain any uncovered acceptance work; tasks 06–13 remain separate milestones.

## Delivery

Pending implementation, independent review and ordinary pull-request merge. The project main ruleset requires PR delivery and resolved review threads, forbids force updates/deletion, and exposes no required CI contexts. No protections will be bypassed. No released version or tag is created by this milestone.
