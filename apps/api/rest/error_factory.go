package rest

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/teathedev/pkg/errors"
)

// SetupErrorFactory configures Huma so that huma.ErrorXXX(...) returns
// *errors.CustomError. Call once from the application (init in this package).
func SetupErrorFactory() {
	huma.NewError = func(status int, msg string, _ ...error) huma.StatusError {
		return errors.NewWithStatus("API", msg, status)
	}
}
