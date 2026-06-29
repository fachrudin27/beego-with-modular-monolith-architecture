// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	// "firstbeegoapi/controllers"

	auth "firstbeegoapi/internal/auth/delivery/api"
	ordering "firstbeegoapi/internal/ordering/delivery/api"
	"firstbeegoapi/internal/shared"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {

	shared.RegisterCORSMiddleware()

	// beego.BConfig.RecoverPanic = true
	beego.InsertFilter("*", beego.BeforeRouter, shared.RequestIDMiddleware)
	beego.InsertFilter("*", beego.BeforeRouter, shared.IPLimiterFilter)
	beego.InsertFilter("/v1/ordering/*", beego.BeforeRouter, shared.JWTAuthMiddleware)

	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/auth",
			beego.NSRouter("/login", &auth.AuthController{}, "post:Login"),
		),
		beego.NSNamespace("/ordering",
			beego.NSRouter("/:objectId", &ordering.OrderingController{}, "get:Get"),
		),
	)
	beego.AddNamespace(ns)
}
