package cache

import "context"

type CacheProvider interface {
	Get(key string) (string, bool)
	Set(key string, value string) error
	Evict(key string) error

	GetContext(ctx context.Context, key string) (string, bool)
	SetContext(ctx context.Context, key string, value string) error
	EvictContext(ctx context.Context, key string) error
}
