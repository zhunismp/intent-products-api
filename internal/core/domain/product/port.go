package product

import (
	"context"
)

type ProductUsecase interface {
	CreateProduct(context.Context, CreateProductCmd) (*Product, error)
	QueryProducts(context.Context, QueryProductCmd) ([]Product, error)
	GetProduct(context.Context, GetProductCmd) (*Product, error)
	DeleteProduct(context.Context, DeleteProductCmd) error
	UpdateCauseStatus(context.Context, UpdateCauseStatusCmd) (*Cause, error)

	BatchGetProduct(context.Context, BatchGetProductCmd) ([]Product, error)
}

type ProductRepository interface {
	CreateProduct(context.Context, Product) (*Product, error)
	QueryProduct(context.Context, QueryProductSpec) ([]Product, error)
	GetProduct(context.Context, OwnerId, ProductId) (*Product, error)
	DeleteProduct(context.Context, OwnerId, ProductId) error

	BatchGetProduct(context.Context, OwnerId, []ProductId) ([]Product, error)
}

type CauseRepository interface {
	CreateCauses(context.Context, ProductId, []Cause) ([]Cause, error)
	GetCauses(context.Context, ProductId) ([]Cause, error)
	UpdateCauseStatus(context.Context, ProductId, CauseStatus) (*Cause, error)
	DeleteCausesByProductID(context.Context, ProductId) error
}
