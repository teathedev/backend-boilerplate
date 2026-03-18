package controller

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/teathedev/backend-boilerplate/apps/api/rest"
	"github.com/teathedev/backend-boilerplate/constants"
	"github.com/teathedev/backend-boilerplate/internal/usecases"
	"github.com/teathedev/backend-boilerplate/types"
	"github.com/teathedev/pkg/logger"
)

// LoginInput is the request body for POST /auth/login.
type LoginInput struct {
	Body types.Login
}

// LoginOutput is the response for POST /auth/login.
type LoginOutput struct {
	Body types.AuthenticationResult
}

// RefreshInput is the request body for POST /auth/refresh.
type RefreshInput struct {
	Body types.RefreshTokenRequest
}

// RefreshOutput is the response for POST /auth/refresh.
type RefreshOutput struct {
	Body types.AuthenticationResult
}

// RegisterInput is the request body for POST /auth/register.
type RegisterInput struct {
	Body types.Register
}

// RegisterOutput is the response for POST /auth/register.
type RegisterOutput struct {
	Body types.AuthenticationResult
}

func init() {
	log := logger.New("AuthenticationController")
	huma.Register(rest.IAMGroup, huma.Operation{
		OperationID: "login",
		Method:      http.MethodPost,
		Path:        "/login",
		Summary:     "Login with identifier and password",
		Responses:   rest.ErrorResponses(),
	}, func(ctx context.Context, in *LoginInput) (*LoginOutput, error) {
		result, err := usecases.Authentication.Login(ctx, &in.Body)
		if err != nil {
			return nil, err
		}
		return &LoginOutput{Body: *result}, nil
	})

	huma.Register(rest.IAMGroup, huma.Operation{
		OperationID: "refresh",
		Method:      http.MethodPost,
		Path:        "/refresh",
		Summary:     "Refresh access token using refresh token",
		Responses:   rest.ErrorResponses(),
	}, func(ctx context.Context, in *RefreshInput) (*RefreshOutput, error) {
		result, err := usecases.Authentication.Refresh(ctx, &in.Body)
		if err != nil {
			return nil, err
		}
		return &RefreshOutput{Body: *result}, nil
	})

	huma.Register(rest.IAMGroup, huma.Operation{
		OperationID: "register",
		Method:      http.MethodPost,
		Path:        "/register",
		Summary:     "Register a new user",
		Responses:   rest.ErrorResponses(),
	}, func(ctx context.Context, in *RegisterInput) (*RegisterOutput, error) {
		log.Info(
			"Registering user",
			logger.LogParams{
				"RequestID": ctx.Value(constants.RequestID),
				"Body":      in.Body,
			},
		)
		result, err := usecases.Authentication.Register(ctx, &in.Body)
		if err != nil {
			log.Error(
				"Failed to register user",
				logger.LogParams{
					"RequestID": ctx.Value(constants.RequestID),
					"Error":     err,
				},
			)
			return nil, err
		}
		return &RegisterOutput{Body: *result}, nil
	})
}
