package transport

import (
	"bytes"
	"net/http"

	"github.com/Chayakorn2002/pms-classroom-backend/utils/transport"
)

func ErrorHandlingMiddleware() TransportMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			crw := &transport.CustomResponseWriter{
				ResponseWriter: w,
				Body:           new(bytes.Buffer),
				StatusCode:     http.StatusOK,
			}

			next.ServeHTTP(crw, r)

			if crw.StatusCode >= http.StatusBadRequest {
				http.Error(w, crw.Body.String(), crw.StatusCode)
			}
		})
	}
}
