package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mrflick72/onlyone-portal/account/account-api/domain/account"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
)

type VauthenticatorAccountRepository struct {
	Client  *http.Client
	BaseUrl string
	logger  *logging.Logger
}

type vauthenticatorUserInfoResponse struct {
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Birthdate  string `json:"birthdate"`
	Email      string `json:"email"`
	Phone      string `json:"phone_number"`
}

func NewVauthenticatorAccountRepository(client *http.Client, baseURL string) account.AccountRepository {
	return &VauthenticatorAccountRepository{
		Client:  client,
		BaseUrl: baseURL,
		logger:  logging.GetLoggerInstanceForComponentByType(&VauthenticatorAccountRepository{}),
	}
}

func (r *VauthenticatorAccountRepository) FindAnAccount(ctx context.Context) (*account.Account, error) {
	req, err := r.newAuthorizedRequest(ctx, http.MethodGet, fmt.Sprintf("%s/userinfo", r.BaseUrl), nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.Client.Do(req)
	if err != nil {
		r.logger.LogErrorfFor("error while calling vauthenticator userinfo endpoint: %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, r.unexpectedStatusError(resp, "userinfo endpoint")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		r.logger.LogErrorfFor("error while reading vauthenticator userinfo response: %s", err)
		return nil, err
	}

	var userInfo vauthenticatorUserInfoResponse
	if err = json.Unmarshal(body, &userInfo); err != nil {
		r.logger.LogErrorfFor("error while unmarshalling vauthenticator userinfo response: %s", err)
		return nil, err
	}

	return &account.Account{
		FirstName: userInfo.GivenName,
		LastName:  userInfo.FamilyName,
		BirthDate: userInfo.Birthdate,
		Email:     userInfo.Email,
		Phone:     userInfo.Phone,
	}, nil
}

func (r *VauthenticatorAccountRepository) Save(ctx context.Context, account *account.Account) error {
	body, err := json.Marshal(account)
	if err != nil {
		r.logger.LogErrorfFor("error while marshalling account payload: %s", err)
		return err
	}

	req, err := r.newAuthorizedRequest(ctx, http.MethodPut, fmt.Sprintf("%s/api/accounts", r.BaseUrl), body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.Client.Do(req)
	if err != nil {
		r.logger.LogErrorfFor("error while calling vauthenticator account update endpoint: %s", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return r.unexpectedStatusError(resp, "account update endpoint")
	}

	return nil
}

func (r *VauthenticatorAccountRepository) newAuthorizedRequest(ctx context.Context, method string, url string, body []byte) (*http.Request, error) {
	var requestBody io.Reader
	if body != nil {
		requestBody = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, requestBody)
	if err != nil {
		r.logger.LogErrorfFor("error while creating request for vauthenticator: %s", err)
		return nil, err
	}

	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		r.logger.LogErrorfFor("error while getting current user: %s", err)
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *user.AccessToken))
	return req, nil
}

func (r *VauthenticatorAccountRepository) unexpectedStatusError(resp *http.Response, endpoint string) error {
	r.logger.LogErrorfFor("vauthenticator %s returned unexpected status: %d", endpoint, resp.StatusCode)
	return fmt.Errorf("vauthenticator %s returned status %d", endpoint, resp.StatusCode)
}
