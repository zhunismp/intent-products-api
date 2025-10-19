package repositories

import (
	"context"

	"github.com/zhunismp/intent-products-api/internal/core/dtos"
	"github.com/zhunismp/intent-products-api/internal/core/entities"
)

type ProductRepository interface {
	CreateProduct(context.Context, entities.Product) (*entities.Product, error)

	// This should not return entities object at all.
	// Just for development purpose
	// TODO: return output dto instead of entities
	QueryProduct(context.Context, dtos.QueryProductInput) ([]entities.Product, error)
	
	DeleteProduct(context.Context, dtos.DeleteProductInput) error
}
