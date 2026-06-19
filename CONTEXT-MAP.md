# Context Map

This is a polyglot monorepo of independent services; domain documentation is split per context. Contexts are added to this map lazily, as their `CONTEXT.md` is created.

## Contexts

- [golang-web-framework](./core-services/golang-web-framework/CONTEXT.md) — shared Go library providing cross-cutting concerns (HTTP bootstrap, JWT auth, caching, OTel, AWS/HTTP clients) consumed by every Gin-based service
- [Tagging](./tagging/tag-api/CONTEXT.md) — per-user catalog of tags used to categorize budget expenses
- [Budget Expense](./budget/budget-api/CONTEXT.md) — tracks a user's budget expenses and revenue, each optionally categorized by tags

## Relationships

- **golang-web-framework → consumers**: `budget-api`, `tag-api`, `account-api`, and `plan-api` each pull this module in via a local path `replace` directive (no version pinning), so any breaking change here breaks all four immediately.
- **Tagging → Budget Expense**: Budget Expense references tags by `(Key, Value)`, denormalized onto each expense at write time. The `UNKNOWN` Sentinel Tag is a convention duplicated as a literal string in both services — there is no shared code enforcing it stays in sync.
