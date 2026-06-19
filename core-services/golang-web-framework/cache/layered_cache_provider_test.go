package cache

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stubCacheProvider is a hand-rolled CacheProvider used to assert
// orchestration rules (call order/counts) without a mocking framework.
type stubCacheProvider struct {
	store        map[string]string
	getCalls     int
	setCalls     int
	evictCalls   int
	name         string
	callSequence *[]string
}

func newStub(name string, seq *[]string) *stubCacheProvider {
	return &stubCacheProvider{store: map[string]string{}, name: name, callSequence: seq}
}

func (s *stubCacheProvider) Get(key string) (string, bool) { return s.GetContext(context.TODO(), key) }
func (s *stubCacheProvider) Set(key, value string) error {
	return s.SetContext(context.TODO(), key, value)
}
func (s *stubCacheProvider) Evict(key string) error { return s.EvictContext(context.TODO(), key) }

func (s *stubCacheProvider) GetContext(ctx context.Context, key string) (string, bool) {
	s.getCalls++
	*s.callSequence = append(*s.callSequence, s.name+":Get")
	v, ok := s.store[key]
	return v, ok
}

func (s *stubCacheProvider) SetContext(ctx context.Context, key, value string) error {
	s.setCalls++
	*s.callSequence = append(*s.callSequence, s.name+":Set")
	s.store[key] = value
	return nil
}

func (s *stubCacheProvider) EvictContext(ctx context.Context, key string) error {
	s.evictCalls++
	*s.callSequence = append(*s.callSequence, s.name+":Evict")
	delete(s.store, key)
	return nil
}

func TestGet_L1Hit_DoesNotTouchL2(t *testing.T) {
	var seq []string
	l1, l2 := newStub("L1", &seq), newStub("L2", &seq)
	l1.store["key"] = "value"
	provider := NewLayeredCacheProvider(l1, l2)

	value, found := provider.Get("key")
	assert.True(t, found)
	assert.Equal(t, "value", value)
	assert.Equal(t, 0, l2.getCalls)
}

func TestGet_L1Miss_L2Hit_BackfillsL1(t *testing.T) {
	var seq []string
	l1, l2 := newStub("L1", &seq), newStub("L2", &seq)
	l2.store["key"] = "value"
	provider := NewLayeredCacheProvider(l1, l2)

	value, found := provider.Get("key")
	assert.True(t, found)
	assert.Equal(t, "value", value)
	assert.Equal(t, "value", l1.store["key"])
	assert.Equal(t, 1, l1.setCalls)
}

func TestGet_BothMiss_ReturnsFalse(t *testing.T) {
	var seq []string
	l1, l2 := newStub("L1", &seq), newStub("L2", &seq)
	provider := NewLayeredCacheProvider(l1, l2)

	value, found := provider.Get("missing")
	assert.False(t, found)
	assert.Equal(t, "", value)
}

func TestSet_WriteOrderIsL2ThenL1(t *testing.T) {
	var seq []string
	l1, l2 := newStub("L1", &seq), newStub("L2", &seq)
	provider := NewLayeredCacheProvider(l1, l2)

	require.NoError(t, provider.Set("key", "value"))
	assert.Equal(t, []string{"L2:Set", "L1:Set"}, seq)
	assert.Equal(t, "value", l1.store["key"])
	assert.Equal(t, "value", l2.store["key"])
}

func TestEvict_OrderIsL1ThenL2(t *testing.T) {
	var seq []string
	l1, l2 := newStub("L1", &seq), newStub("L2", &seq)
	l1.store["key"] = "value"
	l2.store["key"] = "value"
	provider := NewLayeredCacheProvider(l1, l2)

	require.NoError(t, provider.Evict("key"))
	assert.Equal(t, []string{"L1:Evict", "L2:Evict"}, seq)
}

func TestSet_AlwaysReturnsNil(t *testing.T) {
	var seq []string
	provider := NewLayeredCacheProvider(newStub("L1", &seq), newStub("L2", &seq))
	assert.NoError(t, provider.Set("key", "value"))
}

func TestEvict_AlwaysReturnsNil(t *testing.T) {
	var seq []string
	provider := NewLayeredCacheProvider(newStub("L1", &seq), newStub("L2", &seq))
	assert.NoError(t, provider.Evict("key"))
}

func TestNonContextMethods_DelegateToContextMethods(t *testing.T) {
	var seq []string
	l1, l2 := newStub("L1", &seq), newStub("L2", &seq)
	provider := NewLayeredCacheProvider(l1, l2)

	require.NoError(t, provider.Set("key", "value"))
	value, found := provider.Get("key")
	assert.True(t, found)
	assert.Equal(t, "value", value)

	require.NoError(t, provider.Evict("key"))
	_, found = provider.Get("key")
	assert.False(t, found)
}
