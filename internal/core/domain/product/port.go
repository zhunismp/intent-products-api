package product

import (
	"context"
)

type ProductUsecase interface {
	CreateProduct(context.Context, CreateProductCmd) (*Product, error)
	GetProduct(context.Context, GetProductCmd) (*Product, error)
	GetProductByStatus(context.Context, GetProductByStatusCmd) ([]Product, error)
	DeleteProduct(context.Context, DeleteProductCmd) error
	UpdateCauseStatus(context.Context, UpdateCauseStatusCmd) (*Cause, error)
}

type ProductRepository interface {
	CreateProduct(context.Context, Product) (*Product, error)
	GetProduct(context.Context, OwnerId, ProductId) (*Product, error)
	GetProductByStatus(context.Context, OwnerId, Status) ([]Product, error)
	DeleteProduct(context.Context, OwnerId, ProductId) error
}

type CauseRepository interface {
	CreateCauses(context.Context, ProductId, []Cause) ([]Cause, error)
	GetCauses(context.Context, ProductId) ([]Cause, error)
	UpdateCauseStatus(context.Context, ProductId, CauseStatus) (*Cause, error)
	DeleteCausesByProductID(context.Context, ProductId) error
}
