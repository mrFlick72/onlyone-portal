package testutils

import (
	"context"
	"fmt"

	"github.com/mrflick72/budget/budget-api/domain/money"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
)

func NewUserContext() context.Context {
	UserName := "A_USER_NAME"
	user := security.User{UserName: &UserName, Authorities: nil}
	return context.WithValue(context.TODO(), "user", user)
}

func SafeDateFor(dateStr string) date.Date {
	d, err := date.DateFor(dateStr)
	if err != nil {
		fmt.Printf("Error during date parsing data input is: %s", dateStr)
		panic(err)
	}
	return *d
}

func SafeMoneyFor(moneyStr string) money.Money {
	m, err := money.MoneyFor(moneyStr)
	if err != nil {
		fmt.Printf("Error during date parsing data input is: %s", moneyStr)
		panic(err)
	}
	return *m
}
