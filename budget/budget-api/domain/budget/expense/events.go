package expense

import (
	"context"
	"sync"
)

/*
InternalEventBus

	This is intended to make possible async event communication at intra class level.
	ex: notify to teh UpdateExpenseAction that at repository layer we noticed that one search tags is been removed to it has to be updated at storage level and
	notified to teh anaytics api service
*/
type InternalEvent struct {
	Payload *BudgetExpense
	Ctx     context.Context
}
type InternalEventBus = chan InternalEvent

var (
	once    sync.Once
	channel InternalEventBus
)

func NewEventBus() InternalEventBus {
	once.Do(func() {
		channel = make(chan InternalEvent)
	})

	return channel
}

type BudgetExpenseEventPublisher interface {
	CreateBudgetExpense(ctx context.Context, expense BudgetExpense) error
	UpdateBudgetExpense(ctx context.Context, expense BudgetExpense) error
	DeleteBudgetExpense(ctx context.Context, expense BudgetExpense) error
}
