# Scope is filtered with a FilterExpression, not a GSI, until existing tags are backfilled

We're adding `Scope` to `Tag` so clients can ask for tags belonging to one domain (e.g. "expense", "revenue"). The natural long-term shape is a GSI keyed so `Scope` queries are indexed instead of scanned. We rejected building that GSI in this iteration: DynamoDB GSIs are sparse, so every `Tag` written before this field existed — the entire current dataset — would simply be invisible to any scope-filtered query until each one is resaved or backfilled, and we don't yet have a backfill plan or the appetite to write one for this change.

Instead, `FindTagsByScope` runs the existing `Query` keyed on `user_name` (unchanged from `FindAllTags`) and applies a `FilterExpression` on the normalized `scope` attribute afterward. This costs the same read capacity as fetching the full per-user tag set regardless of how selective the filter is, but a user's tag catalog is small, so the inefficiency is negligible. `Scope` itself is optional on write (`PUT /api/tags`), so existing clients are unaffected, and a `Tag` with no `Scope` simply never matches a scope-filtered read.

The GSI-backed index is deferred to a later iteration, once there's a backfill story for existing tags. Don't add the GSI without also resolving what happens to pre-existing scope-less tags.

## Considered Options

- GSI on `(user_name, scope)` with a Projection, built now — rejected: sparse-index semantics would silently drop every pre-existing tag from scoped queries with no backfill in place.
- FilterExpression on the existing user_name-keyed Query (chosen) — no index/migration work now, acceptable cost given small per-user catalogs, indexing revisited once backfill exists.
