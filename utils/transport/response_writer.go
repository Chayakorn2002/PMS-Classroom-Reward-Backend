package transport

import (
	"bytes"
	"net/http"
)

// CustomResponseWriter is a wrapper around http.ResponseWriter that captures the response body and status code.
type CustomResponseWriter struct {
	http.ResponseWriter
	Body       *bytes.Buffer
	StatusCode int
	Error      error
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *CustomResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *CustomResponseWriter) WriteError(b []byte, err error) {
	w.Error = err
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}