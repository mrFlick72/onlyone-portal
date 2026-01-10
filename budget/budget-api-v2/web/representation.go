package web

type BudgetSearchCriteriaRepresentation struct {
	Month         string
	Year          int8
	SearchTagList []string
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
	Total                            string
	DailyBudgetExpenseRepresentation DailyBudgetExpenseRepresentation
	totalDetailList                  []TotalBySearchTagDetail
}

type DailyBudgetExpenseRepresentation struct {
	BudgetExpenseRepresentationList []BudgetExpenseRepresentation
	Date                            string
	Total                           string
}

type TotalBySearchTagDetail struct {
	SearchTagKey   string
	SearchTagValue string
	Total          string
}
