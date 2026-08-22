package web

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mrflick72/onlyone-portal/account/account-api/domain/mfa"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupMfaRouter(repo mfa.MfaRepository) *gin.Engine {
	r := gin.Default()
	RegisterMfaEndpoints(r, repo, &server.GinContextToPlainContextFactory{})
	return r
}

type mockMfaRepository struct {
	mock.Mock
}

func (m *mockMfaRepository) FindAll(ctx context.Context) ([]mfa.MfaDevice, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]mfa.MfaDevice), args.Error(1)
}

func TestFindAllMfaDevices(t *testing.T) {
	aDevice := mfa.MfaDevice{
		UserName:    "j***n@e***.com",
		MfaMethod:   "EMAIL_MFA_METHOD",
		MfaChannel:  "j***n@e***.com",
		MfaDeviceId: "a-device-id",
		Default:     true,
	}
	repo := &mockMfaRepository{}
	repo.On("FindAll", mock.Anything).Return([]mfa.MfaDevice{aDevice}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/account/mfa", nil)
	setupMfaRouter(repo).ServeHTTP(w, req)

	expected, _ := json.Marshal([]mfa.MfaDevice{aDevice})
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(expected), strings.TrimSpace(w.Body.String()))
	repo.AssertExpectations(t)
}

func TestFindAllMfaDevicesWhenEmpty(t *testing.T) {
	repo := &mockMfaRepository{}
	repo.On("FindAll", mock.Anything).Return([]mfa.MfaDevice{}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/account/mfa", nil)
	setupMfaRouter(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[]", strings.TrimSpace(w.Body.String()))
	repo.AssertExpectations(t)
}

func TestFindAllMfaDevicesWhenRepositoryFails(t *testing.T) {
	repo := &mockMfaRepository{}
	repo.On("FindAll", mock.Anything).Return(nil, assert.AnError)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/account/mfa", nil)
	setupMfaRouter(repo).ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	repo.AssertExpectations(t)
}
