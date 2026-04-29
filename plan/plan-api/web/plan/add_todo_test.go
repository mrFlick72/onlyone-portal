package plan

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/plan"
	"github.com/stretchr/testify/assert"
)

func TestAddTodo(t *testing.T) {
	p := aTestPlan()
	router := setupRouter(&mockRepo{plans: []*plan.Plan{p}})

	body, _ := json.Marshal(todoRepresentation{Date: "2026-04-29", Content: "do something"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/plan/"+p.Id+"/todo", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}
