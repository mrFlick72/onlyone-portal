package dynamodb

import (
	"fmt"
	"strings"

	"github.com/mrflick72/budget/budget-api/domain/budget/attachment"
)

type DynamoDbAttachmentIdProvider struct {
	UuidGenerator func() string
}

func (provider *DynamoDbAttachmentIdProvider) GenerateAttachmentIdFor(att *attachment.Attachment) string {
	if att.AttachmentId != "" {
		return att.AttachmentId
	}
	return provider.UuidGenerator()
}

func (provider *DynamoDbAttachmentIdProvider) PartitionKeyFor(att *attachment.Attachment) string {
	return fmt.Sprintf("%s_%s", att.BudgetId, strings.ToUpper(att.BudgetType))
}
