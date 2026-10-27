---
type: "Reference"
title: "Operator requirements and design authority"
description: "Operator requirements and design authority."
status: "draft"
factory:
  format: 1
  assertion:
    frame: "observed"
    certainty: "supported"
    basis: "authoritative-statement"
sources:
  -
    id: "interview"
    resource: "Operator statements in the factoryctl design interview in this session"
    author: "human:operator"
    title: "Factoryctl design interview"
---

# Operator requirements and design authority

This is an LLM-authored extraction of the operator's statements in the factoryctl design interview. The quoted sentences below are direct statements; the grouped requirements are a faithful synthesis, not independent external verification. The conversation is the primary source and is not publicly addressable from this repository. The operator has authorized research and a complete design-document package; application implementation is the subsequent stage.

## Direct authoritative statements

> AI agent always handles git operations.

> Can we remove it as a requirement from the project, and remove it from this repo, and only add it back in once we can prove there is a valid use for it?

> I want it to be possible to run a project even without a .factory/ folder at all if you wanted.

> Also note, I want everything in the factory/ knowledge graph OKF bundle to be human readable and navigatable, it should double as documentation for the desired state of the project as well as the current state of the project.

## Confirmed requirements, synthesized from the interview

1. A knowledge-work factory interface replacing earlier experiments, with CLI, a launchable GUI, and MCP tools for an external AI runtime. Factoryctl supplies structure and interaction; agents interpret instructions and automate work.
2. Git repositories are mandatory. Projects may share a monorepo or occupy separate repositories. Discover project roots from CWD or an explicit path, with only immediate subfolders scanned by default (recursion limit 1). All interfaces use a project collection even with one project; the three required layouts are specified below.
3. No project configuration file is required or read. Optional, versioned `.factory/` holds templates and UI assets. Only `factory/` is the OKF bundle: project metadata, definitions, profiles, prompts, specifications, tasks, decisions and documentation. Source code, deliverables and asset files are outside this graph.
4. Human-readable files hold durable factory state. No factoryctl database, project configuration store, session store or persistent index. Contributor-specific runtime preferences may live outside Git in the user home directory; they are not project knowledge. Git owns historical revisions: change a current field in place instead of adding transition logs. External runtimes retain their own normal storage.
5. Three initial domain types: agent profiles, tasks and decisions. Users can define entirely new types or disable supplied ones. Definitions specify fields, validation, relationships and views. Everything in the bundle is a valid OKF concept; no index files.
6. Every type has list and detail views. Provide a project overview, cross-project views, task kanban, backlinks, direct editing and interactive bidirectional conversations. Prompt templates define GUI actions, may compose templates, and start fresh sessions followed by continuing conversation. Multiple conversations may run concurrently.
7. Latest Codex CLI is the first required adapter implementation; runtime selection belongs to each contributor and is not committed to the project. Profiles remain OKF nodes; do not generate Codex-specific files. Runtime owns conversation persistence. Agents handle Git operations, scheduling and factory loops.
8. Custom presentation supports arbitrary HTML, CSS and JavaScript. GUI serves one trusted local user, or sits behind external infrastructure. Distributed collaboration is separate checkouts exchanging Git changes, not a shared-user database.
9. Single portable compiled binary; desktop, server/Kubernetes and CI usage. Target macOS, Linux and Windows on ARM64 and x86-64; execution testing is initially on the operator's current Mac only. No unsupported-platform execution claim.
10. Use factoryctl's own record format to plan and manage building it. Supply an optional software-development starter set. Maintain living specifications; Git and task records track changes.
11. Clearly distinguish desired state from observations, fact from assumption, and authoritative statements from LLM-derived conclusions. Attribution must survive reading the Markdown without the application.
12. Acceptance examples include research and maintaining a knowledge graph, and discovering, pursuing, submitting and tracking government contract bids through award. Initial design target: 100 repositories, 10,000 total nodes.

13. Factoryctl validates documents as they are written: valid links, mandatory fields, correct field types and correct file placement. Structural validation is distinct from agent-owned workflow policy.

## Authority and proposal boundary

These statements establish product requirements. Technology choices, exact field names, commands, error codes, performance thresholds and the implementation sequence elsewhere are LLM design proposals unless identified as operator-confirmed. A requirement is not a claim that software implements it. Publication of this design does not establish implementation success or human verification of every proposed detail.

Read the [product specification](../specifications/product.md) for traceable requirements, the [architecture](../designs/architecture.md) for proposed mechanisms, and the [design task](../tasks/design-package.md) for current scope and verification.

## Versioning directive — direct operator source

> From now on please ensure this project uses conventional commits and from that semantic versioning. We are starting at v0.1.0, and until I authorize the bump up to and release of v1.0.0 we don't care about the concept of backwards compatibility.

The [release policy](../decisions/release-policy.md) separates these explicit requirements from agent-selected operational conventions.

## Tracking and personal runtime directives — direct operator sources

> By default I want factoryctl to store all project status and tracking documents in git as versioned files

