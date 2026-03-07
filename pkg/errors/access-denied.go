package errors

func NewAccessDeniedError(module string) *CustomError {
	return &CustomError{
		Module:  module,
		Message: "Access Denied",
		Params:  nil,
		Status:  403,
	}
}
