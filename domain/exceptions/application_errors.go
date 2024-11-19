package exceptions

import "net/http"

type GlobalErrors struct {
	ErrInternal                *ExceptionError
	ErrValidation              *ExceptionError
	ErrUnableToProceed         *ExceptionError
	ErrNotFound                *ExceptionError
	ErrInvalidKeyPair          *ExceptionError
	ErrJwtInvalidSigningMethod *ExceptionError
	ErrJwtInvalidSignature     *ExceptionError
	ErrJwtInvalidToken         *ExceptionError
}

func NewGlobalErrors() *GlobalErrors {
	return &GlobalErrors{
		ErrInternal:   NewExceptionError(5000, 000000, "Internal error", http.StatusInternalServerError),
		ErrValidation: NewExceptionError(4000, 000001, "Invalid Request", http.StatusBadRequest),
		ErrNotFound:   NewExceptionError(4000, 000003, "Not found", http.StatusNotFound),
	}
}

type ApplicationError struct {
	ErrUnauthorized      *ExceptionError
	ErrNotFound          *ExceptionError
	ErrBadRequest        *ExceptionError
	ErrUserAlreadyExists *ExceptionError
	ErrEmailNotFound     *ExceptionError
	ErrInvalidCredential *ExceptionError
	ErrAlreadyRedeeemed  *ExceptionError
	ErrInternal          *ExceptionError
}

func NewApplicationError() *ApplicationError {
	return &ApplicationError{
		ErrUnauthorized:      NewExceptionError(4000, 100000, "Unauthorized", http.StatusUnauthorized),
		ErrNotFound:          NewExceptionError(4000, 100001, "Not Found", http.StatusBadRequest),
		ErrBadRequest:        NewExceptionError(4000, 100002, "Bad Request", http.StatusBadRequest),
		ErrUserAlreadyExists: NewExceptionError(4000, 100003, "User already exists", http.StatusBadRequest),
		ErrEmailNotFound:     NewExceptionError(4000, 100004, "Email not found", http.StatusBadRequest),
		ErrInvalidCredential: NewExceptionError(4000, 100005, "Invalid Credential", http.StatusBadRequest),
		ErrAlreadyRedeeemed:  NewExceptionError(4000, 100006, "Already Redeemed", http.StatusBadRequest),
		ErrInternal:          NewExceptionError(5000, 109999, "Internal Server Error", http.StatusInternalServerError),
	}
}
