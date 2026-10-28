# factoryctl

A local document explorer for Git-backed Markdown knowledge bundles, with reusable document services and automation commands. It discovers projects, reads and validates their `factory/` documents, and makes revision-checked edits. Factoryctl never runs Git, a shell, an agent, or workflow instructions.

This is a partial product with document services and a read-only browser explorer, not a complete product release. The living [specifications](factory/specifications/product.md) describe the larger intended product.

## Build and try it

Requires Go 1.27 or later to build. The executable needs no Go installation, network connection, agent, configuration file, or `.factory/` folder to operate.

```sh
CGO_ENABLED=0 go build -o /tmp/factoryctl ./cmd/factoryctl
/tmp/factoryctl --path . web --listen 127.0.0.1:8080
```

Open `http://127.0.0.1:8080` to explore the planning bundle. The browser is read-only; document automation commands remain available below.

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

## Validation and storage

The parser uses YAML and CommonMark ASTs. Untrusted document, metadata, ignore and template-asset reads reject special files before opening and recheck the opened handle; Unix opens also use nonblocking/no-follow flags to guard replacement races. Nonfinite YAML numbers in generic metadata remain visible as `.nan`/`.inf`/`-.inf` strings in JSON with a warning; original bytes and values are retained for validation. Output encoding or write failures return exit 4. It rejects duplicate keys, aliases, custom tags, invalid UTF-8, conflict markers, excessive nesting and oversized documents. Definitions check structural field kinds, required values, enum membership, placement and typed references. Links inside code examples are ignored. Heading anchors use lowercase Unicode letters/numbers, retain hyphens/underscores, replace whitespace with hyphens, and suffix duplicates. Explicit HTML IDs are recognized without executing HTML. External URLs are not fetched; `--require-complete-links` fails when any links remain unchecked. `okf` provides the implemented syntax/base-metadata checks, not a claim of exhaustive upstream optional-metadata conformance.

Every mutation validates a proposed graph, including incoming links and affected definitions/templates. Unrelated existing errors do not block a healthy edit. Local definition overrides remain isolated per project; disabled types remain readable and structurally checked. Embedded base definitions are virtual until initialization materializes them. Templates accept typed inputs and literal `include` paths with bounded depth/output; no functions can execute code, fetch URLs or read arbitrary files.

Writes use the stable `control/write-coordination.md` native advisory lock, SHA-256 revision comparison, root-scoped filesystem operations, a flushed temporary sibling, and atomic replacement (exclusive publication for creation). Symlinks, traversal and special files are rejected. The anchor cannot be edited/deleted through ordinary record operations. Existing bundles need that anchor before mutation; init previews all conflicting defaults and never overwrites them. Interrupted init reports created paths; it is not a multi-file transaction.

External editors and Git do not participate in advisory locks. Dependency rechecks reduce races but cannot provide transactions with non-cooperating writers or across project bundles. Backlinks and validation cover loaded projects only; an unknown external bundle URL cannot be distinguished from an ordinary external URL. No absent-checkout registry exists. Do not use concurrent Git checkout/reset while writing.

Limits: 1 MiB per document, 256 KiB frontmatter, YAML nesting 32, template include depth 16, include budget 1024, prompt output 128 KiB. Documents exceeding limits remain diagnostic entries. Strict bundle validation reports non-Markdown files, including Finder `.DS_Store` files; the CLI never removes them.

## Remaining work

GUI editing/boards/conversations, MCP, runtime launch/questions, watcher refresh, external network link checks, custom view export, extension tool registration and release qualification are unimplemented. Related CLI entry points return `unsupported_capability`. Reads/commands reload disk; there is no persistent index or background process.

Additional bounded gaps: full Git wildmatch/global ignore configuration; complete optional OKF provenance/footnote semantics; dependency-cycle reporting; precise source coordinates for every structural diagnostic; registry-based identification of absent cross-project checkouts; interactive input collection. Templates reject dynamic include paths and inline `define`/`template` constructs; use literal `include` instead. Factory identity changes should be made externally and followed by rediscovery. Software starter defaults are opt-in. The read-only browser assets are embedded. No throughput, full release, Windows/Linux execution, or full-GUI qualification claim is made.

## Checks

```sh
go test ./...
go test -race ./...
go vet ./...
CGO_ENABLED=0 go build -o /tmp/factoryctl ./cmd/factoryctl
```

Tests cover S1–S3 collections, independent definitions, mapped cross-project links, safe target selection, parser hazards, source preservation, template bounds, proposed-graph rejection, external edits, cooperating processes and lock release after process termination. Only execution on the current macOS ARM64 host is qualified in this milestone; other target builds are compilation evidence.

## Milestone review gate

Every milestone delivery requires a current Playwright walkthrough reconciled against the complete [requirements matrix](factory/specifications/requirements-traceability.md). Follow [the project instructions](AGENTS.md) and [demonstration contract](factory/specifications/release-evidence.md). `make review-check` checks row accounting and evidence freshness; `make review-gate` additionally requires every row to pass end to end. The current partial product has unmet full-GUI/MCP/runtime and other requirements, so full-matrix qualification remains incomplete even when its review artifacts are valid.

Review the [milestone demo](/Users/Shared/factoryctl/artifacts/demo/latest.mp4): one browser walkthrough followed by fresh CLI regression scenes, with requirement IDs and outcome explanations on the actual scenes. This same link is updated at every milestone. The [row-by-row index](artifacts/milestone-review/20261028/MATRIX.md) documents remaining gaps; full-matrix qualification is incomplete.

Generated `artifacts/` content and videos stay local and are not included in Git history. After cloning, run `make install-hooks` to enable the commit/push guards. `make repository-check` validates the index and guard tests; CI also checks full history. Generate local evidence before running `make review-check`. Checker unit tests use self-contained synthetic fixtures.

## Explore the current plan in your browser

```sh
make build
./bin/factoryctl --path . web --listen 127.0.0.1:8080
```

Open `http://127.0.0.1:8080`. Browse real Markdown documents, filter by type or recorded intent state, inspect sources and authorization scope, follow same-checkout local links and switch to raw Markdown. “Accepted” is a claim recorded in a document, not a decision made by the application. Missing provenance stays unresolved. The server is loopback-only and read-only; it takes a startup snapshot, so restart it after document changes.

This slice does not implement document editing, conversations, MCP Apps, canonical cross-project URL navigation/backlinks or acceptance submission. Use [the current plan](factory/designs/implementation-plan.md) and [accepted working method](factory/decisions/intent-and-review.md) to distinguish current behavior from proposed continuation. Review videos and generated evidence remain local.
