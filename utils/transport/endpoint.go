package transport

type Endpoint[T, R any] func() (fn Service[T, R])

// NewEndpoint is a factory function that returns an endpoint function. The endpoint function
// returns a service function that has been wrapped with the provided middlewares.
// The middlewares are applied in the order they are provided.
//
// Ex: NewEndpoint(svc, Logging, Claim) returns an endpoint function that returns a service function
// that has been wrapped with Logging and Claim middlewares.
// The order will be Claim -> Logging -> svc
func NewEndpoint[T, R any](svc Service[T, R], middlewares ...EndpointMiddleware[T, R]) func() Endpoint[T, R] {
	return func() Endpoint[T, R] {
		return func() Service[T, R] {
			var newSvc Service[T, R]
			init := true
			for _, m := range middlewares {
				if init {
					newSvc = m(svc)
					init = false
				} else {
					newSvc = m(newSvc)
				}
			}

			if len(middlewares) == 0 {
				return svc
			}

			return newSvc
		}
	}
}
