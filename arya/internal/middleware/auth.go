package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/CodefriendOrg/arya/internal/pb/github.com/CodefriendOrg/kingsguard/pkg/kingsguard"
	"github.com/CodefriendOrg/arya/internal/pkg/auth"
	"github.com/CodefriendOrg/arya/internal/pkg/helper/kingsguard_helper"
)

// AuthMiddleware .
func AuthMiddleware(ctx context.Context, kingsguardHelper *kingsguard_helper.Service, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		var user *kingsguard.User
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString != authHeader {
			user = kingsguardHelper.GetUserFromAccessTokenWoError(ctx, tokenString)
		}

		ctxWithUser := auth.WithUser(r.Context(), user)

		next.ServeHTTP(w, r.WithContext(ctxWithUser))
	})
}
