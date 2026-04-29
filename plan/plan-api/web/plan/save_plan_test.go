package plan

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSavePlan(t *testing.T) {
	router := setupRouter(&mockRepo{})

	body, _ := json.Marshal(planRepresentation{Title: "new plan", Date: "2026-04-29"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/plan", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	_, err := uuid.Parse(resp["id"])
	assert.NoError(t, err, "response id should be a valid UUID")
}
