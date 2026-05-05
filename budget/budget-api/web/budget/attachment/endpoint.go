package attachment

import (
	"github.com/gin-gonic/gin"
	"github.com/mrflick72/budget/budget-api/domain/budget/attachment"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
)

func RegisterExpenseEndpoints(
	r *gin.Engine,
	facade attachment.AttachmentActionsFacade,
) *gin.Engine {

	var logger = logging.GetLoggerInstanceForComponentByTypeName("attachment.RegisterExpenseEndpoints")

	r.PUT("/api/attachment", func(context *gin.Context) {
		logger.LogInfofFor("")
	})

	r.GET("/api/attachment", func(context *gin.Context) {

	})

	r.GET("/api/attachment/:attachmentId/content", func(context *gin.Context) {

	})
	r.DELETE("/api/attachment/:attachmentId", func(context *gin.Context) {

	})

	return r
}
