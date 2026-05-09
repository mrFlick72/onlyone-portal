package attachment

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrflick72/budget/budget-api/domain/budget/attachment"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
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
		budgetType := c.PostForm("budgetType")
		rawDate := c.PostForm("date")
		if budgetType == "" {
			logger.LogErrorFor("missing 'budgetType' field")
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing 'budgetType' field"})
			return
		}
		if budgetId == "" {
			logger.LogErrorFor("missing 'budgetId' field")
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing 'budgetId' field"})
			return
		}
		if rawDate == "" {
			logger.LogErrorFor("missing 'date' field")
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing 'date' field"})
			return
		}
		parsedDate, err := date.DateFor(rawDate)
		if err != nil {
			logger.LogErrorfFor("invalid 'date' field: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'date' field"})
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
				BudgetType:  budgetType,
				Date:        *parsedDate,
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

	r.GET("/api/attachment/metadata/:budgetType/:budgetId", func(c *gin.Context) {
		ctx := contextFactoryConverter.CreateContextFromGin(c)

		budgetType := c.Param("budgetType")
		budgetId := c.Param("budgetId")

		logger.LogInfofFor("budgetType: %v", budgetType)
		logger.LogInfofFor("budgetId: %v", budgetId)
		attachments, err := facade.GetAttachments(ctx, budgetId, attachment.BudgetType(budgetType))
		if err != nil {
			logger.LogErrorfFor("error getting attachments metadata: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		representation := DomainModelToRepresentationModel(attachments)
		c.JSON(http.StatusOK, representation)
	})

	r.GET("/api/attachment/:attachmentId/content", func(c *gin.Context) {
		ctx := contextFactoryConverter.CreateContextFromGin(c)

		attachmentId := c.Param("attachmentId")

		logger.LogInfofFor("attachmentId: %v", attachmentId)
		att, err := facade.GetAttachmentBy(ctx, attachmentId)
		if err != nil {
			logger.LogErrorfFor("error getting attachment content for id %s: %v\n", attachmentId, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		contentType := att.ContentType
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, att.FineName))
		c.Data(http.StatusOK, contentType, att.Content)
	})

	r.DELETE("/api/attachment/:attachmentId", func(context *gin.Context) {

	})
}
