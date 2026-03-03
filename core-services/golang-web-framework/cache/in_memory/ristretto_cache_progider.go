package inmemory

import (
	"errors"
	"time"

	"github.com/dgraph-io/ristretto"
)
	// Get(key string) (string, error)
	// Set(key string, value string) error
	// Evict(key string) error
type RistrettoCacheProvider struct {
	Cache *ristretto.Cache
}

func (provider *RistrettoCacheProvider) Get(key string) (string, bool) {
	val, found := provider.Cache.Get(key)
	if !found {
		return "", false
	}
	return val.(string), true
}

func (provider *RistrettoCacheProvider) Set(key string, value string) error {
	result := provider.Cache.SetWithTTL(key, value, 1, 1*time.Minute)
	if !result {
		return errors.New("The key value " + key + ":" + value + " was not stored in the cache")
	}
	provider.Cache.Wait()
	return nil
}

func (provider *RistrettoCacheProvider) Evict(key string) error {
	provider.Cache.Del(key)
	provider.Cache.Wait()
	return nil
}
