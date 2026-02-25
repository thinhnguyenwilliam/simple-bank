package util

import "context"

type userKeyType struct{}

var userKey = userKeyType{}

func ContextWithUser(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, userKey, username)
}

func UserFromContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(userKey).(string)
	return username, ok
}
