package errors

func NewUnauthorizedError(module string) *CustomError {
	return &CustomError{
		Module:  module,
		Message: "Unauthorized",
		Params:  nil,
		Status:  401,
	}
}
