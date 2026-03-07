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

func NewBadInput(module string, params []BadInputField) *CustomError {
	return &CustomError{
		Module:  module,
		Message: "Validation Failed",
		Params:  params,
		Status:  400,
	}
}
