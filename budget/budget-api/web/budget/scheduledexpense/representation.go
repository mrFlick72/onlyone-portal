package scheduledexpense

import tagRep "github.com/mrflick72/budget/budget-api/web/tags"

type ScheduledExpenseRepresentation struct {
	Id          string                           `json:"id"`
	Description string                           `json:"description"`
	Amount      string                           `json:"amount"`
	Notes       string                           `json:"notes"`
	Tags        []tagRep.SearchTagRepresentation `json:"tags"`
	Day         int                              `json:"day"`
	Month       *int                             `json:"month,omitempty"`
	EndDate     *string                          `json:"endDate,omitempty"`
	Status      string                           `json:"status"`
}

type ScheduledExpenseListRepresentation struct {
	ScheduledExpenses []ScheduledExpenseRepresentation `json:"scheduledExpenses"`
}
