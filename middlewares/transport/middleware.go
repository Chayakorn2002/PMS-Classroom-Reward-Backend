package transport

import "net/http"

// TransportMiddleware defines a middleware function for HTTP transport.
type TransportMiddleware func(http.Handler) http.Handler

// CreateStack creates a middleware stack from a list of TransportMiddleware.
// The middlewares are applied in the order they are provided.
func CreateStack(middlewares ...TransportMiddleware) TransportMiddleware {
	return func(handler http.Handler) http.Handler {
		// Apply middlewares in reverse order
		for i := len(middlewares) - 1; i >= 0; i-- {
			handler = middlewares[i](handler)
		}
		return handler
	}
}

