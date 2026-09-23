package scheduledexpense

import (
	"errors"

	domainscheduledexpense "github.com/mrflick72/budget/budget-api/domain/budget/scheduledexpense"
	"github.com/mrflick72/budget/budget-api/domain/money"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	tagRep "github.com/mrflick72/budget/budget-api/web/tags"
)

func ScheduledExpenseRepresentationToDomainModel(rep ScheduledExpenseRepresentation) (*domainscheduledexpense.ScheduledExpense, error) {
	if rep.Description == "" {
		return nil, errors.New("description is required")
	}
	if rep.Day < 1 || rep.Day > 31 {
		return nil, errors.New("day must be between 1 and 31")
	}
	if rep.Month != nil && (*rep.Month < 1 || *rep.Month > 12) {
		return nil, errors.New("month must be between 1 and 12")
	}

	m, err := money.MoneyFor(rep.Amount)
	if err != nil {
		return nil, err
	}

	searchTags := make([]tags.SearchTag, 0, len(rep.Tags))
	for _, tag := range rep.Tags {
		searchTags = append(searchTags, tags.SearchTag{Key: tag.Key, Value: tag.Value})
	}

	domainModel := &domainscheduledexpense.ScheduledExpense{
		Id:          rep.Id,
		Description: rep.Description,
		Amount:      m,
		Notes:       rep.Notes,
		Tags:        searchTags,
		Day:         rep.Day,
	}

	if rep.Month != nil {
		month := *rep.Month
		domainModel.Month = &month
	}

	if rep.EndDate != nil {
		d, err := date.DateFor(*rep.EndDate)
		if err != nil {
			return nil, err
		}
		domainModel.EndDate = d
	}

	return domainModel, nil
}

func ScheduledExpenseDomainToRepresentationModel(se *domainscheduledexpense.ScheduledExpense) ScheduledExpenseRepresentation {
	searchTags := make([]tagRep.SearchTagRepresentation, 0, len(se.Tags))
	for _, tag := range se.Tags {
		searchTags = append(searchTags, tagRep.SearchTagRepresentation{Key: tag.Key, Value: tag.Value})
	}

	rep := ScheduledExpenseRepresentation{
		Id:          string(se.Id),
		Description: se.Description,
		Amount:      se.Amount.StringifyAmount(),
		Notes:       se.Notes,
		Tags:        searchTags,
		Day:         se.Day,
		Status:      string(se.Status),
	}

	if se.Month != nil {
		month := *se.Month
		rep.Month = &month
	}

	if se.EndDate != nil {
		formatted := se.EndDate.GetFormattedDate()
		rep.EndDate = &formatted
	}

	return rep
}

func ScheduledExpenseListDomainToRepresentationModel(list []domainscheduledexpense.ScheduledExpense) ScheduledExpenseListRepresentation {
	scheduledExpenses := make([]ScheduledExpenseRepresentation, 0, len(list))
	for i := range list {
		scheduledExpenses = append(scheduledExpenses, ScheduledExpenseDomainToRepresentationModel(&list[i]))
	}
	return ScheduledExpenseListRepresentation{ScheduledExpenses: scheduledExpenses}
}
