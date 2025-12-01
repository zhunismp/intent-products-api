package product

import (
	"context"
)

type ProductUsecase interface {
	CreateProduct(ctx context.Context, ownerID uint, title string, price float64, link string, reasons []string) error
	GetProduct(ctx context.Context, ownerID uint, productID uint) (*Product, error)
	GetProductByStatus(ctx context.Context, ownerID uint, status string) ([]*Product, error)
	Move(ctx context.Context, ownerID uint, productID uint, productAfterID *uint) error
	DeleteProduct(ctx context.Context, ownerID uint, productID uint) error
}

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *Product) (uint, error)
	GetProduct(ctx context.Context, ownerID uint, productID uint) (*Product, error)
	GetProductByStatus(ctx context.Context, ownerID uint, status string) ([]*Product, error)
	DeleteProduct(ctx context.Context, ownerID uint, productID uint) error

	GetFirstPosition(ctx context.Context, ownerID uint) (string, error)
	GetPositionByProductID(ctx context.Context, ownerID uint, productID uint) (string, error)
	GetNextPosition(ctx context.Context, ownerID uint, position string) (string, error)
	UpdatePosition(ctx context.Context, ownerID uint, productID uint, position string) error
}
