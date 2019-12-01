package e

import (
	"errors"
	"github.com/sirupsen/logrus"
)

var (
	// Error for all internal error should not be exposed
	ErrInternalError = errors.New("the request failed by internal error")
)

type ResponseError struct {
	Error string `json:"error"`
}

// function to create a response error
func CreateErr(err error) *ResponseError {
	res := &ResponseError{Error: err.Error()}
	return res
}

// function to return internal 500 error
// takes the actual error for logging and return the generic error as response
func InternalError(err error) *ResponseError {
	logrus.Error(err)

	res := &ResponseError{Error: ErrInternalError.Error()}
	return res
}
