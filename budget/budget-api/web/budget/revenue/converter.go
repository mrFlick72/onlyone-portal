package revenue

import (
	"errors"
	"strings"

	domainrevenue "github.com/mrflick72/budget/budget-api/domain/budget/revenue"
	"github.com/mrflick72/budget/budget-api/domain/money"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	tagRep "github.com/mrflick72/budget/budget-api/web/tags"
)

func RevenueRepresentationToDomainModel(rep RevenueRepresentation) (*domainrevenue.Revenue, error) {
	d, err := date.DateFor(rep.Date)
	if err != nil {
		return nil, err
	}
	m, err := money.MoneyFor(rep.Amount)
	if err != nil {
		return nil, err
	}

	searchTags := make([]tags.SearchTag, 0, len(rep.Tags))
	for _, tag := range rep.Tags {
		searchTags = append(searchTags, tags.SearchTag{Key: tag.Key, Value: tag.Value})
	}

	return &domainrevenue.Revenue{
		Id:     domainrevenue.RevenueId(rep.Id),
		Date:   *d,
		Amount: m,
		Note:   rep.Note,
		Tags:   searchTags,
	}, nil
}

func RevenueDomainToRepresentationModel(r *domainrevenue.Revenue) RevenueRepresentation {
	searchTags := make([]tagRep.SearchTagRepresentation, 0, len(r.Tags))
	for _, tag := range r.Tags {
		searchTags = append(searchTags, tagRep.SearchTagRepresentation{Key: tag.Key, Value: tag.Value})
	}

	return RevenueRepresentation{
		Id:     string(r.Id),
		Date:   r.Date.GetFormattedDate(),
		Amount: r.Amount.StringifyAmount(),
		Note:   r.Note,
		Tags:   searchTags,
	}
}

func RevenueListDomainToRepresentationModel(list []domainrevenue.Revenue) RevenueListRepresentation {
	revenues := make([]RevenueRepresentation, 0, len(list))
	total := money.Zero()
	for i := range list {
		revenues = append(revenues, RevenueDomainToRepresentationModel(&list[i]))
		total = total.Plus(list[i].Amount)
	}
	return RevenueListRepresentation{
		Revenues: revenues,
		Total:    total.StringifyAmount(),
	}
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
