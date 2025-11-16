package product

import (
	"time"
)

// model
type Sorting struct {
	Field     string `json:"field" validate:"required,oneof=title price created_at"`
	Direction string `json:"direction" validate:"required,oneof=asc desc"`
}

// request
type QueryProductRequest struct {
	Start  *time.Time `json:"start" validate:"omitempty"`
	End    *time.Time `json:"end" validate:"omitempty,date_after_opt=Start"`
	Status *string    `json:"status" validate:"omitempty,oneof=staging valid bought"`
	Sort   *Sorting   `json:"sort" validate:"omitempty,dive"`
}

type CreateProductRequest struct {
	Title   string   `json:"title" validate:"required"`
	Price   float64  `json:"price" validate:"min=1"`
	Link    *string  `json:"link" validate:"omitempty,url"`
	Reasons []string `json:"reasons" validate:"omitempty,dive,required"`
}

type UpdateCauseStatusRequest struct {
	ProductID string `json:"productId" validate:"required"`
	CauseID   string `json:"causeId" validate:"required"`
	Status    bool   `json:"status" validate:"required"`
}
