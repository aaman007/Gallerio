package context

import (
	"context"
	"gallerio/accounts"
)

var (
	userKey privateKey = "user"
)

type privateKey string

func WithUser(ctx context.Context, user *accounts.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *accounts.User {
	if temp := ctx.Value(userKey); temp != nil {
		if user, ok := temp.(*accounts.User); ok {
			return user
		}
	}
	return nil
}
