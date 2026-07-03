package api

import (
	"net/mail"
	"strings"

	"firstbeegoapi/internal/auth/domain"
	"firstbeegoapi/internal/shared"
)

func validateLoginRequest(request domain.LoginRequest) error {
	email := strings.TrimSpace(request.Email)
	if email == "" {
		return shared.NewValidationError("missing_email", "email is required")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return shared.NewValidationError("invalid_email", "email must be valid")
	}
	if request.Password == "" {
		return shared.NewValidationError("missing_password", "password is required")
	}

	return nil
}
