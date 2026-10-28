# factoryctl

A local CLI, browser and MCP interface for Git-backed Markdown knowledge bundles, with shared document services and opt-in native conversations. It discovers projects, reads and validates their `factory/` documents, and makes revision-checked edits. Document operations never launch Git, a shell or an agent. Explicit conversation mode can start the selected installed Codex runtime and attach project-scoped MCP tools; workflow prose does not authorize actions.

The completion candidate at `a7c22ba847d58e8f4ece6a22ee9bb73302a64c02` implements document editing, boards, runtime-loaded tools and custom views alongside the existing conversations. Independent evidence review supports the implemented scope; exact-candidate approval and delivery remain pending. Strict static linkage on Darwin/Windows is unresolved under U045, so full qualification and project completion are not claimed. The living [specifications](factory/specifications/product.md) remain the desired-state authority.

## Build and try it

Requires Go 1.27 or later to build. Document operations need no Go installation, network connection, agent, configuration file, or `.factory/` folder. Conversations additionally require the selected installed runtime and its own working authentication/network access.

```sh
CGO_ENABLED=0 go build -o /tmp/factoryctl ./cmd/factoryctl
/tmp/factoryctl --path . web --listen 127.0.0.1:8080
```

Open `http://127.0.0.1:8080` to browse and edit the planning bundle, inspect backlinks, and use enum boards. Document operations work without a runtime; conversations require explicit runtime selection.

Discovery scans the selected directory and its immediate children only. A monorepo's projects remain separate even when they share one Git root. `--project ID` selects a project; duplicate identities additionally require the `--checkout` token reported by `projects list`. CWD within a source subtree can fall back to its nearest ancestor project; explicit `--path` never does.

Initialize an existing Git working tree, or an explicitly selected directory inside one:

```sh
factoryctl --path /path/to/project init --project-id my-project --starter software
factoryctl --path /path/to/project records create --type Task --slug first-task \
  --input 'title=First task' --input 'summary=Describe the work'
factoryctl --path /path/to/project records show tasks/first-task.md --json
factoryctl --path /path/to/project records update tasks/first-task.md \
  --expected-revision HASH_FROM_SHOW --set task_state=done
```

Inputs accept JSON values and otherwise remain literal strings. All inputs must be supplied explicitly; missing input fails without prompting, including in `--non-interactive` mode. A valid state enum change is accepted regardless of workflow prose. Structural validity does not establish approval or completion.

Use `records replace PATH --expected-revision HASH --file FILE` for complete-document edits (`--file -` reads stdin), and `records delete PATH --expected-revision HASH` for deletion. Reads expose exact source, metadata, SHA-256 revisions, loaded outgoing links and backlinks. Unknown fields and types remain readable. Scalar field updates splice source bytes without rewriting unrelated content; complex YAML representations require explicit whole-document replacement and return `lossless_edit_unavailable`.

Other implemented operations: `templates diff`, `actions list --record PATH`, `doctor`, list filtering/pagination, and `validate --level okf|factory`. Global flags may precede or follow commands. `--json` puts structured results/errors on stdout; human diagnostics go to stderr. Exit codes are 0 success, 1 invalid structure, 2 invalid arguments/selection, 3 stale revision/existing destination, 4 I/O or unsupported capability, and 5 incomplete validation when requested.

## Conversations and questions

Start the browser with an explicit runtime selection:

```sh
/tmp/factoryctl --path . web --runtime codex --listen 127.0.0.1:8080
```

Open the **Conversation access** URL printed by the process. Its fragment grants access to the trusted conversation controls; keep that URL private. Select a project, optionally include the selected document, and send a message. Questions accept text or choices with custom answers; follow-up messages keep the same native thread. Stop turn interrupts the runtime. Generated views require an explicit trust action and retain a plain-text fallback. A view can draft an answer, but you review and submit it through the trusted host controls. Browser document editing uses the shared revision-checked structural writer. Starting `web` without `--runtime` or `--preferences` enables browsing and document editing without reading runtime preferences.

The chat composer includes **Model** and **Reasoning** selectors using the installed runtime's advertised capabilities. Choose before starting a chat or while an existing conversation is idle; the choice applies to the next message. Enter sends and Shift+Enter inserts a newline.

