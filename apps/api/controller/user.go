package controller

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/teathedev/backend-boilerplate/apps/api/rest"
	"github.com/teathedev/pkg/validation"
)

// CreateUserInput is the request body for POST /users.
type CreateUserInput struct {
	Body struct {
		Email string `json:"email" validate:"required,email" doc:"User email address" format:"email" maxLength:"255"`
		Name  string `json:"name"  validate:"required,min=2" doc:"Display name" minLength:"2" maxLength:"255"`
	}
}

// CreateUserOutput is the response for POST /users.
type CreateUserOutput struct {
	Body struct {
		ID string `json:"id" doc:"Created user ID" format:"uuid"`
	} `json:"body"`
}

func init() {
	huma.Register(rest.IAMServer.API, huma.Operation{
		OperationID: "create-user",
		Method:      http.MethodPost,
		Path:        "/users",
		Summary:     "Create a user",
		Responses:   rest.ErrorResponses(rest.IAMBadInputSchema, rest.IAMErrorSchema),
	}, func(ctx context.Context, in *CreateUserInput) (*CreateUserOutput, error) {
		if err := validation.ValidateStruct(&in.Body); err != nil {
			return nil, err
		}

		// TODO: wire to usecase and persist user; for now return placeholder id
		out := &CreateUserOutput{}
		out.Body.ID = uuid.New().String()
		return out, nil
	})
}
