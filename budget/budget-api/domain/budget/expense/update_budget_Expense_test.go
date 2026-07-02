package expense

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/domain/money"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/stretchr/testify/mock"
)

func TestWhenABudgetExpenseUpdateSucceed(t *testing.T) {

	mockedRepository := new(BudgetExpenseRepositoryMock)
	mockedPublisher := new(BudgetExpenseEventPublisherMock)
	uut := UpdateBudgetExpense{
		Repository:     mockedRepository,
		EventPublisher: mockedPublisher,
	}

	aDate, _ := date.IsoDateFor("2018-01-01")
	anAmount, _ := money.MoneyFor("1.00")
	aBudgetExpense := BudgetExpense{
		Id:       "A_BUDGET_ID",
		UserName: "A_USER_NAME",
		Date:     *aDate,
		Amount:   anAmount,
		Note:     "A_NOTE",
		Tags:     []tags.SearchTag{{Key: "super-market", Value: "super-market"}},
	}

	ctx := testutils.NewUserContext()

	foundBudgetExpense := &BudgetExpense{
		Id:       "A_BUDGET_ID",
		UserName: "A_USER_NAME",
	}
	mockedRepository.On("FindFor", ctx, "A_BUDGET_ID").Return(foundBudgetExpense, nil)
	mockedRepository.On("Save", ctx, &aBudgetExpense).Return(nil)
	mockedPublisher.On("UpdateBudgetExpense", ctx, mock.Anything).Return(nil)

	err := uut.Execute(ctx, &aBudgetExpense)

	assert.Equal(t, nil, err)

	mockedRepository.AssertCalled(t, "FindFor", ctx, "A_BUDGET_ID")
	mockedRepository.AssertCalled(t, "Save", ctx, &aBudgetExpense)
	mockedPublisher.AssertCalled(t, "UpdateBudgetExpense", ctx, aBudgetExpense)
}
func TestListenConsumesAnInternalEventThenPersistsAndPublishes(t *testing.T) {

	mockedRepository := new(BudgetExpenseRepositoryMock)
	mockedPublisher := new(BudgetExpenseEventPublisherMock)

	bus := NewEventBus()
	defer bus.Close()
	uut := UpdateBudgetExpense{
		Repository:     mockedRepository,
		EventPublisher: mockedPublisher,
		EventBus:       bus,
		Logger:         logging.GetLoggerInstanceForComponentByTypeName("UpdateBudgetExpenseTest"),
	}

	aBudgetExpense := &BudgetExpense{
		Id:       "A_BUDGET_ID",
		UserName: "A_USER_NAME",
		Tags:     []tags.SearchTag{tags.UnknownSentinel()},
	}

	// The listener detaches the request context via context.WithoutCancel, so the
	// context reaching Save/UpdateBudgetExpense is a derived value, not the one
	// published — hence mock.Anything on the context argument.
	done := make(chan struct{})
	mockedRepository.On("Save", mock.Anything, aBudgetExpense).Return(nil)
	mockedPublisher.On("UpdateBudgetExpense", mock.Anything, *aBudgetExpense).Return(nil).
		Run(func(args mock.Arguments) { close(done) })

	go uut.Listen()
	bus.Publish(InternalEvent{Payload: aBudgetExpense, Ctx: testutils.NewUserContext()})

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for Listen to consume the internal event")
	}

	mockedRepository.AssertCalled(t, "Save", mock.Anything, aBudgetExpense)
	mockedPublisher.AssertCalled(t, "UpdateBudgetExpense", mock.Anything, *aBudgetExpense)
}

