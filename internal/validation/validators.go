package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/teathedev/pkg/validation"
)

var usernameRe = regexp.MustCompile(`^[a-zA-Z0-9_]{3,32}$`)

// Init registers custom validation rules for the application.
func init() {
	// Username: 3-32 chars, letters, numbers, underscore.
	_ = validation.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		v, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		return usernameRe.MatchString(v)
	})

	// User identifier: allow non-empty; detailed checks are done in usecases.
	_ = validation.RegisterValidation("user_identifier", func(fl validator.FieldLevel) bool {
		v, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		return v != ""
	})

	// E164 phone number: basic length + leading '+' check.
	_ = validation.RegisterValidation("e164", func(fl validator.FieldLevel) bool {
		v, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		l := len(v)
		if l < 8 || l > 15 {
			return false
		}
		if v[0] != '+' {
			return false
		}
		return true
	})

	// Password: delegate to any policy you like; minimal length here.
	_ = validation.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		v, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		return len(v) >= 8
	})
}
