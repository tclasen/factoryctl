// scale qualifies the actual shared loader against a disposable deterministic corpus.
// It is development tooling, never included in the factoryctl executable.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/tclasen/factoryctl/internal/app"
	"github.com/tclasen/factoryctl/internal/project"
)

func put(path, data string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(data), 0644)
}

func recordPath(n int) string {
	if n%2 == 0 {
		return fmt.Sprintf("research/fact-%03d.md", n)
	}
	return fmt.Sprintf("opportunities/bid-%03d.md", n)
}

func corpus(root string, repos, nodes, size int) error {
	if repos < 1 || nodes < repos*7 || nodes%repos != 0 {
		return fmt.Errorf("nodes must be divisible by repositories, with at least seven per repository")
	}
	if _, err := os.Lstat(root); !os.IsNotExist(err) {
		return fmt.Errorf("corpus destination must not exist (preserve previous attempts): %s", root)
	}
	contents := map[string]string{}
	per := nodes / repos
	for r := 0; r < repos; r++ {
		base := filepath.Join(root, fmt.Sprintf("repository-%03d", r))
		if err := put(filepath.Join(base, ".git/HEAD"), "ref: refs/heads/main\n"); err != nil {
			return err
		}
		for _, dir := range []string{"objects", "refs/heads"} {
			if err := os.MkdirAll(filepath.Join(base, ".git", dir), 0755); err != nil {
				return err
			}
		}
		contents[filepath.Join(base, "factory/factory.md")] = fmt.Sprintf("---\ntype: Factory\ntitle: Synthetic knowledge work %03d\nproject_id: synthetic-%03d\n---\n# Synthetic knowledge work\n\nScratch data only. No real research or procurement claims.\n", r, r)
		contents[filepath.Join(base, "factory/control/write-coordination.md")] = "---\ntype: Factory Control\ntitle: Synthetic writer coordination\n---\n# Writer coordination\n\nCooperative writers lock this anchor.\n"
		contents[filepath.Join(base, "factory/prompts/review-evidence.md")] = "---\ntype: Prompt\ntitle: Review synthetic evidence\napplies_to: [Research Fact, Opportunity]\ninputs: {}\nconversation: fresh\n---\n# Review synthetic evidence\n\nPreviewing this prompt does not launch an agent or contact anyone.\n\n# Template\n\n```factory-prompt\nReview {{ .Record.Title }} in {{ .Project.Name }}. Separate cited synthetic claims, assumptions and inference. Ask the human a clarifying question before changing the record. Do not submit bids or contact external parties.\n{{ .Record.Body }}\n```\n"
		for _, typ := range []struct{ name, dir, field, values string }{{"Research Fact", "research", "certainty", "[cited, assumption, inference]"}, {"Opportunity", "opportunities", "stage", "[discovered, pursuing, submitted, tracking, awarded]"}} {
			contents[filepath.Join(base, "factory/definitions", typ.dir+".md")] = fmt.Sprintf("---\ntype: Record Definition\ntitle: %s definition\ndefines: %s\ndirectory: %s\nfields:\n  %s:\n    kind: enum\n    required: true\n    values: %s\nviews:\n  board:\n    field: %s\n---\n# %s\n\nThe harness interprets workflow prose. Stage is a recorded claim, never an approval.\n", typ.name, typ.name, typ.dir, typ.field, typ.values, typ.field, typ.name)
		}
		count := per - 5
		for n := 0; n < count; n++ {
			typ, field, value := "Research Fact", "certainty", []string{"cited", "assumption", "inference"}[(n/2)%3]
			if n%2 == 1 {
				typ, field, value = "Opportunity", "stage", []string{"discovered", "pursuing", "submitted", "tracking", "awarded"}[(n/2)%5]
			}
			next, prev := recordPath((n+1)%count), recordPath((n+count-1)%count)
			text := fmt.Sprintf("---\ntype: %s\ntitle: Synthetic %s %03d\n%s: %s\nsources:\n  - id: synthetic\n    resource: https://example.invalid/synthetic-source\n    author: human:synthetic-fixture\nfactory:\n  format: 1\n---\n# Evidence\n\nSynthetic fixture %03d. Cited fact: fixture budget is 42 units (invented source). Assumption: review is possible. LLM inference: inspect eligibility. These are synthetic claims, not external facts.\n\n[Next](../%s#evidence) and [Previous](../%s#evidence).\n\nHarness action: compare the source and assumption; ask the human before any external contact. No actual bid, submission or award occurred.\n", typ, field, n, field, value, n, next, prev)
			contents[filepath.Join(base, "factory", recordPath(n))] = text
		}
	}
	keys := make([]string, 0, len(contents))
	total := 0
	for p, s := range contents {
		keys = append(keys, p)
		total += len(s)
	}
	sort.Strings(keys)
	if total > size {
		return fmt.Errorf("declared corpus size %d smaller than metadata %d", size, total)
	}
	remain := size - total
	for i, p := range keys {
		extra := remain / (len(keys) - i)
		remain -= extra
		if len(contents[p])+extra > 1<<20 {
			return fmt.Errorf("document exceeds declared 1 MiB maximum")
		}
		if err := put(p, contents[p]+strings.Repeat("x", extra)); err != nil {
			return err
		}
	}
	return nil
}

