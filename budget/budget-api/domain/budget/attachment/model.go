package attachment

type AttachmentMetadata struct {
	AttachmentId string
	BudgetId     string // it can be the budget resource in which such attachment is attached to. It is something like, expense or revenue
	Owner        string
	FineName     string
	ContentType  string
	Metadata     map[string]string
}
type Attachment struct {
	AttachmentMetadata
	Content []byte
}
