package expense

import (
	"context"
	"errors"
	"time"

	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
)

// reclassificationTimeout bounds the async durable-reclassify work so a stuck
// Save or event publish cannot block the listener goroutine indefinitely. It is
// layered on top of the request context detached via context.WithoutCancel, so
// the work keeps the request's values (user identity, trace) but gains its own
// independent deadline.
const reclassificationTimeout = 30 * time.Second

type CreateBudgetExpense struct {
	Repository     BudgetExpenseRepository
	EventPublisher BudgetExpenseEventPublisher
	Logger         *logging.Logger
}

func (action *CreateBudgetExpense) Execute(ctx context.Context, budgetExpense *BudgetExpense) error {
	applyDefaultTagIfMissing(budgetExpense)
	user, _ := security.GetCurrentUser(ctx)
	budgetExpense.UserName = *user.UserName
	if err := action.Repository.Save(ctx, budgetExpense); err != nil {
		action.Logger.LogErrorfFor("create save failed for budget expense %s: %v", budgetExpense.Id, err)
		return err
	}
	if err := action.EventPublisher.CreateBudgetExpense(ctx, *budgetExpense); err != nil {
		action.Logger.LogErrorfFor("create event publish failed for budget expense %s: %v", budgetExpense.Id, err)
	}
	return nil
}

// applyDefaultTagIfMissing defaults an expense with no tags to the UNKNOWN sentinel tag,
// which tag-api always includes in its catalog without requiring per-user onboarding
// (see tagging/tag-api/docs/adr/0001-unknown-sentinel-tag-not-persisted.md). The
// sentinel is defined once in the tags package and shared with revenue.
func applyDefaultTagIfMissing(budgetExpense *BudgetExpense) {
	if len(budgetExpense.Tags) == 0 {
		budgetExpense.Tags = []tags.SearchTag{tags.UnknownSentinel()}
	}
}

type FindSpentBudget struct {
	BudgetExpenseRepository BudgetExpenseRepository
	SearchTagRepository     tags.SearchTagRepository
	Logger                  *logging.Logger
}

func (action *FindSpentBudget) Execute(ctx context.Context, month date.Month, year date.Year, searchTagKeys []tags.SearchTagKey) (*SpentBudget, error) {
	firstDate, err := date.FirstDateOfMonth(month, year)
	if err != nil {
		return nil, err
	}
	lastDate, err := date.LastDateOfMonth(month, year)
	if err != nil {
		return nil, err
	}

	budgetByDateRange, err := action.BudgetExpenseRepository.FindByDateRange(ctx, firstDate, lastDate, searchTagKeys)
	if err != nil {
		return nil, err
	}

	searchTags, err := action.getAllSearchTagFor(ctx, budgetByDateRange)
	return NewSpentBudget(budgetByDateRange, searchTags), err
}

func (action *FindSpentBudget) getAllSearchTagFor(ctx context.Context, budgetExpenses []BudgetExpense) ([]tags.SearchTag, error) {
	var searchTags []tags.SearchTag
	seen := make(map[string]bool)
	for _, expense := range budgetExpenses {
		for _, tag := range expense.Tags {
			if !seen[tag.Key] {
				seen[tag.Key] = true
				searchTag, err := action.SearchTagRepository.GetTagBy(ctx, tag.Key)
				if err != nil {
					return nil, err
				}
				searchTags = append(searchTags, *searchTag)
			}
		}

	}
	return searchTags, nil
}

type UpdateBudgetExpense struct {
	Repository     BudgetExpenseRepository
	EventPublisher BudgetExpenseEventPublisher
	EventBus       *InternalEventBus
	Logger         *logging.Logger
}

// Listen consumes reclassification events until the bus is closed, durably
// persisting each fix and re-emitting the UPDATE event. It is started as a
// single goroutine by the composition root and returns when the bus is closed
// at shutdown.
func (action *UpdateBudgetExpense) Listen() {
	for {
		select {
		case event := <-action.EventBus.Events():
			action.reclassify(event)
		case <-action.EventBus.Done():
			action.drainAndStop()
			return
		}
	}
}

