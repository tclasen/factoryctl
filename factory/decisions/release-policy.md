---
type: Decision
title: Conventional Commits and pre-1.0 versioning
description: Commit and release rules for developing factoryctl.
status: draft
factory:
  format: 1
  assertion:
    frame: desired
    certainty: mixed
    basis: mixed
sources:
  - id: operator
    resource: ../sources/operator-brief.md
    title: Operator versioning directive
    author: human:operator
decision_state: accepted
---

# Conventional Commits and pre-1.0 versioning

## Operator-authorized policy

Use Conventional Commits and derive semantic versions from them. The initial release target is **0.1.0**, tagged **v0.1.0**. No backwards compatibility commitment applies during 0.x development. Do not add compatibility shims, deprecation periods or migrations solely to preserve earlier 0.x behavior. This does not authorize unrelated data destruction or weaken validation.

The operator must explicitly authorize both advancing to and releasing **1.0.0**. Neither breaking changes nor release automation may cross that boundary automatically. This policy does not assert that 0.1.0 has been built, tagged or published.

## Operational conventions — agent-selected implementation of the directive

Follow [Conventional Commits 1.0.0](https://www.conventionalcommits.org/en/v1.0.0/): `type(scope)!: description`, with optional scope and breaking marker. Use `feat` for features, `fix` for fixes, and appropriate `docs`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `style` or `revert` types for other work. Mark incompatible interface or file-format changes with `!` or a `BREAKING CHANGE:` footer even during 0.x; recording the change does not promise compatibility.

All new authored commits, PR titles and merge commit titles must follow this convention. Preserve existing published history. When preserving commits through a merge, choose a conventional merge title such as `chore: merge release policy (#2)`; do not let the integration commit introduce a duplicate feature or fix signal.

After the initial 0.1.0 release, calculate the next release from all commits reachable from the release candidate but not the previous release tag. Use the highest applicable increment once per release:

| Changes since the last release | Increment during 0.x |
| --- | --- |
| Any breaking marker, regardless of commit type | Minor; reset patch to zero |
| Any `feat` without a breaking marker | Minor; reset patch to zero |
| Any `fix` or `perf` without a breaking marker | Patch |
| Only other types without breaking markers | No automatic version change or release |

For example, fixes after 0.1.0 produce 0.1.1; a feature or breaking change produces 0.2.0. The first release remains 0.1.0 regardless of accumulated development commits. Before the authorized 1.0.0 transition, define the stable public contract and switch the release calculation to normal post-1.0 SemVer rules.

[Semantic Versioning 2.0.0](https://semver.org/) explicitly treats major version zero as unstable. The particular 0.x increment mapping above is a project convention, not a requirement imposed by SemVer. Released tags and artifacts must not be rewritten.

## Enforcement and current limits

Agents apply this policy immediately when committing and merging. Automated commit checks, version calculation and release guards are acceptance requirements for [release qualification](../tasks/11-release-qualification.md); they are not implemented by this document. Release tooling must test the initial-version case, aggregate increment precedence, breaking changes on any type, no-release changes and rejection of an unauthorized major version of 1 or greater. No release is published by adopting this policy.

## Mandatory milestone review gate

Every milestone delivery must include a current Playwright demo exercising the complete requirements traceability matrix end to end, as explicitly directed by the operator. Apply [the demonstration gate](../specifications/release-evidence.md#milestone-review-delivery) before declaring delivery complete. Reconcile every atomic row with current executable evidence, video scene timestamps and exact gaps. A partial review artifact may be shared, but no missing or unsupported capability may be reported as passing, and no video may be deferred merely because the milestone is not a tagged release.
