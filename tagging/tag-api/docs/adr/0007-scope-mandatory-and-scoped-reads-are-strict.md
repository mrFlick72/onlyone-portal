# Scope is mandatory and scoped reads are strict

Supersedes [0003](./0003-scope-always-persisted-empty-string-default.md) and [0006](./0006-scoped-query-includes-unscoped-tags.md).

[0003](./0003-scope-always-persisted-empty-string-default.md) and [0006](./0006-scoped-query-includes-unscoped-tags.md) were both written for a **migration window**: nothing produced real scopes yet, so an empty `Scope` was the universal "no real Scope" value, an empty-`Scope` query meant "return everything," and a scoped query was *inclusive* of unscoped/legacy tags so `budget-api` could ask for `expense` tags and still see the unscoped catalog it had in production. Those rules existed precisely so scoping could be adopted **before** any producer wrote scopes and **without** a data backfill.

That window is now closing deliberately. The frontend is gaining a per-scope tag-management UI (Expense / Revenue tabs) that writes a real, non-empty `Scope` on every create, and the existing data is being backfilled to `scope = "expense"` by the maintainer. Once scopes are authoritative, the migration-window leniency becomes a liability: it would show every legacy tag under *both* the Expense and Revenue tabs, blurring the per-scope separation the UI exists to provide, and it leaves an "unscoped" bucket that no longer has a legitimate producer.

Starting with this iteration, `Scope` is **authoritative**:

- **Write — mandatory, non-blank.** `PUT /api/tags` rejects a tag whose `NormalizeScope(scope)` is empty with `400`. Every persisted tag carries a real, non-empty `Scope`. (Reverses the "default empty string, written unconditionally" half of [0003](./0003-scope-always-persisted-empty-string-default.md).)
- **Read — strict.** `FindAllTags(ctx, scope)` filters on exactly `#scope = :scope`. A scoped read returns only tags of that normalized scope; unscoped/legacy tags and foreign-scope tags are excluded. (Reverses the inclusive `#scope = :scope OR #scope = :empty OR attribute_not_exists(#scope)` filter of [0006](./0006-scoped-query-includes-unscoped-tags.md).)
- **No unfiltered path.** The `if scope != ""` branch that skipped the `FilterExpression` and returned the whole catalog is removed, as is the unscoped `GET /api/tags` route that was its only caller. There is no longer any "give me everything" read. `PUT /api/tags` (create) is unaffected and remains the only write path.
- **`UNKNOWN` sentinel unchanged.** It is still synthesized and appended to every scoped read regardless of scope ([0001](./0001-unknown-sentinel-tag-not-persisted.md) stands) so every scoped consumer keeps a catch-all. It has no `Scope` of its own and is never persisted, so "mandatory scope on write" does not apply to it. The tag-management UI hides it client-side, since it is a technical tag, not a user-authored catalog entry.

**Consequence — deploy ordering is now load-bearing.** `budget-api` already calls `GET /api/tags/scope/expense` and, under [0006](./0006-scoped-query-includes-unscoped-tags.md), received its unscoped production tags through the inclusive filter. Under strict scoping it receives only `expense`-scoped tags, so the maintainer's backfill of existing tags to `scope = "expense"` **must precede** this tag-api deploy — otherwise expense tags vanish from the budget UI until the backfill runs. `budget-api` itself needs no code change.

## Considered Options

- **Keep the migration-window rules** ([0003](./0003-scope-always-persisted-empty-string-default.md) + [0006](./0006-scoped-query-includes-unscoped-tags.md)) — rejected: with a real per-scope producer arriving, "unscoped shows under every scope" defeats the Expense/Revenue separation the UI is being built for, and an unfiltered "return everything" path is a footgun once scope is meant to be authoritative.
- **Strict reads, but leave scope optional on write** (tolerate blank scope) — rejected: a blank-scope tag would be unreachable under strict reads (it matches no scope), so it could be created but never listed. Mandatory non-blank scope on write is the only consistent pairing with strict reads.
- **Scope mandatory on write + strict reads + no unfiltered path (chosen)** — makes `Scope` authoritative end to end, gives the UI clean per-scope catalogs, and removes the unfiltered escape hatch. Cost: a real data backfill, with deploy ordering as a hard prerequisite.