Open **Defaults** to explicitly save the selected pair globally or for this project. Personal defaults live in `~/.factoryctl/chat-settings.json`, outside project Git. Project overrides use canonical local project roots, so separate projects with the same ID do not share defaults. Reset removes the chosen override; new chats inherit the remaining defaults. Existing conversations keep their own selection. Choosing a model or sending a message does not save defaults. Factoryctl creates the private configuration file only when you choose Save or Reset, and reports stale-write conflicts rather than overwriting a newer save.

The retained interactive terminal uses the same conversation services:

```sh
/tmp/factoryctl --path . chat --runtime codex --prompt 'Help me explore the next milestone.'
```

Run `chat` in an actual terminal: piped stdin, `--json` and `--non-interactive` are rejected. Enter text or displayed choice IDs (comma-separated for multiple choices), or a custom answer. `/decline` declines a question, `/cancel` interrupts the active turn, and `/quit` closes the conversation. Generated views have a text fallback here. Ordinary document commands never prompt.

Both `web` and `chat` accept `--runtime-executable /absolute/path/to/codex`, `--model NAME`, `--reasoning-effort VALUE`, `--preferences FILE` and `--chat-settings /absolute/path/to/file.json`. The last flag selects an alternate personal settings file for services or isolated testing. Explicit conversation selections override invocation flags, then saved project defaults, saved global defaults, legacy personal preferences and native defaults. A model flag clears inherited reasoning unless an effort flag is also supplied; an effort-only flag keeps the inherited model. Runtime options require `--runtime` or an explicit `--preferences` file. Only the Codex adapter is implemented; unsupported selections fail without substitution.

On an explicit runtime action, the default preference location is the current OS user's `~/.factoryctl/preferences.md`. Factoryctl does not create or update it. An optional file looks like:

```yaml
---
version: 1
runtime: codex
---
Personal runtime selection; keep outside project Git.
```

Optional frontmatter fields are `executable` (absolute path), `model` and `reasoning_effort`. Missing default preferences are allowed; an explicitly selected missing file or malformed preferences fail before launch. Credentials stay with the runtime.

For an MCP client, configure a stdio server invocation:

```sh
/tmp/factoryctl --path /absolute/path/to/project mcp
```

MCP stdout/stdin carry JSON-RPC, never terminal answers; do not add `--json`. Select exactly one initialized project. Document tools expose scoped reads, validation and revision-checked writes. `factory_user_ask`, `factory_user_wait` and `factory_ui_present` require an attached trusted conversation surface; standalone MCP returns `interaction_unavailable` for those tools. Factoryctl's conversation launch attaches its own scoped MCP child through invocation-only configuration, without writing global runtime configuration. Tools cannot submit user answers.

See the [interaction capability report](factory/sources/interaction-capabilities.md) for observed runtime and terminal evidence, first-party sandbox limits and unqualified external MCP Apps interoperability. The interactive CLI remains available; no retirement decision is implied.

## Validation and storage

The parser uses YAML and CommonMark ASTs. Untrusted document, metadata, ignore and template-asset reads reject special files before opening and recheck the opened handle; Unix opens also use nonblocking/no-follow flags to guard replacement races. Nonfinite YAML numbers in generic metadata remain visible as `.nan`/`.inf`/`-.inf` strings in JSON with a warning; original bytes and values are retained for validation. Output encoding or write failures return exit 4. It rejects duplicate keys, aliases, custom tags, invalid UTF-8, conflict markers, excessive nesting and oversized documents. Definitions check structural field kinds, required values, enum membership, placement and typed references. Links inside code examples are ignored. Heading anchors use lowercase Unicode letters/numbers, retain hyphens/underscores, replace whitespace with hyphens, and suffix duplicates. Explicit HTML IDs are recognized without executing HTML. External URLs are not fetched; `--require-complete-links` fails when any links remain unchecked. `okf` provides the implemented syntax/base-metadata checks, not a claim of exhaustive upstream optional-metadata conformance.

Every mutation validates a proposed graph, including incoming links and affected definitions/templates. Unrelated existing errors do not block a healthy edit. Local definition overrides remain isolated per project; disabled types remain readable and structurally checked. Embedded base definitions are virtual until initialization materializes them. Templates accept typed inputs and literal `include` paths with bounded depth/output; no functions can execute code, fetch URLs or read arbitrary files.

