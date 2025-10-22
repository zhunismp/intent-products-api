package product

import (
	"context"
	"time"
)

type ProductUsecase interface {
	CreateProduct(context.Context, CreateProductCmd) (*Product, error)
	QueryProduct(context.Context, QueryProductCmd) ([]Product, error)
	DeleteProduct(context.Context, DeleteProductCmd) error
}

type ProductRepository interface {
	CreateProduct(context.Context, Product) (*Product, error)
	QueryProduct(context.Context, QueryProductSpec) ([]Product, error)
	DeleteProduct(context.Context, string, string) error
}

// usecase command
type CreateProductCmd struct {
	OwnerID string   `validate:"required"`
	Title   string   `validate:"required"`
	Price   float64  `validate:"required,gt=0"`
	Link    *string  `validate:"omitempty,url"`
	Reasons []string `validate:"omitempty,dive"`
}

type DeleteProductCmd struct {
	OwnerID   string `validate:"required"`
	ProductID string `validate:"required"`
}

type QueryProductCmd struct {
	OwnerID string             `validate:"required"`
	Filters QueryProductFilter `validate:"required"`
}

type QueryProductFilter struct {
	Start  *time.Time `validate:"omitempty"`
	End    *time.Time `validate:"omitempty"`
	Status *string    `validate:"omitempty,oneof=staging valid bought"`
}

// repo query
type QueryProductSpec struct {
	OwnerID string
	Start   *time.Time
	End     *time.Time
	Status  *string
	Sort    *Sorting
}

type Sorting struct {
	SortField     string
	SortDirection int
}
