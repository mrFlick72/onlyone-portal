# Tagging

Owns the per-user catalog of tags used to categorize budget expenses. Each user has their own private set of tags.

## Language

**Tag**:
A user-owned `(Key, Value, Scope)` tuple in the catalog. `Key` is a server-generated UUID, opaque and never chosen by the client; `Value` is the human-readable label the user picks (e.g. "Groceries"); `Scope` is optional.
_Avoid_: Category, label

**Scope**:
A free-text label on a `Tag` naming the single domain it applies to (e.g. `"expense"`, `"revenue"`) — one `Tag` has at most one `Scope`, never several. Not a controlled vocabulary: any caller-supplied string is accepted, but tag-api normalizes it (trim + lower-case) before persisting and before matching, so `"Expense"` and `"expense"` are the same Scope. Optional — a `Tag` with no `Scope` simply doesn't match any scope-filtered query. Every `Tag` that existed before this field was introduced has no `Scope` until its owner resaves it.
_Avoid_: Category, type, kind

**Sentinel Tag**:
A `Tag` that is never persisted to storage and never appears in `SaveTag`/`PUT /api/tags` — it is synthesized on every read of the catalog (`FindAllTags`, `FindTagsByScope`). `UNKNOWN` (`Key: "UNKNOWN", Value: "UNKNOWN"`) is the one sentinel tag in this system: it exists so that any expense an upstream service marks as untagged resolves to a known catalog entry, without requiring every user to be onboarded with it ahead of time. It has no `Scope` of its own but is synthesized into every scope-filtered result regardless of the requested `Scope`, so a scope-filtered caller always has a catch-all to fall back to.
_Avoid_: Default tag, placeholder tag
