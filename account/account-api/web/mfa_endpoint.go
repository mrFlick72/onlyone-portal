package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrflick72/onlyone-portal/account/account-api/domain/mfa"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
)

const MFA_ENDPOINT_PREFIX = "/api/account/mfa"

func RegisterMfaEndpoints(
	r *gin.Engine,
	MfaRepository mfa.MfaRepository,
	contextFactoryConverter server.ContextFactoryConverter,
) *gin.Engine {

	var logger = logging.GetLoggerInstanceForComponentByTypeName("web.RegisterMfaEndpoints")

	r.GET(MFA_ENDPOINT_PREFIX, func(c *gin.Context) {
		ctx := contextFactoryConverter.CreateContextFromGin(c)
		devices, err := MfaRepository.FindAll(ctx)
		if err != nil {
			logger.LogErrorFor(err)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		c.JSON(http.StatusOK, devices)
	})

	return r
}
