---
type: "Specification"
title: "Codex conversation and MCP integration"
description: "Codex conversation and MCP integration."
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

# Codex conversation and MCP integration

**Desired behavior:** bidirectional GUI conversation with current Codex CLI, fresh sessions per action, concurrent conversations, runtime-owned persistence and factoryctl MCP tools. Codex is the first adapter implementation; contributor selection and optional home preferences follow [personal runtime preferences](personal-runtime.md), not repository configuration. These requirements come from the [operator](../sources/operator-brief.md). The adapter design below is proposed. [Observed protocol evidence](../sources/environment.md) is narrower than a working integration.

## Process and protocol boundary

When the contributor selects Codex, launch their selected installed `codex app-server` as a child process with stdio JSON-RPC. One child per canonical checkout/personal runtime selection snapshot isolates conversations; multiple Codex threads may run in it. Spawn with an argument array, never `sh -c`, using the selected repository as CWD. Keep protocol stdout separate from diagnostic stderr. Discover and validate supported protocol features at startup; record the detected runtime version in diagnostic output, not a durable session database.

Initialize the protocol, create a thread with `thread/start`, then submit the rendered action through `turn/start`. Stream item/turn notifications into the GUI. Send follow-ups on the same thread; `turn/interrupt` cancels active work. Handle server-initiated input and approval requests with their request IDs, displaying the runtime's reason and exact available choices. Send only the user's choice; do not fabricate approvals or translate errors into permission grants. Unsupported request variants display an actionable diagnostic and leave the action unapproved.

Read native schemas/capabilities rather than assuming all model effort values or experimental fields exist. Contributor-selected `model` and `reasoning_effort` are optional and never read from shared profile metadata: validate against runtime capabilities, map to supported thread/turn fields, and otherwise explain the unsupported setting. Omitted values defer to native runtime defaults. Do not change the user's global configuration or create `.codex/`, `AGENTS.md` or equivalent harness files as part of initialization.

## Supplying factory tools

Start the app-server with invocation-only overrides defining an MCP stdio server whose command is the absolute factoryctl executable and arguments are `mcp --path <selected-root> --project <id>`. Encode override values with a proper TOML serializer passed as separate `-c` argv entries; do not interpolate shell strings. This uses Codex's own per-invocation configuration mechanism, not a generated config file. The exact accepted override keys and readiness/error events must be qualified by the first runtime spike before committing the adapter implementation.

Codex is the MCP client; factoryctl is the MCP server. MCP supplies data/tools to the agent. App-server supplies conversation/approval control to the GUI. These are separate directions. The agent may also use its normal shell/file tools; factoryctl is not a policy sandbox.

Scope the MCP child to the selected project by default. Read tools may include other explicitly loaded projects in a multi-project invocation; writes always require an explicit target project within the launched scope. Cross-project references do not grant filesystem access. If mandatory MCP setup fails, fail the action before presenting a functioning factory conversation; do not silently continue with missing tools.

## Context assembly

Pass the selected profile's instructions, rendered action, selected record path and revision, project root, relevant definitions, epistemic/source labels and known reference paths. Inline only a bounded selection; expose full retrieval through tools. Mark truncation and unresolved sources explicitly. The agent receives instructions to reread before writing and supply current revision tokens. Never describe the conversation scope as an enforced security restriction when the harness has broader filesystem tools.

Starting a conversation does not execute a workflow, schedule a job or change task state. A future `start factory loop` action is just a prompt/profile invocation. Agent planning, subagent creation and Git operations stay in the external runtime.

## Lifecycle and failures

GUI conversation states are in-memory projections: starting, running, awaiting_input, idle, interrupted, failed, disconnected. They are not persisted task transitions. Distinguish `threadId`, `turnId`, item IDs and JSON-RPC request IDs. Do not retry a mutating turn automatically after an ambiguous disconnect: inspect runtime thread state, then let the user retry or continue. Correlate responses and tolerate interleaved events from concurrent threads. Serialize competing turns per thread or expose the runtime's queue semantics explicitly; different threads may run concurrently.

Closing a browser tab does not kill its thread while the factoryctl process is running. On server shutdown, stop accepting new actions, request cancellation of owned active turns, drain responses for a bounded interval, then close child transports/processes. Report unresolved effects; no durable factoryctl recovery queue. If a runtime terminates, browsing and direct record edits remain available.

Session persistence belongs to Codex. A restarted factoryctl may list/resume native threads through supported runtime APIs filtered to the selected CWD; it does not promise factory-specific ticket associations after restart or mirror transcripts. An unavailable runtime history response is not an empty-history claim. Credentials remain with Codex; no login/password storage in factoryctl.

## Qualification

The first implementation slice must demonstrate a fresh thread, streaming, follow-up, interruption, one native approval/input round trip, a factoryctl MCP read/write with stale-revision rejection, and two concurrent conversations. Pin the tested Codex version in the qualification evidence. Treat future unknown schemas as compatibility diagnostics, not proof of support merely because the executable is newer. See [runtime task](../tasks/07-codex-adapter.md) and [OpenAI App Server architecture](https://openai.com/index/unlocking-the-codex-harness/).
