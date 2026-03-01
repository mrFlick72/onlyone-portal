//go:build test

package tags

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type SearchTagRepositoryMock struct {
	mock.Mock
}

func (mock *SearchTagRepositoryMock) GetTagBy(ctx context.Context, searchTagKey SearchTagKey) (*SearchTag, error) {
	args := mock.Called(ctx, searchTagKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SearchTag), args.Error(1)
}
