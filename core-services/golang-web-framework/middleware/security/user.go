package security

import (
	"context"
	"errors"
)

type UserName = string


type User struct {
	UserName *UserName
	Authorities *[]string
	AccessToken *string
}

func GetCurrentUser(ctx *context.Context) (*User, error) {
	if v := (*ctx).Value("user"); v != nil {
		user := v.(User)
		return &user, nil
	} else {
		return nil, errors.New("no current user in context")
	}
}

