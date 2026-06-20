# Tagging

Owns the per-user catalog of tags used to categorize budget records. Each user has their own private set of tags, and each tag is scoped to the single domain it applies to (`expense` today; `revenue` catalog supported, with revenue record tagging forthcoming).

## Language

**Tag**:
A user-owned `(Key, Value, Scope)` tuple in the catalog. `Key` is a server-generated UUID, opaque and never chosen by the client; `Value` is the human-readable label the user picks (e.g. "Groceries"); `Scope` is mandatory — every `Tag` belongs to exactly one `Scope`.
_Avoid_: Category, label

**Scope**:
A label on a `Tag` naming the single domain it applies to (`"expense"`, `"revenue"`) — one `Tag` has exactly one `Scope`, never several and never none. Not a controlled vocabulary at the type level: any caller-supplied string is accepted as long as it is non-blank, but tag-api normalizes it (trim + lower-case) before persisting and before matching, so `"Expense"` and `"expense"` are the same Scope. `Scope` is **mandatory on write**: `SaveTag` / `PUT /api/tags` reject a tag whose normalized `Scope` is empty. `Scope` is **authoritative on read**: a query for a `Scope` returns only `Tag`s of exactly that normalized `Scope` — there is no inclusive-of-unscoped read and no "return everything" read. A `Tag` therefore lives under one `Scope` and is reachable from that `Scope` alone. (Earlier iterations treated `Scope` as optional with an empty-string default, an empty-`Scope` query meaning "everything," and a scoped query including unscoped tags — a migration-window model now retired; see [ADR 0007](./docs/adr/0007-scope-mandatory-and-scoped-reads-are-strict.md), which superseded [0003](./docs/adr/0003-scope-always-persisted-empty-string-default.md) and [0006](./docs/adr/0006-scoped-query-includes-unscoped-tags.md). Existing unscoped tags are backfilled to a real `Scope` out of band.)
_Avoid_: Category, type, kind

**Sentinel Tag**:
A `Tag` that is never persisted to storage and never appears in `SaveTag`/`PUT /api/tags` — it is synthesized on every read of the catalog (`FindAllTags`, regardless of the `Scope` it's called with). `UNKNOWN` (`Key: "UNKNOWN", Value: "UNKNOWN"`) is the one sentinel tag in this system: it exists so that any expense an upstream service marks as untagged resolves to a known catalog entry, without requiring every user to be onboarded with it ahead of time. It has no `Scope` of its own — the mandatory-`Scope` rule applies to persisted `Tag`s, not to this synthesized sentinel — yet it is appended to every scope-filtered result regardless of the requested `Scope`, so every scoped caller always has a catch-all to fall back to. Consumers that present the catalog for editing rather than selection (the tag-management UI) hide it, since it is a technical tag, not a user-authored one.
_Avoid_: Default tag, placeholder tag
