package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
)

type RestSearchTagRepository struct {
	Client  *http.Client
	BaseURL string
	Scope   string
	logger  *logging.Logger
}

func NewRestSearchTagRepository(client *http.Client, baseURL string, scope string) tags.SearchTagRepository {
	return &RestSearchTagRepository{
		Client:  client,
		BaseURL: baseURL,
		Scope:   scope,
		logger:  logging.GetLoggerInstanceForComponentByType(&RestSearchTagRepository{}),
	}
}

func (repository *RestSearchTagRepository) GetTagBy(ctx context.Context, key string) (*tags.SearchTag, error) {
	searchTags, err := repository.GetAllTags(ctx)
	if err != nil {
		return nil, err
	}

	for _, searchTag := range searchTags {
		if searchTag.Key == key {
			return &searchTag, nil
		}
	}
	// key isn't in the catalog — most likely because the tag was deleted
	// (tag-api ADR 0008). Degrade to the UNKNOWN sentinel in place rather
	// than erroring; see
	// docs/adr/0003-deleted-tag-references-resolve-to-unknown-in-place.md.
	unknown := tags.UnknownSentinel()
	return &unknown, nil
}

func (repository *RestSearchTagRepository) GetAllTags(ctx context.Context) ([]tags.SearchTag, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/tags/scope/%s", repository.BaseURL, repository.Scope), nil)
	if err != nil {
		repository.logger.LogErrorfFor("Error while calling tag API: %s", err)
		return nil, err
	}
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		repository.logger.LogErrorfFor("Error while getting current user: %s", err)
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *user.AccessToken))

	resp, err := repository.Client.Do(req)
	if err != nil {
		repository.logger.LogErrorfFor("Error while calling tag API: %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// A non-2xx response must never fall through to the unmarshal path below:
		// an empty or stub body on a degraded/misrouted response would otherwise
		// unmarshal into a valid, empty []SearchTag — indistinguishable from "this
		// user genuinely has zero tags in this scope" and, via GetTagBy's
		// not-found-in-catalog fallback (ADR 0003), silently drives durable
		// UNKNOWN reclassification for every tag key looked up in this window
		// (see docs/adr/0004-tag-catalog-fetch-must-validate-http-status-before-trusting-empty-result.md).
		repository.logger.LogErrorfFor("Unexpected status code %d from tag API", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status code %d from tag API", resp.StatusCode)
	}

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		repository.logger.LogErrorfFor("Error while calling tag API: %s", err)
		return nil, err
	}
	var searchTags []tags.SearchTag
	err = json.Unmarshal(body, &searchTags)
	if err != nil {
		repository.logger.LogErrorfFor("Error while un marshalling tag API response: %s", err)
		return nil, err
	}
	return searchTags, nil
}
