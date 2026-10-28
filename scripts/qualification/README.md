# Local qualification tooling

Run on the current Mac from the repository root:

```sh
python3 -m unittest discover -s scripts/release -p 'test_*.py'
CGO_ENABLED=0 GOCACHE=/private/tmp/factoryctl-go-cache GOTMPDIR=/private/tmp go test ./scripts/qualification/...
python3 scripts/qualification/run.py --output artifacts/completion-qualification/0.1.0/attempt-001
```

The output path must be new. Failed attempts, logs, corpus, target binaries and
reports stay local under the ignored artifacts tree. Reusable corpora remain
in the report's unique `/private/tmp/factoryctl-qualification-*` path: placing
them under ignored artifacts would correctly exclude them from discovery.
No publication or version
tagging occurs. The report binds all measured source inputs before and after;
if another worker changes them, rerun after source freeze. `--skip-scale` and
`--skip-builds` support diagnosed partial reruns, whose scope remains explicit.

The scale runner generates exactly 100 repositories, 10,000 Markdown nodes and
100 MiB of UTF-8 Markdown by default, including metadata and two custom type
definitions per project. Every domain record links to two others. Definitions
provide research certainty and synthetic opportunity pipeline boards. Sources,
assumptions, inference and fictional awards are visibly attributed. These are
scratch Git repositories with no commits, not real research or procurement.
Generate a smaller reusable browser corpus with:

```sh
CGO_ENABLED=0 GOCACHE=/private/tmp/factoryctl-go-cache GOTMPDIR=/private/tmp go run ./scripts/qualification/scale --corpus /private/tmp/factoryctl-breadth-001 --repositories 2 --nodes 30 --bytes 131072 --generate-only
```

Measurements call the actual project discovery and shared application loader.
Cold means the first application load in that process; OS disk caches are not
evicted. Query p95 covers 100 in-memory raw-document scans using the actual
shared CLI/MCP `app.QueryMatcher`, including matcher compilation per query,
excluding transport/load/serialization. Refresh uses production `LoadCached`
after one external edit, comparing current bytes and revalidating. A fresh full
reload follows as a separate comparison. Browser debounce is not measured.
RSS uses instantaneous `ps` samples in bytes, not peak RSS; Go heap
allocation is reported separately. Threshold booleans report observed targets
without converting unmeasured watcher behavior into a pass.

Six CGO-disabled builds embed their application assets. The inspector uses Go
ELF/Mach-O/PE readers, target build metadata and embedded-marker checks. Linux
must have no ELF interpreter or dynamic imports. Darwin system interfaces and
Windows platform DLLs are reported explicitly; these are not claims of full
static linkage. Foreign binaries are never executed. Packaging evidence does
not qualify browser behavior, server shutdown, runtime interactions or release
acceptance.

Release policy tooling consumes an explicit JSON array of full commit messages:

```sh
python3 scripts/release/version.py messages.json --previous 0.1.0 --proposed 0.2.0
```

Omit `--previous` for the initial 0.1.0 candidate. After that, breaking markers
on any type or `feat` yield a minor increment, `fix`/`perf` a patch, other types
no release. Aggregate precedence is applied once. Nonconventional messages,
incorrect proposals and every major version 1+ are rejected. There is no bypass
flag for major releases; an authorized transition requires a separately reviewed
policy/tool change. Caller owns the exact reachable commit range and publication.
