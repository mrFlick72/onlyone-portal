package attachment

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/domain/budget/attachment"
)

func TestDomainModelToRepresentationModelMapsFields(t *testing.T) {
	input := []attachment.AttachmentMetadata{
		{
			AttachmentId: "a-1",
			BudgetId:     "b-1",
			BudgetType:   "expense",
			Owner:        "alice",
			FineName:     "alice.pdf",
		},
		{
			AttachmentId: "a-2",
			BudgetId:     "b-2",
			BudgetType:   "revenue",
			Owner:        "bob",
			FineName:     "bob.png",
		},
	}

	got := DomainModelToRepresentationModel(input)

	assert.Equal(t, 2, len(got))
	assert.Equal(t, AttachmentMetadataRepresentation{
		FileName:     "alice.pdf",
		Owner:        "alice",
		BudgetId:     "b-1",
		BudgetType:   "expense",
		AttachmentId: "a-1",
	}, got[0])
	assert.Equal(t, AttachmentMetadataRepresentation{
		FileName:     "bob.png",
		Owner:        "bob",
		BudgetId:     "b-2",
		BudgetType:   "revenue",
		AttachmentId: "a-2",
	}, got[1])
}

func TestDomainModelToRepresentationModelWithEmptyInputReturnsEmptySlice(t *testing.T) {
	got := DomainModelToRepresentationModel(nil)

	assert.Equal(t, 0, len(got))
}
