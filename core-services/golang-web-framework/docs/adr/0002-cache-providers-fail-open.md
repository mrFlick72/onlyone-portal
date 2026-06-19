---
status: accepted
---

# Cache providers fail open, in the provider itself

A `CacheProvider` must never be the source of a caller-facing error. `Get`/`GetContext` log internal failures and return a miss; `Set`/`Evict`/their `*Context` variants log the failure and always return `nil`. This is a deliberate behavior change from `RistrettoCacheProvider.Set`, which previously returned an error when `SetWithTTL` was rejected. We chose to bake this contract into each concrete provider (Ristretto, Redis) rather than relying on `LayeredCacheProvider` to add resilience on top, so every provider remains safe to use standalone — not only when wrapped by the composite. Only the real delegate behind a cache-aside call (the DB/HTTP call made on a cache miss) is allowed to return an error to the application.

## Considered Options

- Let `Set`/`Evict` keep returning real errors and have `LayeredCacheProvider` swallow them at the composite boundary — rejected because it would leave `RedisCacheProvider` and `RistrettoCacheProvider` unsafe to use directly, contradicting the goal of additive, independently-usable providers.
- Add a circuit breaker around the Redis provider on top of fail-open — deferred; not needed for the current scale and adds operational complexity not yet justified.

## Consequences

Callers that previously could (in principle) branch on a `Set` error from `RistrettoCacheProvider` will never see one again. The only real consumer at the time of this decision (`budget-api`'s `RistrettoCachedSearchTagRepository`) already only logs that error and continues, so this is not an observable behavior change for it.
