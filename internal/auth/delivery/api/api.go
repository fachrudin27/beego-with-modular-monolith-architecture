package api

import (
	"encoding/json"
	"io"

	"firstbeegoapi/internal/auth/app"
	"firstbeegoapi/internal/auth/domain"
	"firstbeegoapi/internal/shared"

	beego "github.com/beego/beego/v2/server/web"
)

type AuthController struct {
	beego.Controller
	app.AuthService
}

// @Title Login
// @Description login user and return jwt token
// @Param	body	body	domain.LoginRequest	true	"login payload"
// @Success 200 {object} domain.LoginResponse
// @Failure 400 invalid request
// @router /login [post]
func (a *AuthController) Login() {
	var request domain.LoginRequest
	body := a.Ctx.Input.RequestBody
	if len(body) == 0 && a.Ctx.Request.Body != nil {
		var err error
		body, err = io.ReadAll(a.Ctx.Request.Body)
		if err != nil {
			shared.ZapLogger("error", "Login POST API Log", "auth", shared.RequestID(a.Ctx), a.Ctx.Request.URL.Path, a.Ctx.Input.RequestBody, []byte(err.Error()))
			shared.WriteError(a.Ctx, shared.NewValidationError("invalid_request_body", "request body is invalid"))
			return
		}
	}
	if err := json.Unmarshal(body, &request); err != nil {
		shared.ZapLogger("error", "Login POST API Log", "auth", shared.RequestID(a.Ctx), a.Ctx.Request.URL.Path, a.Ctx.Input.RequestBody, []byte(err.Error()))
		shared.WriteError(a.Ctx, shared.NewValidationError("invalid_request_body", "request body must be valid json"))
		return
	}
	if err := validateLoginRequest(request); err != nil {
		shared.ZapLogger("info", "Login POST API Log", "auth", shared.RequestID(a.Ctx), a.Ctx.Request.URL.Path, a.Ctx.Input.RequestBody, []byte(err.Error()))
		shared.WriteError(a.Ctx, err)
		return
	}

	requestID := shared.RequestID(a.Ctx)
	serviceCtx := shared.WithLogContext(a.Ctx.Request.Context(), shared.LogContext{
		Service:     "auth",
		RequestID:   requestID,
		URL:         a.Ctx.Request.URL.Path,
		RequestBody: a.Ctx.Input.RequestBody,
	})

	response, err := a.AuthService.LoginAct(serviceCtx, request)
	if err != nil {
		shared.ZapLogger("info", "Login POST API Log", "auth", shared.RequestID(a.Ctx), a.Ctx.Request.URL.Path, a.Ctx.Input.RequestBody, []byte(err.Error()))
		shared.WriteError(a.Ctx, err)
		return
	}

	responseJson, _ := json.Marshal(response)
	shared.ZapLogger("info", "Login POST API Log", "auth", shared.RequestID(a.Ctx), a.Ctx.Request.URL.Path, a.Ctx.Input.RequestBody, []byte(responseJson))
	shared.WriteSuccess(a.Ctx, 200, "login success", response)
}