// drainAndStop processes any events still buffered when the bus is closed, so a
// shutdown does not silently drop reclassifications already queued. The bus is
// closed only after the HTTP server has drained (see the composition root), so
// no new events are published during the drain and the buffer only shrinks; the
// non-blocking default returns once it is empty.
func (action *UpdateBudgetExpense) drainAndStop() {
	for {
		select {
		case event := <-action.EventBus.Events():
			action.reclassify(event)
		default:
			action.Logger.LogInfofFor("reclassification listener stopped")
			return
		}
	}
}

// reclassify durably persists a read-time UNKNOWN reclassification and, only on
// a successful write, re-emits the UPDATE event that reindexes the analytics
// projection — so analytics is never told a reclassification happened that did
// not land in DynamoDB. The originating request context is already cancelled by
// the time this runs asynchronously, so it is detached with
// context.WithoutCancel — the request's cancellation is dropped while its values
// (the user identity the repository's user_name condition needs, trace context)
// are preserved — and then given its own reclassificationTimeout so a stuck
// write or publish cannot block the listener goroutine forever.
func (action *UpdateBudgetExpense) reclassify(event InternalEvent) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(event.Ctx), reclassificationTimeout)
	defer cancel()
	if err := action.Repository.Save(ctx, event.Payload); err != nil {
		action.Logger.LogErrorfFor("reclassification save failed for budget expense %s: %v", event.Payload.Id, err)
		return
	}
	if err := action.EventPublisher.UpdateBudgetExpense(ctx, *event.Payload); err != nil {
		action.Logger.LogErrorfFor("reclassification update publish failed for budget expense %s: %v", event.Payload.Id, err)
	}
}

func (action *UpdateBudgetExpense) Execute(ctx context.Context, budgetExpense *BudgetExpense) error {
	applyDefaultTagIfMissing(budgetExpense)
	userName, err := security.GetCurrentUser(ctx)
	if err != nil {
		return err
	}

	existingBudgetExpense, err := action.Repository.FindFor(ctx, budgetExpense.Id)
	if err != nil {
		action.Logger.LogErrorfFor("update lookup failed for budget expense %s: %v", budgetExpense.Id, err)
		return err
	}
	if existingBudgetExpense != nil && existingBudgetExpense.UserName == *userName.UserName {
		if err := action.Repository.Save(ctx, budgetExpense); err != nil {
			action.Logger.LogErrorfFor("update save failed for budget expense %s: %v", budgetExpense.Id, err)
			return err
		}
		if err := action.EventPublisher.UpdateBudgetExpense(ctx, *budgetExpense); err != nil {
			action.Logger.LogErrorfFor("update event publish failed for budget expense %s: %v", budgetExpense.Id, err)
		}
		return nil
	}
	return errors.New("budget expense not found or user not authorized to update it")
}

type DeleteBudgetExpense struct {
	Repository     BudgetExpenseRepository
	EventPublisher BudgetExpenseEventPublisher
	Logger         *logging.Logger
}

func (action *DeleteBudgetExpense) Execute(ctx context.Context, id BudgetExpenseId) error {
	userName, err := security.GetCurrentUser(ctx)
	if err != nil {
		return err
	}

	existingBudgetExpense, err := action.Repository.FindFor(ctx, id)
	if err != nil {
		action.Logger.LogErrorfFor("delete lookup failed for budget expense %s: %v", id, err)
		return err
	}

	if existingBudgetExpense != nil && existingBudgetExpense.UserName == *userName.UserName {
		if err := action.Repository.Delete(ctx, id); err != nil {
			action.Logger.LogErrorfFor("delete failed for budget expense %s: %v", id, err)
			return err
		}
		if err := action.EventPublisher.DeleteBudgetExpense(ctx, *existingBudgetExpense); err != nil {
			action.Logger.LogErrorfFor("delete event publish failed for budget expense %s: %v", existingBudgetExpense.Id, err)
		}
		return nil
	}
	return errors.New("budget expense not found or user not authorized to delete it")
}
