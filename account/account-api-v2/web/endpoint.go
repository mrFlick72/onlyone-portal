package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrflick72/onlyone-portal/account/account-api/domain/account"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/pkg/time/date"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
)

/*
	    app.get(ENDPOINT_PREFIX, async (req: Request, res: Response) => {
	        let account = await accountRepository.findAnAccount();
	        const formattedDate = getFormattedDate(account);

	        res.status(200)
	            .json(
	                {
	                    firstName: account.firstName,
	                    lastName: account.lastName,
	                    birthDate: formattedDate,
	                    email: account.email,
	                    phone: account.phone
	                }
	            )
	            .end()
	    });

	    app.put(ENDPOINT_PREFIX, (req: Request, res: Response) => {
	        const account = req.body as Account
	        const username = SecurityContextHolder.getStore()?.userName!!

	        const parsedBirthDate = moment(account.birthDate, "DD/MM/YYYY") || "";
	        const accountToBeStored = {
	            firstName: account.firstName,
	            lastName: account.lastName,
	            birthDate: moment(parsedBirthDate).format("YYYY-MM-DD") || "",
	            email: username,
	            phone: account.phone
	        };

	        updateAccount.execute(accountToBeStored)
	            .then(_ => {
	                res.status(204).end()
	            }).catch(err => {
	                console.error(err)
	                res.status(500).end()
	            });
	    });


	    const getFormattedDate = (account: Account): string => {
	        let formattedDate: string = "";
	        try {
	            formattedDate = moment(account.birthDate, "YYYY-MM-DD").format("DD/MM/YYYY")
	        } catch (e) {
	        }
	        return formattedDate;
	    }
	}
*/
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
			c.JSON(http.StatusInternalServerError, nil)
		}

		formattedDate := getFormattedDate(userAccount)
		c.JSON(http.StatusInternalServerError, account.Account{
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

		formattedIsoDate := getFormattedIsoDate(&userAccount)
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

func getFormattedDate(account *account.Account) string {

	formattedDate, err := date.IsoDateFor(account.BirthDate)
	if err != nil {
		formattedDate.GetFormattedDate()
	}
	return ""
}

func getFormattedIsoDate(account *account.Account) string {

	formattedDate, err := date.DateFor(account.BirthDate)
	if err != nil {
		formattedDate.GetIsoFormattedDate()
	}
	return ""
}
