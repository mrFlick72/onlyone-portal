package inmemory

import (
	"errors"
	"time"

	"github.com/dgraph-io/ristretto"
)

type RistrettoCacheProvider struct {
	Cache *ristretto.Cache
}

func (provider *RistrettoCacheProvider) Get(key string) (string, error) {
	val, found := provider.Cache.Get(key)
	if !found {
		return "", nil
	}
	return val.(string), nil
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
