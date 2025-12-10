package domain

import (
	"context"
	"errors"
)

type UserName = string

type User struct {
	UserName *UserName
}

func GetCurrentUser(ctx *context.Context) (*User, error) {
	if v := (*ctx).Value("current_user"); v != nil {
		return &User{UserName: v.(*UserName)}, nil
	} else {
		return nil, errors.New("no current user in context")
	}
}

func SetCurrentUser(user *User, ctx *context.Context) (*context.Context, error) {
	newCtx := context.WithValue(*ctx, "current_user", user.UserName)
	return &newCtx, nil
}
