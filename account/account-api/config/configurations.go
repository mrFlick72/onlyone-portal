package config

import (
	"net/http"
	"sync"

	"github.com/mrflick72/onlyone-portal/account/account-api/adapter/rest"
	"github.com/mrflick72/onlyone-portal/account/account-api/domain/account"
	"github.com/mrflick72/onlyone-portal/account/account-api/domain/mfa"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/config"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/httpclient"
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

var mfaRepository mfa.MfaRepository

var mfaOnce sync.Once

func NewVauthenticatorMfaRepository() mfa.MfaRepository {
	mfaOnce.Do(func() {
		mfaRepository = &rest.VauthenticatorMfaRepository{
			Client:  httpclient.NewHTTPClient(),
			BaseUrl: configurationManager.GetConfigFor("idp.base-url"),
			Logger:  logging.GetLoggerInstanceForComponentByType(&rest.VauthenticatorMfaRepository{}),
		}
	})

	return mfaRepository
}
