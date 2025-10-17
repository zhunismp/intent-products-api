package errors

import (
	"fmt"
)

type AppError struct {
	// Code    string `json:"code"`
	Status  StatusCode `json:"status"`
	Message string     `json:"message"`
	Err     error      `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err == nil {
		return e.Message
	}
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func (e *AppError) Unwrap() error { return e.Err }

// constructor
func BadRequest(msg string, err error) *AppError {
	return &AppError{Status: StatusBadRequest, Message: msg, Err: err}
}

func Internal(msg string, err error) *AppError {
	return &AppError{Status: StatusInternalServerError, Message: msg, Err: err}
}

func NotFound(msg string, err error) *AppError {
	return &AppError{Status: StatusNotFound, Message: msg, Err: err}
}
