package dynamodb

import "github.com/mrflick72/budget/budget-api/domain/budget/scheduledexpense"

type DynamoDbScheduledExpenseIdProvider struct {
	UuidGenerator func() string
}

func (provider *DynamoDbScheduledExpenseIdProvider) GenerateIdFor(se *scheduledexpense.ScheduledExpense) scheduledexpense.ScheduledExpenseId {
	if se.Id != "" {
		return se.Id
	}
	return provider.UuidGenerator()
}
