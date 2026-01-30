package money

import "github.com/shopspring/decimal"

var SCALE_PRECISION = 2
var ZERO, _ = MoneyFor("0.00")
var ONE, _ = MoneyFor("1.00")

type Money struct {
	content decimal.Decimal
}

func MoneyFor(amount string) (*Money, error) {
	var result, err = decimal.NewFromString(amount)
	if err != nil {
		return nil, err
	}
	result = result.Round(2)
	return &Money{
		content: result,
	}, nil
}

func (m *Money) Plus(money *Money) *Money {
	return &Money{content: m.content.Add(money.content)}
}

func (m *Money) StringifyAmount() string {
	return m.content.StringFixed(2)
}
