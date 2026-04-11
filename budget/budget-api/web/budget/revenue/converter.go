package revenue

import (
	"errors"
	"strings"

	domainrevenue "github.com/mrflick72/budget/budget-api/domain/budget/revenue"
	"github.com/mrflick72/budget/budget-api/domain/money"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
)

func RevenueRepresentationToDomainModel(rep RevenueRepresentation) *domainrevenue.Revenue {
	d, _ := date.DateFor(rep.Date)
	m, _ := money.MoneyFor(rep.Amount)
	return &domainrevenue.Revenue{
		Id:     domainrevenue.RevenueId(rep.Id),
		Date:   *d,
		Amount: m,
		Note:   rep.Note,
		
	}
}

func RevenueDomainToRepresentationModel(r *domainrevenue.Revenue) RevenueRepresentation {
	return RevenueRepresentation{
		Id:     string(r.Id),
		Date:   r.Date.GetFormattedDate(),
		Amount: r.Amount.StringifyAmount(),
		Note:   r.Note,
	}
}

func RevenueListDomainToRepresentationModel(list []domainrevenue.Revenue) []RevenueRepresentation {
	result := make([]RevenueRepresentation, 0, len(list))
	for i := range list {
		result = append(result, RevenueDomainToRepresentationModel(&list[i]))
	}
	return result
}

// parseYearQueryParam preserves the Python revenue-api wire format "year=2023"
// (optionally followed by further ";key=value" segments — only the first is read).
func parseYearQueryParam(q string) (date.Year, error) {
	if q == "" {
		return date.Year{}, errors.New("missing query parameter")
	}
	first := strings.Split(q, ";")[0]
	pair := strings.SplitN(first, "=", 2)
	if len(pair) != 2 || pair[0] != "year" || pair[1] == "" {
		return date.Year{}, errors.New("invalid query parameter, expected year=YYYY")
	}
	return date.NewYearFor(pair[1]), nil
}
