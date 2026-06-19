package domain

import "context"

type FindAllTags struct {
	Repository TagRepository
}

func (action *FindAllTags) Execute(ctx context.Context) ([]Tag, error) {
	tags, err := action.Repository.FindAllTags(ctx)
	if err != nil {
		return nil, err
	}
	return append(tags, Tag{Key: "UNKNOWN", Value: "UNKNOWN"}), nil
}
