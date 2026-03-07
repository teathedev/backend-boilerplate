// Package rest provides generic REST/Huma helpers for the application.
// Application-specific setup (routers, servers, Swagger UI) lives in apps/api/rest.
package rest

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/teathedev/backend-boilerplate/pkg/errors"
)

// SetupErrorFactory configures Huma so that huma.ErrorXXX(...) returns
// *errors.CustomError. Call once from the application (e.g. apps/api/rest init).
func SetupErrorFactory() {
	huma.NewError = func(status int, msg string, _ ...error) huma.StatusError {
		return errors.NewWithStatus("API", msg, status)
	}
}
