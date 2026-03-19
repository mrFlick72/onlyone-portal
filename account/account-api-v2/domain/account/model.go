package account

type Account struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	BirthDate string `json:"birthDate"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}
