package shared

import (
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
	beecontext "github.com/beego/beego/v2/server/web/context"
)

func CORSMiddleware() beego.HandleFunc {
	return func(ctx *beecontext.Context) {
		ctx.Output.Header("Access-Control-Allow-Origin", "*")
		ctx.Output.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		ctx.Output.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, Authorization")
		ctx.Output.Header("Access-Control-Expose-Headers", "Content-Length")
		ctx.Output.Header("Access-Control-Max-Age", "86400")

		if ctx.Request.Method == http.MethodOptions {
			ctx.Output.SetStatus(http.StatusNoContent)
			_ = ctx.Output.Body(nil)
			return
		}
	}
}

func RegisterCORSMiddleware() {
	beego.InsertFilter("*", beego.BeforeRouter, CORSMiddleware(), beego.WithReturnOnOutput(true))
}
