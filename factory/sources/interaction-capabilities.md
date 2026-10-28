---
type: Reference
title: Interaction runtime and surface capabilities
description: Observed Codex runtime and terminal evidence, first-party view boundaries and CLI disposition.
status: draft
factory:
  format: 1
  assertion:
    frame: observed
    certainty: mixed
    basis: direct-observation
sources:
  - id: interaction-contract
    resource: ../specifications/user-interaction.md
    title: User interaction requirements
  - id: runtime-contract
    resource: ../specifications/agent-runtime.md
    title: Native runtime contract
  - id: active-task
    resource: ../tasks/13-user-interaction.md
    title: Current bounded milestone and acceptance scope
---

# Interaction capabilities

This report separates implementation, bounded live observations and remaining qualification. It is not a full requirements-matrix pass, a milestone completion declaration or an external MCP Apps interoperability claim. Local evidence below belongs to the 2026-10-28 qualification run on macOS ARM64 using installed `codex-cli 0.164.0`; generated artifacts remain outside Git.

## Capability table

| Capability | Evidence and boundary |
| --- | --- |
| Native app-server startup and invocation-only MCP attachment | Observed live through stdio initialization, native thread/turn events and a synthetic MCP echo tool. The Go adapter independently received `FACTORY_PROBE:go-adapter-roundtrip`. No global runtime configuration write is implemented; the probe did not measure before/after hashes of global configuration. |
| Native model/effort selection | Implemented validation against `model/list`, with deterministic adapter tests. Live conversations used runtime defaults; arbitrary model/effort combinations are not qualified. |
| Factoryctl terminal questions | Observed live in an actual PTY with the real factoryctl MCP bridge: text answer, custom answer to a single-choice question and agent acknowledgement of both. Inputs were automated test submissions, not a human acceptance review. |
| Same-thread follow-up | Observed live: the agent repeated both supplied answers on the next turn; native identities show all three turns used one thread. |
| Cancellation | Observed live: a pending unanswered question was interrupted, with matching native `turn/interrupt` and `turn/completed` status `interrupted`; terminal exit was zero. This is not a user answer or approval. |
| Native runtime input | Observed live in an adapter probe using native plan mode: `item/tool/requestUserInput` request ID `0`, an automated host response `Violet`, then completed agent output `NATIVE_INPUT Violet`. This qualifies adapter routing, not browser rendering of that native request. |
| Native approvals | Schema support and host validation exist. No live human command/file/permission approval round-trip is qualified by these probes. Unknown request variants remain unsupported and unapproved. |
| Browser conversations | Implemented opt-in host controls, scoped questions, follow-ups and interruption. Runtime/PTY evidence here does not substitute for the coordinator's separate browser capture and review. |
| Generated browser views | Implemented first-party HTML/CSS/JavaScript presentation after explicit trust, sandboxed frame and plain fallback. Frame drafts require a separate trusted-host submit. This report does not claim that source inspection alone qualifies every browser isolation requirement. |
| External MCP Apps | Not qualified. First-party presentation and limited host messaging are not evidence of interoperable external Apps handshake/resources/bridge support. Standalone MCP advertises empty resource lists; interaction tools require an attached trusted surface. |
| Isolation and failure handling | Deterministic tests cover scoped broker credentials, invalid/repeated native replies, disconnection and bounded adapter overflow. They are distinct from live multi-project or all-scenarios qualification. |
| Persistence and portability | Factoryctl keeps conversations/questions in memory; restart recovery is absent. Native persistence is runtime-owned. Execution evidence is macOS ARM64 only; other platform builds are not execution qualification. |

## Reproducible local evidence

Evidence is under `artifacts/interaction-runtime/20261028/`, which is ignored and not distributed in a clone. `schema/` contains the installed runtime's generated experimental schema. `probe-summary.json`, `qualification.json` and `go-adapter-probe.json` record the original runtime/MCP observations. Raw native event files may include account or configuration metadata and must remain private local diagnostics.

