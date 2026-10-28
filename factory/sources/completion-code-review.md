---
type: Reference
title: Provisional completion code review
description: Independent review findings and tests before source freeze and current media qualification.
status: draft
factory:
  format: 1
  assertion:
    frame: observed
    certainty: supported
    basis: llm-inference
---

# Provisional completion code review

Reviewer `/root/completion_design_review`, run `completion-review-design` follow-up; no product authorship. Project-only inspection and review evidence. HEAD remains baseline `71c05d172e7a7a83628257acd2b2d58b3eb1f18c` with concurrent uncommitted implementation. This is not frozen-candidate, video, final-product or delivery acceptance. One independent perspective; actual model/cost/full native trace unknown.

Source-fit is unchanged from [design review](completion-design-review.md): the supplied explicit instruction authorizes completing existing remaining scope; original-message URL remains unavailable. New engineering schemas do not become operator-selected syntax. Product remains a structural document interface, browser host and harness launcher; no factory Git, backend extension execution or policy engine is authorized.

Read current completion task/brief and acceptance A01–A25, plus changed web/app/knowledge/MCP/CLI code, targeted tests and safeio. Inspected template export explicitly. Qualification/release scripts, all historical implementations, live Codex behavior, browser journeys and media were not comprehensively reviewed in this pass. All required sources supplied for this bounded follow-up were accessible.

## Delivery F01 recheck

**Design F01 closed:** completion task now records the exact schema, storage purpose, one-hop delivery_ref, manifest/criteria/assessment binding, branch-target limitation and absence of a publication endpoint. `delivery.Resolve` now rereads original and owner bytes before returning a matching local result; incomplete mappings return unverified. Local matching label explicitly separates recorded acceptance from content facts. Optional ignored observations remain evidence rather than mandatory configuration/database. There is no policy/Git/network implementation in the inspected resolver.

This resolves the design finding and the two previously requested narrow implementation repairs. It does not qualify all delivery scenarios: branch/remote verification remains explicitly unavailable and A21 remains incomplete until its required disposition/evidence is reconciled. Current tests cover local manifest freshness, shared plan/task result, deletion marker, store exclusion, no done-field inference and a static symlink. Concurrent reference changes, unavailable storage, branch observations and incomplete mappings need executed negative coverage at qualification.

## Raw product findings

Findings were sent to Coordinator while workers continued. They are observations against the inspected live files below, not claims about a later repaired candidate. Owners: web worker for C01–C03/C06–C09; extension/app worker for C04/C05 and snapshot part of C07; Coordinator integrates and obtains recheck. All remain open pending repair evidence.

| Label | Severity | Location | Concrete finding / acceptance effect |
| --- | --- | --- | --- |
| C01 | Major | `assets/app.js` renderList; `assets/editor.js` board | List renders `rows.slice(0,200)` without page/load-more control. Board truncates each type to first 200 before grouping. Record 201 cannot be reached by normal list/board traversal; filtering is not pagination. Add bounded navigable pages/virtualization and prove final-page records can be viewed/moved. A09/A13. |
| C02 | Major | `assets/editor.js` board | Cards render only when field equals a declared enum value. External missing/invalid state cards disappear rather than occupying the specified diagnostic column. Preserve visibility and allow repair without treating an invalid state as a legal transition. A03/A04/A09. |
| C03 | Major | `custom.go` customBundle near definition lookup | Reads `views.detail_template` without rejecting `!Enabled`, `Invalid` or `ExtensionInvalid`. A disabled type still offers/trusts/runs its specialized custom presentation; generic access should remain, specialized views should be disabled. Add negative tests with existing assets and disabled/invalid definition. A04/A10/A14. |
| C04 | Major | `app/create.go` CreateChecked; `app/store.go` Mutation | Template rendering occurs before writer lock. Mutation pins the definition hash, which includes the reference string but not referenced template/include bytes. A template-only edit between rendering and locked reload can persist old output that remains structurally valid; the locked dependency recheck compares the new snapshot to itself. Pin/render dependencies inside the locked operation or compare exact dependency revisions under lock. Test an intervening template-only/include-only edit and intact original/nonexistent destination. A05/A06/A14. |
| C05 | Major remaining work | `mcp/server.go` factory_nodes_list; tools schema | List supports limit but no cursor, filters or continuation. It truncates to 200 and sorts path rather than the shared contract ordering. Required MCP/CLI pagination parity is incomplete, independently of successful bounded list tests. Add snapshot-bound continuation and changed-snapshot refusal. A07/A13. |
| C06 | Major containment risk | `custom.go` readAsset | Checks ancestors with Lstat, then calls unrooted `safeio.ReadFile`. safeio no-follows the final file only; an ancestor swapped to an outside symlink after checking can redirect the read outside the project. Use project-root-scoped OS opens as in delivery/store in addition to checks, and retain static symlink refusal. This is a code-level race finding, not an executed exploit claim. A10/A14. |
| C07 | Major | `documents.go` refresh; `app/snapshot.go` snapshot fingerprint | Refresh returns early when document-byte snapshot hash is unchanged. Adding a non-Markdown file or rejected symlink under factory changes structural diagnostics without adding fingerprinted documents, leaving browser health stale. Include diagnostic/discovery changes in refresh identity or refresh diagnostics independently. A03 external revalidation. |
| C08 | Minor API contract | `documents.go` DELETE response | Mutate(delete) returns the removed old document, so `d != nil`; handler then looks up its now-removed details key and returns a zero Document instead of `{removed: ref}`. Return the documented removed identity explicitly and assert response, not just HTTP 200/disk absence. A07. |
| C09 | Major collection gap | `web.go` renderHandler.resolve | Rejects all absolute URLs before consulting snapshot-resolved metadata links. A canonical cross-project bundle URL mapping to an already-loaded project still renders unavailable; backlink derivation alone does not repair body/source navigation. Reuse `snap.Links` resolution and preserve ambiguous/missing targets explicitly. A01/A16. |

