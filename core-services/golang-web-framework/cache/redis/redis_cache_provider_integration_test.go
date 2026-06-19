//go:build integration

package redis

import (
	"context"
	"testing"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	redistest "github.com/testcontainers/testcontainers-go/modules/redis"
)

func newRealRedisProvider(t *testing.T, namespace string, ttl time.Duration) *RedisCacheProvider {
	ctx := context.Background()
	container, err := redistest.Run(ctx, "redis:7-alpine")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = testcontainers.TerminateContainer(container)
	})

	connStr, err := container.ConnectionString(ctx)
	require.NoError(t, err)

	opts, err := goredis.ParseURL(connStr)
	require.NoError(t, err)
	client := goredis.NewClient(opts)
	t.Cleanup(func() { _ = client.Close() })

	return NewCacheProvider(client, namespace, ttl)
}

func TestIntegration_SetThenGet_RoundTrips(t *testing.T) {
	ctx := context.Background()
	provider := newRealRedisProvider(t, "ns", time.Minute)

	require.NoError(t, provider.SetContext(ctx, "key", "value"))

	value, found := provider.GetContext(ctx, "key")
	assert.True(t, found)
	assert.Equal(t, "value", value)
}

func TestIntegration_Evict_RemovesKey(t *testing.T) {
	ctx := context.Background()
	provider := newRealRedisProvider(t, "ns", time.Minute)
	require.NoError(t, provider.SetContext(ctx, "key", "value"))

	require.NoError(t, provider.EvictContext(ctx, "key"))

	_, found := provider.GetContext(ctx, "key")
	assert.False(t, found)
}

func TestIntegration_Namespace_PrefixesKeys(t *testing.T) {
	ctx := context.Background()
	provider := newRealRedisProvider(t, "myservice", time.Minute)
	require.NoError(t, provider.SetContext(ctx, "key", "value"))

	value, found := provider.GetContext(ctx, "key")
	assert.True(t, found)
	assert.Equal(t, "value", value)
}

func TestIntegration_Get_MissReturnsFalse(t *testing.T) {
	ctx := context.Background()
	provider := newRealRedisProvider(t, "ns", time.Minute)

	_, found := provider.GetContext(ctx, "missing")
	assert.False(t, found)
}
