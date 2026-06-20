# 0001 — Expense-scoped tag lookup, hardcoded at the wiring layer

- Status: Accepted
- Date: 2026-06-20

## Context

`tag-api` now supports filtering tags by `Scope` via `GET /api/tags/scope/:scope` (tag-api PR #14 / its ADR 0004).
`budget-api` is the only backend consumer of tag-api's tag listing, and within budget-api only **expense** tracking uses
tagging at all — revenue does not. Previously `RestSearchTagRepository` called the unscoped `GET /api/tags` and ignored
`Scope`, fetching every tag across every scope and discarding the ones it did not need.

We want budget-api to ask tag-api for exactly the tags it needs (`expense`-scoped).

## Decision

- `RestSearchTagRepository` calls `GET /api/tags/scope/expense` instead of `GET /api/tags`.
- The `"expense"` scope is a literal constant defined **once, at the wiring layer** (`config.NewSearchTagRepository`),
  and passed as a constructor argument into `rest.NewRestSearchTagRepository` (drives the request path) and on into
  `rest.NewRistrettoCachedSearchTagRepository` (drives the cache key).
- `domain/tags.SearchTagRepository` (`GetAllTags(ctx)`, `GetTagBy(ctx, key)`) and the `SearchTag` value object are
  **unchanged**. `Scope` stays purely an adapter/wiring concern and never reaches budget-api's domain — its `SearchTag`
  type has no `Scope` field. `GetTagBy` already delegates to `GetAllTags` and filters in memory, so it inherits the
  narrowed result set for free.
- The Ristretto cache key changes from `search_tags_user_<userName>` to `search_tags_user_<userName>_scope_expense`, so
  the cached entry's identity matches what it actually holds and a future second-scope caller sharing this code path
  cannot silently collide with it.

## Why hardcode the scope rather than make it configurable or thread it through the domain

- **Not configurable**: the scope is not an operational knob. budget-api tags expenses; that is a property of the code,
  not of a deployment. A config key would invite misconfiguration (e.g. pointing at `revenue` and silently breaking
  expense tagging) with no legitimate use for the flexibility.
- **Not in the domain**: budget-api's domain has no concept of tag scope and no `Scope` field on `SearchTag`. Threading
  scope through `SearchTagRepository` would leak a tag-api filtering concept into a domain port that does not need it,
  widening the interface for every implementation and test for no domain benefit. Keeping it at the wiring layer holds
  the concern where the binding to tag-api's contract already lives.

## Consequences

- budget-api fetches a smaller, relevant result set from tag-api.
- The cache key uniquely identifies the expense-scoped result set; mixing scopes in one process is collision-safe.
- tag-api's unscoped `GET /api/tags` is **not** removed — the frontend's `SearchTagsPage` still depends on it for an
  unfiltered tag-catalog listing, a structurally different need. Removing that endpoint requires a UX redesign of that
  page first and is out of scope here.