### Template export and other bounded observations

Export is additive: records/prompts create through Mutate; view export uses root-scoped exclusive create, leaves existing custom files intact, and checks asset ancestors/final regularity. No overwrite or automatic runtime launch found. Existing-record skip uses Lstat without a revision token because it performs no update; no stale overwrite follows. Root-scoped opens bound outside-root races, though check/open separation does not promise immunity from every inside-root external rename. Tests prove basic preservation; ancestor-replacement negative coverage is not established here.

MCP selected-project checks reject explicit mismatched ref identities, extension names cannot choose new roots, mappings cannot override target selectors/type, and write envelopes require revision/definition tokens. Namespaced list/call reload is request-time; reconnect guidance is honest and notifications are not claimed. Generic structural checks remain active for disabled types, which is consistent with the knowledge-format contract; disabled specialized interactions must remain distinct. Tool annotations remain advisory. A complete MCP schema/bounds fuzz or protocol compliance review was not run.

## Independent checks

Executed with `GOCACHE=/private/tmp/factoryctl-go-cache GOTMPDIR=/private/tmp` on current host:

| Command | Actual result |
| --- | --- |
| `go test ./internal/app ./internal/knowledge ./internal/mcp ./internal/delivery ./cmd/factoryctl` | Exit 0; all five packages passed (some cached). |
| `go test ./internal/web` | Exit 0; passed in 1.264s. Escalation approved for local loopback HTTP tests. |
| `go test -race ./internal/app ./internal/mcp ./internal/delivery` | Exit 0; all three packages passed, cached results reported. |

No failed test attempt or automatic approval rejection occurred. An exploratory source read named nonexistent `internal/safeio/safeio.go`; corrected to discovered `read.go`, with no effects. No new tests were written and no product/Git files changed by Reviewer. Tests ran against a live mutable tree; they are provisional evidence, not exact frozen-candidate qualification. Ordinary 120-second/test 10-minute bounds and two diagnosed retries were not exhausted.

## Acceptance accounting for this review

