# A scoped FindAllTags query is inclusive of unscoped tags

[0003](./0003-scope-always-persisted-empty-string-default.md) kept unscoped tags (`scope == ""`, or no `scope` attribute at all for pre-feature rows) reachable through the *empty-scope* read path — `FindAllTags(ctx, "")` skips the filter and returns everything. But a *non-empty* scope applied a strict `#scope = :scope` filter, so a scoped caller saw only tags whose scope matched exactly and nothing unscoped.

That strict behavior broke `budget-api` the moment it started asking for expense-scoped tags. `budget-api` now calls `GET /api/tags/scope/expense` to resolve the tags stored on each expense. In production every one of those tags is **unscoped**: the frontend still creates tags through the unscoped `PUT /api/tags` (no scope sent → `scope == ""`), and legacy tags predate the field entirely. The strict filter returned only the `UNKNOWN` sentinel, so `budget-api`'s per-tag `GetTagBy` failed for every real tag, every expense was dropped in `FindByDateRange`, and single-expense reads 500'd. Nothing in the system actually produces `expense`-scoped tags yet, so strict scoping had no data to match.

Starting with this iteration, a non-empty scope is an **inclusive** filter:

```
#scope = :scope OR #scope = :empty OR attribute_not_exists(#scope)
```

A scoped query returns tags matching that scope **plus** all unscoped tags (`scope == ""` or no attribute). Only tags scoped to a *different* value (e.g. `revenue` when `expense` was requested) are excluded. This is the same "shared/legacy tags stay reachable without a backfill" guarantee [0003](./0003-scope-always-persisted-empty-string-default.md) gave the empty-scope path, now extended to scoped callers — which is what makes scoping safe to adopt incrementally, before any producer writes real scopes.

A consequence: an unscoped tag is visible under *every* scope. That is intended for the migration window. If/when scopes become authoritative (the frontend writes them and existing tags are backfilled), this inclusiveness can be revisited.

## Considered Options

- Keep the strict `#scope = :scope` filter and backfill existing tags to `scope = "expense"` plus teach the frontend to write scopes — rejected for now: a multi-service, data-migration-sized change to fix a runtime breakage, and the unscoped `PUT /api/tags` path is deliberately retained, so unscoped tags keep being created regardless.
- Revert `budget-api` to the unscoped `GET /api/tags` — rejected: abandons the expense-scoping intent instead of making it work.
- Inclusive scoped filter (chosen) — one change at the tag-api read layer, preserves [0003](./0003-scope-always-persisted-empty-string-default.md)'s no-backfill guarantee, keeps `budget-api`'s scoped lookup working against today's unscoped data while still excluding foreign scopes.
