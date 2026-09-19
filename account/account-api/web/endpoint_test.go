package web

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mrflick72/onlyone-portal/account/account-api/domain/account"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockAccountRepository struct {
	mock.Mock
}

func (m *mockAccountRepository) FindAnAccount(ctx context.Context) (*account.Account, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*account.Account), args.Error(1)
}

func (m *mockAccountRepository) Save(ctx context.Context, acc *account.Account) error {
	args := m.Called(ctx, acc)
	return args.Error(0)
}

func setupAccountRouter(repo account.AccountRepository) *gin.Engine {
	r := gin.Default()
	RegisterEndpoints(r, &account.UpdateAccount{AccountRepository: repo}, repo, &server.GinContextToPlainContextFactory{})
	return r
}

func TestFindAnAccountReturnsLocale(t *testing.T) {
	repo := new(mockAccountRepository)
	repo.On("FindAnAccount", mock.Anything).Return(&account.Account{
		FirstName: "Mario",
		Email:     "mario@example.com",
		Locale:    "it",
	}, nil)

	router := setupAccountRouter(repo)
	request := httptest.NewRequest(http.MethodGet, ENDPOINT_PREFIX, nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
	var response account.Account
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	assert.Equal(t, "it", response.Locale)
}

func TestSaveAccountForwardsLocale(t *testing.T) {
	repo := new(mockAccountRepository)
	repo.On("Save", mock.Anything, mock.MatchedBy(func(acc *account.Account) bool {
		return acc.Locale == "en"
	})).Return(nil)

	router := setupAccountRouter(repo)
	body, _ := json.Marshal(account.Account{FirstName: "Mario", Locale: "en"})
	request := httptest.NewRequest(http.MethodPut, ENDPOINT_PREFIX, bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusNoContent, recorder.Code)
	repo.AssertExpectations(t)
}