| Criteria | Provisional result and limitation |
| --- | --- |
| A01 | Incomplete: C09, no complete discovery/live topology regression in this pass. |
| A02 | Unverified broadly; source boundary preserved in inspected changes. |
| A03 | Incomplete: C02/C07; focused structural rejection tests passed. |
| A04 | Incomplete: C02/C03; core extension fixtures passed. |
| A05 | Incomplete: C04; rendering/export tests passed within their narrower coverage. |
| A06 | Incomplete: C04; app mutation/concurrency tests passed. |
| A07 | Incomplete: C05/C08; focused MCP/web/CLI tests passed. |
| A08 | Unverified live native runtime in this pass. |
| A09 | Incomplete: C01/C02; no independent browser journey yet. |
| A10 | Incomplete: C03/C06; static containment/digest web tests passed. |
| A11 | Unverified synthetic breadth/browser/harness journey. |
| A12 | Unverified platform/server/linkage qualification here; source contains no new Git/backend extension execution. |
| A13 | Incomplete: C01/C05; scale benchmark/release evidence excluded here. |
| A14 | Incomplete: C03/C04/C06; focused reload/prose-boundary tests passed, actual browser/harness journey pending. |
| A15 | Core fixture passed; full interface/harness/browser proof unverified. |
| A16 | Incomplete: C09; S1–S3 full regression unverified. |
| A17 | Live actual user/harness dialogue unverified here; preserved trust boundary inspected. |
| A18 | Current external Apps remains unqualified; no richer capability or CLI retirement inferred. |
| A19 | Unverified; no current media reviewed. |
| A20 | Unverified full current matrix; prior row counts are not acceptance. |
| A21 | Design ready, focused local tests passed; all required qualification scenarios/dispositions incomplete. |
| A22 | Regression tests in web passed; live selection/default journeys unverified here. |
| A23 | Existing request sequence guard/title assignment inspected; current browser/video proof unverified. |
| A24 | Existing keyboard composition/repeat guards inspected; current browser/video proof unverified. |
| A25 | Source-fit reviewed with original-link limitation; no application policy engine requested or added. Historical missing grants stay missing. |

## Inspected-source identity checkpoint

SHA-256 captured at provisional review close. Concurrent source changes mean this is a per-file checkpoint, not an atomic tree/final candidate or proof that every test used these exact bytes. Repairs invalidate affected observations and require recheck.

| Source | SHA-256 |
| --- | --- |
| Completion task/brief | `82d9df6aea53f1c21cd566f25dbf9aa7ea599e7cb45dc54a66db2ee1c889e769` |
| delivery/delivery.go | `b8058cfa053f8ffa40721bc11faeaeb79e6147098486605b98dede7b68934ff1` |
| web/delivery.go | `d6501820905201e23f392d267ddfdfc257574a1221b89af701457c5f9f7b1649` |
| web/web.go | `64bef8e26ee3c5049ca7328fe6e5de77f1e3f46c79eaad36cfc49c476291ee82` |
| web/documents.go | `766872be80de41628e7807ec6196e044ed0f2c41a05b70c5921e86aa469d285a` |
| web/custom.go | `a54a23ab15bca919722adf915f372c49e78483c0d6819b51cd97eb35ede61893` |
| web/assets/app.js | `3fb426ded293d7abcfa4c10c7085fdf87efa7b5115e27ad9dcd2893f71ca0da7` |
| web/assets/editor.js | `788218352f511ed65fc3ecf2f76cf9b87bf84267b26b8a0628105c4af47ea8e6` |
| app/create.go | `80becee76288aa89bad52fb23012b7a1a6ffa47d428d1fe4fccda211b43c1f42` |
| app/store.go | `82ec94eb8e26721dd4fc614e3548e200f02be3551efa134d64f896c70981c75a` |
| app/snapshot.go | `bbaa74e91d00fac030eb69cb3974a2d4314f356dbda60906810f70d4437929e2` |
| app/extensions.go | `7d3bf74a9a8334126df51fecf0b48f972dc28565e5a1f1ac9a9c1e016d3af100` |
| app/actions.go | `a0e2d7401f2ebecb270aa581d96b60cdd851d6fe8e6883c94019f28dea422fc3` |
| app/export.go | `a0a84d56c3e1990a851f8ae5f28ba3152255810192e001bc7d737138a702f886` |
| knowledge/extensions.go | `c9e2577ba6e10bcf11416f8ba3feb423d8291811c02fd35cea9b2662a6213dbf` |
| knowledge/schema.go | `77db6ae3413e0b38d6658caaf9b6fcc7fd44e674feb8d830b59aeeac17e8a09f` |
| mcp/server.go | `7367d604bb687f5feeee1c69c403fa22a21e698f628966c52481cceca3bb0a03` |
| mcp/tools.go | `8ac8cafe4bcdaca024bb14659538c4f7214d6ea1ada2b72eebe5a9507c5fc5be` |
| cmd/factoryctl/main.go | `b14e75053ca51ade7a86896d70b842eeefda13ee7153c0de39d2918cfbe5badf` |

