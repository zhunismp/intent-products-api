package product

import (
	"time"
)

type QueryProductRequest struct {
	Start  *time.Time `json:"start" validate:"omitempty"`
	End    *time.Time `json:"end" validate:"omitempty,date_after_opt=Start"`
	Status *string    `json:"status" validate:"omitempty,oneof=staging valid bought"`
	Sort   *struct {
		Field     string `json:"field" validate:"required,oneof=title price created_at"`
		Direction string `json:"direction" validate:"required,oneof=asc desc"`
	} `json:"sort" validate:"omitempty,dive"`
}

type CreateProductRequest struct {
	Title   string   `json:"title" validate:"required"`
	Price   float64  `json:"price" validate:"min=1"`
	Link    *string  `json:"link" validate:"omitempty,url"`
	Reasons []string `json:"reasons" validate:"omitempty,dive,required"`
}

type UpdateProductRequest struct{}

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
