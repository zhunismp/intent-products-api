package dto

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/zhunismp/intent-products-api/internal/core/domain/shared/apperrors"
)

type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type ValidationErrorResponse struct {
	ErrorMessage string            `json:"errorMessage"`
	ErrorFields  map[string]string `json:"errorFields"`
}

func HandleError(c fiber.Ctx, err error) error {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		return c.Status(apperrors.MapToHttpCode(appErr.Code)).JSON(ErrorResponse{ErrorMessage: appErr.Message})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{ErrorMessage: "something went wrong"})
}

func HandleResponse(c fiber.Ctx, status int, message string, data any) error {
	return c.Status(status).JSON(
		SuccessResponse{
			Message: message,
			Data:    data,
		},
	)
}