Next safe action: assigned implementation owners repair/recheck C01–C09, Coordinator reconciles remaining acceptance/media requirements, freezes the candidate, and requests independent final review. This review artifact and product changes remain pending normal Coordinator-owned PR/main delivery. No acceptance, deployment, release or media publication follows from this provisional review.

## Narrow repair and cache follow-up

Coordinator subsequently changed CreateChecked to capture complete snapshot identity plus static template asset dependencies (including inactive include branches), compare asset maps before/after rendering, and pass ExpectedSnapshot/ExpectedAssets to Mutate. Mutate compares these under the writer lock and rechecks assets before publication. This resolves C04 by source inspection; the new app tests exercise changed Markdown source graph and changed asset pins with absent destination after rejection. Conservative whole-graph pinning may reject unrelated concurrent changes, but does not allow stale rendering. Existing noncooperative-editor race limits still apply.

MCP now accepts cursor, binds it to snapshot/checkout/type/query and validates offset, returning next_cursor. C05's inaccessible continuation is resolved by inspection; shared ordering and scalar filter parity remain incomplete against the written target. Those differences must be reconciled rather than silently called full A07 parity.

Reviewed new LoadCached and Document.BodyInfo changes: unchanged bytes, not timestamps, control shallow immutable parsed-document reuse; link checks use cached MarkdownInfo outputs, and existing consumers inspected only read those outputs. Schema is mutable validation metadata and initially survived definition removal through the shallow copy (new C10, major cache correctness finding). Returned C10 to Coordinator; a subsequent source check observed `d.Schema = "unavailable"` before each validation pass, resolving the specific stale-schema mechanism. Definition-removal cold-versus-cached equivalence test remains a useful required regression. Prior snapshot bytes/identity stay separate; broad cold/cache corpus equivalence and scale claims remain unverified here.

New `go test -race ./internal/app ./internal/knowledge ./internal/templates ./internal/mcp ./internal/delivery` ran: app passed 5.669s, knowledge 1.179s, templates 1.524s and delivery 2.016s; MCP failed because sandbox denied the bridge fixture's tcp6 loopback listen. This was an environment failure, not a product assertion failure. One diagnosed retry with approved listener access, `go test -race ./internal/mcp`, passed in 2.313s. Preserve this failure alongside success; earlier no-failure statement applies only to the initial checks above.

Repair checkpoint hashes (before the subsequent observed Schema reset, not an atomic frozen tree): create `8b6bf49e523635704b695db1f988fe126c99d107193618fe936545903f9a5a4f`; store `30226bc4e63c067d1dd1987038d9754783dbfdb312927bd0340c8ae89e594af6`; snapshot `814d5c760de93770264fa7b755edd4b6a6dc257240520069cf4e4b81ac7be525`; templates `52eb232f21905d7622de2242bd4c8b025a8a163dcf0ae40a88275bcb08d1dd95`; knowledge/document `df914ed376f572d5a5582c189a264dd01da18d4b7ed297ebaf28c9969c6b4ac0`; app/links `13db27ea8002a56416143283d619119ef78da889dff502808c429c28f45f8e33`; MCP server `86ec9b7e7f94951c8ca75720f74ed5fcef4915e7651ec221a64840228644ca9c`; MCP tools `d2c39dd44ae18201da474582f3ea03d635001607dba9dfa2c903ac90e9c56df6`.

Workers continue repairing GUI findings. Newly observed Show more/diagnostic-column and diagnostic snapshot identity changes are not promoted to final acceptance by this partial inspection. Full source freeze, meaningful regression, complete criteria/evidence and independent actual media review remain required.

## Source-freeze recheck before capture

Coordinator declared web source frozen and requested independent recheck of all raw findings, current CLI/app/knowledge/templates/MCP fixes and tests. Inspected those sources again; no media/evidence qualification is implied. Parent may reopen the source freeze for performance or repairs; any changed sources invalidate affected recheck hashes.

