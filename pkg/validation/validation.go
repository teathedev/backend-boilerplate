// Package validation provides custom validation functions and utilities for struct validation.
// It includes validators for user roles, states, subscriptions, passwords, and other domain-specific validations.
package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/teathedev/backend-boilerplate/pkg/errors"

	_validator "github.com/go-playground/validator/v10"
)

var validate = _validator.New()

// ValidateStruct validates the given struct using struct tags and returns a CustomError with
// validation field details if validation fails. Uses json tag names in error responses.
func ValidateStruct[T any](data *T) *errors.CustomError {
	errParams := []errors.BadInputField{}
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	for _, err := range err.(_validator.ValidationErrors) {
		element := errors.BadInputField{
			Field:     err.Field(),
			Condition: errors.BadInputConditionNotValid,
			Value:     fmt.Sprintf("%v", err.Value()),
		}
		errParams = append(errParams, element)
	}

	return errors.NewBadInput("Validation", errParams)
}

func init() {
	// use json tags as err.Field()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
}
