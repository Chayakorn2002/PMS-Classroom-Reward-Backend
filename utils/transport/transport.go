package transport

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Chayakorn2002/pms-classroom-backend/utils/ctxkey"
)

// NewTransport is a transport layer that decodes the incoming request and calls the endpoint
// with the decoded request. It also encodes the response and sends it back to the client.
func NewTransport[T, R any](req any, endpoint func() Endpoint[T, R], middlewares ...EndpointMiddleware[T, R]) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode the request
		requestBody, err := readRequestBody(r)
		if err != nil {
			JsonResponse(w, err, http.StatusBadRequest)
			return
		}

		// Unmarshal the request body
		if len(requestBody) != 0 {
			err = json.Unmarshal(requestBody, &req)
			if err != nil {
				JsonResponse(w, err, http.StatusBadRequest)
				return
			}
		}

		typedReq, ok := req.(T)
		if !ok {
			JsonResponse(w, "invalid request type", http.StatusInternalServerError)
			return
		}

		// Call the endpoint and handle the response
		resp, serviceError := callEndpoint(endpoint, r, typedReq)
		if serviceError != nil {
			ctx := ctxkey.WithError(r.Context(), serviceError)
			r = r.WithContext(ctx)
			// fmt.Println("r.Context() in transport", r.Context())
			ErrorHandler(r, w, serviceError)
		} else {
			JsonResponse(w, resp, http.StatusOK)
		}
	}
}

// readRequestBody reads the request body and returns it as a byte slice.
func readRequestBody(r *http.Request) ([]byte, error) {
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

// callEndpoint calls the endpoint with the given request and returns the response and any error.
func callEndpoint[T, R any](endpoint func() Endpoint[T, R], r *http.Request, req T) (R, error) {
	return endpoint()()(r.Context(), req)
}
