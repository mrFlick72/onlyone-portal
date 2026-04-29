package plan

type todoRepresentation struct {
	Id       string `json:"id"`
	UserName string `json:"user_name"`
	Date     string `json:"date"`
	Content  string `json:"content"`
}

type planRepresentation struct {
	Id       string               `json:"id"`
	UserName string               `json:"user_name"`
	Title    string               `json:"title"`
	Date     string               `json:"date"`
	Todos    []todoRepresentation `json:"todos"`
}
