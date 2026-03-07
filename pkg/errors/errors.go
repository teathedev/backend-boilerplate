// Package errors provides a custom error type for the application.
package errors

import (
	"encoding/json"
	"fmt"
)

type ICustomError interface {
	Error() string
	Log()
	serializeParams() string
}

type CustomError struct {
	Module  string `json:"-"`
	Message string `json:"message"`
	Params  any    `json:"error"`
	Status  int    `json:"-"`
}

func (e *CustomError) serializeParams() (string, error) {
	byt, err := json.Marshal(e.Params)

	return string(byt), err
}

func (e *CustomError) serialize() (string, error) {
	byt, err := json.Marshal(e)

	return string(byt), err
}

func (e CustomError) Error() string {
	msg, err := e.serialize()
	if err != nil {
		return "Serialize error failed"
	}
	return msg
}

func (e *CustomError) Log() string {
	params, err := e.serializeParams()
	if err != nil {
		return "Serialize params failed"
	}
	return fmt.Sprintf("[%s] # [%s] : [%s]", e.Module, e.Message, params)
}

func New(module string, message string) *CustomError {
	return &CustomError{
		Module:  module,
		Message: message,
		Params:  nil,
		Status:  500,
	}
}
