package domain

import "github.com/shopspring/decimal"

var SCALE_PRECISION = 2
var ZERO, _ = MoneyFor("0.00")
var ONE, _ = MoneyFor("1.00")

type Money struct {
	content decimal.Decimal
}

func MoneyFor(amount string) (*Money, error) {
	return nil, nil
}

func (m *Money) Plus(money *Money) *Money {
	return nil
}

func (m *Money) StringifyAmount() string {
	return ""
}