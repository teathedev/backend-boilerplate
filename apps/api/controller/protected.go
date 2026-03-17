package controller

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/teathedev/backend-boilerplate/apps/api/rest"
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
	Body struct {
		Message string `json:"message" doc:"Static response for authenticated user"`
	} `json:"body"`
}

// AdminOutput is the response for GET /protected/admin (role protected).
type AdminOutput struct {
	Body struct {
		Message string `json:"message" doc:"Static response for admin users"`
	} `json:"body"`
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

	// Auth required - requires valid Bearer token
	huma.Register(rest.ProtectedGroup, huma.Operation{
		OperationID: "me",
		Method:      http.MethodGet,
		Path:        "/me",
		Summary:     "Get current user info (requires authentication)",
		Middlewares: huma.Middlewares{rest.AuthMiddleware},
		Responses:   rest.ErrorResponses(),
	}, func(ctx context.Context, _ *struct{}) (*MeOutput, error) {
		out := &MeOutput{}
		out.Body.Message = "You are authenticated"
		return out, nil
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
