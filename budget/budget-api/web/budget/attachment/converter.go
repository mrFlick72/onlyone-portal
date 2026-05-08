package attachment

import (
	"github.com/mrflick72/budget/budget-api/domain/budget/attachment"
)

func DomainModelToRepresentationModel(model []attachment.AttachmentMetadata) []AttachmentMetadataRepresentation {
	result := make([]AttachmentMetadataRepresentation, len(model))

	for _, item := range model {
		result = append(result, AttachmentMetadataRepresentation{
			FileName:     item.FineName,
			Owner:        item.Owner,
			BudgetId:     item.BudgetId,
			BudgetType:   item.BudgetType,
			AttachmentId: item.AttachmentId,
		})
	}
	return result
}
