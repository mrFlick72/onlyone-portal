package security

import (
	"context"
	"fmt"
	"time"

	"github.com/lestrrat-go/httprc/v3"
	"github.com/lestrrat-go/jwx/v3/jwk"
)

const (
	jwkRefreshInterval  = 15 * time.Minute
	jwkBootFetchTimeout = 10 * time.Second
)

// NewCachedJwkSet returns a jwk.Set whose contents are kept fresh by a
// background refresh goroutine bound to ctx. Cancel ctx on shutdown to stop
// the goroutine.
//
// Register waits for the initial fetch to complete before returning; the wait
// is bounded by jwkBootFetchTimeout so an unreachable JWKS endpoint fails
// boot fast rather than hanging indefinitely.
func NewCachedJwkSet(ctx context.Context, jwksURL string) (jwk.Set, error) {
	cache, err := jwk.NewCache(ctx, httprc.NewClient())
	if err != nil {
		jwk_logger.LogErrorfFor("jwk: new cache: %v", err)
		return nil, fmt.Errorf("jwk: new cache: %w", err)
	}

	bootCtx, cancel := context.WithTimeout(ctx, jwkBootFetchTimeout)
	defer cancel()
	if err := cache.Register(bootCtx, jwksURL, jwk.WithMinInterval(jwkRefreshInterval)); err != nil {
		jwk_logger.LogErrorfFor("jwk: register %q: %v", jwksURL, err)
		return nil, fmt.Errorf("jwk: register %q: %w", jwksURL, err)
	}

	cachedSet, err := cache.CachedSet(jwksURL)
	if err != nil {
		jwk_logger.LogErrorfFor("jwk: cached set %q: %v", jwksURL, err)
		return nil, fmt.Errorf("jwk: cached set %q: %w", jwksURL, err)
	}
	return cachedSet, nil
}
