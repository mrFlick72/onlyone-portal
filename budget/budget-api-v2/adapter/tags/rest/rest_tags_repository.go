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

var logger = logging.GetLoggerInstance()

type RestSearchTagRepository struct {
	Client  *http.Client
	BaseURL string
}

func (repository *RestSearchTagRepository) GetTagBy(ctx context.Context, key string) (*tags.SearchTag, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/tags", repository.BaseURL), nil)
	if err != nil {
		logger.LogErrorfFor("Error while calling tag API: %s", err)
		return nil, err
	}
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		logger.LogErrorfFor("Error while getting current user: %s", err)
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *user.AccessToken))

	resp, err := repository.Client.Do(req)
	if err != nil {
		logger.LogErrorfFor("Error while calling tag API: %s", err)
		return nil, err
	}
	defer resp.Body.Close()
	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.LogErrorfFor("Error while calling tag API: %s", err)
		return nil, err
	}
	var searchTags []tags.SearchTag
	err = json.Unmarshal(body, &searchTags)
	if err != nil {
		logger.LogErrorfFor("Error while un marshalling tag API response: %s", err)
		return nil, err
	}
	// todo: this is not efficient, we should have an endpoint to get a tag by key
	logger.LogInfofFor("Received search tags from tag API: %v", searchTags)
	for _, searchTag := range searchTags {
		if searchTag.Key == key {
			return &searchTag, nil
		}
	}
	return nil, fmt.Errorf("tag with key '%s' not found", key)

}
