# Context Map

This is a polyglot monorepo of independent services; domain documentation is split per context. Contexts are added to this map lazily, as their `CONTEXT.md` is created.

## Contexts

- [golang-web-framework](./core-services/golang-web-framework/CONTEXT.md) — shared Go library providing cross-cutting concerns (HTTP bootstrap, JWT auth, caching, OTel, AWS/HTTP clients) consumed by every Gin-based service
- [Tagging](./tagging/tag-api/CONTEXT.md) — per-user catalog of tags, each scoped to one domain (`expense` / `revenue`), used to categorize budget records
- [Budget Expense](./budget/budget-api/CONTEXT.md) — tracks a user's budget expenses and revenue, each optionally categorized by tags

## Relationships

- **golang-web-framework → consumers**: `budget-api`, `tag-api`, `account-api`, and `plan-api` each pull this module in via a local path `replace` directive (no version pinning), so any breaking change here breaks all four immediately.
- **Tagging → Budget Expense**: Both Budget Expense *and* Budget Revenue records carry tags — expense tags from the `expense` Scope, revenue tags from the `revenue` Scope. Each record stores tag **keys** only and resolves each key to its current value from Tagging on read, via that aggregate's scoped query. Scoped reads are strict, so an aggregate sees only its own scope's tags plus the `UNKNOWN` sentinel. The `UNKNOWN` Sentinel Tag is a convention duplicated as a literal string in both services — there is no shared code enforcing it stays in sync (budget-api keeps a single `tags.UnknownSentinel` used by both aggregates). Tags are maintained per scope through the frontend tag-management UI (Expense / Revenue tabs). Revenue tagging covers storage, resolution and display only — revenue emits no events and has no by-tag totals, so Tagging does not yet flow into any revenue analytics.
