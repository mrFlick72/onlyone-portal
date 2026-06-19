---
status: accepted
---

# Cross-replica L1 staleness is accepted, not solved

`LayeredCacheProvider` does not invalidate other replicas' L1 entries when one replica writes or evicts a key — it only updates its own process's L1 and the shared L2. If a value changes via one replica, every other replica keeps serving its stale L1 copy until that entry's own L1 TTL expires. We accepted this bound (staleness ≤ L1 TTL) for this iteration rather than building a real-time invalidation mechanism, because the consumers in scope tolerate brief staleness and a pub/sub invalidation channel is meaningful additional infrastructure with no concrete need yet.

## Considered Options

- Pub/sub invalidation: have each write/evict publish a message (e.g. via Redis Pub/Sub) that other replicas' `LayeredCacheProvider` instances subscribe to and use to evict their own L1 entry — deferred as a follow-up; real complexity (delivery guarantees, reconnect handling) not yet justified by a concrete consumer need.
- Skip L1 entirely and read L2 (Redis) on every call — rejected for this story; it would erase the latency benefit that motivated adding an in-process tier in the first place.

## Consequences

Anyone deploying a multi-replica consumer of `LayeredCacheProvider` must choose L1 TTL deliberately: it is the de-facto upper bound on how long a stale read can survive after another replica's write. A future story may need to add cross-replica invalidation if a consumer's staleness tolerance shrinks below what L1 TTL can practically provide.
