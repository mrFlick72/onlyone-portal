package scheduledexpense

import (
	"errors"
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

	r.GET("/api/budget/scheduled-expense/:id", func(c *gin.Context) {
		ctx := ContextFactoryConverter.CreateContextFromGin(c)
		scheduledExpense, err := facade.FindScheduledExpense(ctx, c.Param("id"))
		if err != nil {
			logger.LogErrorfFor("Error finding scheduled expense: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if scheduledExpense == nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, ScheduledExpenseDomainToRepresentationModel(scheduledExpense))
	})

	r.PUT("/api/budget/scheduled-expense/:id", func(c *gin.Context) {
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
		domainModel.Id = c.Param("id")
		if err := facade.UpdateScheduledExpense(ctx, domainModel); err != nil {
			if errors.Is(err, scheduledexpense.ErrScheduledExpenseNotFound) {
				c.Status(http.StatusNotFound)
				return
			}
			logger.LogErrorfFor("Error updating scheduled expense: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	})

	r.DELETE("/api/budget/scheduled-expense/:id", func(c *gin.Context) {
		ctx := ContextFactoryConverter.CreateContextFromGin(c)
		if err := facade.DeleteScheduledExpense(ctx, c.Param("id")); err != nil {
			if errors.Is(err, scheduledexpense.ErrScheduledExpenseNotFound) {
				c.Status(http.StatusNotFound)
				return
			}
			logger.LogErrorfFor("Error deleting scheduled expense: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
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
		if err := facade.CreateScheduledExpense(ctx, domainModel); err != nil {
			logger.LogErrorfFor("Error creating scheduled expense: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusCreated)
	})

	return r
}
