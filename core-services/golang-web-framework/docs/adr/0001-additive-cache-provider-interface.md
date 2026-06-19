---
status: accepted
---

# Widen CacheProvider additively instead of a breaking change

`budget-api`, `tag-api`, `account-api`, and `plan-api` all consume `golang-web-framework` through a local path `replace` directive with no version pinning, so any signature change to an existing `CacheProvider` method breaks their build the moment this module changes — even when the change is scoped to framework-only work. To add Redis/context support we widened the interface by adding three new methods (`GetContext`, `SetContext`, `EvictContext`) alongside the original `Get`/`Set`/`Evict`, rather than changing the originals to take a `context.Context`. The non-context methods are permanent one-line shims to the `*Context` versions via `context.TODO()` — they are not deprecated cruft pending removal, they are load-bearing backward compatibility for as long as those services use the `replace` directive.

## Considered Options

- Change `Get`/`Set`/`Evict` to take `context.Context` directly and update every consumer in lockstep — rejected because it requires a coordinated multi-repo change across services with no version pinning between them.
- Add a second, parallel interface (e.g. `ContextCacheProvider`) instead of widening `CacheProvider` — rejected because it would force every concrete provider and every consumer to juggle two interface types for the same concept.
