package i18n-api

import (
"context"
"errors"
)

var ErrBundleNotFound = errors.New("message bundle not found")

type BundleRepository interface {
	GetBundle(ctx context.Context, key BundleKey) (*MessageBundle, error)
}
