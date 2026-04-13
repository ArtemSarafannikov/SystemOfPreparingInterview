package middleware

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/CodefriendOrg/arya/internal/pkg/logger"
	"github.com/CodefriendOrg/arya/internal/pkg/user_error"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.uber.org/zap"
)

// ErrorMiddleware .
func ErrorMiddleware(ctx context.Context, err error) *gqlerror.Error {
	gqlErr := graphql.DefaultErrorPresenter(ctx, err)
	if gqlErr == nil {
		return nil
	}

	var userErr *user_error.Error
	switch {
	case errors.As(err, &userErr) && !errors.Is(err, user_error.InternalError):
		gqlErr.Message = userErr.ErrorType()
	default:
		logger.Errorf(ctx, "[INTERNAL]", zap.Error(err))
		gqlErr.Message = user_error.InternalError.Error()
	}

	return gqlErr
}

// PanicMiddleware .
func PanicMiddleware(ctx context.Context, err any) error {
	logger.Errorf(ctx, "[PANIC]", zap.Any("error", err))
	return user_error.InternalError
}
