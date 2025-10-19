package dtos

import "time"

type QueryProductInput struct {
	OwnerID    string             `validate:"required"`
	Pagination *Pagination        `validate:"omitempty"`
	Filters    QueryProductFilter `validate:"required"`
}

type QueryProductFilter struct {
	Start  *time.Time `validate:"omitempty"`
	End    *time.Time `validate:"omitempty"`
	Status *string    `validate:"omitempty,oneof=staging valid bought"`
}

type Pagination struct {
	Page int `validate:"required,gt=0"`
	Size int `validate:"required,gt=0"`
}
