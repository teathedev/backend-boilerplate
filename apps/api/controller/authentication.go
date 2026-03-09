package controller

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/teathedev/backend-boilerplate/apps/api/rest"
	"github.com/teathedev/backend-boilerplate/internal/usecases"
	"github.com/teathedev/pkg/validation"
	"github.com/teathedev/backend-boilerplate/types"
)

// LoginInput is the request body for POST /auth/login.
type LoginInput struct {
	Body types.Login
}

// LoginOutput is the response for POST /auth/login.
type LoginOutput struct {
	Body types.AuthenticationResult `json:"body"`
}

// RegisterInput is the request body for POST /auth/register.
type RegisterInput struct {
	Body types.Register
}

// RegisterOutput is the response for POST /auth/register.
type RegisterOutput struct {
	Body types.AuthenticationResult `json:"body"`
}

func init() {
	huma.Register(rest.IAMServer.API, huma.Operation{
		OperationID: "login",
		Method:      http.MethodPost,
		Path:        "/login",
		Summary:     "Login with identifier and password",
		Responses:   rest.ErrorResponses(rest.IAMBadInputSchema, rest.IAMErrorSchema),
	}, func(ctx context.Context, in *LoginInput) (*LoginOutput, error) {
		result, err := usecases.Authentication.Login(ctx, &in.Body)
		if err != nil {
			return nil, err
		}
		return &LoginOutput{Body: *result}, nil
	})

	huma.Register(rest.IAMServer.API, huma.Operation{
		OperationID: "register",
		Method:      http.MethodPost,
		Path:        "/register",
		Summary:     "Register a new user",
		Responses:   rest.ErrorResponses(rest.IAMBadInputSchema, rest.IAMErrorSchema),
	}, func(ctx context.Context, in *RegisterInput) (*RegisterOutput, error) {
		if err := validation.ValidateStruct(&in.Body); err != nil {
			return nil, err
		}
		result, err := usecases.Authentication.Register(ctx, &in.Body)
		if err != nil {
			return nil, err
		}
		return &RegisterOutput{Body: *result}, nil
	})
}
