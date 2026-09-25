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

// ScheduledExpenseStatusRepresentation is the PATCH /:id body. Status is
// validated against ACTIVE/PAUSED by the handler, not by a binding tag, so an
// unknown value and a missing one both get the same 400.
type ScheduledExpenseStatusRepresentation struct {
	Status string `json:"status"`
}
