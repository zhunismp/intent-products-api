package usecases

import (
	"context"

	"github.com/zhunismp/intent-products-api/internal/core/dtos"
	"github.com/zhunismp/intent-products-api/internal/core/entities"
)

type ProductUsecase interface {
	CreateProduct(context.Context, dtos.CreateProductInput) (*entities.Product, error)
	QueryProduct(context.Context, dtos.QueryProductInput) ([]entities.Product, error)
	DeleteProduct(context.Context, dtos.DeleteProductInput) error
}