// On a persistence failure the listener logs and does NOT re-emit the UPDATE
// event, so the analytics projection is never told a reclassification happened
// that did not durably land in DynamoDB.
func TestListenDoesNotPublishWhenPersistenceFails(t *testing.T) {

	mockedRepository := new(BudgetExpenseRepositoryMock)
	mockedPublisher := new(BudgetExpenseEventPublisherMock)

	bus := NewEventBus()
	defer bus.Close()
	uut := UpdateBudgetExpense{
		Repository:     mockedRepository,
		EventPublisher: mockedPublisher,
		EventBus:       bus,
		Logger:         logging.GetLoggerInstanceForComponentByTypeName("UpdateBudgetExpenseTest"),
	}

	aBudgetExpense := &BudgetExpense{
		Id:       "A_BUDGET_ID",
		UserName: "A_USER_NAME",
		Tags:     []tags.SearchTag{tags.UnknownSentinel()},
	}

	saved := make(chan struct{})
	mockedRepository.On("Save", mock.Anything, aBudgetExpense).Return(fmt.Errorf("save failed")).
		Run(func(args mock.Arguments) { close(saved) })

	go uut.Listen()
	bus.Publish(InternalEvent{Payload: aBudgetExpense, Ctx: testutils.NewUserContext()})

	select {
	case <-saved:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for the listener to attempt the save")
	}

	// Give the listener a window to (incorrectly) publish if the guard regressed.
	time.Sleep(100 * time.Millisecond)
	mockedRepository.AssertCalled(t, "Save", mock.Anything, aBudgetExpense)
	mockedPublisher.AssertNotCalled(t, "UpdateBudgetExpense", mock.Anything, mock.Anything)
}

// On shutdown the listener drains events already buffered before returning, so
// queued reclassifications are not silently dropped (review finding #3).
func TestListenDrainsBufferedEventsOnShutdown(t *testing.T) {

	mockedRepository := new(BudgetExpenseRepositoryMock)
	mockedPublisher := new(BudgetExpenseEventPublisherMock)

	bus := NewEventBus()
	uut := UpdateBudgetExpense{
		Repository:     mockedRepository,
		EventPublisher: mockedPublisher,
		EventBus:       bus,
		Logger:         testLogger,
	}

	mockedRepository.On("Save", mock.Anything, mock.Anything).Return(nil)
	mockedPublisher.On("UpdateBudgetExpense", mock.Anything, mock.Anything).Return(nil)

	// Queue several events and stop the bus before the listener runs; every
	// buffered event must still be processed before Listen returns.
	for i := 0; i < 3; i++ {
		bus.Publish(InternalEvent{Payload: &BudgetExpense{Id: "A_BUDGET_ID"}, Ctx: testutils.NewUserContext()})
	}
	bus.Close()

	done := make(chan struct{})
	go func() {
		uut.Listen()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Listen did not drain buffered events and stop")
	}

	mockedRepository.AssertNumberOfCalls(t, "Save", 3)
	mockedPublisher.AssertNumberOfCalls(t, "UpdateBudgetExpense", 3)
}

func TestWhenABudgetExpenseUpdateWithNoTagsDefaultsToUnknown(t *testing.T) {

	mockedRepository := new(BudgetExpenseRepositoryMock)
	mockedPublisher := new(BudgetExpenseEventPublisherMock)
	uut := UpdateBudgetExpense{
		Repository:     mockedRepository,
		EventPublisher: mockedPublisher,
	}

	aDate, _ := date.IsoDateFor("2018-01-01")
	anAmount, _ := money.MoneyFor("1.00")
	aBudgetExpense := BudgetExpense{
		Id:       "A_BUDGET_ID",
		UserName: "A_USER_NAME",
		Date:     *aDate,
		Amount:   anAmount,
		Note:     "A_NOTE",
		Tags:     []tags.SearchTag{},
	}

	ctx := testutils.NewUserContext()

	foundBudgetExpense := &BudgetExpense{
		Id:       "A_BUDGET_ID",
		UserName: "A_USER_NAME",
	}
	mockedRepository.On("FindFor", ctx, "A_BUDGET_ID").Return(foundBudgetExpense, nil)
	mockedRepository.On("Save", ctx, &aBudgetExpense).Return(nil)
	mockedPublisher.On("UpdateBudgetExpense", ctx, mock.Anything).Return(nil)

	err := uut.Execute(ctx, &aBudgetExpense)

	assert.Equal(t, nil, err)
	assert.Equal(t, []tags.SearchTag{{Key: "UNKNOWN", Value: "UNKNOWN"}}, aBudgetExpense.Tags)
}

