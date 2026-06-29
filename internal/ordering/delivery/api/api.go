package api

import (
	"encoding/json"
	"firstbeegoapi/internal/ordering/app"
	"firstbeegoapi/internal/ordering/domain"
	"firstbeegoapi/internal/shared"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
)

type OrderingController struct {
	beego.Controller
	app.OrderingService
}

// @Title Get
// @Description find object by objectid
// @Param	objectId		path 	string	true		"the objectid you want to get"
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router /:objectId [get]
func (o *OrderingController) Get() {

	objectId := o.Ctx.Input.Param(":objectId")
	if objectId == "" {
		shared.ZapLogger("warn", "Ordering Get API Log", "ordering", "/api", shared.RequestID(o.Ctx), o.Ctx.Request.URL.Path, o.Ctx.Input.RequestBody, []byte("object id is required"))
		shared.WriteError(o.Ctx, shared.NewValidationError("missing_object_id", "object id is required"))
		return
	}

	id, err := strconv.Atoi(objectId)
	if err != nil {
		shared.ZapLogger("warn", "Ordering Get API Log", "ordering", "/api", shared.RequestID(o.Ctx), o.Ctx.Request.URL.Path, o.Ctx.Input.RequestBody, []byte(err.Error()))
		shared.WriteError(o.Ctx, shared.NewValidationError("invalid_object_id", "object id must be a number"))
		return
	}

	requestID := shared.RequestID(o.Ctx)
	serviceCtx := shared.WithLogContext(o.Ctx.Request.Context(), shared.LogContext{
		Service:     "ordering",
		Position:    "/api",
		RequestID:   requestID,
		URL:         o.Ctx.Request.URL.Path,
		RequestBody: o.Ctx.Input.RequestBody,
	})

	service, err := o.OrderingService.CheckOrderByProductIdAct(serviceCtx, domain.CheckOrderByProductIdRequest{
		ProductId: int64(id),
	})
	if err != nil {
		shared.ZapLogger("error", "Ordering Get API Log", "ordering", "/api", shared.RequestID(o.Ctx), o.Ctx.Request.URL.Path, o.Ctx.Input.RequestBody, []byte(err.Error()))
		shared.WriteError(o.Ctx, err)
		return
	}

	responseJson, _ := json.Marshal(service)

	shared.ZapLogger("info", "Ordering Get API Log", "ordering", "/api", shared.RequestID(o.Ctx), o.Ctx.Request.URL.Path, o.Ctx.Input.RequestBody, responseJson)
	shared.WriteSuccess(o.Ctx, 200, "success get data", service)
}
