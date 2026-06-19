package rest

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/ristretto"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/cache"
	inmemory "github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/cache/in_memory"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
)

type RistrettoCachedSearchTagRepository struct {
	CacheProvider cache.CacheProvider
	Delegate      tags.SearchTagRepository
	logger        *logging.Logger
}

func NewRistrettoCachedSearchTagRepository(delegate tags.SearchTagRepository) tags.SearchTagRepository {
	logger := logging.GetLoggerInstanceForComponentByType(&RistrettoCachedSearchTagRepository{})
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		logger.LogErrorfFor("Error while creating Ristretto cache: %s", err)
		panic(err)
	}
	return &RistrettoCachedSearchTagRepository{
		CacheProvider: &inmemory.RistrettoCacheProvider{
			Cache: cache,
		},
		Delegate: delegate,
		logger:   logger,
	}
}

func (repository *RistrettoCachedSearchTagRepository) GetTagBy(ctx context.Context, key string) (*tags.SearchTag, error) {
	searchTags, err := repository.GetAllTags(ctx)
	if err != nil {
		return nil, err
	}

	for _, searchTag := range searchTags {
		if searchTag.Key == key {
			return &searchTag, nil
		}
	}
	return nil, fmt.Errorf("tag with key '%s' not found", key)
}

func (repository *RistrettoCachedSearchTagRepository) GetAllTags(ctx context.Context) ([]tags.SearchTag, error) {
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		repository.logger.LogErrorfFor("Error fetching current user: %s", err)
		return nil, err
	}

	cacheKey := fmt.Sprintf("search_tags_user_%s", *user.UserName)

	cachedSearchTags, found := repository.getFromTheCache(ctx, cacheKey)
	if found {
		return cachedSearchTags, nil
	}

	return repository.fetchFromDelegate(ctx, cacheKey, *user.UserName)
}

func (repository *RistrettoCachedSearchTagRepository) getFromTheCache(ctx context.Context, cacheKey string) ([]tags.SearchTag, bool) {
	if cachedValue, found := repository.CacheProvider.GetContext(ctx, cacheKey); found {
		repository.logger.LogDebugfFor("Cache hit for key: %s", cacheKey)
		var searchTags []tags.SearchTag
		err := json.Unmarshal([]byte(cachedValue), &searchTags)
		if err != nil {
			repository.logger.LogErrorfFor("Error un marshalling cached value for key: %s, error: %s", cacheKey, err)
			return nil, false
		}
		return searchTags, true
	}
	repository.logger.LogDebugfFor("Cache miss for key: %s", cacheKey)
	return nil, false
}

func (repository *RistrettoCachedSearchTagRepository) fetchFromDelegate(ctx context.Context, cacheKey, userName string) ([]tags.SearchTag, error) {
	repository.logger.LogDebugfFor("Cache miss for search tags, fetching from delegate for user: %s", userName)
	searchTags, err := repository.Delegate.GetAllTags(ctx)
	repository.logger.LogDebugfFor("Fetched search tags from delegate for user: %s, number of tags: %d", userName, len(searchTags))

	if err != nil {
		repository.logger.LogErrorfFor("Error fetching search tags from delegate: %s", err)
		return nil, err
	}
	repository.setToTheCache(ctx, cacheKey, searchTags)
	return searchTags, nil
}

func (repository *RistrettoCachedSearchTagRepository) setToTheCache(ctx context.Context, cacheKey string, searchTags []tags.SearchTag) {
	jsonData, err := json.Marshal(searchTags)
	if err != nil {
		repository.logger.LogErrorfFor("Error marshalling search tags for key: %s, error: %s", cacheKey, err)
		return
	}

	repository.CacheProvider.SetContext(ctx, cacheKey, string(jsonData))
	repository.logger.LogDebugfFor("Cache set for key: %s", cacheKey)
}
