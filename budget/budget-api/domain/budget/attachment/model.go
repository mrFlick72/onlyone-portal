package attachment

import "github.com/mrflick72/budget/budget-api/domain/time/date"

//todo rename BudgetType into resourceType and BudgetId into resourceId

type AttachmentMetadata struct {
	AttachmentId string
	BudgetId     string // it can be the budget resource in which such attachment is attached to. It is something like, expense or revenue
	BudgetType   string
	Date         date.Date
	Owner        string
	FineName     string
	ContentType  string
	Metadata     map[string]string
}
type Attachment struct {
	AttachmentMetadata
	Content []byte
}
type BudgetType string

const (
	Expense BudgetType = "expense"
	Revenue BudgetType = "revenue"
)
