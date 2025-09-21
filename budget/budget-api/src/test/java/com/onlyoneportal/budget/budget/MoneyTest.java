package com.onlyoneportal.budget.budget;


import com.onlyoneportal.budget.Money;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;

import java.math.BigDecimal;
import java.math.RoundingMode;

public class MoneyTest {

    @Test
    public void manoeyScaleIsCorrect() throws Exception {
        Money money = Money.moneyFor("12.506");
        BigDecimal expectedValue = new BigDecimal(12.51).setScale(2, RoundingMode.HALF_DOWN);
        BigDecimal actualValue = money.amount();
        Assertions.assertEquals(actualValue, expectedValue);
    }

    @Test
    public void addOperation() throws Exception {
        Money firstAddendum = Money.moneyFor("12.50");
        Money secondAddendum = Money.moneyFor("10.22");

        Money expectedValue = Money.moneyFor("22.72");

        Assertions.assertEquals(firstAddendum.plus(secondAddendum), expectedValue);

    }
}