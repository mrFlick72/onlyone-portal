# 0001 — Prefer `type` over `interface` for object shapes

- Status: Accepted
- Date: 2026-06-22

## Context

The frontend declares object shapes — component `Props`, domain models, message bundles, select options — inconsistently. A survey of `src/` found the split roughly even: **32** files declare props with `interface XProps {…}` and **22** with `type XProps = {…}`. The non-props layer already leans `type`: domain models (`BudgetExpense`, `BudgetRevenue`), the message-bundle types, and `SelectOption` are all `type` aliases. So the inconsistency is concentrated in component props, where both forms appear with no rule deciding which.

`interface` and `type` are interchangeable for plain object shapes. `interface` offers two things `type` does not:

- **Declaration merging** — re-opening an interface by re-declaring it.
- **`extends`** — though `type` expresses the same via intersection (`A & B`).

A survey of our own code found that **neither is used**: no `…Props` interface `extends` anything, and there is no declaration merging anywhere in `src/`. The historical counter-argument for `interface` — marginally better type-checker caching and error messages for large/deep object hierarchies — does not apply to small component-props objects.

## Decision

**`type` is the default for object shapes in the frontend.**

- Declare component `Props`, domain models, and other object shapes as `type X = { … }`.
- `interface` is reserved for the two cases where it earns its keep:
  1. **Declaration merging** is genuinely needed (e.g. augmenting a global or a third-party module).
  2. **Extending a third-party `interface`** that is itself exported as an `interface` (so `extends` reads naturally against the upstream type).
- Everything else uses `type`. Unions and intersections (e.g. `MenuItemProps & { scope?: TagScope }`) are expressed directly, which `type` handles uniformly and `interface` cannot.

All existing `interface` object shapes are converted to `type` in one standalone, behavior-preserving change (nothing `extends`/merges, so the transform is mechanical), validated by `tsc --noEmit` and the production build.

## Why not keep both / standardize on `interface`

- **Both (status quo)** — rejected: the even 32/22 split is exactly the inconsistency this resolves; "either is fine" is what produced it.
- **Standardize on `interface`** — rejected: it would fight the union/intersection props the codebase already has, contradict the `type`-based domain/bundle layer, and buy nothing, since the interface-only features (merging, `extends`) are unused here.
- **`type` (chosen)** — one rule, consistent with the existing non-props layer, no loss of any capability we actually use, and no accidental declaration merging.

## Enforcement (deferred)

This ADR is **documentation-enforced only** at acceptance. There is no linter in the frontend today (only `tsc --noEmit`), so nothing automatically prevents a stray `interface …Props` from reappearing — there is a known drift window. Introducing ESLint with `@typescript-eslint/consistent-type-definitions: ['error', 'type']` (whose `--fix` would also have performed this conversion) is tracked as a separate follow-up, deliberately kept out of this change so a convention codemod does not smuggle in a from-scratch linting toolchain and a broader ruleset decision.
