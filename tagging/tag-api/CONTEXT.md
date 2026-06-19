# Tagging

Owns the per-user catalog of tags used to categorize budget expenses. Each user has their own private set of tags.

## Language

**Tag**:
A user-owned `(Key, Value)` pair in the catalog. `Key` is a server-generated UUID, opaque and never chosen by the client; `Value` is the human-readable label the user picks (e.g. "Groceries").
_Avoid_: Category, label

**Sentinel Tag**:
A `Tag` that is never persisted to storage and never appears in `SaveTag`/`PUT /api/tags` — it is synthesized on every read of the catalog (`FindAllTags`). `UNKNOWN` (`Key: "UNKNOWN", Value: "UNKNOWN"`) is the one sentinel tag in this system: it exists so that any expense an upstream service marks as untagged resolves to a known catalog entry, without requiring every user to be onboarded with it ahead of time.
_Avoid_: Default tag, placeholder tag
