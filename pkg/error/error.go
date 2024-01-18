package error

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// Error stores information about an error.
type Error struct {
	Code    int
	Message string
	err     error
	cause   error
	stack   string
}

// NewError creates a new error.
func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		err:     errors.New(message),
		stack:   message + "\n" + createStackTrace(),
	}
}

// WrapError creates a new error wrapping an existing error.
func WrapError(code int, message string, err error) *Error {
	switch cErr := err.(type) {
	case *Error:
		return &Error{
			Code:    code,
			Message: message,
			err:     errors.New(message + ": " + cErr.err.Error()),
			cause:   cErr,
			stack:   message + "\n" + createStackTrace() + "Caused by: " + cErr.stack,
		}
	default:
		return &Error{
			Code:    code,
			Message: message,
			err:     errors.New(message + ": " + cErr.Error()),
			cause:   cErr,
			stack:   message + "\n" + cErr.Error() + "\n" + createStackTrace(),
		}
	}
}

// Error returns the string representation.
func (e *Error) Error() string {
	return e.err.Error()
}

// Cause returns the wrapped error or nil if error is route cause.
func (e *Error) Cause() error {
	return e.cause
}

// StackTrace returns the stack trace.
func (e *Error) StackTrace() string {
	return e.stack
}

// IsAuthError returns true if error is a authentication/authorization error.
func (e *Error) IsAuthError() bool {
	return e.Code <= -100 && e.Code > -200
}

// IsPermissionError returns true if error is a permission error.
func (e *Error) IsPermissionError() bool {
	return e.Code <= -200 && e.Code > -300
}

// IsValidationError returns true if error is a validation error.
func (e *Error) IsValidationError() bool {
	return e.Code <= -300 && e.Code > -400
}

// IsLogicError returns true if error is a logic error.
func (e *Error) IsLogicError() bool {
	return e.Code <= -400 && e.Code > -500
}

// IsSystemError returns true if error is a system error.
func (e *Error) IsSystemError() bool {
	return e.Code <= -500
}

// --- Helper functions ---

func createStackTrace() string {
	pcs := make([]uintptr, 32)
	n := runtime.Callers(3, pcs)
	pcs = pcs[:n]
	frames := runtime.CallersFrames(pcs)

	var str strings.Builder
	for {
		frame, more := frames.Next()
		if strings.Contains(frame.File, "runtime/") {
			break
		}
		str.WriteString(fmt.Sprintf("  %s\n    %s:%d\n", frame.Function, frame.File, frame.Line))
		if !more {
			break
		}
	}
	return str.String()
}
