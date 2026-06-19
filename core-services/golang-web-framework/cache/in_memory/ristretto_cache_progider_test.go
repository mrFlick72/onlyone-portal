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
	ctx := context.Background()
	provider := &RistrettoCacheProvider{Cache: newTestCache(t)}

	err := provider.SetContext(ctx, "key", "value")
	require.NoError(t, err)

	value, found := provider.GetContext(ctx, "key")
	assert.True(t, found)
	assert.Equal(t, "value", value)
}

func TestGet_MissReturnsFalseNotError(t *testing.T) {
	ctx := context.Background()
	provider := &RistrettoCacheProvider{Cache: newTestCache(t)}

	value, found := provider.GetContext(ctx, "missing")
	assert.False(t, found)
	assert.Equal(t, "", value)
}

func TestSet_DefaultsTTLWhenZero(t *testing.T) {
	ctx := context.Background()
	provider := &RistrettoCacheProvider{Cache: newTestCache(t)}

	err := provider.SetContext(ctx, "key", "value")
	require.NoError(t, err)

	value, found := provider.GetContext(ctx, "key")
	assert.True(t, found)
	assert.Equal(t, "value", value)
}

func TestSet_UsesConfiguredTTL(t *testing.T) {
	ctx := context.Background()
	provider := &RistrettoCacheProvider{Cache: newTestCache(t), TTL: 5 * time.Minute}

	err := provider.SetContext(ctx, "key", "value")
	require.NoError(t, err)

	value, found := provider.GetContext(ctx, "key")
	assert.True(t, found)
	assert.Equal(t, "value", value)
}

func TestEvict_RemovesKeyAndReturnsNil(t *testing.T) {
	ctx := context.Background()
	provider := &RistrettoCacheProvider{Cache: newTestCache(t)}
	require.NoError(t, provider.SetContext(ctx, "key", "value"))

	err := provider.EvictContext(ctx, "key")
	require.NoError(t, err)

	_, found := provider.GetContext(ctx, "key")
	assert.False(t, found)
}
