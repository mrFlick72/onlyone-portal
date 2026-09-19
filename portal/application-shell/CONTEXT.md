# Application Shell

The frontend SPA consuming every backend context. Its own domain vocabulary should match the backend's (see `CONTEXT-MAP.md`); this file exists to flag where it currently doesn't.

## Language

**Locale**:
Account's `Locale` (see `account/account-api/CONTEXT.md`) is an open BCP-47 tag as far as vauthenticator and account-api are concerned. The Application Shell narrows this at the edit surface to a closed set of the languages it actually ships message bundles for (`it`, `en` today) — the account can only ever be set to one of these through this UI, even though the underlying field accepts anything.
_Avoid_: Language (as a distinct concept from Locale — this app treats them as the same field, just constrained at the edit boundary)

**SearchTag** *(naming collision — not a single concept)*:
Two unrelated types share this name. `budget/search-tags/domain/SearchTag.ts` defines the real Tag catalog entry (`{key, value}`, mirroring tag-api's `Tag`). `budget/expense/domain/BudgetExpense.ts` separately defines its own local `SearchTag` (`{tagKey, tagValue}`) for a tag reference attached to an expense — a different shape, same name, no import relationship between them. Tracked as a known follow-up from issue #20 ("Rename `SearchTag` → `Tag`, disambiguating from the expense filter-by-tag feature"), not yet done.
_Avoid_: treating the two `SearchTag`s as interchangeable
