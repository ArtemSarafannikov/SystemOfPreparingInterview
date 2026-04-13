package directives

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/CodefriendOrg/arya/internal/pkg/auth"
	"github.com/CodefriendOrg/arya/internal/pkg/user_error"
)

// Auth .
func Auth(ctx context.Context, _ any, next graphql.Resolver) (any, error) {
	user := auth.GetUserFromContext(ctx)
	if user == nil {
		return nil, user_error.WithoutLoggerMessage(user_error.Unauthorized)
	}

	return next(ctx)
}
