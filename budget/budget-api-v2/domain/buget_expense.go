package domain

type SpentBudget struct {
	BudgetExpenseList BudgetExpense
	SearchTags        map[string]string
}

type BudgetExpense struct {
	Id       BudgetExpenseId
	UserName UserName
	Date     Date
	Amount   Money
	Note     string
	Tag      string
}

type BudgetExpenseId = string
type UserName = string
