package scheduledexpense

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrflick72/budget/budget-api/domain/budget/scheduledexpense"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
)

func RegisterScheduledExpenseEndpoints(
	r *gin.Engine,
	ContextFactoryConverter server.ContextFactoryConverter,
	facade scheduledexpense.ScheduledExpenseActions,
) *gin.Engine {

	var logger = logging.GetLoggerInstanceForComponentByTypeName("scheduledexpense.RegisterScheduledExpenseEndpoints")

	r.GET("/api/budget/scheduled-expense", func(c *gin.Context) {
		ctx := ContextFactoryConverter.CreateContextFromGin(c)
		scheduledExpenses, err := facade.FindScheduledExpenses(ctx)
		if err != nil {
			logger.LogErrorfFor("Error finding scheduled expenses: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, ScheduledExpenseListDomainToRepresentationModel(scheduledExpenses))
	})

	r.POST("/api/budget/scheduled-expense", func(c *gin.Context) {
		var representation ScheduledExpenseRepresentation

		ctx := ContextFactoryConverter.CreateContextFromGin(c)
		if err := c.ShouldBindJSON(&representation); err != nil {
			logger.LogErrorfFor("Error binding JSON: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		domainModel, err := ScheduledExpenseRepresentationToDomainModel(representation)
		if err != nil {
			logger.LogErrorfFor("Error converting scheduled expense representation: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		facade.CreateScheduledExpense(ctx, domainModel)
		c.Status(http.StatusCreated)
	})

	return r
}
