package tags

import "context"

type SearchTagKey = string
type SearchTagValue = string

type SearchTag struct {
	Key   SearchTagKey   `json:"key"`
	Value SearchTagValue `json:"value"`
}

type SearchTagRepository interface {
	GetTagBy(ctx context.Context, key string) (*SearchTag, error)
	GetAllTags(ctx context.Context) ([]SearchTag, error)
}
