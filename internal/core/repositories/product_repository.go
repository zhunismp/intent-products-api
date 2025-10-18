package repositories

import (
	"context"

	"github.com/zhunismp/intent-products-api/internal/core/entities"
)

type ProductRepository interface {
	CreateProduct(context.Context, entities.Product) (*entities.Product, error)
}