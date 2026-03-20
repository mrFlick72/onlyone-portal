package account

import "context"

type UpdateAccount struct {
	AccountRepository AccountRepository
}

func (action *UpdateAccount) Execute(ctx context.Context, account *Account) error {
	return action.AccountRepository.Save(ctx, account)
}
