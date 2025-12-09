package domain

type UserName = string

type User struct {
	UserName *UserName
}

type UserRepository interface {
	GetCurrentUser() (*User, error)
	SetCurrentUser(user *User) error
}
