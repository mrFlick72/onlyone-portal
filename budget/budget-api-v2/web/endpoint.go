package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrflick72/budget/budget-api/domain/budget/expense"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
)

func RegisterEndpoints(
	r *gin.Engine,
	facade expense.BudgetExpenseActions,
) *gin.Engine {

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
	r.PUT("/api/budget/expense/:id", func(c *gin.Context) {})
	r.PUT("/api/budget/expense", func(c *gin.Context) {})
	r.POST("/api/budget/expense", func(c *gin.Context) {
		var representation BudgetExpenseRepresentation

		ctx := server.CopyGinKeysToRequestContext(c)

		if err := c.ShouldBindJSON(&representation); err != nil {
			fmt.Printf("Error binding JSON: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Printf("Received representation: %+v\n", representation)
		domainModel := RepresentationModelToDomainModel(representation)
		fmt.Printf("Received domain model: %+v\n", domainModel)
		facade.CreateBudgetExpense(ctx, domainModel)
		c.Status(http.StatusCreated)
	})
	r.DELETE("/api/budget/expense/:id", func(c *gin.Context) {})

	return r
}
