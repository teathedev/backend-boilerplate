package errors

func NewUnsupportedFileFormatError(module string) *CustomError {
	return &CustomError{
		Module:  module,
		Message: "Unsupported File Format",
		Params:  nil,
		Status:  415,
	}
}
