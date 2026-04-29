package plan

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeletePlan(t *testing.T) {
	p := aTestPlan()
	repo := &mockRepo{}
	repo.On("DeletePlan", p.Id, testUser).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/plan/"+p.Id, nil)
	setupRouter(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	repo.AssertExpectations(t)
}
