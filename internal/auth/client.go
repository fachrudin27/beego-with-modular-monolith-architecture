package auth

import (
	"context"

	"firstbeegoapi/internal/auth/app"
	"firstbeegoapi/internal/auth/domain"
)

type AuthContract interface {
	LoginAct(ctx context.Context, request domain.LoginRequest) (domain.LoginResponse, error)
}

type authApi struct {
	app.AuthService
}

func NewModuleApi() AuthContract {
	return &authApi{}
}

func (a *authApi) LoginAct(ctx context.Context, request domain.LoginRequest) (domain.LoginResponse, error) {
	return a.AuthService.LoginAct(ctx, request)
}
