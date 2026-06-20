package domain

import "context"

type FindAllTags struct {
	Repository TagRepository
}

// Execute returns the user's tags plus the UNKNOWN sentinel. An empty scope
// returns every tag unfiltered; a non-empty scope returns only tags whose
// normalized Scope matches it.
func (action *FindAllTags) Execute(ctx context.Context, scope string) ([]Tag, error) {
	tags, err := action.Repository.FindAllTags(ctx, NormalizeScope(scope))
	if err != nil {
		return nil, err
	}
	return append(tags, Tag{Key: "UNKNOWN", Value: "UNKNOWN"}), nil
}
