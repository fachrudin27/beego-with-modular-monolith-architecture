package shared

import beecontext "github.com/beego/beego/v2/server/web/context"

func GetRequestIdContext(ctx *beecontext.Context) string {
	return RequestID(ctx)
}
