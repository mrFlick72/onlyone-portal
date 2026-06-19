package cache

import "context"

type CacheProvider interface {
	GetContext(ctx context.Context, key string) (string, bool)
	SetContext(ctx context.Context, key string, value string) error
	EvictContext(ctx context.Context, key string) error
}
