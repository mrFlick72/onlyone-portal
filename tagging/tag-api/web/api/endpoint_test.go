package api

import (
	"context"
	"net/http"
	"net/http/httptest"
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

func TestFindAllTagsGiveEmptyTags(t *testing.T) {
	router := setupRouter()

	// simple in-memory mock implementing domain.TagRepository
	mock := &MockRepo{tags: []domain.Tag{}}
	var repo domain.TagRepository = mock
	RegisterEndpoints(router, &server.GinContextToPlainContextFactory{}, repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/tags", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[]", w.Body.String())
}

func TestFindAllTagsGiveNotEmptyTags(t *testing.T) {
	router := setupRouter()

	// simple in-memory mock implementing domain.TagRepository
	mock := &MockRepo{tags: []domain.Tag{
		{
			Key:   "tag1",
			Value: "value1",
		},
		{
			Key:   "tag2",
			Value: "value2",
		},
	}}

	var repo domain.TagRepository = mock
	RegisterEndpoints(router, &server.GinContextToPlainContextFactory{}, repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/tags", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[{\"key\":\"tag1\",\"value\":\"value1\"},{\"key\":\"tag2\",\"value\":\"value2\"}]", w.Body.String())
}


// MockRepo is a simple in-memory implementation of domain.TagRepository for tests.
type MockRepo struct {
	tags []domain.Tag
}

func (m *MockRepo) SaveTag(ctx context.Context, tag *domain.Tag) error {
	m.tags = append(m.tags, *tag)
	return nil
}

func (m *MockRepo) GetTagBy(ctx context.Context, key string) (*domain.Tag, error) {
	for _, t := range m.tags {
		if t.Key == key {
			return &t, nil
		}
	}
	return nil, nil
}

func (m *MockRepo) FindAllTags(ctx context.Context) ([]domain.Tag, error) {
	return m.tags, nil
}
