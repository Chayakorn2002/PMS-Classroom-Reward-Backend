package transport

import (
	"context"
)

// Service is a function that performs the actual work of a service.
// It process the request and returns the response and error.
// error return from Service should be an agree upon type errors.
// we will not inject the error into the response as it will be handle by endpoint
type Service[T, R any] func(ctx context.Context, req T) (resp R, err error)