func TestWhenABudgetExpenseUpdateDoesDoneNothingBecauseTheBudgetExpenseDoesNotExist(t *testing.T) {

	mockedRepository := new(BudgetExpenseRepositoryMock)
	uut := UpdateBudgetExpense{
		Repository: mockedRepository,
		Logger:     testLogger,
	}

	aDate, _ := date.IsoDateFor("2018-01-01")
	anAmount, _ := money.MoneyFor("1.00")
	aBudgetExpense := BudgetExpense{
		Id:       "A_BUDGET_ID",
		UserName: "A_USER_NAME",
		Date:     *aDate,
		Amount:   anAmount,
		Note:     "A_NOTE",
		Tags:     []tags.SearchTag{{Key: "super-market", Value: "super-market"}},
	}

	ctx := testutils.NewUserContext()

	mockedRepository.On("FindFor", ctx, "A_BUDGET_ID").Return(nil, fmt.Errorf("budget expense with the id %s was not found", "A_BUDGET_ID"))

	err := uut.Execute(ctx, &aBudgetExpense)

	assert.NotEqual(t, nil, err)

	mockedRepository.AssertCalled(t, "FindFor", ctx, "A_BUDGET_ID")
	mockedRepository.AssertNotCalled(t, "Save", ctx, &aBudgetExpense)

}
func TestWhenABudgetExpenseUpdateFails(t *testing.T) {

	mockedRepository := new(BudgetExpenseRepositoryMock)
	uut := UpdateBudgetExpense{
		Repository: mockedRepository,
		Logger:     testLogger,
	}

	aDate, _ := date.IsoDateFor("2018-01-01")
	anAmount, _ := money.MoneyFor("1.00")
	aBudgetExpense := BudgetExpense{
		Id:       "A_BUDGET_ID",
		UserName: "A_USER_NAME",
		Date:     *aDate,
		Amount:   anAmount,
		Note:     "A_NOTE",
		Tags:     []tags.SearchTag{{Key: "super-market", Value: "super-market"}},
	}

	ctx := testutils.NewUserContext()

	foundBudgetExpense := &BudgetExpense{
		Id:       "A_BUDGET_ID",
		UserName: "A_USER_NAME",
	}
	mockedRepository.On("FindFor", ctx, "A_BUDGET_ID").Return(foundBudgetExpense, nil)
	mockedRepository.On("Save", ctx, &aBudgetExpense).Return(fmt.Errorf("budget expense save with the id %s failed", "A_BUDGET_ID"))

	err := uut.Execute(ctx, &aBudgetExpense)

	assert.NotEqual(t, nil, err)

	mockedRepository.AssertCalled(t, "FindFor", ctx, "A_BUDGET_ID")
	mockedRepository.AssertCalled(t, "Save", ctx, &aBudgetExpense)
}

func TestWhenABudgetExpenseUpdateFailsBecauseUserOwnership(t *testing.T) {
	mockedRepository := new(BudgetExpenseRepositoryMock)
	uut := UpdateBudgetExpense{
		Repository: mockedRepository,
		Logger:     testLogger,
	}

	ctx := testutils.NewUserContext()
	foundBudgetExpense := &BudgetExpense{
		Id:       "A_BUDGET_ID",
		UserName: "A_DIFFERENT_USER_NAME",
	}
	aBudgetExpense := &BudgetExpense{
		Id:       "A_BUDGET_ID",
		UserName: "A_USER_NAME",
	}
	mockedRepository.On("FindFor", ctx, "A_BUDGET_ID").Return(foundBudgetExpense, nil)

	err := uut.Execute(ctx, aBudgetExpense)

	assert.NotEqual(t, nil, err)

	mockedRepository.AssertCalled(t, "FindFor", ctx, "A_BUDGET_ID")
	mockedRepository.AssertNotCalled(t, "Save", ctx, aBudgetExpense)
}
