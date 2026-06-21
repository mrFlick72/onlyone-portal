package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mrFlick72/onlyone-portal/tagging/tag-api/domain"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	return r
}

func TestFindAllTagsByScopeGiveOnlyTheSentinelWhenScopeEmpty(t *testing.T) {
	router := setupRouter()

	// simple in-memory mock implementing domain.TagRepository
	mock := &MockRepo{tags: []domain.Tag{}}
	var repo domain.TagRepository = mock
	RegisterEndpoints(router, &server.GinContextToPlainContextFactory{}, repo, &domain.FindAllTags{Repository: repo})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/tags/scope/expense", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[{\"key\":\"UNKNOWN\",\"value\":\"UNKNOWN\"}]", w.Body.String())
}

func TestFindAllTagsByScopeGiveMatchingTagsPlusSentinel(t *testing.T) {
	router := setupRouter()

	// simple in-memory mock implementing domain.TagRepository
	mock := &MockRepo{tags: []domain.Tag{
		{Key: "tag1", Value: "value1", Scope: "expense"},
		{Key: "tag2", Value: "value2", Scope: "expense"},
	}}

	var repo domain.TagRepository = mock
	RegisterEndpoints(router, &server.GinContextToPlainContextFactory{}, repo, &domain.FindAllTags{Repository: repo})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/tags/scope/expense", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[{\"key\":\"tag1\",\"value\":\"value1\",\"scope\":\"expense\"},{\"key\":\"tag2\",\"value\":\"value2\",\"scope\":\"expense\"},{\"key\":\"UNKNOWN\",\"value\":\"UNKNOWN\"}]", w.Body.String())
}

func TestPutTagRejectsBlankScope(t *testing.T) {
	router := setupRouter()

	mock := &MockRepo{tags: []domain.Tag{}}
	var repo domain.TagRepository = mock
	RegisterEndpoints(router, &server.GinContextToPlainContextFactory{}, repo, &domain.FindAllTags{Repository: repo})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/tags", strings.NewReader(`{"value":"Groceries","scope":"   "}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Empty(t, mock.tags, "no tag should be persisted when scope is blank")
}

func TestPutTagAcceptsNonBlankScope(t *testing.T) {
	router := setupRouter()

	mock := &MockRepo{tags: []domain.Tag{}}
	var repo domain.TagRepository = mock
	RegisterEndpoints(router, &server.GinContextToPlainContextFactory{}, repo, &domain.FindAllTags{Repository: repo})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/tags", strings.NewReader(`{"value":"Groceries","scope":"expense"}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Len(t, mock.tags, 1)
	assert.Equal(t, "expense", mock.tags[0].Scope)
	assert.NotEmpty(t, mock.tags[0].Key, "a server-side UUID key should be assigned")
}

func TestFindAllTagsByScopeFiltersOutNonMatchingAndScopelessTags(t *testing.T) {
	router := setupRouter()

	mock := &MockRepo{tags: []domain.Tag{
		{Key: "tag1", Value: "Groceries", Scope: "Expense"},
		{Key: "tag2", Value: "Salary", Scope: "revenue"},
		{Key: "tag3", Value: "Legacy", Scope: ""},
	}}

	var repo domain.TagRepository = mock
	RegisterEndpoints(router, &server.GinContextToPlainContextFactory{}, repo, &domain.FindAllTags{Repository: repo})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/tags/scope/EXPENSE", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[{\"key\":\"tag1\",\"value\":\"Groceries\",\"scope\":\"Expense\"},{\"key\":\"UNKNOWN\",\"value\":\"UNKNOWN\"}]", w.Body.String())
}

func TestFindAllTagsByScopeAlwaysIncludesUnknownSentinel(t *testing.T) {
	router := setupRouter()

	mock := &MockRepo{tags: []domain.Tag{
		{Key: "tag1", Value: "Groceries", Scope: "expense"},
	}}

	var repo domain.TagRepository = mock
	RegisterEndpoints(router, &server.GinContextToPlainContextFactory{}, repo, &domain.FindAllTags{Repository: repo})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/tags/scope/revenue", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[{\"key\":\"UNKNOWN\",\"value\":\"UNKNOWN\"}]", w.Body.String())
}

// MockRepo is a simple in-memory implementation of domain.TagRepository for tests.
type MockRepo struct {
	tags []domain.Tag
}

func (m *MockRepo) SaveTag(ctx context.Context, tag *domain.Tag) error {
	m.tags = append(m.tags, *tag)
	return nil
}

func (m *MockRepo) FindAllTags(ctx context.Context, scope string) ([]domain.Tag, error) {
	// Strict scope filtering — mirrors the real repository contract: only tags
	// whose normalized scope equals the requested scope are returned.
	var matching []domain.Tag
	for _, t := range m.tags {
		if domain.NormalizeScope(t.Scope) == scope {
			matching = append(matching, t)
		}
	}
	return matching, nil
}
