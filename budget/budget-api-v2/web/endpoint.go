package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrflick72/budget/budget-api/domain/budget/expense"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
)

func RegisterEndpoints(
	r *gin.Engine,
	ContextFactoryConverter server.ContextFactoryConverter,
	facade expense.BudgetExpenseActions,
) *gin.Engine {

	var logger = logging.GetLoggerInstance()
	/*
			@GetMapping
		    public ResponseEntity getBudgetExpenseList(@RequestParam("q") BudgetSearchCriteriaRepresentation budgetExpenseRequest)

		    @PutMapping
		    public ResponseEntity getBudgetExpenseListBy(@RequestBody BudgetSearchCriteriaRepresentation budgetExpenseRequest)

		    @PutMapping("/{id}")
		    public ResponseEntity updateBudgetExpense(@PathVariable("id") String id, @RequestBody BudgetExpenseRepresentation request)

		    @PostMapping
		    public ResponseEntity newBudgetExpense(@RequestBody BudgetExpenseRepresentation budgetExpenseRepresentation)

		    @DeleteMapping("/{id}")
		    public ResponseEntity deleteBudgetExpense(@PathVariable("id") String id)
	*/
	r.GET("/api/budget/expense", func(c *gin.Context) {})
	r.PUT("/api/budget/expense/:id", func(c *gin.Context) {
		var representation BudgetExpenseRepresentation

		ctx := ContextFactoryConverter.CreateContextFromGin(c)

		fmt.Println(ctx)
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
	r.PUT("/api/budget/expense", func(c *gin.Context) {})
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
	r.DELETE("/api/budget/expense/:id", func(c *gin.Context) {})

	return r
}
