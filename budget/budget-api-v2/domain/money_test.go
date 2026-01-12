package domain

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

// @Test
// public void manoeyScaleIsCorrect() throws Exception {
//     Money money = Money.moneyFor("12.506");
//     BigDecimal expectedValue = new BigDecimal(12.51).setScale(2, RoundingMode.HALF_DOWN);
//     BigDecimal actualValue = money.amount();
//     Assertions.assertEquals(actualValue, expectedValue);
// }

// @Test
// public void addOperation() throws Exception {
//     Money firstAddendum = Money.moneyFor("12.50");
//     Money secondAddendum = Money.moneyFor("10.22");

//     Money expectedValue = Money.moneyFor("22.72");

//     Assertions.assertEquals(firstAddendum.plus(secondAddendum), expectedValue);

// }

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
