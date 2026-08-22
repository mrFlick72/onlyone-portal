package web

import (
	"bytes"
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

func (m *mockMfaRepository) StartEnrollment(ctx context.Context, mfaMethod string, mfaChannel string) (string, error) {
	args := m.Called(ctx, mfaMethod, mfaChannel)
	return args.String(0), args.Error(1)
}

func (m *mockMfaRepository) Associate(ctx context.Context, ticket string, code string) error {
	args := m.Called(ctx, ticket, code)
	return args.Error(0)
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

func TestStartMfaEnrollment(t *testing.T) {
	repo := &mockMfaRepository{}
	repo.On("StartEnrollment", mock.Anything, "EMAIL_MFA_METHOD", "jane@example.com").Return("a-ticket", nil)

	w := httptest.NewRecorder()
	body := bytes.NewBufferString(`{"mfaMethod":"EMAIL_MFA_METHOD","mfaChannel":"jane@example.com"}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/account/mfa/enrollment", body)
	setupMfaRouter(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, `{"ticket":"a-ticket"}`, w.Body.String())
	repo.AssertExpectations(t)
}

func TestStartMfaEnrollmentWhenRepositoryFails(t *testing.T) {
	repo := &mockMfaRepository{}
	repo.On("StartEnrollment", mock.Anything, "EMAIL_MFA_METHOD", "jane@example.com").Return("", assert.AnError)

	w := httptest.NewRecorder()
	body := bytes.NewBufferString(`{"mfaMethod":"EMAIL_MFA_METHOD","mfaChannel":"jane@example.com"}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/account/mfa/enrollment", body)
	setupMfaRouter(repo).ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	repo.AssertExpectations(t)
}

func TestAssociateMfaEnrollment(t *testing.T) {
	repo := &mockMfaRepository{}
	repo.On("Associate", mock.Anything, "a-ticket", "123456").Return(nil)

	w := httptest.NewRecorder()
	body := bytes.NewBufferString(`{"ticket":"a-ticket","code":"123456"}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/account/mfa/associate", body)
	setupMfaRouter(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	repo.AssertExpectations(t)
}

func TestAssociateMfaEnrollmentWhenRepositoryFails(t *testing.T) {
	repo := &mockMfaRepository{}
	repo.On("Associate", mock.Anything, "a-ticket", "000000").Return(assert.AnError)

	w := httptest.NewRecorder()
	body := bytes.NewBufferString(`{"ticket":"a-ticket","code":"000000"}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/account/mfa/associate", body)
	setupMfaRouter(repo).ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	repo.AssertExpectations(t)
}
