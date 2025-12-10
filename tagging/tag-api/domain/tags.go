package domain

import "context"

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type TagRepository interface {
	SaveTag(ctx *context.Context, tag *Tag) error
	GetTagBy(ctx *context.Context, key string) (*Tag, error)
	FindAllTags(ctx *context.Context) (*[]Tag, error)
}
