package revenue

import (
	"context"
	"errors"
	"sort"

	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
)

// applyDefaultTagIfMissing defaults a revenue with no tags to the shared UNKNOWN
// sentinel — the exact analogue of the expense behavior — so every revenue has
// at least one tag once saved.
func applyDefaultTagIfMissing(revenue *Revenue) {
	if len(revenue.Tags) == 0 {
		revenue.Tags = []tags.SearchTag{tags.UnknownSentinel()}
	}
}

type CreateRevenue struct {
	Repository RevenueRepository
}

func (action *CreateRevenue) Execute(ctx context.Context, revenue *Revenue) error {
	applyDefaultTagIfMissing(revenue)
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		return err
	}
	revenue.UserName = *user.UserName
	return action.Repository.Save(ctx, revenue)
}

type UpdateRevenue struct {
	Repository RevenueRepository
}

func (action *UpdateRevenue) Execute(ctx context.Context, revenue *Revenue) error {
	applyDefaultTagIfMissing(revenue)
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		return err
	}

	existingRevenue, err := action.Repository.FindFor(ctx, revenue.Id)
	if err != nil {
		return err
	}
	if existingRevenue != nil && existingRevenue.UserName == *user.UserName {
		return action.Repository.Save(ctx, revenue)
	}
	return errors.New("revenue not found or user not authorized to update it")
}

type FindRevenue struct {
	Repository RevenueRepository
}

func (action *FindRevenue) Execute(ctx context.Context, year date.Year) ([]Revenue, error) {
	firstDate, err := date.FirstDateOfMonth(date.JANUARY(), year)
	if err != nil {
		return nil, err
	}
	lastDate, err := date.LastDateOfMonth(date.DECEMBER(), year)
	if err != nil {
		return nil, err
	}
	revenues, err := action.Repository.FindByDateRange(ctx, firstDate, lastDate)
	if err != nil {
		return nil, err
	}
	sort.Slice(revenues, func(i, j int) bool {
		return revenues[i].Date.GetTime().Before(revenues[j].Date.GetTime())
	})
	return revenues, nil
}

type DeleteRevenue struct {
	Repository RevenueRepository
}

func (action *DeleteRevenue) Execute(ctx context.Context, id RevenueId) error {
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		return err
	}

	existingRevenue, err := action.Repository.FindFor(ctx, id)
	if err != nil {
		return err
	}

	if existingRevenue != nil && existingRevenue.UserName == *user.UserName {
		return action.Repository.Delete(ctx, id)
	}
	return errors.New("revenue not found or user not authorized to delete it")
}