| Raw finding | Current disposition | Independent evidence / remaining action |
| --- | --- | --- |
| C01 | Partially fixed, **open major** | List/board Show more controls now expose records after 200. Both increment an unbounded limit and render prefix `[0,limit)`, so repeated use eventually renders all 10,000 DOM records. The GUI explicitly requires bounded rendering/virtualization; use fixed-size pages or virtualized windows. |
| C02 | Partially fixed, **open major** | Diagnostic column now appears and offers valid destination selection. Its identity is the literal label `Unassigned / invalid`, which is also a legal enum value. That legitimate value takes the diagnostic branch and valid cards disappear. Use a distinct internal flag/sentinel and preserve arbitrary valid enum strings. |
| C03 | Resolved, bounded code/test verification | `customBundle` refuses disabled/invalid/ExtensionInvalid definitions. `TestCustomDisabledDefinitionUsesGenericFallback` passed in web race suite. Generic record access retained. |
| C04 | Resolved, bounded code/test verification | Snapshot/assets captured around render and rechecked under lock/pre-publication; `TestTemplateSourceAndAssetPinsInsideWriter` passed in app race suite. Template dependency walker covers static inactive includes. No full concurrent external-editor guarantee inferred. |
| C05 | Lost-page and ordering mechanisms resolved; scalar field filters remain unimplemented target contract | MCP next cursor binds snapshot, checkout, type and query; offset parsing is strict; sorting now type/title/path within selected project. `TestMCPPaginationBindsFiltersAndSnapshot` passed. Core truncation issue closed. Exact scalar-filter API target remains a coverage gap for full parity, not silently implemented. |
| C06 | Resolved by containment inspection; static negative tests passed | `readAsset` now uses `os.OpenRoot` and `safeio.ReadRoot` after ancestor checks, bounding outside-root replacement races. Existing static symlink/external-asset tests passed. No executed adversarial race exploit/denial proof claimed. |
| C07 | Backend fixed, **browser part open major** | `snapshotIdentity` includes diagnostics and server test for diagnostic-only invalidation passed. Browser refresh reroutes viewed detail only when its own revision changes and never updates snapshot-note. A linked target/definition change can alter current record diagnostics/effective definition without changing record bytes, leaving displayed health and edit identity stale. Refresh derived detail on relevant snapshot changes while preserving active drafts, and give definition-stale reconciliation. |
| C08 | Resolved by source inspection | DELETE branch now explicitly returns removed ref regardless of old-document return. Existing CRUD race test passed; it still does not assert exact returned removed identity. |
| C09 | Resolved, bounded code/test verification | Browser resolution now consults valid snapshot resolvedLinks before local fallback. `TestCanonicalCrossProjectLinksAndBacklinks` passed; it asserts actual rewritten body target plus backlink identity. Full S1–S3 browser journey remains separate. |
| C10 | Resolved, bounded code/test verification | Validation resets Schema before resolving current definitions. `TestCachedDefinitionRemovalMatchesColdLoad` passed and asserts cached/cold equivalence plus prior snapshot immutability. Same-length/timestamp-restored edit and delete tests passed. |

New C11 (**major HTTP pagination contract**): `/api/v1/nodes` cursor binds only snapshot; reusing a page cursor with another checkout/type/query silently applies an offset to a different result set. Numeric suffix Atoi error is discarded, so malformed offsets become zero. Bind scope/filter identity and reject malformed offsets, matching MCP behavior. Reported directly to Coordinator before capture.

No new execution or policy scope expansion found in app.QueryMatcher: regexp.QuoteMeta keeps input literal, case-insensitive flag supplies Unicode matching, compiled once per query and applied to existing bytes. New literal/bracket/Greek-case test passed. Cached parser/link metadata consumers inspected remain read-only; immutable reuse is bounded by content hash and full structural revalidation.

Independent command against declared source freeze:

`GOCACHE=/private/tmp/factoryctl-go-cache GOTMPDIR=/private/tmp go test -race ./internal/app ./internal/knowledge ./internal/templates ./internal/mcp ./internal/delivery ./internal/web ./cmd/factoryctl`

Approved listener-access execution; exit 0. App 5.050s, MCP 2.630s, web 3.648s, CLI 5.383s; knowledge/templates/delivery returned cached passes. No failed/retried operation in this recheck. This is meaningful unit/integration regression, not a substitute for browser/media/performance/remote-main qualification.

### Recheck identities

