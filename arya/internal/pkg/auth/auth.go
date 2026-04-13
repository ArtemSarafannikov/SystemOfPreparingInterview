package auth

import (
	"context"

	"github.com/CodefriendOrg/arya/internal/pb/github.com/CodefriendOrg/kingsguard/pkg/kingsguard"
)

type contextKey string

// userCtxKey .
const userCtxKey = contextKey("user")

// User .
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// GetUserFromContext .
func GetUserFromContext(ctx context.Context) *User {
	user, _ := ctx.Value(userCtxKey).(*User)
	return user
}

// WithUser .
func WithUser(ctx context.Context, user *kingsguard.User) context.Context {
	if user == nil {
		return ctx
	}

	return context.WithValue(ctx, userCtxKey, &User{
		ID:       user.Id,
		Username: user.Username,
	})
}
