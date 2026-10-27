---
type: "Reference"
title: "Observed design environment and evidence limits"
description: "Observed design environment and evidence limits."
status: "draft"
factory:
  format: 1
  assertion:
    frame: "observed"
    certainty: "supported"
    basis: "direct-observation"
sources:
  -
    id: "local-tools"
    resource: "Local git inspection, go version, uname -m, codex --version and generated app-server JSON schema on 2026-10-27"
    title: "Authoring-session tool observations"
---

# Observed design environment and evidence limits

## Direct observations

During design preparation on 2026-10-27, the local checkout had no commits or application files. Its remote identified `tclasen/factoryctl`. The local host reported `arm64`; installed tools reported Go `1.27.3` for `darwin/arm64` and Codex CLI `0.164.0`.

`codex app-server --help` exposed stdio and schema-generation support. Generating the non-experimental JSON schema locally produced 320 JSON schema files in temporary scratch storage. The schema exposes `thread/start` fields including cwd, model, config and developerInstructions, and `turn/start` fields including threadId, input and effort. This is protocol-surface evidence, not an executed conversation or successful MCP handshake.

## Limits

No factoryctl executable exists. No Codex session was launched by factoryctl, no runtime approval round trip was tested, no cloud process was deployed and no Linux/Windows executable was run. Installed versions are observations, not a claim that they remain latest. Implementation must record the version it actually qualifies.

Temporary protocol/schema and validation files are authoring tools outside the repository. They are not application persistence or required by a user of factoryctl. [Release acceptance](../specifications/acceptance.md) defines the later tests.
