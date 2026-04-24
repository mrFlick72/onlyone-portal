package revenue

type RevenueRepresentation struct {
	Id     string `json:"id"`
	Date   string `json:"date"`
	Amount string `json:"amount"`
	Note   string `json:"note"`
}

type RevenueListRepresentation struct {
	Revenues []RevenueRepresentation `json:"revenues"`
	Total    string                  `json:"total"`
}