func percentile(samples []float64) float64 {
	sort.Float64s(samples)
	return samples[(95*len(samples)+99)/100-1]
}
func elapsed(start time.Time) float64 { return float64(time.Since(start)) / float64(time.Millisecond) }
func rss() int64 {
	// ps reports KiB on the execution-qualified Mac. This is an instantaneous
	// process RSS sample, not peak RSS and not the Go heap estimate.
	out, err := exec.Command("ps", "-o", "rss=", "-p", strconv.Itoa(os.Getpid())).Output()
	if err != nil {
		return -1
	}
	n, err := strconv.ParseInt(strings.TrimSpace(string(out)), 10, 64)
	if err != nil {
		return -1
	}
	return n * 1024
}

func run() error {
	root := flag.String("corpus", "", "new scratch corpus path (must not exist)")
	repos := flag.Int("repositories", 100, "scratch Git repositories")
	nodes := flag.Int("nodes", 10000, "total Markdown nodes including definitions and Factory nodes")
	size := flag.Int("bytes", 100<<20, "exact total Markdown bytes")
	generate := flag.Bool("generate-only", false, "generate reusable browser fixtures without measurement")
	flag.Parse()
	if *root == "" {
		return fmt.Errorf("--corpus is required")
	}
	if err := corpus(*root, *repos, *nodes, *size); err != nil {
		return err
	}
	result := map[string]any{"repositories": *repos, "nodes": *nodes, "markdown_bytes": *size, "max_file_bytes": 1 << 20, "links_per_domain_record": 2, "go_version": runtime.Version(), "host": runtime.GOOS + "/" + runtime.GOARCH, "corpus": *root, "execution_scope": "current host only", "cold_definition": "first application load in this process; OS disk cache not evicted"}
	if *generate {
		return json.NewEncoder(os.Stdout).Encode(result)
	}
	start := time.Now()
	collection, err := project.Discover(*root, true)
	if err != nil {
		return err
	}
	result["discovery_ms"] = elapsed(start)
	if len(collection.Projects) != *repos || len(collection.Diagnostics) != 0 {
		return fmt.Errorf("discovery got %d projects, %d diagnostics", len(collection.Projects), len(collection.Diagnostics))
	}
	start = time.Now()
	snap, err := app.Load(collection.Projects)
	if err != nil {
		return err
	}
	result["cold_load_ms"] = elapsed(start)
	result["rss_after_load_bytes"] = rss()
	count := 0
	for _, b := range snap.Bundles {
		count += len(b.Docs)
	}
	if count != *nodes {
		return fmt.Errorf("loaded %d nodes, expected %d", count, *nodes)
	}
	if len(snap.Diagnostics) != 0 {
		return fmt.Errorf("qualification corpus has diagnostics: %v", snap.Diagnostics[0])
	}
	validLinks, uncheckedLinks := 0, 0
	for _, link := range snap.Links {
		if link.Status == "valid" {
			validLinks++
		}
		if link.Status == "unchecked" {
			uncheckedLinks++
		}
	}
	result["valid_links"] = validLinks
	result["unchecked_external_source_links"] = uncheckedLinks
	// Exercise the shared CLI/MCP cached raw-document query matcher directly.
	// Transport discovery/load/serialization and browser title-only queries are
	// excluded and must not be inferred from this measurement.
	samples := []float64{}
	matches := 0
	for i := 0; i < 100; i++ {
		start = time.Now()
		matches = 0
		query := []string{"synthetic", "fixture budget", "absent-query-987654"}[i%3]
		matcher := app.DocumentMatcher(query)
		for _, b := range snap.Bundles {
			for _, d := range b.Docs {
				if matcher(d) {
					matches++
				}
			}
		}
		samples = append(samples, elapsed(start))
		expected := []int{*nodes - 2*(*repos), *nodes - 5*(*repos), 0}[i%3]
		if matches != expected {
			return fmt.Errorf("query %q matched %d documents, expected %d", query, matches, expected)
		}
	}
	result["query_p95_ms"] = percentile(samples)
	result["query_samples"] = len(samples)
	result["query_last_match_count"] = matches
	result["query_scope"] = "actual app.DocumentMatcher used by CLI/MCP over cached raw documents, matcher compilation once/query included; excludes load/HTTP/serialization"
	result["rss_after_queries_bytes"] = rss()
	// One external edit then the production cached shared-loader refresh.
	// A full reload is measured separately afterward; no browser debounce is included.
	changed := filepath.Join(collection.Projects[0].Root, "factory", recordPath(0))
	data, err := os.ReadFile(changed)
	if err != nil {
		return err
	}
	data = []byte(strings.Replace(string(data), "fixture budget is 42", "fixture budget is 43", 1))
	start = time.Now()
	if err := os.WriteFile(changed, data, 0644); err != nil {
		return err
	}
	refreshed, err := app.LoadCached(collection.Projects, snap)
	if err != nil {
		return err
	}
	result["incremental_refresh_ms"] = elapsed(start)
	result["rss_after_incremental_refresh_bytes"] = rss()
	runtime.KeepAlive(snap)
	if snap.Revision == refreshed.Revision {
		return fmt.Errorf("external edit not observed")
	}
	result["refresh_scope"] = "production LoadCached after one external disk edit, compares current bytes and revalidates; no watcher debounce/HTTP timing measured"
	start = time.Now()
	full, err := app.Load(collection.Projects)
	if err != nil {
		return err
	}
	result["full_reload_comparison_ms"] = elapsed(start)
	if full.Revision != refreshed.Revision {
		return fmt.Errorf("cached refresh differs from fresh full reload")
	}
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	result["go_heap_alloc_bytes"] = mem.HeapAlloc
	result["targets"] = map[string]any{"cold_load_under_5000ms": result["cold_load_ms"].(float64) < 5000, "cached_query_p95_under_200ms": result["query_p95_ms"].(float64) < 200, "incremental_loader_refresh_under_1000ms": result["incremental_refresh_ms"].(float64) < 1000, "sampled_rss_under_512MiB": result["rss_after_load_bytes"].(int64) > 0 && result["rss_after_load_bytes"].(int64) < 512<<20 && result["rss_after_incremental_refresh_bytes"].(int64) > 0 && result["rss_after_incremental_refresh_bytes"].(int64) < 512<<20}
	result["limitations"] = []string{"OS cache not evicted", "instantaneous RSS samples, not peak; excludes browser/harness", "no HTTP/watcher or debounce timing", "non-host platforms not executed", "generated scratch Git repositories have no commits"}
	return json.NewEncoder(os.Stdout).Encode(result)
}
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
