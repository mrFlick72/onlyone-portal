package revenue

import (
	"github.com/mrflick72/budget/budget-api/domain/money"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
)

type Revenue struct {
	Id       RevenueId
	UserName UserName
	Date     date.Date
	Amount   money.Money
	Note     string
	Tags     []tags.SearchTag
}

type RevenueId = string
type UserName = string

type RevenueIdProvider interface {
	GenerateIdFor(revenue *Revenue) RevenueId
}
