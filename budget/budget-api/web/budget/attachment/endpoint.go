package attachment

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrflick72/budget/budget-api/domain/budget/attachment"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
)

func RegisterAttachmentEndpoints(
	r *gin.Engine,
	contextFactoryConverter server.ContextFactoryConverter,
	facade attachment.AttachmentActions,
) {

	var logger = logging.GetLoggerInstanceForComponentByTypeName("attachment.RegisterAttachmentEndpoints")

	r.POST("/api/attachment", func(c *gin.Context) {
		fileHeader, err := c.FormFile("file")
		if err != nil {
			logger.LogErrorfFor("missing 'file' part: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing 'file' part"})
			return
		}

		budgetId := c.PostForm("budgetId")
		attachmentId := c.PostForm("attachmentId")
		if budgetId == "" {
			logger.LogErrorFor("missing 'budgetId' field")
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing 'budgetId' field"})
			return
		}

		opened, err := fileHeader.Open()
		if err != nil {
			logger.LogErrorfFor("cannot open uploaded file: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer opened.Close()

		content, err := io.ReadAll(opened)
		if err != nil {
			logger.LogErrorfFor("cannot read uploaded file: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx := contextFactoryConverter.CreateContextFromGin(c)
		domainAttachment := &attachment.Attachment{
			AttachmentMetadata: attachment.AttachmentMetadata{
				BudgetId:    budgetId,
				FineName:    fileHeader.Filename,
				ContentType: fileHeader.Header.Get("Content-Type"),
			},
			Content: content,
		}

		if attachmentId != "" {
			domainAttachment.AttachmentId = attachmentId
		}

		if err := facade.SaveAttachment(ctx, domainAttachment); err != nil {
			logger.LogErrorfFor("error adding attachment: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	})

	r.GET("/api/attachment", func(context *gin.Context) {

	})

	r.GET("/api/attachment/:attachmentId/content", func(context *gin.Context) {

	})

	r.DELETE("/api/attachment/:attachmentId", func(context *gin.Context) {

	})
}
