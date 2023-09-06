package core

import (
	"errors"
	"fmt"
)

// Application error codes
const (
	ECONFLICT       = "conflict"
	EINTERNAL       = "internal"
	EINVALID        = "invalid"
	ENOTFOUND       = "not_found"
	ENOTIMPLEMENTED = "not_implemented"
	EUNAUTHORIZED   = "unauthorized"
)

type Error struct {
	// Machine-readable error code
	Code string

	// Human-readable error message
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("app error: code=%s message=%v", e.Code, e.Message)
}

func ErrorCode(err error) string {
	var e *Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Code
	}
	return EINTERNAL
}

func ErrorMessage(err error) string {
	var e *Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Message
	}
	return "Internal error."
}

func WrapErr(err error) error {
	if err == nil {
		return nil
	}
	var appError *Error
	if errors.As(err, &appError) {
		return appError
	} else {
		return Errorw(EINTERNAL, err)
	}
}

func Errorf(code string, format string, args ...any) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

func Errorw(code string, err error) *Error {
	return &Error{
		Code:    code,
		Message: err.Error(),
	}
}
