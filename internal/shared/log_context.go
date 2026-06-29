package shared

import "context"

type logContextKey struct{}

type LogContext struct {
	Service     string
	Position    string
	RequestID   string
	URL         string
	RequestBody []byte
}

func WithLogContext(ctx context.Context, logCtx LogContext) context.Context {
	return context.WithValue(ctx, logContextKey{}, logCtx)
}

func LogContextFrom(ctx context.Context) (LogContext, bool) {
	logCtx, ok := ctx.Value(logContextKey{}).(LogContext)
	return logCtx, ok
}
