package redis

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newMiniredisProvider(t *testing.T, namespace string, ttl time.Duration) (*RedisCacheProvider, *miniredis.Miniredis) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(mr.Close)

	client := goredis.NewClient(&goredis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = client.Close() })

	return NewCacheProvider(client, namespace, ttl), mr
}

func TestSetThenGet_RoundTrips(t *testing.T) {
	ctx := context.Background()
	provider, _ := newMiniredisProvider(t, "ns", time.Minute)

	require.NoError(t, provider.SetContext(ctx, "key", "value"))

	value, found := provider.GetContext(ctx, "key")
	assert.True(t, found)
	assert.Equal(t, "value", value)
}

func TestGet_MissReturnsFalseNotError(t *testing.T) {
	ctx := context.Background()
	provider, _ := newMiniredisProvider(t, "ns", time.Minute)

	value, found := provider.GetContext(ctx, "missing")
	assert.False(t, found)
	assert.Equal(t, "", value)
}

func TestNamespace_PrefixesKeysInRedis(t *testing.T) {
	ctx := context.Background()
	provider, mr := newMiniredisProvider(t, "myservice", time.Minute)

	require.NoError(t, provider.SetContext(ctx, "key", "value"))

	value, err := mr.Get("myservice:key")
	require.NoError(t, err)
	assert.Equal(t, "value", value)
}

func TestEvict_RemovesKey(t *testing.T) {
	ctx := context.Background()
	provider, _ := newMiniredisProvider(t, "ns", time.Minute)
	require.NoError(t, provider.SetContext(ctx, "key", "value"))

	require.NoError(t, provider.EvictContext(ctx, "key"))

	_, found := provider.GetContext(ctx, "key")
	assert.False(t, found)
}

func TestTTL_ExpiresViaFastForward(t *testing.T) {
	ctx := context.Background()
	provider, mr := newMiniredisProvider(t, "ns", 1*time.Second)
	require.NoError(t, provider.SetContext(ctx, "key", "value"))

	mr.FastForward(2 * time.Second)

	_, found := provider.GetContext(ctx, "key")
	assert.False(t, found)
}

func TestGet_DegradesToMissWhenConnectionBroken(t *testing.T) {
	ctx := context.Background()
	provider, mr := newMiniredisProvider(t, "ns", time.Minute)
	mr.Close()

	value, found := provider.GetContext(ctx, "key")
	assert.False(t, found)
	assert.Equal(t, "", value)
}

func TestSet_NeverReturnsErrorWhenConnectionBroken(t *testing.T) {
	ctx := context.Background()
	provider, mr := newMiniredisProvider(t, "ns", time.Minute)
	mr.Close()

	err := provider.SetContext(ctx, "key", "value")
	assert.NoError(t, err)
}

func TestEvict_NeverReturnsErrorWhenConnectionBroken(t *testing.T) {
	ctx := context.Background()
	provider, mr := newMiniredisProvider(t, "ns", time.Minute)
	mr.Close()

	err := provider.EvictContext(ctx, "key")
	assert.NoError(t, err)
}
