package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrflick72/onlyone-portal/account/account-api/domain/account"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/pkg/time/date"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
)

const ENDPOINT_PREFIX = "/api/account/user-account"

func RegisterEndpoints(
	r *gin.Engine,
	AccountUpdate *account.UpdateAccount,
	AccountRepository account.AccountRepository,
	contextFactoryConverter server.ContextFactoryConverter,
) *gin.Engine {

	var logger = logging.GetLoggerInstanceForComponentByTypeName("web.RegisterEndpoints")

	r.GET(ENDPOINT_PREFIX, func(c *gin.Context) {
		ctx := contextFactoryConverter.CreateContextFromGin(c)
		userAccount, err := AccountRepository.FindAnAccount(ctx)
		if err != nil {
			logger.LogErrorFor(err)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		formattedDate := getFormattedDate(userAccount, *logger)
		c.JSON(http.StatusOK, account.Account{
			FirstName: userAccount.FirstName,
			LastName:  userAccount.FirstName,
			BirthDate: formattedDate,
			Email:     userAccount.Email,
			Phone:     userAccount.Phone,
		})

	})

	r.PUT(ENDPOINT_PREFIX, func(c *gin.Context) {
		var userAccount account.Account
		ctx := contextFactoryConverter.CreateContextFromGin(c)

		if err := c.ShouldBindJSON(&userAccount); err != nil {
			logger.LogErrorfFor("Error binding JSON: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		formattedIsoDate := getFormattedIsoDate(&userAccount, *logger)
		userAccountToBeStored := &account.Account{
			FirstName: userAccount.FirstName,
			LastName:  userAccount.LastName,
			BirthDate: formattedIsoDate,
			Email:     userAccount.Email,
			Phone:     userAccount.Phone,
		}
		AccountRepository.Save(ctx, userAccountToBeStored)
		c.JSON(http.StatusNoContent, nil)
	})

	return r

}

func getFormattedDate(account *account.Account, logger logging.Logger) string {
	formattedDate, err := date.IsoDateFor(account.BirthDate)
	if err != nil {
		logger.LogErrorfFor("date parsing for %s is in error: %v", account.BirthDate, err)
		return ""
	} else {
		return formattedDate.GetFormattedDate()
	}
}

func getFormattedIsoDate(account *account.Account, logger logging.Logger) string {
	formattedDate, err := date.DateFor(account.BirthDate)
	if err != nil {
		logger.LogErrorfFor("date parsing for %s is in error: %v", account.BirthDate, err)
		return ""
	} else {
		return formattedDate.GetIsoFormattedDate()
	}
}
