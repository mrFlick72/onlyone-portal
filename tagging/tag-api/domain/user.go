package domain

import (
	"context"
	"errors"
)

type UserName = string


type User struct {
	UserName *UserName
	Authorities *[]string
}

func GetCurrentUser(ctx *context.Context) (*User, error) {
	if v := (*ctx).Value("user"); v != nil {
		user := v.(User)
		return &user, nil
	} else {
		return nil, errors.New("no current user in context")
	}
}

// todo to delete
func SetCurrentUser(user User, ctx *context.Context) (*context.Context, error) {
	newCtx := context.WithValue(*ctx, "user", user)
	return &newCtx, nil
}
