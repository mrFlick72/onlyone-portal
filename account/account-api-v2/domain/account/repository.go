package account

import "context"

type AccountRepository interface {
	FindAnAccount(ctx context.Context) (*Account, error)

	Save(ctx context.Context, account *Account) error
}
