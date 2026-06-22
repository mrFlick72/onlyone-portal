package revenue

import tagRep "github.com/mrflick72/budget/budget-api/web/tags"

type RevenueRepresentation struct {
	Id     string                           `json:"id"`
	Date   string                           `json:"date"`
	Amount string                           `json:"amount"`
	Note   string                           `json:"note"`
	Tags   []tagRep.SearchTagRepresentation `json:"tags"`
}

type RevenueListRepresentation struct {
	Revenues []RevenueRepresentation `json:"revenues"`
	Total    string                  `json:"total"`
}
