package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mrflick72/onlyone-portal/account/account-api/domain/mfa"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
)

type VauthenticatorMfaRepository struct {
	Client  *http.Client
	BaseUrl string
	Logger  *logging.Logger
}

func (r *VauthenticatorMfaRepository) FindAll(ctx context.Context) ([]mfa.MfaDevice, error) {
	req, err := r.newAuthorizedRequest(ctx, http.MethodGet, fmt.Sprintf("%s/api/mfa/enrollment", r.BaseUrl), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := r.Client.Do(req)
	if err != nil {
		r.Logger.LogErrorfFor("error while calling vauthenticator mfa enrollment endpoint: %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, r.unexpectedStatusError(resp, "mfa enrollment list endpoint")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		r.Logger.LogErrorfFor("error while reading vauthenticator mfa enrollment response: %s", err)
		return nil, err
	}

	var devices []mfa.MfaDevice
	if err = json.Unmarshal(body, &devices); err != nil {
		r.Logger.LogErrorfFor("error while unmarshalling vauthenticator mfa enrollment response: %s", err)
		return nil, err
	}

	return devices, nil
}

type enrollmentRequest struct {
	MfaChannel string `json:"mfaChannel"`
	MfaMethod  string `json:"mfaMethod"`
}

type enrollmentResponse struct {
	Ticket string `json:"ticket"`
}

type associationRequest struct {
	Ticket string `json:"ticket"`
	Code   string `json:"code"`
}

func (r *VauthenticatorMfaRepository) StartEnrollment(ctx context.Context, mfaMethod string, mfaChannel string) (string, error) {
	body, err := json.Marshal(enrollmentRequest{MfaChannel: mfaChannel, MfaMethod: mfaMethod})
	if err != nil {
		r.Logger.LogErrorfFor("error while marshalling mfa enrollment payload: %s", err)
		return "", err
	}

	req, err := r.newAuthorizedRequest(ctx, http.MethodPost, fmt.Sprintf("%s/api/mfa/enrollment", r.BaseUrl), body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.Client.Do(req)
	if err != nil {
		r.Logger.LogErrorfFor("error while calling vauthenticator mfa enrollment endpoint: %s", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", r.unexpectedStatusError(resp, "mfa enrollment start endpoint")
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		r.Logger.LogErrorfFor("error while reading vauthenticator mfa enrollment response: %s", err)
		return "", err
	}

	var enrollment enrollmentResponse
	if err = json.Unmarshal(responseBody, &enrollment); err != nil {
		r.Logger.LogErrorfFor("error while unmarshalling vauthenticator mfa enrollment response: %s", err)
		return "", err
	}

	return enrollment.Ticket, nil
}

func (r *VauthenticatorMfaRepository) Associate(ctx context.Context, ticket string, code string) error {
	body, err := json.Marshal(associationRequest{Ticket: ticket, Code: code})
	if err != nil {
		r.Logger.LogErrorfFor("error while marshalling mfa association payload: %s", err)
		return err
	}

	req, err := r.newAuthorizedRequest(ctx, http.MethodPost, fmt.Sprintf("%s/api/mfa/associate", r.BaseUrl), body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.Client.Do(req)
	if err != nil {
		r.Logger.LogErrorfFor("error while calling vauthenticator mfa associate endpoint: %s", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return r.unexpectedStatusError(resp, "mfa associate endpoint")
	}

	return nil
}

func (r *VauthenticatorMfaRepository) newAuthorizedRequest(ctx context.Context, method string, url string, body []byte) (*http.Request, error) {
	var requestBody io.Reader
	if body != nil {
		requestBody = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, requestBody)
	if err != nil {
		r.Logger.LogErrorfFor("error while creating request for vauthenticator: %s", err)
		return nil, err
	}

	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		r.Logger.LogErrorfFor("error while getting current user: %s", err)
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *user.AccessToken))
	return req, nil
}

func (r *VauthenticatorMfaRepository) unexpectedStatusError(resp *http.Response, endpoint string) error {
	r.Logger.LogErrorfFor("vauthenticator %s returned unexpected status: %d", endpoint, resp.StatusCode)
	return fmt.Errorf("vauthenticator %s returned status %d", endpoint, resp.StatusCode)
}
