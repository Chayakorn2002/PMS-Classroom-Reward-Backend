package ctxkey

import (
	"context"
)

var (
	CTX_KEY_ERR = ContextKey{Name: "error"}
)

func WithError(ctx context.Context, err error) context.Context {
	return context.WithValue(ctx, CTX_KEY_ERR, err)
}

func GetError(ctx context.Context) error {
	if err, ok := ctx.Value(CTX_KEY_ERR).(error); ok {
		return err
	}

	return nil
}
