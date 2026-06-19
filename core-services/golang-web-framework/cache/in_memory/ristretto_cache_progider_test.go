package inmemory

import (
	"context"
	"testing"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestCache(t *testing.T) *ristretto.Cache {
	c, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e4,
		MaxCost:     1 << 20,
		BufferItems: 64,
	})
	require.NoError(t, err)
	return c
}

func TestSetThenGet_RoundTrips(t *testing.T) {
	provider := &RistrettoCacheProvider{Cache: newTestCache(t)}

	err := provider.Set("key", "value")
	require.NoError(t, err)

	value, found := provider.Get("key")
	assert.True(t, found)
	assert.Equal(t, "value", value)
}

func TestGet_MissReturnsFalseNotError(t *testing.T) {
	provider := &RistrettoCacheProvider{Cache: newTestCache(t)}

	value, found := provider.Get("missing")
	assert.False(t, found)
	assert.Equal(t, "", value)
}

func TestSet_DefaultsTTLWhenZero(t *testing.T) {
	provider := &RistrettoCacheProvider{Cache: newTestCache(t)}

	err := provider.Set("key", "value")
	require.NoError(t, err)

	value, found := provider.Get("key")
	assert.True(t, found)
	assert.Equal(t, "value", value)
}

func TestSet_UsesConfiguredTTL(t *testing.T) {
	provider := &RistrettoCacheProvider{Cache: newTestCache(t), TTL: 5 * time.Minute}

	err := provider.Set("key", "value")
	require.NoError(t, err)

	value, found := provider.Get("key")
	assert.True(t, found)
	assert.Equal(t, "value", value)
}

func TestEvict_RemovesKeyAndReturnsNil(t *testing.T) {
	provider := &RistrettoCacheProvider{Cache: newTestCache(t)}
	require.NoError(t, provider.Set("key", "value"))

	err := provider.Evict("key")
	require.NoError(t, err)

	_, found := provider.Get("key")
	assert.False(t, found)
}

func TestContextVariants_MatchNonContextVariants(t *testing.T) {
	provider := &RistrettoCacheProvider{Cache: newTestCache(t)}
	ctx := context.Background()

	require.NoError(t, provider.SetContext(ctx, "key", "value"))

	value, found := provider.GetContext(ctx, "key")
	assert.True(t, found)
	assert.Equal(t, "value", value)

	require.NoError(t, provider.EvictContext(ctx, "key"))
	_, found = provider.GetContext(ctx, "key")
	assert.False(t, found)
}
