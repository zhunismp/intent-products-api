package dtos

import "time"

type CreateProductInput struct {
	OwnerID string   `validate:"required"`
	Title   string   `validate:"required"`
	Price   float64  `validate:"required,gt=0"`
	Link    *string  `validate:"omitempty,url"`
	Reasons []string `validate:"omitempty,dive"`
}

type QueryProductInput struct {
	OwnerID string             `validate:"required"`
	Filters QueryProductFilter `validate:"required"`
}

type QueryProductFilter struct {
	Start  *time.Time `validate:"omitempty"`
	End    *time.Time `validate:"omitempty"`
	Status *string    `validate:"omitempty,oneof=staging valid bought"`
}
