package cache

import "context"

// layeredCacheProvider composes two CacheProviders into a two-level cache.
// It depends only on the CacheProvider interface — it has no knowledge of
// the concrete L1/L2 implementations (e.g. Ristretto, Redis).
type layeredCacheProvider struct {
	l1 CacheProvider
	l2 CacheProvider
}

// NewLayeredCacheProvider composes l1 (fast, typically in-process) and l2
// (typically shared/networked) into a single CacheProvider. Reads check l1
// first and backfill l1 from l2 on an l1 miss / l2 hit. Writes go to l2 then
// l1 (write-through, L2-before-L1 order). Evicts go to l1 then l2 (mirror
// order). TTLs are owned entirely by whichever concrete l1/l2 providers are
// passed in — LayeredCacheProvider has no TTL of its own. Cross-replica L1
// staleness is not addressed by this type.
func NewLayeredCacheProvider(l1, l2 CacheProvider) CacheProvider {
	return &layeredCacheProvider{l1: l1, l2: l2}
}

func (provider *layeredCacheProvider) Get(key string) (string, bool) {
	return provider.GetContext(context.TODO(), key)
}

func (provider *layeredCacheProvider) Set(key string, value string) error {
	return provider.SetContext(context.TODO(), key, value)
}

func (provider *layeredCacheProvider) Evict(key string) error {
	return provider.EvictContext(context.TODO(), key)
}

func (provider *layeredCacheProvider) GetContext(ctx context.Context, key string) (string, bool) {
	if value, found := provider.l1.GetContext(ctx, key); found {
		return value, true
	}

	value, found := provider.l2.GetContext(ctx, key)
	if !found {
		return "", false
	}

	_ = provider.l1.SetContext(ctx, key, value)
	return value, true
}

func (provider *layeredCacheProvider) SetContext(ctx context.Context, key string, value string) error {
	_ = provider.l2.SetContext(ctx, key, value)
	_ = provider.l1.SetContext(ctx, key, value)
	return nil
}

func (provider *layeredCacheProvider) EvictContext(ctx context.Context, key string) error {
	_ = provider.l1.EvictContext(ctx, key)
	_ = provider.l2.EvictContext(ctx, key)
	return nil
}
