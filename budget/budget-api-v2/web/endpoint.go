package web

import "github.com/gin-gonic/gin"

func RegisterEndpoints(r *gin.Engine) *gin.Engine {

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
	r.POST("/api/budget/expense", func(c *gin.Context) {})
	r.DELETE("/api/budget/expense/:id", func(c *gin.Context) {})

	return r
}
