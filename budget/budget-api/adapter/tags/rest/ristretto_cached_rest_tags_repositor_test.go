//go:build test

package rest

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
	"github.com/stretchr/testify/mock"
)

type CacheProviderMock struct {
	mock.Mock
}

func (m *CacheProviderMock) GetContext(ctx context.Context, key string) (string, bool) {
	args := m.Called(ctx, key)
	return args.String(0), args.Bool(1)
}

func (m *CacheProviderMock) SetContext(ctx context.Context, key string, value string) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (m *CacheProviderMock) EvictContext(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func userContextFor(userName string) context.Context {
	u := userName
	return context.WithValue(context.Background(), "user", security.User{UserName: &u})
}

func TestGetAllTagsReturnsCachedValueOnCacheHitWithoutCallingDelegate(t *testing.T) {
	ctx := userContextFor("A_USER")
	cacheKey := "search_tags_user_A_USER_scope_expense"
	cachedTags := []tags.SearchTag{{Key: "tagKey", Value: "tagValue"}}
	cachedJSON, _ := json.Marshal(cachedTags)

	cacheProviderMock := new(CacheProviderMock)
	cacheProviderMock.On("GetContext", ctx, cacheKey).Return(string(cachedJSON), true)

	delegateMock := new(tags.SearchTagRepositoryMock)

	repository := &RistrettoCachedSearchTagRepository{
		CacheProvider: cacheProviderMock,
		Delegate:      delegateMock,
		Scope:         "expense",
		logger:        logging.GetLoggerInstanceForComponentByType(&RistrettoCachedSearchTagRepository{}),
	}

	result, err := repository.GetAllTags(ctx)

	assert.Equal(t, nil, err)
	assert.Equal(t, cachedTags, result)
	cacheProviderMock.AssertExpectations(t)
	delegateMock.AssertNotCalled(t, "GetAllTags", mock.Anything)
	cacheProviderMock.AssertNotCalled(t, "Get", mock.Anything)
}

func TestGetAllTagsFetchesFromDelegateAndPopulatesCacheOnCacheMiss(t *testing.T) {
	ctx := userContextFor("A_USER")
	cacheKey := "search_tags_user_A_USER_scope_expense"
	fetchedTags := []tags.SearchTag{{Key: "tagKey", Value: "tagValue"}}
	fetchedJSON, _ := json.Marshal(fetchedTags)

	cacheProviderMock := new(CacheProviderMock)
	cacheProviderMock.On("GetContext", ctx, cacheKey).Return("", false)
	cacheProviderMock.On("SetContext", ctx, cacheKey, string(fetchedJSON)).Return(nil)

	delegateMock := new(tags.SearchTagRepositoryMock)
	delegateMock.On("GetAllTags", ctx).Return(fetchedTags, nil)

	repository := &RistrettoCachedSearchTagRepository{
		CacheProvider: cacheProviderMock,
		Delegate:      delegateMock,
		Scope:         "expense",
		logger:        logging.GetLoggerInstanceForComponentByType(&RistrettoCachedSearchTagRepository{}),
	}

	result, err := repository.GetAllTags(ctx)

	assert.Equal(t, nil, err)
	assert.Equal(t, fetchedTags, result)
	cacheProviderMock.AssertExpectations(t)
	delegateMock.AssertExpectations(t)
	cacheProviderMock.AssertNotCalled(t, "Get", mock.Anything)
	cacheProviderMock.AssertNotCalled(t, "Set", mock.Anything, mock.Anything)
}

func TestGetTagByReturnsUnknownSentinelWhenKeyNotInCachedCatalog(t *testing.T) {
	ctx := userContextFor("A_USER")
	cacheKey := "search_tags_user_A_USER_scope_expense"
	cachedTags := []tags.SearchTag{{Key: "other", Value: "otherValue"}}
	cachedJSON, _ := json.Marshal(cachedTags)

	cacheProviderMock := new(CacheProviderMock)
	cacheProviderMock.On("GetContext", ctx, cacheKey).Return(string(cachedJSON), true)

	delegateMock := new(tags.SearchTagRepositoryMock)

	repository := &RistrettoCachedSearchTagRepository{
		CacheProvider: cacheProviderMock,
		Delegate:      delegateMock,
		Scope:         "expense",
		logger:        logging.GetLoggerInstanceForComponentByType(&RistrettoCachedSearchTagRepository{}),
	}

	result, err := repository.GetTagBy(ctx, "deletedTagKey")

	assert.Equal(t, nil, err)
	assert.Equal(t, &tags.SearchTag{Key: "UNKNOWN", Value: "UNKNOWN"}, result)
	delegateMock.AssertNotCalled(t, "GetAllTags", mock.Anything)
}