| Current inspected source | SHA-256 |
| --- | --- |
| cmd/factoryctl/main.go | `86960a958974b9bcdb09152574523120ccb3becf634b8b8bbfd2a4f562e1abf4` |
| app/create.go | `8b6bf49e523635704b695db1f988fe126c99d107193618fe936545903f9a5a4f` |
| app/store.go | `30226bc4e63c067d1dd1987038d9754783dbfdb312927bd0340c8ae89e594af6` |
| app/snapshot.go | `228dbec41f33404f521019a0e0b2c1a55a4bf2dcfe1c14df77705fccbd964af5` |
| app/refresh_test.go | `d864b86b5fd009b68c83fe8815c759046d5c4857ba12c9c86f77fcc3b5b7f678` |
| app/search.go | `a348f42d9801be1ad6c850c5935d22a020e98267c4c9943a3f6beffad16324bf` |
| knowledge/document.go | `df914ed376f572d5a5582c189a264dd01da18d4b7ed297ebaf28c9969c6b4ac0` |
| templates/templates.go | `52eb232f21905d7622de2242bd4c8b025a8a163dcf0ae40a88275bcb08d1dd95` |
| mcp/server.go | `f0b6972869d90ac2b0f027c282f720272bcaa850e3370b8f7984a89170bf6935` |
| mcp/tools.go | `d2c39dd44ae18201da474582f3ea03d635001607dba9dfa2c903ac90e9c56df6` |
| mcp/pagination_test.go | `75859ceeecde93ed4d3a49de680b5e5ab8f55bc724e5096d965b117e870a116f` |
| web/documents.go | `d2b14a10fb2b054873badccfb2ccba30a264ab5e6c6ad4503ab0beb5e9520d3b` |
| web/custom.go | `5c45b0b3c5480ac3ff375ad5bb89d39f08232d098ac13a1f452032c1e5b9f392` |
| web/web.go | `3bd1cf6ca36457e0430c7cea23e24a306397d43d9b1621b1257e15e6cb6498e6` |
| web/documents_test.go | `31344b2f5da2ea878a9d4d588c70afbd4ba3c6c65ad1e9f829ed81bb17f1b4dc` |
| web/assets/app.js | `3fb426ded293d7abcfa4c10c7085fdf87efa7b5115e27ad9dcd2893f71ca0da7` |
| web/assets/editor.js | `788218352f511ed65fc3ecf2f76cf9b87bf84267b26b8a0628105c4af47ea8e6` |
| delivery/delivery.go | `b8058cfa053f8ffa40721bc11faeaeb79e6147098486605b98dede7b68934ff1` |

Technical acceptance remains **not ready** for the declared full scope while C01/C02/C07/C11 and remaining parity/qualification gaps are open. Ready to review narrowly repaired sources, then actual final evidence/media after recapture against the accepted frozen source. Coordinator retains all Git and delivery ownership; Reviewer changes only this report.

## Candidate eba8723 repair recheck and performance criteria reading

Observed candidate commit `eba8723417217326eee147cf7c69ec8e612ccacf`; product working tree clean at inspection, with only preserved untracked `.DS_Store` files. Subsequent product repairs invalidate the affected candidate/recheck. Media has not been reviewed.

Latest C01 is resolved by source inspection: list uses a 200-record offset window with Previous/Next; board combines scoped groups and selects a 200-card global page, clears old card snapshots, and exposes Previous/Next. C02 is resolved by a state object with independent `invalid` boolean, preserving a legitimate enum value named “Unassigned / invalid.” C11 is resolved by snapshot plus checkout/project/type/query hash and strict numeric offset parsing; `TestHTTPNodeCursorBindsSelectionAndRejectsMalformedOffset` passed in the independent race suite.

C07's derived-detail refresh is repaired: changed snapshot refresh updates snapshot-note and reroutes current detail even when its document revision is unchanged; the same-record draft/base remain in memory. **One residual recovery issue remains open:** save-error handling enters reconciliation only for `stale_revision`, not `stale_definition`. An active draft retains its old base definition token after a definition-only external edit, correctly receives a conflict, but gets no current-token/reconcile controls. Returned this precise case to Coordinator before capture; extend the same recovery branch to definition conflict and recheck. C03/C04/C06/C08/C09/C10 previously resolved mechanisms remain present. C05 continuation/order repairs remain present; scalar-filter contract coverage remains unimplemented/unverified rather than silently qualified.

