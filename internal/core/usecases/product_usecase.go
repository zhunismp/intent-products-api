package usecases

import (
	"context"

	"github.com/zhunismp/intent-products-api/internal/adapters/http/transport"
)

type ProductUsecase interface {
	CreateProduct(context.Context, transport.CreateProductRequest) error
}