Writes use the stable `control/write-coordination.md` native advisory lock, SHA-256 revision comparison, root-scoped filesystem operations, a flushed temporary sibling, and atomic replacement (exclusive publication for creation). Symlinks, traversal and special files are rejected. The anchor cannot be edited/deleted through ordinary record operations. Existing bundles need that anchor before mutation; init previews all conflicting defaults and never overwrites them. Interrupted init reports created paths; it is not a multi-file transaction.

External editors and Git do not participate in advisory locks. Dependency rechecks reduce races but cannot provide transactions with non-cooperating writers or across project bundles. Backlinks and validation cover loaded projects only; an unknown external bundle URL cannot be distinguished from an ordinary external URL. No absent-checkout registry exists. Do not use concurrent Git checkout/reset while writing.

Limits: 1 MiB per document, 256 KiB frontmatter, YAML nesting 32, template include depth 16, include budget 1024, prompt output 128 KiB. Documents exceeding limits remain diagnostic entries. Strict bundle validation reports non-Markdown files, including Finder `.DS_Store` files; the CLI never removes them.

## Qualification status and limits

The current [local review package](artifacts/completion-review/a7c22ba-final-v2/requirements-coverage.json) proposes 105 passing atomic rows and leaves U045 blocked. This is prepared evidence, not independent final acceptance. Its source-bound run records 169 named race-test passes; all six macOS/Linux/Windows ARM64/x86-64 targets were built and inspected. Only the current Mac executed the product. Linux static linkage was observed; Darwin imports system libraries and Windows imports platform DLLs. Whether that packaging satisfies the operator's single-static-binary requirement remains unanswered, with no inferred exception.

The declared 100-repository, 10,000-node, 100 MiB benchmark measured loader refresh at about 2.545 seconds, missing the proposed one-second target. That target is an LLM engineering proposal, not an operator-approved performance limit. The measurement excludes browser/harness memory and end-to-end HTTP/polling visibility; it is not a guaranteed GUI refresh latency.

External MCP Apps interoperability remains unqualified; first-party generated views do not establish it, and the interactive CLI is retained. External network link checks are not provided, and backlinks cover loaded projects only. Conversation/question state remains process-local; native history stays with the runtime. Factory identity changes require an external edit followed by rediscovery. Templates intentionally require literal include paths and bounded expansion. No foreign-platform execution, final full-matrix acceptance or release publication is claimed.

## Checks

```sh
go test ./...
go test -race ./...
go vet ./...
CGO_ENABLED=0 go build -o /tmp/factoryctl ./cmd/factoryctl
```

Tests cover S1–S3 collections, independent definitions, mapped cross-project links, safe target selection, parser hazards, source preservation, template bounds, proposed-graph rejection, external edits, cooperating processes and lock release after process termination. Only execution on the current macOS ARM64 host is qualified in this milestone; other target builds are compilation evidence.

## Milestone review gate

Every milestone delivery requires a current Playwright walkthrough reconciled against the complete [requirements matrix](factory/specifications/requirements-traceability.md). Follow [the project instructions](AGENTS.md) and [demonstration contract](factory/specifications/release-evidence.md). `make review-check` checks row accounting and evidence freshness; `make review-gate` additionally requires every row to pass end to end. The prepared completion matrix retains U045 as blocked, and exact-candidate delivery review remains separate. Row accounting and valid artifacts cannot establish full qualification while that mandatory requirement remains unresolved.

Review the [milestone demo](/Users/Shared/factoryctl/artifacts/demo/latest.mp4): the stable local handoff link for the combined walkthrough. The versioned [completion package](artifacts/completion-review/a7c22ba-final-v2/results.json) and [current row accounting](artifacts/completion-review/a7c22ba-final-v2/requirements-coverage.json) bind the prepared film, executed checks and remaining gap to the candidate. Earlier milestone artifacts remain historical evidence. Independent media review supports this partial package; full-matrix qualification remains blocked at U045.

Generated `artifacts/` content and videos stay local and are not included in Git history. After cloning, run `make install-hooks` to enable the commit/push guards. `make repository-check` validates the index and guard tests; CI also checks full history. Generate local evidence before running `make review-check`. Checker unit tests use self-contained synthetic fixtures.

## Explore the current plan in your browser

```sh
make build
./bin/factoryctl --path . web --listen 127.0.0.1:8080
```

Open `http://127.0.0.1:8080`. Browse real Markdown documents, filter by type or recorded intent state, inspect sources and authorization scope, follow resolved local or loaded canonical cross-project links and switch to raw Markdown. “Accepted” is a claim recorded in a document, not a decision made by the application. Missing provenance stays unresolved. The server is loopback-only. It reloads external document and definition changes without a restart, while preserving active editor drafts.

