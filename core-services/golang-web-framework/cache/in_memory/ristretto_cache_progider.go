package inmemory

import (
	"context"
	"sync"
	"time"

	"github.com/dgraph-io/ristretto"

	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
)

const defaultTTL = 1 * time.Minute

//todo move the logger as RistrettoCacheProvider field and then init it in the 
type RistrettoCacheProvider struct {
	Cache *ristretto.Cache
	TTL   time.Duration
}

var loggerOnce = sync.OnceValue(func() *logging.Logger {
	return logging.GetLoggerInstanceForComponentByType(&RistrettoCacheProvider{})
})

func (provider *RistrettoCacheProvider) GetContext(ctx context.Context, key string) (string, bool) {
	val, found := provider.Cache.Get(key)
	if !found {
		return "", false
	}
	return val.(string), true
}

// SetContext never returns an error to the caller: a CacheProvider must never
// be the source of a caller-facing error, only the real delegate behind a
// cache-aside call should ever surface one. Internal failures are logged and
// swallowed.
func (provider *RistrettoCacheProvider) SetContext(ctx context.Context, key string, value string) error {
	ttl := provider.TTL
	if ttl <= 0 {
		ttl = defaultTTL
	}

	if !provider.Cache.SetWithTTL(key, value, 1, ttl) {
		loggerOnce().LogErrorfFor("Ristretto cache rejected SetWithTTL for key: %s", key)
		return nil
	}
	provider.Cache.Wait()
	return nil
}

func (provider *RistrettoCacheProvider) EvictContext(ctx context.Context, key string) error {
	provider.Cache.Del(key)
	provider.Cache.Wait()
	return nil
}
