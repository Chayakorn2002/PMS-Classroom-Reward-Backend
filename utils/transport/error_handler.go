package transport

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Chayakorn2002/pms-classroom-backend/domain/dto"
	"github.com/Chayakorn2002/pms-classroom-backend/domain/exceptions"
)

func ErrorHandler(r *http.Request, w http.ResponseWriter, err error) {
	// fmt.Println("r.Context() in error handler", r.Context())

	// Retrieve neccessary details
	// Status code defaults to 500
	httpCode := http.StatusInternalServerError

	// Use response code from ExceptionError
	var cErr *exceptions.ExceptionError

	switch {
	// Error is an ExceptionError
	case errors.As(err, &cErr):
		fmt.Println("cErr.DebugMessage", cErr.DebugMessage)
		httpCode = cErr.HttpStatusCode

		resp := &dto.CommonErrorResponse{
			Status: cErr.Code,
			Error: &dto.CommonError{
				Code:    cErr.APIStatusCode,
				Message: cErr.GlobalMessage,
			},
		}
		JsonResponseError(w, resp, err, httpCode)

	// Default error response
	default:
		resp := &dto.CommonErrorResponse{
			Status: exceptions.NewGlobalErrors().ErrInternal.APIStatusCode,
			Error: &dto.CommonError{
				Code:    cErr.APIStatusCode,
				Message: exceptions.NewGlobalErrors().ErrInternal.GlobalMessage,
			},
		}
		JsonResponseError(w, resp, err, httpCode)
	}
}
