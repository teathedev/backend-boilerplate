package errors

type BadInputCondition string

const (
	BadInputConditionNotValid BadInputCondition = "not_valid"
	BadInputConditionNotFound BadInputCondition = "not_found"
)

type BadInputField struct {
	Field     string            `json:"field"`
	Condition BadInputCondition `json:"condition"`
	Value     string            `json:"value"`
}

// APIErrorResponse is the documented shape of 400 (bad input / validation) error bodies (OpenAPI).
// Use only for 400 responses so the spec shows message + error array.
type APIErrorResponse struct {
	Message string          `json:"message" doc:"Error message"`
	Error   []BadInputField `json:"error,omitempty" doc:"Validation errors (only for 400)"`
}

// APIErrorSimple is the documented shape of 401/403/404/500 error bodies (OpenAPI).
// Message only; no validation error array.
type APIErrorSimple struct {
	Message string `json:"message" doc:"Error message"`
}

func NewBadInput(module string, params []BadInputField) *CustomError {
	return &CustomError{
		Module:  module,
		Message: "Validation Failed",
		Params:  params,
		Status:  400,
	}
}
