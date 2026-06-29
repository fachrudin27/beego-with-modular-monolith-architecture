package shared

import (
	"errors"
	"os"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	beecontext "github.com/beego/beego/v2/server/web/context"
	"github.com/golang-jwt/jwt/v5"
)

const (
	authUserContextKey = "auth_user"
	defaultJWTSecret   = "change-this-secret"
	defaultJWTIssuer   = "firstbeegoapi"
)

type JWTUser struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type JWTClaims struct {
	User JWTUser `json:"user"`
	jwt.RegisteredClaims
}

func GenerateJWT(user JWTUser, expiresAt time.Time) (string, error) {
	claims := JWTClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    JWTIssuer(),
			Subject:   user.Email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWTSecret()))
}

func ParseJWT(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(JWTSecret()), nil
	}, jwt.WithIssuer(JWTIssuer()))
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func JWTAuthMiddleware(ctx *beecontext.Context) {
	authHeader := ctx.Input.Header("Authorization")
	if authHeader == "" {
		WriteError(ctx, NewUnauthorizedError("missing_authorization", "authorization header is required"))
		return
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		WriteError(ctx, NewUnauthorizedError("invalid_authorization", "authorization header must use Bearer token"))
		return
	}

	tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
	if tokenString == "" {
		WriteError(ctx, NewUnauthorizedError("missing_token", "bearer token is required"))
		return
	}

	claims, err := ParseJWT(tokenString)
	if err != nil {
		WriteError(ctx, NewUnauthorizedError("invalid_token", "token is invalid or expired"))
		return
	}

	ctx.Input.SetData(authUserContextKey, claims.User)
}

func CurrentJWTUser(ctx *beecontext.Context) (JWTUser, bool) {
	user, ok := ctx.Input.GetData(authUserContextKey).(JWTUser)
	return user, ok
}

func JWTSecret() string {
	if secret := strings.TrimSpace(os.Getenv("JWT_SECRET")); secret != "" {
		return secret
	}

	return beego.AppConfig.DefaultString("jwtsecret", defaultJWTSecret)
}

func ValidateJWTConfig(runMode string) error {
	if runMode == "dev" {
		return nil
	}

	if JWTSecret() == defaultJWTSecret {
		return errors.New("jwtsecret must be configured for non-dev runmode")
	}

	return nil
}

func JWTIssuer() string {
	if issuer := strings.TrimSpace(os.Getenv("JWT_ISSUER")); issuer != "" {
		return issuer
	}

	return beego.AppConfig.DefaultString("jwtissuer", defaultJWTIssuer)
}

func JWTExpiresIn() time.Duration {
	seconds := beego.AppConfig.DefaultInt64("jwtexpiresin", 3600)
	if seconds <= 0 {
		seconds = 3600
	}

	return time.Duration(seconds) * time.Second
}