The operator suggested GitHub Issues/Projects, Linear, Trello and Jira as possible alternative tracking backends, and project configuration as a possible way to select them. This is a future possibility, not an instruction to implement or configure those backends now.

> Codex being the runtime that is used is individual contributor specific implementation details and shouldn't be tracked in the repo.

The operator suggested that personal runtime selection could live in the running user's home directory. The exact optional preference path and precedence in [personal runtime preferences](../specifications/personal-runtime.md) are agent-designed mechanisms, not operator-selected syntax.

## Runtime-extension clarification — operator authority

In a subsequent design discussion, the operator clarified:

> Factoryctl itself needs to be a single static binary. It shouldn't be shelling out to git, but it can shell out to AI agents, which may then shell out to git.

The operator described factoryctl as a thin CLI/Web-GUI wrapper over an existing harness such as Codex, Claude Code or pi-coding-agent. Humans navigate clean rendered Markdown; GUI action buttons invoke the selected harness to execute instructions. Factoryctl exposes MCP tools to govern document reads and writes. Different document types need their own lifecycle behavior and tools, decoupled from the binary and extensible by users at runtime.

The operator authorized documenting and merging the design, then explicitly excluded custom code execution inside factoryctl. The earlier proposed WebAssembly mechanism is rejected, not deferred. Arbitrary code belongs in the browser UI (for example a custom drag-and-drop kanban board), or in work performed by an AI harness after a prompt. Document types use declarative rules and prompt-backed lifecycle actions. This authorizes design documentation and normal Git delivery, not a claim of implemented features. Go remains the selected design direction; exact metadata syntax, dependency versions and invocation flags require implementation qualification. Codex remains the first required adapter; naming Claude Code and pi establishes the adapter boundary, not a claim that either is integrated.

This clarification strengthens R04 to prohibit even read-only Git subprocesses and adds R26–R28. Core scheduling and autonomous agent dispatch remain excluded. The later minimal-interface directive below narrows type-specific validation to document structure; lifecycle interpretation belongs to the harness. MCP tools enforce their own operations; restricting bypass requires qualified harness permissions or OS isolation. A prompt alone cannot force an agent to use MCP.

## Minimal interface directive — controlling operator requirement

The operator emphasized that factoryctl must be simple and minimal: all complexity of what a software factory or another knowledge-work factory means belongs to the AI harness interpreting Markdown supplied by factoryctl. Factoryctl makes that knowledge viewable and exposes tools that keep documents standardized. The operator requested that this principle be documented and the repository updated.

This supersedes the earlier proposal for factoryctl-enforced lifecycle transitions and state-conditioned requirements. Standardization means document syntax, required fields, field types/enums, placement, links and safe file updates. Workflow meaning, readiness, review sufficiency, approvals and permission to advance are interpreted by the agent from Markdown. A schema enum may list state labels; factoryctl does not decide which transition is appropriate. Prompts and instructions are content, not a rule language executed by the backend.

The [minimal-interface decision](../decisions/minimal-interface.md) is the governing scope test for future implementation. Browser code for custom UI and harness execution through prompts remain supported. No workflow engine, policy evaluator, business-rule DSL or backend custom-code execution is authorized by this directive.

## Discovery depth directive — direct operator source

> For the project discovery, keep the recursion limit to 1 by default. Meaning only look in the folders in the folder you are currently in, not further levels of recusion.

The default scan inspects immediate subfolders only. The selected directory is depth 0; its immediate subfolders are depth 1. Existing selection of the current Git working tree or an explicit path remains available.

## Project collections and user interaction — operator requirements

The operator requires three exercised use cases throughout implementation and testing: launch inside a single-project Git repository; inside a monorepo containing multiple projects with independent factory configuration/state and cross-links; and inside a parent folder containing independent Git repositories with their own factory configuration/state and cross-links. All use the same UI, including a collection containing one project in the first case. “Configuration/state” is represented by each project's existing Markdown definitions, records and optional assets; this request does not prescribe a new configuration file.

The CLI and web UI must let the selected coding harness ask the user questions through a tool, supporting free text, multiple choice and back-and-forth conversation. The operator requested MCP Apps exploration so the LLM can present interactive UI. They subsequently proposed that sufficiently strong MCP Apps support could justify deprecating the CLI and supporting only a rich web user interface, including UI generated at runtime rather than only predefined elements. This is a conditional product direction to qualify, not an assertion of protocol/harness support or immediate CLI retirement. A web launcher and MCP server entry point are distinct from an interactive terminal UI.

## Release demonstration and traceability — operator requirements

For every release, provide a Playwright video demonstrating all implemented features, including older features to prove they still work; omit only features deprecated and removed. The recording should resemble human use, with smooth scrolling and enough pauses for orientation. Overlay what is being demonstrated/evaluated and highlight the element or area needing attention.

Maintain a full requirements traceability matrix covering everything the operator declared a requirement, linking it to implementation evidence and ongoing evidence that it continues to function. Historical evidence alone is not current-release proof.
