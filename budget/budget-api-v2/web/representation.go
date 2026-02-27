package web

import "github.com/mrflick72/budget/budget-api/domain/tags"

type BudgetSearchCriteriaRepresentation struct {
	Month         string `json:"month"`
	Year          string `json:"year"`
	SearchTagList []tags.SearchTagKey `json:"searchTagList"`
}

type BudgetExpenseRepresentation struct {
	Id       string `json:"id"`
	Date     string `json:"date"`
	Amount   string `json:"amount"`
	Note     string `json:"note"`
	TagKey   string `json:"tagKey"`
	TagValue string `json:"tagValue"`
}

type SpentBudgetRepresentation struct {
	Total                            string `json:"total"`
	DailyBudgetExpenseRepresentation []DailyBudgetExpenseRepresentation `json:"dailyBudgetExpenseRepresentationList"`
	TotalDetailList                  []TotalBySearchTagDetail `json:"totalDetailList"`
}

type DailyBudgetExpenseRepresentation struct {
	BudgetExpenseRepresentationList []BudgetExpenseRepresentation `json:"budgetExpenseRepresentationList"`
	Date                            string `json:"date"`
	Total                           string `json:"total"`
}

type TotalBySearchTagDetail struct {
	SearchTagKey   string `json:"searchTagKey"`
	SearchTagValue string `json:"searchTagValue"`
	Total          string `json:"total"`
}
