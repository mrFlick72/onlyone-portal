package plan

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/plan"
	"github.com/stretchr/testify/assert"
)

func TestDeletePlan(t *testing.T) {
	p := aTestPlan()
	router := setupRouter(&mockRepo{plans: []*plan.Plan{p}})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/plan/"+p.Id, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
