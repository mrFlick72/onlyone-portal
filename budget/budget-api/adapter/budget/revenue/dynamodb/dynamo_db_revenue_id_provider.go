package dynamodb

import (
	"encoding/base64"
	"fmt"

	"github.com/mrflick72/budget/budget-api/domain/budget/revenue"
)

type DynamoDbRevenueIdProvider struct {
	SaltGenerator func() string
}

func (provider *DynamoDbRevenueIdProvider) GenerateIdFor(revenue *revenue.Revenue) revenue.RevenueId {
	return fmt.Sprintf("%s-%s",
		provider.partitionKeyFrom(revenue.Date.GetYear(), revenue.UserName),
		provider.rangeKeyFrom(revenue.Date.GetMonth(), revenue.Date.GetDay()))
}

func (provider *DynamoDbRevenueIdProvider) partitionKeyFrom(year int, userName string) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d_%s", year, userName)))
}

func (provider *DynamoDbRevenueIdProvider) rangeKeyFrom(month, day int) string {
	rangeKey := fmt.Sprintf("%d_%d_%s", month, day, provider.SaltGenerator())
	return base64.StdEncoding.EncodeToString([]byte(rangeKey))
}
