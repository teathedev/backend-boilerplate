package rest

import (
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/teathedev/backend-boilerplate/constants"
	"github.com/teathedev/backend-boilerplate/internal/usecases"
	"github.com/teathedev/backend-boilerplate/types"
)

// AuthMiddleware validates the Bearer token and injects the user into context.
// Returns 401 if the token is missing or invalid.
func AuthMiddleware(ctx huma.Context, next func(huma.Context)) {
	auth := ctx.Header("Authorization")
	if auth == "" {
		huma.WriteErr(APIInstance, ctx, http.StatusUnauthorized, "Missing Authorization header")
		return
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(auth, prefix) {
		huma.WriteErr(APIInstance, ctx, http.StatusUnauthorized, "Invalid Authorization format")
		return
	}
	token := strings.TrimPrefix(auth, prefix)
	user, err := usecases.Authentication.GetUserByToken(ctx.Context(), token)
	if err != nil {
		huma.WriteErr(APIInstance, ctx, http.StatusUnauthorized, "Invalid or expired token")
		return
	}
	ctx = huma.WithValue(ctx, constants.User, user)
	next(ctx)
}

// RequireRoleMiddleware checks that the authenticated user has one of the allowed roles.
// Must be used after AuthMiddleware. Returns 403 if the user's role is not allowed.
func RequireRoleMiddleware(allowed ...types.UserRoles) func(huma.Context, func(huma.Context)) {
	allowedSet := make(map[types.UserRoles]struct{}, len(allowed))
	for _, r := range allowed {
		allowedSet[r] = struct{}{}
	}
	return func(ctx huma.Context, next func(huma.Context)) {
		userVal := ctx.Context().Value(constants.User)
		if userVal == nil {
			huma.WriteErr(APIInstance, ctx, http.StatusUnauthorized, "Authentication required")
			return
		}
		user, ok := userVal.(*types.User)
		if !ok {
			huma.WriteErr(APIInstance, ctx, http.StatusInternalServerError, "Invalid user context")
			return
		}
		if _, ok := allowedSet[user.Role]; !ok {
			huma.WriteErr(APIInstance, ctx, http.StatusForbidden, "Insufficient permissions")
			return
		}
		next(ctx)
	}
}
