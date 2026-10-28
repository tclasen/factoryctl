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

## Candidate and review checkpoint

Initial implementation commit: `51e95bb493915ab639175f89109c4cfe2a42e378`, tree `ff59e16f41c0b94d646531a89a6915e5b268f7a0`. All six packages passed normal tests, race tests and vet. The developer built CGO-disabled Darwin ARM64, Linux AMD64 and Windows AMD64 binaries; only Darwin was executed. Coordinator independently rebuilt the Darwin binary and repeated both external CLI journeys successfully. Its SHA-256 was `6f7c1ea2a007c014ac914305346545fe4ebdd44ee19d36d11e1f6d434ecfdc89`.

A clean Git archive of that candidate, initialized as a temporary repository, validated with zero errors, zero warnings and 26 unchecked links (`complete: false`). Mach-O inspection reported system `libSystem` and `libresolv` dependencies: this is a self-contained application binary, not a claim of fully static Darwin linkage. Compiler cache-stat warnings were nonfatal; artifacts and exit codes were checked. A governance hash check initially used the wrong working directory; rerunning from the workspace root verified every snapshotted instruction/policy unchanged.

Fresh independent implementation review reproduced two blocking findings on the initial candidate despite passing tests:

- IR01 (P2): opening a FIFO before checking regular-file status can block document validation; analogous template/discovery reads require the same repair. Expected: special files fail promptly and remain diagnostics without being consumed.
- IR02 (P2): nonfinite YAML metadata can make JSON encoding fail while the command returns success and no output. Expected: nonfinite data is handled explicitly and encoding failures cannot produce silent success.

Both findings were assigned to the developer for regression-tested repair. The initial candidate is not approved for merge. Repairs add a shared regular-file reader with pre/post checks and nonblocking/no-follow Unix opens, including document, template, discovery and writer rereads. Generic nonfinite YAML values retain their source bytes and are represented as explicit YAML-token strings in JSON with a diagnostic. JSON serialization and stdout errors now propagate instead of reporting success. New bounded subprocess and output-failure tests cover the reproduced cases. The repair candidate passes `go test ./...`, `go test -race ./...`, and `go vet ./...`; CGO-disabled Darwin ARM64, Linux AMD64 and Windows AMD64 builds also pass. FIFO regressions terminate within their three-second deadlines; NaN/positive and negative infinity outputs parse as JSON without changing source bytes; failing output writers return exit 4 in JSON, human and help modes. The repaired owned-source manifest SHA-256 is `2a5ba135f267156f2a4feaf88adb787d81acd56ff4dbc8ab110fde3f4cf47976`. The independent reviewer must recheck the exact repair candidate before acceptance; its result and final commit/main identities belong to the delivery PR.

## Product scope still pending

Native Codex integration and its runtime spike, MCP/HTTP servers, GUI, browser extensions, user-question broker, long-running filesystem watcher, release video, full scale qualification and release publication are not acceptance claims of this CLI milestone. Any other unimplemented CLI contract must be reported explicitly in the final test evidence and command behavior. Tasks 01–05 retain any uncovered acceptance work; tasks 06–13 remain separate milestones.

## Delivery

Implementation and the proposed milestone closure are on [PR 8](https://github.com/tclasen/factoryctl/pull/8). Closure becomes effective only when the independent repair recheck passes, normal PR requirements are satisfied, the reviewed content merges, and remote main containment/content is verified. The PR records that exact-candidate review and delivery outcome; before those gates this record is a proposal, not a delivered-completion claim. The project main ruleset requires PR delivery and resolved review threads, forbids force updates/deletion, and exposes no required CI contexts. No protections will be bypassed. No released version or tag is created by this milestone.

## Review artifact omission and correction

The initial milestone handoff omitted the Playwright review video. The operator requested it immediately afterward. [The review-video task](../tasks/cli-review-video.md) records the correction and its fresh evidence; treating the implementation as an unreleased milestone did not satisfy the operator’s review expectation. The first CLI source delivery remains historical implementation evidence, separate from video acceptance.
