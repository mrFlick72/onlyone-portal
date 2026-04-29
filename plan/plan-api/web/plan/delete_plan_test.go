package plan

import (
	"github.com/mrflick72/onlyone-portal/plan/plan-service/internal/test"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeletePlan(t *testing.T) {
	aPlan := test.ANewPlan()
	repo := &mockRepo{}
	repo.On("DeletePlan", aPlan.Id, testUser).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/plan/"+aPlan.Id, nil)
	setupRouter(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	repo.AssertExpectations(t)
}
