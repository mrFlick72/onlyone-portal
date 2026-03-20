package config

import (
	"net/http"
	"sync"

	"github.com/mrflick72/onlyone-portal/account/account-api/adapter/rest"
	"github.com/mrflick72/onlyone-portal/account/account-api/domain/account"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/config"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
)

var configurationManager = config.GetConfigurationManagerInstance()

var accountRepository account.AccountRepository

var once sync.Once

func NewVauthenticatorAccountRepository() account.AccountRepository {
	once.Do(func() {
		accountRepository = &rest.VauthenticatorAccountRepository{
			Client:  &http.Client{},
			BaseUrl: configurationManager.GetConfigFor("idp.base-url"),
			Logger:  logging.GetLoggerInstanceForComponentByType(&rest.VauthenticatorAccountRepository{}),
		}
	})

	return accountRepository
}

func NewAccountUpdate() *account.UpdateAccount {
	return &account.UpdateAccount{
		AccountRepository: NewVauthenticatorAccountRepository(),
	}
}
