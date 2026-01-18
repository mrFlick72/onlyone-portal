package tags

import "context"

type SearchTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type TagRepository interface {
	SaveTag(ctx *context.Context, tag *SearchTag) error
	GetTagBy(ctx *context.Context, key string) (*SearchTag, error)
	FindAllTags(ctx *context.Context) (*[]SearchTag, error)
}