Browser editing, boards, canonical cross-project navigation/backlinks and opt-in conversations are implemented. The application displays recorded acceptance claims; viewing a document or submitting a chat answer does not accept the project. External MCP Apps interoperability remains unqualified. Use [the current plan](factory/designs/implementation-plan.md) and [accepted working method](factory/decisions/intent-and-review.md) to distinguish current behavior from proposed continuation. Review videos and generated evidence remain local.


## Web editing, boards and custom views

Run `factoryctl --path . gui` (or `web`) to serve the loopback browser interface. The same project collection supports a single project, monorepo project roots and sibling repositories. Use the project/type/intent filters, overview, linked documents and backlinks. The browser reflects external edits without restart; it retains an unsaved draft when a concurrent edit conflicts.

**Edit fields** changes declared fields through shared structural validation; **Edit Markdown** handles complete documents and complex metadata. Saves carry document and effective-definition hashes. New documents use the selected type's template and hash; deletion refuses referenced records. Enum boards support drag/drop and keyboard selection, including a diagnostic column for invalid/missing values. Next/previous controls keep document lists and boards on fixed-size pages; records beyond the first page remain reachable. These operations work without a runtime and never interpret workflow prose as approval.

Definitions may provide `views.detail_template` and `views.assets` under `.factory/views/`. Enable a repository view explicitly to run its packaged HTML/CSS/JavaScript in an isolated frame. Trust binds the definition and transitive asset content; changes require trust again. The scoped bridge can request validated edits or preview declared prompt actions; sending a prompt to the harness remains an explicit host action. Missing or invalid assets retain the generic view. `templates export --kind record|prompt|view` materializes absent defaults and preserves existing customizations.

## Declarative tools and actions

`extensions list --json` discovers runtime-loaded namespaced document bindings. `extensions call NAME --input KEY=VALUE ... --json` invokes the same service as HTTP/MCP. Bindings support fixed read/create/update/delete operations; no backend code, executable handlers, expressions, state transitions or approval rules are registered. Read the exact schema and examples in [extensions](factory/specifications/extensions.md) and [implementation evidence](factory/sources/completion-extensions-brief.md).

Use `actions list --record PATH`, `actions render PROMPT --record PATH`, or explicitly launch `actions run PROMPT --record PATH --runtime codex` in a terminal. Rendering never starts a runtime. MCP document writers require `expected_definition` in addition to document revision where applicable; obtain it from reads or definition discovery. Tool discovery reloads current definitions; clients caching the tool list must reconnect. MCP and CLI listings expose snapshot-bound continuation cursors. Legacy CLI writers accept optional `--expected-definition` and always validate the current graph inside the bundle lock.

## Delivery observations

Task/plan views show delivery separately from recorded work and intent. Project Markdown owns the explicit target, criteria and attributable acceptance assessment. A plan's `delivery_ref` points to the same task; both display the same result. Optional ignored `.factory/delivery/` observations bind those claims to a scoped content manifest, observer and verification time. A local content match does not authenticate assessor authority. Missing criteria, missing evidence and changed content remain explicit. Branch/remote target observations are externally recorded and remain unverified by factoryctl; no Git or hosting calls are made. See [the completion brief](factory/tasks/project-completion.md#delivery-observation-implementation-brief-and-reviewed-contract) for the exact version-1 schema.

## Qualification and server use

The binary includes its default browser assets and needs no application-side Node/Python installation. For a trusted server or Kubernetes pod, run the same loopback HTTP server behind an operator-configured authenticated proxy/sidecar; mount existing checkouts and optional personal runtime preferences outside project Git. Factoryctl does not deploy infrastructure or provide an accounts database. SIGTERM closes the server and owned conversations. CI uses noninteractive CLI commands.

[Local qualification tools](scripts/qualification/README.md) produce versioned scale, platform packaging and release-policy evidence. Only the current Mac executes product tests; Linux/Windows binaries are cross-built and inspected, not executed. CGO-disabled packaging is not evidence of fully static Darwin/Windows linkage; U045 remains blocked pending the operator's answer. `make check`, `make repository-check`, and the current requirements/video gates remain required. Recordings and all generated artifacts remain local and ignored. No release has been published; the initial release target is 0.1.0, and the current development binary reports 0.1.0-dev.
