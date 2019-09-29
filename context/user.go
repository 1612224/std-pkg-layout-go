package context

import (
	"context"
	"log"
	app "useritem"
)

type contextKey string

const (
	userKey contextKey = "user"
)

// WithUser derives a new context with an user
func WithUser(ctx context.Context, user *app.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// User retrieves an user from context
func User(ctx context.Context) *app.User {
	tmp := ctx.Value(userKey)
	if tmp == nil {
		// user not found
		return nil
	}
	user, ok := tmp.(*app.User)
	if !ok {
		// value is not an user
		// this is a bug
		log.Fatalf("context: user value set incorrectly. type=%T, value=%#v", tmp, tmp)
		return nil
	}
	return user
}
