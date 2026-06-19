---
status: accepted
---

# Remove the non-context Get/Set/Evict shims from CacheProvider

ADR-0001 kept `Get`/`Set`/`Evict` on `CacheProvider` as "permanent" `context.TODO()` shims, on the assumption that some consumer might still be calling them and that auditing every consumer across the monorepo wasn't worth doing at the time. We have since done that audit: `budget-api`'s `RistrettoCachedSearchTagRepository` was the only caller of `cache.CacheProvider` anywhere outside this module, and #8 migrated it to call `GetContext`/`SetContext` exclusively. `tag-api`, `account-api`, and `plan-api` never referenced `cache.CacheProvider` at all. With the actual call-site risk ADR-0001 was hedging against now at zero, we removed `Get`, `Set`, and `Evict` from the interface and from every implementation (`RistrettoCacheProvider`, `RedisCacheProvider`, `layeredCacheProvider`), leaving only the `*Context` methods.

This supersedes ADR-0001's "permanent, not pending removal" framing for the non-context methods specifically — the additive-widening principle ADR-0001 establishes for *future* signature changes still holds.

## Considered Options

- Leave the shims in place indefinitely as ADR-0001 originally decided — rejected once the audit showed they had zero real callers; carrying dead interface surface forever has no offsetting benefit.
- Deprecate via comment only, remove later — rejected as unnecessary ceremony: removal is itself the deprecation, and there's no consumer left to give a transition window to.

## Consequences

This is a breaking signature change to `cache.CacheProvider` and all three concrete providers. It is safe only because every consumer that pulls this module in via the local `replace` directive (`budget-api`, `tag-api`, `account-api`, `plan-api`) was verified to not call the removed methods before merging. Any future provider or consumer must use `GetContext`/`SetContext`/`EvictContext`; there is no non-context fallback anymore.
