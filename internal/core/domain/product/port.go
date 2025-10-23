package product

import (
	"context"
	"time"
)

type ProductUsecase interface {
	CreateProduct(context.Context, CreateProductCmd) (*Product, error)
	QueryProduct(context.Context, QueryProductCmd) ([]Product, error)
	GetProduct(context.Context, GetProductCmd) (*Product, error)
	DeleteProduct(context.Context, DeleteProductCmd) error
}

type ProductRepository interface {
	CreateProduct(context.Context, Product) (*Product, error)
	QueryProduct(context.Context, QueryProductSpec) ([]Product, error)
	GetProduct(context.Context, string, string) (*Product, error)
	DeleteProduct(context.Context, string, string) error
}

// usecase command
type CreateProductCmd struct {
	OwnerID string
	Title   string
	Price   float64
	Link    *string
	Reasons []string
}

type DeleteProductCmd struct {
	OwnerID   string
	ProductID string
}

type GetProductCmd struct {
	OwnerID   string
	ProductID string
}

type QueryProductCmd struct {
	OwnerID string
	Filters QueryProductFilter
}

type QueryProductFilter struct {
	Start  *time.Time
	End    *time.Time
	Status *string
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
