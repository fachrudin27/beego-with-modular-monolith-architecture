package shared

import (
	"errors"
	"net/http"
	"time"

	beecontext "github.com/beego/beego/v2/server/web/context"
)

type ErrorKind string

const (
	ErrorKindValidation   ErrorKind = "validation"
	ErrorKindNotFound     ErrorKind = "not_found"
	ErrorKindConflict     ErrorKind = "conflict"
	ErrorKindUnauthorized ErrorKind = "unauthorized"
	ErrorKindInternal     ErrorKind = "internal"
	ErrorManyRequest      ErrorKind = "too_many_request"
)

type AppError struct {
	Kind    ErrorKind `json:"kind"`
	Code    string    `json:"code"`
	Message string    `json:"message"`
	Err     error     `json:"-"`
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}

	if e.Err != nil {
		return e.Err.Error()
	}

	return e.Message
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Err
}

func NewAppError(kind ErrorKind, code string, message string, err error) *AppError {
	return &AppError{
		Kind:    kind,
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func NewManyRequestError(code string, message string) *AppError {
	return NewAppError(ErrorManyRequest, code, message, nil)
}

func NewValidationError(code string, message string) *AppError {
	return NewAppError(ErrorKindValidation, code, message, nil)
}

func NewNotFoundError(code string, message string) *AppError {
	return NewAppError(ErrorKindNotFound, code, message, nil)
}

func NewConflictError(code string, message string) *AppError {
	return NewAppError(ErrorKindConflict, code, message, nil)
}

func NewUnauthorizedError(code string, message string) *AppError {
	return NewAppError(ErrorKindUnauthorized, code, message, nil)
}

func NewInternalError(code string, message string, err error) *AppError {
	return NewAppError(ErrorKindInternal, code, message, err)
}

type APIResponse[T any] struct {
	Message   string     `json:"message"`
	Data      T          `json:"data,omitempty"`
	Error     *ErrorBody `json:"error,omitempty"`
	RequestId string     `json:"request_id"`
	CreateIn  int64      `json:"create_in"`
}

type ErrorBody struct {
	Kind    ErrorKind `json:"kind"`
	Code    string    `json:"code"`
	Message string    `json:"message"`
}

func WriteSuccess[T any](ctx *beecontext.Context, statusCode int, message string, data T) {
	WriteJSON(ctx, statusCode, APIResponse[T]{
		Message:   message,
		Data:      data,
		RequestId: RequestID(ctx),
		CreateIn:  time.Now().Unix(),
	})
}

func WriteError(ctx *beecontext.Context, err error) {
	appErr := ToAppError(err)
	statusCode := StatusCode(appErr)

	WriteJSON(ctx, statusCode, APIResponse[any]{
		Message: "request failed",
		Error: &ErrorBody{
			Kind:    appErr.Kind,
			Code:    appErr.Code,
			Message: appErr.Message,
		},
		RequestId: RequestID(ctx),
		CreateIn:  time.Now().Unix(),
	})
}

func WriteJSON(ctx *beecontext.Context, statusCode int, payload any) {
	ctx.Output.SetStatus(statusCode)
	_ = ctx.Output.JSON(payload, false, false)
}

func ToAppError(err error) *AppError {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	return NewInternalError("internal_error", "internal server error", err)
}

func StatusCode(err *AppError) int {
	if err == nil {
		return http.StatusOK
	}

	switch err.Kind {
	case ErrorKindValidation:
		return http.StatusBadRequest
	case ErrorKindNotFound:
		return http.StatusNotFound
	case ErrorKindConflict:
		return http.StatusConflict
	case ErrorKindUnauthorized:
		return http.StatusUnauthorized
	case ErrorManyRequest:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}
