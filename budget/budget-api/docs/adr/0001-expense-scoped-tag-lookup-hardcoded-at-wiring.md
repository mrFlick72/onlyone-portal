# budget-api fetches tags scoped to `expense`, hardcoded at the wiring layer, not threaded through the domain

tag-api now exposes `GET /api/tags/scope/:scope` alongside its existing unscoped `GET /api/tags` (see tag-api's [0004-single-find-all-tags-method-with-scope-parameter.md](../../../tagging/tag-api/docs/adr/0004-single-find-all-tags-method-with-scope-parameter.md)). budget-api is presently the only backend consumer of tag-api, and within the Budget Expense context only expense tracking uses tagging at all — revenue tracking has no tagging UI or behavior. We switch budget-api's `RestSearchTagRepository` from the unscoped endpoint to `GET /api/tags/scope/expense`.

The literal value `"expense"` is created once, at the wiring layer (`config/configurations.go`, inside `NewSearchTagRepository`), and passed as a constructor argument into `NewRestSearchTagRepository` and on into `NewRistrettoCachedSearchTagRepository`. `domain.SearchTagRepository`'s interface (`GetAllTags(ctx)`, `GetTagBy(ctx, key)`) is unchanged — `Scope` never becomes a concept the budget-api domain knows about, matching the fact that budget-api's own `SearchTag` type has no `Scope` field. `GetTagBy` already delegates to `GetAllTags` and filters the result in memory, so it gets the narrowed result set for free, with no separate change.

Because the cached value's content (expense-scoped tags) now differs from what the old unscoped endpoint returned (all tags), the Ristretto cache key changes from `search_tags_user_<userName>` to `search_tags_user_<userName>_scope_expense`, so the cached entry's identity matches what it actually holds and a future second-scope caller sharing this code path can't silently collide with it.

This round deliberately stops at budget-api. tag-api's unscoped `GET /api/tags` is **not** removed: the frontend's tag-management page (`SearchTagsPage` / `SearchTagRepository.ts`) depends on it for an unfiltered "all of my tags" listing, which is a structurally different need from a scope-filtered consumer. Removing that endpoint would require redesigning that page first — out of scope here, left to a future issue.

## Considered Options

- Thread `scope` through `domain.SearchTagRepository` (`GetAllTags(ctx, scope)`), with the `FindSpentBudget` domain action passing `"expense"` explicitly — rejected: would make budget-api's domain aware of a concept it has no model for, for a value that never varies today.
- A new YAML config key (e.g. `budget-api.tag-scope`) — rejected: configurability only pays for itself if the value might change per deployment; nothing in this domain suggests it will.
- Collapse tag-api's two routes and redesign the frontend's tag-management page now, so the unscoped route could be deleted in this same round — rejected: turns a backend client swap into a three-repo change including a frontend UX redesign; deferred to a separate issue.
