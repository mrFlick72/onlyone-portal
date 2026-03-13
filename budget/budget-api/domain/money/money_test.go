package money

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestMoneyScaleIsCorrect(t *testing.T) {
	money, _ := MoneyFor("12.506")

	assert.Equal(t, "12.51", money.StringifyAmount())
}

func TestAddOperation(t *testing.T) {
	firstAddendum, _ := MoneyFor("12.50")
	secondAddendum, _ := MoneyFor("10.22")
	expectedValue := "22.72"

	assert.Equal(t, expectedValue, firstAddendum.Plus(secondAddendum).StringifyAmount())
}
