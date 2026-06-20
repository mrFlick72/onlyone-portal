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

type FindTagsByScope struct {
	Repository TagRepository
}

func (action *FindTagsByScope) Execute(ctx context.Context, scope string) ([]Tag, error) {
	tags, err := action.Repository.FindTagsByScope(ctx, NormalizeScope(scope))
	if err != nil {
		return nil, err
	}
	return append(tags, Tag{Key: "UNKNOWN", Value: "UNKNOWN"}), nil
}
