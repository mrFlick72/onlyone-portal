package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrflick72/onlyone-portal/account/account-api/domain/mfa"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
)

const MFA_ENDPOINT_PREFIX = "/api/account/mfa"

type mfaEnrollmentRequest struct {
	MfaMethod  string `json:"mfaMethod"`
	MfaChannel string `json:"mfaChannel"`
}

type mfaEnrollmentResponse struct {
	Ticket string `json:"ticket"`
}

type mfaAssociationRequest struct {
	Ticket string `json:"ticket"`
	Code   string `json:"code"`
}

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

	r.POST(MFA_ENDPOINT_PREFIX+"/enrollment", func(c *gin.Context) {
		var request mfaEnrollmentRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			logger.LogErrorfFor("Error binding JSON: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := contextFactoryConverter.CreateContextFromGin(c)
		ticket, err := MfaRepository.StartEnrollment(ctx, request.MfaMethod, request.MfaChannel)
		if err != nil {
			logger.LogErrorFor(err)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		c.JSON(http.StatusCreated, mfaEnrollmentResponse{Ticket: ticket})
	})

	r.POST(MFA_ENDPOINT_PREFIX+"/associate", func(c *gin.Context) {
		var request mfaAssociationRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			logger.LogErrorfFor("Error binding JSON: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := contextFactoryConverter.CreateContextFromGin(c)
		if err := MfaRepository.Associate(ctx, request.Ticket, request.Code); err != nil {
			logger.LogErrorFor(err)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		c.JSON(http.StatusNoContent, nil)
	})

	return r
}
