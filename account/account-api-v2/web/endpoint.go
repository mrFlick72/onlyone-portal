package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrflick72/onlyone-portal/account/account-api/domain/account"
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

	r.GET(ENDPOINT_PREFIX, func(ctx *gin.Context) {
		appCtx := contextFactoryConverter.CreateContextFromGin(ctx)
		account, err := AccountRepository.FindAnAccount(appCtx)
		if err != nil {

			ctx.JSON(http.StatusInternalServerError, nil)
		}

		ctx.JSON(http.StatusInternalServerError, account)

	})

	r.PUT(ENDPOINT_PREFIX, func(ctx *gin.Context) {
		ctx.JSON(http.StatusNoContent, nil)
	})
	return r
}
