package user_error

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Type string

type Error struct {
	originalMessage string
	errorType       Type
}

func (t Type) Error() string {
	return string(t)
}

func (e *Error) Error() string {
	return e.originalMessage
}

func (e *Error) UserError() error {
	return e.errorType
}

func (e *Error) ErrorType() string {
	return e.errorType.Error()
}

// WithoutLoggerMessage .
func WithoutLoggerMessage(errorType Type) error {
	return &Error{
		originalMessage: "",
		errorType:       errorType,
	}
}

// FromError .
func FromError(errorType Type, err error) error {
	if err == nil {
		return nil
	}

	return &Error{
		originalMessage: err.Error(),
		errorType:       errorType,
	}
}

// FromGRPCError .
func FromGRPCError(err error) error {
	switch status.Code(err) {
	case codes.OK:
		return nil
	case codes.NotFound:
		return &Error{
			errorType: NotFound,
		}
	case codes.Unauthenticated:
		return &Error{
			errorType: Unauthorized,
		}
	default:
		return &Error{
			errorType: InternalError,
		}
	}
}

// New .
func New(errorType Type, message string, args ...any) error {
	return &Error{
		originalMessage: fmt.Sprintf(message, args...),
		errorType:       errorType,
	}
}
