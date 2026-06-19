package redis

import (
	"context"
	"errors"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
)

// defaultOperationTimeout bounds every individual Redis call so a slow or
// dead backend degrades to a cache miss quickly instead of blocking the
// caller.
const defaultOperationTimeout = 200 * time.Millisecond

type RedisCacheProvider struct {
	client    *goredis.Client
	namespace string
	ttl       time.Duration
	timeout   time.Duration
	logger    *logging.Logger
}

// NewCacheProvider builds a RedisCacheProvider. namespace is a required
// key-prefix (e.g. "budget-api") prepended to every key as
// namespace + ":" + key, so multiple services can safely share one Redis
// instance/database. ttl is the L2 expiry applied on every Set.
func NewCacheProvider(client *goredis.Client, namespace string, ttl time.Duration) *RedisCacheProvider {
	return &RedisCacheProvider{
		client:    client,
		namespace: namespace,
		ttl:       ttl,
		timeout:   defaultOperationTimeout,
		logger:    logging.GetLoggerInstanceForComponentByType(&RedisCacheProvider{}),
	}
}

func (provider *RedisCacheProvider) key(key string) string {
	return provider.namespace + ":" + key
}

func (provider *RedisCacheProvider) Get(key string) (string, bool) {
	return provider.GetContext(context.TODO(), key)
}

func (provider *RedisCacheProvider) Set(key string, value string) error {
	return provider.SetContext(context.TODO(), key, value)
}

func (provider *RedisCacheProvider) Evict(key string) error {
	return provider.EvictContext(context.TODO(), key)
}

// GetContext fails open: any Redis error, including a real miss (redis.Nil)
// or a connection/timeout failure, degrades to a cache miss rather than
// surfacing an error to the caller.
func (provider *RedisCacheProvider) GetContext(ctx context.Context, key string) (string, bool) {
	opCtx, cancel := context.WithTimeout(ctx, provider.timeout)
	defer cancel()

	value, err := provider.client.Get(opCtx, provider.key(key)).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			provider.logger.LogDebugfFor("Redis cache miss for key: %s", key)
		} else {
			provider.logger.LogErrorfFor("Redis error on Get for key: %s, error: %s", key, err)
		}
		return "", false
	}
	return value, true
}

// SetContext never returns an error to the caller: a CacheProvider must
// never be the source of a caller-facing error. Internal Redis failures are
// logged and swallowed.
func (provider *RedisCacheProvider) SetContext(ctx context.Context, key string, value string) error {
	opCtx, cancel := context.WithTimeout(ctx, provider.timeout)
	defer cancel()

	if err := provider.client.Set(opCtx, provider.key(key), value, provider.ttl).Err(); err != nil {
		provider.logger.LogErrorfFor("Redis error on Set for key: %s, error: %s", key, err)
	}
	return nil
}

func (provider *RedisCacheProvider) EvictContext(ctx context.Context, key string) error {
	opCtx, cancel := context.WithTimeout(ctx, provider.timeout)
	defer cancel()

	if err := provider.client.Del(opCtx, provider.key(key)).Err(); err != nil {
		provider.logger.LogErrorfFor("Redis error on Evict for key: %s, error: %s", key, err)
	}
	return nil
}

// Close releases the underlying Redis client's connections. Not part of the
// CacheProvider interface — lifecycle is the consuming service's
// responsibility, not wired into WebServerConfigurer.
func (provider *RedisCacheProvider) Close() error {
	return provider.client.Close()
}
