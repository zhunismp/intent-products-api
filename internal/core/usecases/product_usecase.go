package usecases

import (
	"context"

	"github.com/zhunismp/intent-products-api/internal/common/errors"
	"github.com/zhunismp/intent-products-api/internal/core/dtos"
)

type ProductUsecase interface {
	CreateProduct(context.Context, dtos.CreateProductInput) *errors.AppError
}