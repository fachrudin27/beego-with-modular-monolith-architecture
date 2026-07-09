package app

import (
	"context"
	"strings"
	"time"

	"firstbeegoapi/internal/auth/domain"
	"firstbeegoapi/internal/shared"
)

func (a *AuthService) LoginAct(ctx context.Context, request domain.LoginRequest) (domain.LoginResponse, error) {
	email := strings.TrimSpace(strings.ToLower(request.Email))
	if email == "" {
		if logCtx, ok := shared.LogContextFrom(ctx); ok {
			shared.ZapLogger("info", "Login POST API Log", logCtx.Service, logCtx.RequestID, logCtx.URL, logCtx.RequestBody, []byte("email is required"))
		}
		return domain.LoginResponse{}, shared.NewValidationError("missing_email", "email is required")
	}
	if request.Password == "" {
		if logCtx, ok := shared.LogContextFrom(ctx); ok {
			shared.ZapLogger("info", "Login POST API Log", logCtx.Service, logCtx.RequestID, logCtx.URL, logCtx.RequestBody, []byte("password is required"))
		}
		return domain.LoginResponse{}, shared.NewValidationError("missing_password", "password is required")
	}

	user, ok := authenticateUser(email, request.Password)
	if !ok {
		if logCtx, ok := shared.LogContextFrom(ctx); ok {
			shared.ZapLogger("info", "Login POST API Log", logCtx.Service, logCtx.RequestID, logCtx.URL, logCtx.RequestBody, []byte("email or password is invalid"))
		}
		return domain.LoginResponse{}, shared.NewValidationError("invalid_credentials", "email or password is invalid")
	}

	expiresIn := shared.JWTExpiresIn()
	token, err := shared.GenerateJWT(shared.JWTUser{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, time.Now().Add(expiresIn))
	if err != nil {
		if logCtx, ok := shared.LogContextFrom(ctx); ok {
			shared.ZapLogger("error", "Login POST API Log", logCtx.Service, logCtx.RequestID, logCtx.URL, logCtx.RequestBody, []byte(err.Error()))
		}
		return domain.LoginResponse{}, shared.NewInternalError("generate_token_failed", "failed to generate token", err)
	}

	return domain.LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(expiresIn.Seconds()),
		User:        user,
	}, nil
}

func authenticateUser(email string, password string) (domain.User, bool) {
	if email != "admin@example.com" || password != "password" {
		return domain.User{}, false
	}

	return domain.User{
		ID:    1,
		Name:  "Admin",
		Email: email,
	}, true
}
