package attachment

type AttachmentMetadataRepresentation struct {
	FileName     string `json:"fileName"`
	Owner        string `json:"owner"`
	BudgetId     string `json:"budgetId"`
	BudgetType   string `json:"budgetType"`
	AttachmentId string `json:"attachmentId"`
}
