//go:build test

package rest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
)

func userContextWithToken(userName, accessToken string) context.Context {
	u := userName
	token := accessToken
	return context.WithValue(context.Background(), "user", security.User{UserName: &u, AccessToken: &token})
}

func TestGetAllTagsCallsScopedTagEndpoint(t *testing.T) {
	var requestedPath string
	var authorizationHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPath = r.URL.Path
		authorizationHeader = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"key":"tagKey","value":"tagValue"}]`))
	}))
	defer server.Close()

	repository := NewRestSearchTagRepository(server.Client(), server.URL, "expense")

	result, err := repository.GetAllTags(userContextWithToken("A_USER", "A_TOKEN"))

	assert.Equal(t, nil, err)
	assert.Equal(t, "/api/tags/scope/expense", requestedPath)
	assert.Equal(t, "Bearer A_TOKEN", authorizationHeader)
	assert.Equal(t, []tags.SearchTag{{Key: "tagKey", Value: "tagValue"}}, result)
}

func TestGetTagByCallsScopedTagEndpointAndFiltersByKey(t *testing.T) {
	var requestedPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"key":"tagKey","value":"tagValue"},{"key":"other","value":"otherValue"}]`))
	}))
	defer server.Close()

	repository := NewRestSearchTagRepository(server.Client(), server.URL, "expense")

	result, err := repository.GetTagBy(userContextWithToken("A_USER", "A_TOKEN"), "tagKey")

	assert.Equal(t, nil, err)
	assert.Equal(t, "/api/tags/scope/expense", requestedPath)
	assert.Equal(t, &tags.SearchTag{Key: "tagKey", Value: "tagValue"}, result)
}

func TestGetTagByReturnsUnknownSentinelWhenKeyNotInCatalog(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"key":"other","value":"otherValue"}]`))
	}))
	defer server.Close()

	repository := NewRestSearchTagRepository(server.Client(), server.URL, "expense")

	result, err := repository.GetTagBy(userContextWithToken("A_USER", "A_TOKEN"), "deletedTagKey")

	assert.Equal(t, nil, err)
	assert.Equal(t, &tags.SearchTag{Key: "UNKNOWN", Value: "UNKNOWN"}, result)
}
