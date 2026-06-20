# Scope is filtered with a FilterExpression, not a GSI, until existing tags are backfilled

We're adding `Scope` to `Tag` so clients can ask for tags belonging to one domain (e.g. "expense", "revenue"). The natural long-term shape is a GSI keyed so `Scope` queries are indexed instead of scanned. We rejected building that GSI in this iteration: DynamoDB GSIs are sparse, so every `Tag` written before this field existed — the entire current dataset — would simply be invisible to any scope-filtered query until each one is resaved or backfilled, and we don't yet have a backfill plan or the appetite to write one for this change.

Instead, `FindAllTags(ctx, scope)` runs the same `Query` keyed on `user_name` regardless of `scope`, and applies a `FilterExpression` on the normalized `scope` attribute afterward, except when the requested `scope` is empty — see [0003](./0003-scope-always-persisted-empty-string-default.md) for why an empty `scope` skips filtering entirely instead of matching on an empty string, and [0004](./0004-single-find-all-tags-method-with-scope-parameter.md) for why this is one method rather than two. This costs the same read capacity as fetching the full per-user tag set regardless of how selective the filter is, but a user's tag catalog is small, so the inefficiency is negligible. Supplying `Scope` on write (`PUT /api/tags`) is optional, so existing clients are unaffected.

The GSI-backed index is deferred to a later iteration, once there's a backfill story for existing tags. Don't add the GSI without also resolving what happens to pre-existing scope-less tags.

## Considered Options

- GSI on `(user_name, scope)` with a Projection, built now — rejected: sparse-index semantics would silently drop every pre-existing tag from scoped queries with no backfill in place.
- FilterExpression on the existing user_name-keyed Query (chosen) — no index/migration work now, acceptable cost given small per-user catalogs, indexing revisited once backfill exists.
