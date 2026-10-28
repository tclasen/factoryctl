# factoryctl work instructions

Read `factory/factory.md` for the current product state and `factory/sources/operator-brief.md` plus the specifications relevant to the assigned task. Preserve the minimal-interface boundary: document structure and safe access belong in factoryctl; factory workflow interpretation belongs in the external harness.

## Mandatory milestone demo gate

Every milestone delivery, including an unreleased development milestone, must include a Playwright demo exercising the complete requirements traceability matrix end to end. Apply `factory/specifications/release-evidence.md` before claiming completion. Include old and new implemented behavior, current executable assertions, one readable demo video with requirement IDs and outcome explanations on the actual scenes, highlights, and a row-by-row evidence/timestamp mapping for every current atomic requirement.

Run `make review-check REVIEW_EVIDENCE=<current requirements-coverage.json>` to verify accounting, then `make review-gate REVIEW_EVIDENCE=<current requirements-coverage.json>` for full qualification. The checker validates evidence structure, source freshness and observed named-test results; it does not replace independent semantic or visual review. Do not alter statuses or weaken criteria just to pass it.

Complete row accounting is not complete satisfaction. Missing, partially exercised, blocked and unimplemented requirements remain explicit gaps; unsupported-operation errors do not pass the missing capability. A browser evaluation panel is not a product GUI. A partial review package may be shared with full-matrix qualification prominently marked incomplete, but cannot be presented as a fully qualified milestone. Only an explicit operator decision can accept a scoped exception. Do not defer the video because work is not a tagged release.

Keep video, command results and other non-Markdown artifacts outside `factory/`; keep the evidence Reference and complete requirements matrix inside it. Use immutable versioned artifact locations and content hashes. Review the exported recording, retain failed attempts and limitations, and deliver only recorder code and non-video evidence through the normal reviewed PR/main workflow. Follow `factory/decisions/release-policy.md` for commits and versions; no release publication is implied by milestone work.

The final milestone handoff must use this single stable demo link: `/Users/Shared/factoryctl/artifacts/demo/latest.mp4`. Keep an identical versioned archive and verify the canonical file hash after delivery. Do not deliver a separate requirements-only film; retain the complete row/gap matrix as supporting documentation.

Videos must remain local. Never commit, push, upload, attach to GitHub, or otherwise publish recordings without explicit operator authorization. General repository delivery authority does not authorize uploading videos. Preserve local archives and provide the stable local file link above.

The entire `artifacts/` directory and generated recordings are local-only, including JSON, images and hashes stored there. Never stage any of them, even with `git add -f`. Run `make install-hooks` in this checkout and `make repository-check` before delivery. The pre-commit hook checks the index; pre-push and CI check full proposed history. Preserve the guard and require the `repository-content` CI job when restoring repository protections. Hooks are local safeguards, not a guarantee against deliberate bypass. A clean clone does not include local review evidence; generate it locally before running the demo evidence checks.

## Intent and demo-first working method

Apply `factory/decisions/intent-and-review.md`. In each task separate the exact user request/context, assistant-selected approach and bounded execution authorization. A proposal remains a proposal until acceptance has evidence; summaries, commits and technical review cannot upgrade it or change its author. Keep one current plan in `factory/designs/implementation-plan.md`, explicitly superseding affected prior choices and dependencies.

For this operator, use a concise visible-outcome checkpoint and a working local demo/video as the primary review surface; detailed prose and accessible text remain available. Short assent authorizes the clear small slice, not unread options, future scope or publishing artifacts. Ask only material unresolved questions, not routine implementation permission. Independent review checks source-grounded outcome fit before technical correctness. Repeated corrections trigger affected-assumption reassessment before further work. User-visible provenance labels are recorded claims, never inferred approvals or a backend workflow engine.
