package ctxkey

type ContextKey struct {
	Name string
}

var (
	CTX_KEY_ENV             = ContextKey{Name: "env"}
	CTX_KEY_REQUEST_ID      = ContextKey{Name: "request_id"}
	CTX_KEY_ACCEPT_LANGUAGE = ContextKey{Name: "accept_language"}
	CTX_KEY_AUTHORIZATION   = ContextKey{Name: "authorization"}
)
