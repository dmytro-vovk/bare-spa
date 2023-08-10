package errors

import (
	"fmt"
	"github.com/Sergii-Kirichok/DTekSpeachParser/pkg/jsonrpc"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type ErrorType int

const (
	NoType = ErrorType(iota)
	ParsingErr
	BadRequest
	NotFound
	InvalidParams
	ValidationErr
	CommandErr
)

type customError struct {
	code    ErrorType
	err     error
	context map[string]interface{}
}

// Error returns the message of a customError
func (e customError) Error() string {
	return e.err.Error()
}

// New creates a new customError
func (e ErrorType) New(msg string) error { return e.Newf(msg) }

// Newf creates a new customError with formatted message
func (e ErrorType) Newf(msg string, args ...interface{}) error {
	return customError{
		code: e,
		err:  fmt.Errorf(msg, args...),
	}
}

func (e ErrorType) Use(err error) error {
	if err == nil {
		return nil
	}

	return customError{
		code: e,
		err:  err,
	}
}

// Wrap creates a new wrapped error
func (e ErrorType) Wrap(err error, msg string) error { return e.Wrapf(err, msg) }

// Wrapf creates a new wrapped error with formatted message
func (e ErrorType) Wrapf(err error, msg string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return customError{
		code: e,
		err:  errors.Wrapf(err, msg, args...),
	}
}

// New creates a no type error
func New(msg string) error { return Newf(msg) }

// Newf creates a no type error with formatted message
func Newf(msg string, args ...interface{}) error {
	return customError{
		code: NoType,
		err:  fmt.Errorf(msg, args...),
	}
}

func Origin(err error) error {
	if customErr, ok := err.(customError); ok {
		return Origin(customErr.err)
	}

	return err
}

// Wrap wraps an error with a string
func Wrap(err error, msg string) error { return Wrapf(err, msg) }

// Wrapf wraps an error with format string
func Wrapf(err error, msg string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	wrappedErr := errors.Wrapf(err, msg, args...)
	if customErr, ok := err.(customError); ok {
		return customError{
			code:    customErr.code,
			err:     wrappedErr,
			context: customErr.context,
		}
	}

	return customError{
		code: NoType,
		err:  wrappedErr,
	}
}

func Unwrap(err error) error {
	if customErr, ok := err.(customError); ok {
		return errors.Unwrap(customErr.err)
	}

	return errors.Unwrap(err)
}

// Cause gives the original error
func Cause(err error) error {
	if customErr, ok := err.(customError); ok {
		return errors.Cause(customErr.err)
	}

	return errors.Cause(err)
}

// AddErrorContext adds a context to an error
func AddErrorContext(err error, key string, data interface{}) error {
	if customErr, ok := err.(customError); ok {
		if customErr.context == nil {
			customErr.context = map[string]interface{}{key: data}
		} else {
			customErr.context[key] = data
		}

		return customErr
	}

	return customError{code: NoType, err: err, context: map[string]interface{}{key: data}}
}

// GetContext returns the error context
func GetContext(err error) map[string]interface{} {
	if customErr, ok := err.(customError); ok {
		return customErr.context
	}

	return nil
}

// GetType returns the error type
func GetType(err error) ErrorType {
	if customErr, ok := err.(customError); ok {
		return customErr.code
	}

	return NoType
}

func AsJSONRPC(err error, trans ut.Translator) *jsonrpc.Error {
	if err == nil {
		return nil
	}

	code := GetType(err)
	switch code {
	case ParsingErr:
		code = jsonrpc.ParseError
	case BadRequest:
		code = jsonrpc.InvalidRequest
	case NotFound:
		code = jsonrpc.MethodNotFound
	case InvalidParams:
		code = jsonrpc.InvalidParams
	case ValidationErr:
		ret := New("invalid input data")
		for _, vErr := range Origin(err).(validator.ValidationErrors) {
			ret = AddErrorContext(ret, vErr.Field(), vErr.Translate(trans))
		}
		err = ret
	}

	return jsonrpc.NewError(int(code), err, GetContext(err))
}
