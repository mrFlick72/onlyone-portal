# Application Shell

The frontend SPA consuming every backend context. Its own domain vocabulary should match the backend's (see `CONTEXT-MAP.md`); this file exists to flag where it currently doesn't.

## Language

**SearchTag** *(naming collision — not a single concept)*:
Two unrelated types share this name. `budget/search-tags/domain/SearchTag.ts` defines the real Tag catalog entry (`{key, value}`, mirroring tag-api's `Tag`). `budget/expense/domain/BudgetExpense.ts` separately defines its own local `SearchTag` (`{tagKey, tagValue}`) for a tag reference attached to an expense — a different shape, same name, no import relationship between them. Tracked as a known follow-up from issue #20 ("Rename `SearchTag` → `Tag`, disambiguating from the expense filter-by-tag feature"), not yet done.
_Avoid_: treating the two `SearchTag`s as interchangeable
