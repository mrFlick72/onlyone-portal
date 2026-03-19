package rest

import (
	"context"

	"github.com/mrflick72/onlyone-portal/account/account-api/domain/account"
)

type VauthenticatorAccountRepository struct {
	BaseUrl string
}

func (r *VauthenticatorAccountRepository) FindAnAccount(ctx context.Context) (*account.Account, error) {
	return nil, nil
}

func (r *VauthenticatorAccountRepository) Save(ctx context.Context, account *account.Account) error {
	return nil
}
