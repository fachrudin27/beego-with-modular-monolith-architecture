package api

import (
	"encoding/json"
	"firstbeegoapi/internal/ordering/app"
	"firstbeegoapi/internal/shared"

	beego "github.com/beego/beego/v2/server/web"
)

type OrderingController struct {
	beego.Controller
	app.OrderingService
}

var orderingService = app.NewOrderingService(nil)

func SetOrderingService(service *app.OrderingService) {
	if service == nil {
		return
	}

	orderingService = service
}

// @Title Get
// @Description find object by objectid
// @Param	objectId		path 	string	true		"the objectid you want to get"
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router /:objectId [get]
func (o *OrderingController) Get() {

	request, err := validateCheckOrderByProductIDRequest(o.Ctx.Input.Param(":objectId"))
	if err != nil {
		shared.ZapLogger("warn", "Ordering Get API Log", "ordering", shared.RequestID(o.Ctx), o.Ctx.Request.URL.Path, o.Ctx.Input.RequestBody, []byte(err.Error()))
		shared.WriteError(o.Ctx, err)
		return
	}

	requestID := shared.RequestID(o.Ctx)
	serviceCtx := shared.WithLogContext(o.Ctx.Request.Context(), shared.LogContext{
		Service:     "ordering",
		RequestID:   requestID,
		URL:         o.Ctx.Request.URL.Path,
		RequestBody: o.Ctx.Input.RequestBody,
	})

	service, err := orderingService.CheckOrderByProductIdAct(serviceCtx, request)
	if err != nil {
		shared.ZapLogger("error", "Ordering Get API Log", "ordering", shared.RequestID(o.Ctx), o.Ctx.Request.URL.Path, o.Ctx.Input.RequestBody, []byte(err.Error()))
		shared.WriteError(o.Ctx, err)
		return
	}

	responseJson, _ := json.Marshal(service)

	shared.ZapLogger("info", "Ordering Get API Log", "ordering", shared.RequestID(o.Ctx), o.Ctx.Request.URL.Path, o.Ctx.Input.RequestBody, responseJson)
	shared.WriteSuccess(o.Ctx, 200, "success get data", service)
}
