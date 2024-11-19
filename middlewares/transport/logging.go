package transport

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/Chayakorn2002/pms-classroom-backend/utils/ctxkey"
	"github.com/Chayakorn2002/pms-classroom-backend/utils/logger"
	"github.com/Chayakorn2002/pms-classroom-backend/utils/transport"
)

func LoggingMiddleware() TransportMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			requestBody, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody)) // Restore the request body

			crw := &transport.CustomResponseWriter{
				ResponseWriter: w,
				Body:           new(bytes.Buffer),
				StatusCode:     http.StatusOK,
			}

			next.ServeHTTP(crw, r)

			// fmt.Println("r.Context() in logging middleware", r.Context())
			err := ctxkey.GetError(r.Context())

			elapse := time.Since(startTime)
			responseBody := crw.Body.String()
			headers := r.Header

			var fields []any
			// append md log
			fields = append(fields,
				slog.String("logger_name", "canonical"),
				slog.Group("httpserver_md",
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.String("ip", fmt.Sprint(headers["X-Forwarded-For"])),
					slog.String("duration", elapse.String()),
					slog.String("accept-language", convertHeaderAttrToString("Accept-Language", headers)),
					slog.String("x-request-id", convertHeaderAttrToString("X-Request-Id", headers)),
				),
			)

			var level logger.Level
			if crw.StatusCode >= http.StatusBadRequest {
				level = logger.Error
			} else {
				level = logger.Info
			}

			logger.CanonicalLogger(
				r.Context(),
				*logger.Slog,
				level,
				requestBody,
				[]byte(responseBody),
				err,
				logger.CanonicalLog{
					Transport: "http",
					Traffic:   "internal",
					Method:    r.Method,
					Status:    crw.StatusCode,
					Path:      r.URL.Path,
					Duration:  elapse,
				},
				fields,
			)
		})
	}
}

func convertHeaderAttrToString(key string, headers map[string][]string) string {
	if header, ok := headers[key]; ok {
		return header[0]
	}
	return ""
}