`chat-pty-final/probe-pty.py`, `chat-pty-final/stdout.txt`, `chat-pty-final/native-identities.jsonl` and `chat-pty-final/result.json` record the terminal exercise. The script opened a real PTY, verified `isatty` and disabled echo so assertions could not succeed from echoed input. It supplied `juniper-qualification`, then custom choice text `custom amber review`; the agent returned both after `CLI_ROUNDTRIP` and again after `FOLLOWUP_MEMORY`. A third turn asked a question and was cancelled without an answer. Native thread identity was `01a1805b-d864-7363-b438-2d8f727fbf91` throughout.

The PTY binary's before/after SHA-256 was `3e3fd1d45d4145a2f12c826f848e81ca45b00d5b5ba00b56d820fbb48933bd6b`; result SHA-256 was `fd9991dbfca0c1af7b873fd927c8a7a55cab222bc2c5b58deff9e185c2c5786c`. The source inventory pins revision `bb93156a29270e9a0f0a396cffdae2d3488bdf1e`; all 51 recorded source files remained unchanged through the run. This pins that probe candidate, not later implementation changes or the eventual delivery binary. Earlier piped-stdin runs remain historical transport diagnostics and are not interactive terminal qualification. The first attempt encountered native model capacity exhaustion; the failed attempt was retained rather than counted as a pass.

`native-input/probe-native.go` and `native-input/result.json` record the native input exercise. Result SHA-256 was `4765b1cfe1f3f4abf2cc6e31655055c7281b8f719179fc76a4a7bfd60230c6f7`. The native input probe used explicit plan mode; ordinary factoryctl conversations do not thereby gain a claim that native input is available in every runtime mode.

The adapter's deterministic tests and race checks passed during its implementation. Final integrated checks, browser recording, complete matrix dispositions and the stable local demo are owned by the milestone evidence package; this reference does not infer their result. The stable review link remains `/Users/Shared/factoryctl/artifacts/demo/latest.mp4`, updated by the delivery workflow only after the new capture is validated.

## CLI recommendation

Retain the interactive CLI alongside the primary web exploration surface. The PTY exercise establishes a useful terminal answer/follow-up path, while external Apps interoperability and full browser parity remain unqualified. This is a recommendation, not a retirement decision. Removing the CLI requires the explicit operator decision and parity evidence specified by the interaction contract. Document automation and MCP stdio are separate interfaces and do not collect human answers from protocol stdin.


## Integrated qualification follow-up

The [current milestone package](interaction-milestone-evidence.md) adds successful actual browser free-text/form/chart drafts, fresh same-thread recall, generated-code failure fallback and navigation isolation. Live automated host API tests also passed concurrent project/thread separation, cross-scope answer rejection, actual refusal, multiple selections, connection persistence and old-ID rejection after process restart. Their source/result hashes and exact video scenes are retained in that package. These observations do not upgrade native browser approval coverage or external Apps interoperability.


## Current completion candidate — superseding bounded capability observations

The earlier observations above remain historical. Current product revision `a7c22ba847d58e8f4ece6a22ee9bb73302a64c02` adds the [completion evidence](completion-milestone-evidence.md). Actual browser/native runs now cover free-text and custom-choice answers, generated form/chart drafts with explicit host submission, same-thread follow-up, concurrent project/thread isolation and interruption. The model run exercises Astra/low, Sol/low and runtime-default Sol/medium, personal global/project defaults, restart/reset and project isolation. These observations do not qualify every advertised model or future surface.

Current real PTY attempt 3 performs a definition-and-revision-pinned native MCP Research Fact edit and fresh readback in the same native conversation. Approval-only attempt 6 receives the actual native `item/commandExecution/requestApproval` accept/cancel offer; a trusted synthetic cancel produces correlated command-declined/no-execution and turn-interrupted events. The independent trace verifier passes that cancellation criterion. The original whole recorder remains `passed:false` because its termination wait expected an idle prompt after interruption; no successful whole-run rerun is claimed. Invocation-only `approvals_reviewer="user"` is schema-supported; no global runtime configuration change was made. This proves the negative approval route, not human acceptance or every native approval variant. The observed `mcpServer/elicitation/request` variant remains unsupported.

External MCP Apps interoperability remains unqualified. The CLI is retained. Foreign binaries were cross-built and inspected, never executed. Full completion remains blocked by U045's unresolved static-linkage interpretation and final combined evidence/review/delivery; a capability report cannot create operator acceptance.
