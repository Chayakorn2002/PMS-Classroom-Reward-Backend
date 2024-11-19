package transport

import (
	"context"
	"net/http"

	"github.com/Chayakorn2002/pms-classroom-backend/utils/ctxkey"
)

// ClaimMiddleware is a middleware that extracts the required headers from the incoming request and sets them in the context.
func ClaimMiddleware() TransportMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract required headers from the incoming request
			claim := GetRequiredHeadersFromHttp(r)

			// Set the extracted headers in the context
			r = SetUserDataToCtx(r, claim)

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// TODO: To be discussed
type CommonHeaders struct {
	AcceptLanguage string `json:"Accept-Language"`
	Authorization  string `json:"Authorization"`
	RequestId      string `json:"x-request-id"`
}

func GetRequiredHeadersFromHttp(r *http.Request) CommonHeaders {
	claim := CommonHeaders{
		AcceptLanguage: r.Header.Get("Accept-Language"),
		Authorization:  r.Header.Get("Authorization"),
		RequestId:      r.Header.Get("x-request-id"),
	}

	return claim
}

func SetUserDataToCtx(r *http.Request, claim CommonHeaders) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, ctxkey.CTX_KEY_ACCEPT_LANGUAGE, claim.AcceptLanguage)
	ctx = context.WithValue(ctx, ctxkey.CTX_KEY_AUTHORIZATION, claim.Authorization)
	ctx = context.WithValue(ctx, ctxkey.CTX_KEY_REQUEST_ID, claim.RequestId)

	return r.WithContext(ctx)
}
