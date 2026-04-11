package expense

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrflick72/budget/budget-api/domain/budget/expense"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
)

func RegisterExpenseEndpoints(
	r *gin.Engine,
	ContextFactoryConverter server.ContextFactoryConverter, //todo could be private
	facade expense.BudgetExpenseActions,
) *gin.Engine {

	var logger = logging.GetLoggerInstanceForComponentByTypeName("web.RegisterEndpoints")

	r.PUT("/api/budget/expense", func(c *gin.Context) {
		var representation BudgetSearchCriteriaRepresentation

		if err := c.ShouldBindJSON(&representation); err != nil {
			logger.LogErrorfFor("Error binding JSON: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := ContextFactoryConverter.CreateContextFromGin(c)
		spentBudget, err := facade.FindSpentBudget(ctx, date.NewMonthFor(representation.Month), date.NewYearFor(representation.Year), representation.SearchTagList)
		if err != nil {
			logger.LogErrorfFor("Error finding spent budget: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, SpentBudgetDomainToRepresentationModel(spentBudget))
	})

	r.PUT("/api/budget/expense/:id", func(c *gin.Context) {
		var representation BudgetExpenseRepresentation

		ctx := ContextFactoryConverter.CreateContextFromGin(c)
		if err := c.ShouldBindJSON(&representation); err != nil {
			logger.LogErrorfFor("Error binding JSON: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		domainModel := RepresentationModelToDomainModel(representation)
		domainModel.Id = c.Param("id")
		facade.UpdateBudgetExpense(ctx, domainModel)
		c.Status(http.StatusNoContent)
	})

	r.POST("/api/budget/expense", func(c *gin.Context) {
		var representation BudgetExpenseRepresentation

		ctx := ContextFactoryConverter.CreateContextFromGin(c)

		if err := c.ShouldBindJSON(&representation); err != nil {
			logger.LogErrorfFor("Error binding JSON: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		domainModel := RepresentationModelToDomainModel(representation)
		facade.CreateBudgetExpense(ctx, domainModel)
		c.Status(http.StatusCreated)
	})

	r.DELETE("/api/budget/expense/:id", func(c *gin.Context) {
		ctx := ContextFactoryConverter.CreateContextFromGin(c)
		facade.DeleteBudgetExpense(ctx, c.Param("id"))
		c.Status(http.StatusNoContent)
	})

	return r
}
