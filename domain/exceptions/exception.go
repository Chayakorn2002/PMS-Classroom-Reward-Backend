package exceptions

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

type ExceptionError struct {
	Code           int
	HttpStatusCode int
	APIStatusCode  int
	GlobalMessage  string
	DebugMessage   string
	StackErrors    error
	StackCaller    []byte
}

func NewExceptionError(code int, apiStatusCode int, globalMessage string, httpStatusCode int) *ExceptionError {
	return &ExceptionError{
		Code:           code,
		APIStatusCode:  apiStatusCode,
		GlobalMessage:  globalMessage,
		HttpStatusCode: httpStatusCode,
	}
}

func (cErr *ExceptionError) Error() string {
	return cErr.GlobalMessage
}

func (cErr *ExceptionError) WithDebugMessage(debugMessage string) *ExceptionError {
	cErr.DebugMessage = debugMessage
	return cErr
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type errorField struct {
	Kind    string `json:"kind"`
	Stack   string `json:"stack"`
	Message string `json:"message"`
}

func GetStackField(err error) errorField {
	var stack string

	if serr, ok := err.(stackTracer); ok {
		st := serr.StackTrace()
		stack = fmt.Sprintf("%+v", st)
		if len(stack) > 0 && stack[0] == '\n' {
			stack = stack[1:]
		}
	}
	return errorField{
		Kind:    reflect.TypeOf(err).String(),
		Stack:   stack,
		Message: err.Error(),
	}
}
