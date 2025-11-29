package product

import (
	"context"
)

type ProductUsecase interface {
	CreateProduct(ctx context.Context, cmd CreateProductCmd) error
	GetProduct(ctx context.Context, cmd GetProductCmd) (*Product, error)
	GetProductByStatus(ctx context.Context, cmd GetProductByStatusCmd) ([]*Product, error)
	UpdatePriority(ctx context.Context, cmd UpdatePriorityCmd) error
	DeleteProduct(ctx context.Context, cmd DeleteProductCmd) error
}

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *Product) (uint, error)
	GetProduct(ctx context.Context, OwnerId, productID uint) (*Product, error)
	GetProductByStatus(ctx context.Context, ownerID uint, status string) ([]*Product, error)
	BulkGetProducts(ctx context.Context, ownerID uint, productIDs []uint) ([]*Product, error)
	GetLastProductPriority(ctx context.Context, ownerID uint) (int64, error)
	DeleteProduct(ctx context.Context, ownerID uint, productID uint) error
}
