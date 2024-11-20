package transport

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/Chayakorn2002/pms-classroom-backend/domain/dto"
	"github.com/Chayakorn2002/pms-classroom-backend/domain/exceptions"
)

var (
	logFile *os.File
	logger  *log.Logger
)

func init() {
	var err error
	logFile, err = os.OpenFile("log/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	logger = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

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
		// write the debug message to log file
		logger.Printf("Debug message: %s, Error: %v", cErr.DebugMessage, err)

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
