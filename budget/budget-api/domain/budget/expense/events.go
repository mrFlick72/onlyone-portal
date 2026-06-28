package expense

import (
	"context"
	"sync"
)

type EventBus = chan string

var (
	once    sync.Once
	channel EventBus
)

func NewEventBus() EventBus {
	once.Do(func() {
		channel = make(chan string)
	})

	return channel
}

type BudgetExpenseEventPublisher interface {
	CreateBudgetExpense(ctx context.Context, expense BudgetExpense) error
	UpdateBudgetExpense(ctx context.Context, expense BudgetExpense) error
	DeleteBudgetExpense(ctx context.Context, expense BudgetExpense) error
}