Inspected new bounded-memory matcher and root-scoped refresh: ASCII fact is recorded from parsed bytes; ASCII Horspool matching folds both pattern/text, Unicode cases use the quoted literal regexp matcher. Randomized 300-case reference-equivalence and Kelvin/Greek/accent folding tests passed. Cached unchanged-document reads stream SHA-256 through a shared 32 KiB buffer, bounded to maximum+1, through root-scoped regular opens; bytes that differ are reread and reparsed. No timestamp trust, permanent cache, backend extension execution or policy interpretation was introduced by these optimizations. A snapshot still cannot promise a transaction against arbitrary external writers.

Independent approved-listener race suite against this repair state exited 0: app 5.273s; knowledge 1.200s; templates 1.231s; MCP 2.872s; delivery 1.370s; web 3.732s; CLI 5.702s. Command is the same seven-package `go test -race` recorded above. No failure/retry this round. Reviewer did not run a browser concurrently with Coordinator's exclusive capture.

Latest inspected hashes: app/search `990bc7b81b96e35e9edf096a6081067bfc49e6ab93690cfabbcb9573c3359230`; app/snapshot `4fdb3b338e20287911db77a80435ca46f8b5fa8ec35763eacd8b78c06e37c73a`; knowledge/document `a683e26fbd85bd8cebcd07714403b3405b0f49cfcde4e12748a840a98ccc5a24`; web/app.js `ee0dc6899a837bd13664585c95e62be28dc3b92c65a9c104df22d18871874dfb`; web/editor.js `769a0bf446a44e08a80491550fba6cf2e91a09e850ab936005676dd1ad4fe568`; web/documents `dfa9a17eee67dee5f65deb60a2438c271ac0100d30d06f7ea0a83fea3356b64d`.

### Mandatory workload versus proposed timing targets

Independent controlling-source reading:

- Operator brief confirmed requirement 12 names an initial design target of 100 repositories/10,000 total nodes. Its authority/proposal boundary explicitly classifies performance thresholds as LLM design proposals unless operator-confirmed.
- R21 requires designing for that workload and separating current-Mac execution from other-platform build claims. U057 requires “Design and benchmark for 100 repositories and 10,000 nodes.” Neither states an operator-selected numerical latency or memory floor.
- A13 requires the actual 100-repository/10,000-node benchmark with a declared 100 MiB corpus, cold/p95 query/incremental refresh/memory measurements against proposed architecture thresholds, and hardware/version/result reporting. It also requires current dogfood and requirement/task coverage.
- Architecture describes <5s cold, <200ms cached queries, <1s edit visibility after debounce, <512MiB RSS as “proposed thresholds to measure,” and explicitly says they are not achieved results. Acceptance introduction and matrix closing note preserve that attribution.

Thus a measured timing target miss must remain a target miss, but does not alone demonstrate failure of an operator-mandated latency bound or require inventing a scoped exception. Qualification of U057 may establish real workload design/measurement coverage while explicitly reporting engineering targets unmet. Do not mark target attainment passed, equate loader timing with end-to-end GUI visibility, reinterpret a complete benchmark as all of A13, or remove the misses from visible handoff. This is source interpretation, not a waiver, requirement weakening or operator acceptance of degraded behavior.

Inspected optional local `artifacts/completion-qualification/0.1.0/attempt-007/report.json` and scale runner. It declares 100 repositories, 10,000 loaded nodes, exactly 104857600 Markdown bytes, 19,000 valid links and 9,500 unchecked synthetic external citations, 100 query samples, current Mac/Go/hardware and OS-cache limitation. Recorded cold 13729.426ms, query p95 18.088ms, incremental shared-loader refresh 2530.037ms, sampled RSS after load 476381184 bytes and after refresh 506200064 bytes. Cold and incremental-loader targets are false; query and sampled-RSS targets true. Peak RSS, HTTP/debounce/browser edit visibility and foreign-platform execution are explicitly not measured. Runner query-scope text still names app.QueryMatcher although it now directly invokes app.DocumentMatcher; update that description for precision. Before/after source maps are supplied; benchmark predates the final app.js action-draft fix and does not independently qualify that browser change. Quiet rerun and current complete evidence remain Coordinator-owned.

At Coordinator request, Reviewer will hold further `factory/` writes after this checkpoint until capture is declared complete; any findings during that period are returned through messages. Review/delivery remains pending and no media or full-product acceptance is issued.
