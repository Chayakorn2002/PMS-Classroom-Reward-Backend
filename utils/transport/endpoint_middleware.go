package transport

// Middleware is a function that wraps a service with additional behavior.
// It is a common pattern for adding cross-cutting concerns like logging, tracing, etc.
//
// Ex: Logging, Claim, etc.
type EndpointMiddleware[T, R any] func(Service[T, R]) Service[T, R]
