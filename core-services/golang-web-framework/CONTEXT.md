# golang-web-framework

Shared Go library providing the cross-cutting concerns (HTTP server bootstrap, JWT auth, caching, OpenTelemetry, AWS/HTTP clients) consumed by every Gin-based service in OnlyOne Portal.

## Language

### Caching

**Layered Cache**:
A `CacheProvider` composed of a fast in-process L1 tier and a shared L2 tier via `NewLayeredCacheProvider`. Reads check L1 first and backfill it from L2 on an L1 miss; writes go to L2 then L1.
_Avoid_: Near cache, multi-level cache, tiered cache, two-level cache.

**L1 / L2**:
The naming convention for the two tiers of a Layered Cache. L1 is the fast, per-replica, in-process tier (e.g. Ristretto). L2 is the shared tier visible to every replica (e.g. Redis).
_Avoid_: local/remote cache, near/far tier, fast/slow tier.

**Backfill**:
Populating L1 with a value just read from L2 after an L1 miss, so the next read for that key hits L1 directly.
_Avoid_: cache warming, promotion.

**Fail-open**:
The contract that a `CacheProvider` never surfaces an internal failure (a broken connection, a rejected write) as an error to its caller. `Get`/`GetContext` degrade to a miss; `Set`/`Evict`/their `*Context` variants log the failure and always return `nil`. Only the real delegate behind a cache-aside call (the DB/HTTP call made on a cache miss) is allowed to return an error to the application.
_Avoid_: graceful degradation, soft failure.

**Namespace**:
The mandatory key-prefix (`namespace + ":" + key`) every Redis-backed cache provider is constructed with, preventing key collisions when multiple services share one Redis instance.
_Avoid_: prefix, scope.
