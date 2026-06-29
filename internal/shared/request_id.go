package shared

import (
	"crypto/rand"
	"encoding/base64"

	beecontext "github.com/beego/beego/v2/server/web/context"
)

const requestIDContextKey = "request_id"

func RequestIDMiddleware(ctx *beecontext.Context) {
	requestID := ctx.Input.Header("X-Request-Id")
	if requestID == "" {
		requestID = newRequestID()
	}

	ctx.Input.SetData(requestIDContextKey, requestID)
	ctx.Output.Header("X-Request-Id", requestID)
}

func RequestID(ctx *beecontext.Context) string {
	if ctx == nil || ctx.Input == nil {
		return ""
	}

	requestID, ok := ctx.Input.GetData(requestIDContextKey).(string)
	if !ok || requestID == "" {
		requestID = newRequestID()
		ctx.Input.SetData(requestIDContextKey, requestID)
	}

	return requestID
}

func newRequestID() string {
	b := make([]byte, 30)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
