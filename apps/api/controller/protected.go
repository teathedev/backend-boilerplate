package controller

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/teathedev/backend-boilerplate/apps/api/rest"
	"github.com/teathedev/backend-boilerplate/constants"
	"github.com/teathedev/backend-boilerplate/internal/usecases"
	"github.com/teathedev/backend-boilerplate/types"
)

// PingOutput is the response for GET /protected/ping (public).
type PingOutput struct {
	Body struct {
		Message string `json:"message" doc:"Static ping response"`
	} `json:"body"`
}

// MeOutput is the response for GET /protected/me (auth required).
type MeOutput struct {
	Body types.User `json:"body"`
}

// UpdateMeInput is the request body for PATCH /protected/me.
type UpdateMeInput struct {
	Body types.UpdateMe `json:"body"`
}

// UpdatePasswordInput is the request body for PATCH /protected/me/password.
type UpdatePasswordInput struct {
	Body types.UpdatePassword `json:"body"`
}

// AdminOutput is the response for GET /protected/admin (role protected).
type AdminOutput struct {
	Body struct {
		Message string `json:"message" doc:"Static response for admin users"`
	} `json:"body"`
}

func userFromContext(ctx context.Context) *types.User {
	u := ctx.Value(constants.User)
	if u == nil {
		return nil
	}
	user, ok := u.(*types.User)
	if !ok {
		return nil
	}
	return user
}

func init() {
	// Public endpoint - no auth required
	huma.Register(rest.ProtectedGroup, huma.Operation{
		OperationID: "ping",
		Method:      http.MethodGet,
		Path:        "/ping",
		Summary:     "Public ping endpoint",
		Responses:   rest.ErrorResponses(),
	}, func(ctx context.Context, _ *struct{}) (*PingOutput, error) {
		out := &PingOutput{}
		out.Body.Message = "pong"
		return out, nil
	})

	// Auth required - get authenticated user (middleware added per-endpoint)
	huma.Register(rest.ProtectedGroup, huma.Operation{
		OperationID: "me",
		Method:      http.MethodGet,
		Path:        "/me",
		Summary:     "Get current authenticated user",
		Middlewares: huma.Middlewares{rest.AuthMiddleware},
		Responses:   rest.ErrorResponses(),
	}, func(ctx context.Context, _ *struct{}) (*MeOutput, error) {
		u := userFromContext(ctx)
		if u == nil {
			return nil, nil // AuthMiddleware ensures user exists; defensive
		}
		return &MeOutput{Body: *u}, nil
	})

	// Auth required - update authenticated user (middleware added per-endpoint)
	huma.Register(rest.ProtectedGroup, huma.Operation{
		OperationID: "updateMe",
		Method:      http.MethodPatch,
		Path:        "/me",
		Summary:     "Update current user profile",
		Middlewares: huma.Middlewares{rest.AuthMiddleware},
		Responses:   rest.ErrorResponses(),
	}, func(ctx context.Context, in *UpdateMeInput) (*MeOutput, error) {
		u := userFromContext(ctx)
		if u == nil {
			return nil, nil
		}
		updated, err := usecases.User.UpdateMe(ctx, u.ID, &in.Body)
		if err != nil {
			return nil, err
		}
		return &MeOutput{Body: *updated}, nil
	})

	// Auth required - update password (middleware added per-endpoint)
	huma.Register(rest.ProtectedGroup, huma.Operation{
		OperationID: "updatePassword",
		Method:      http.MethodPatch,
		Path:        "/me/password",
		Summary:     "Update current user password",
		Middlewares: huma.Middlewares{rest.AuthMiddleware},
		Responses:   rest.ErrorResponses(),
	}, func(ctx context.Context, in *UpdatePasswordInput) (*struct{}, error) {
		u := userFromContext(ctx)
		if u == nil {
			return nil, nil
		}
		if err := usecases.User.UpdatePassword(ctx, u.ID, &in.Body); err != nil {
			return nil, err
		}
		return &struct{}{}, nil
	})

	// Role protected - requires SuperUser role
	huma.Register(rest.ProtectedGroup, huma.Operation{
		OperationID: "admin",
		Method:      http.MethodGet,
		Path:        "/admin",
		Summary:     "Admin-only endpoint (requires SuperUser role)",
		Middlewares: huma.Middlewares{rest.AuthMiddleware, rest.RequireRoleMiddleware(types.UserRolesSuperUser)},
		Responses:   rest.ErrorResponses(),
	}, func(ctx context.Context, _ *struct{}) (*AdminOutput, error) {
		out := &AdminOutput{}
		out.Body.Message = "Welcome, admin"
		return out, nil
	})
}
