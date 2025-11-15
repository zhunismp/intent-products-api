package apperrors

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

const (
	ErrCodeNotFound     = "NOT_FOUND"
	ErrCodeValidation   = "VALIDATION_ERROR"
	ErrCodeUnauthorized = "UNAUTHORIZED"
	ErrCodeForbidden    = "FORBIDDEN"
	ErrCodeInternal     = "INTERNAL_ERROR"
)

func MapToHttpCode(code string) int {
	switch code {
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeValidation:
		return http.StatusUnprocessableEntity
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func MapToGrpcStatus(code string) codes.Code {
	switch code {
	case ErrCodeNotFound:
		return codes.NotFound
	case ErrCodeValidation:
		return codes.InvalidArgument
	case ErrCodeUnauthorized:
		return codes.Unauthenticated
	case ErrCodeForbidden:
		return codes.PermissionDenied
	default:
		return codes.Internal
	}
}
