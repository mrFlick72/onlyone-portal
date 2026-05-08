package attachment

type AttachmentMetadataRepresentation struct {
	FileName     string `json:"fileName"`
	Owner        string `json:"owner"`
	BudgetId     string `json:"type"`
	BudgetType   string `json:"budgetId"`
	AttachmentId string `json:"attachmentId"`
}
