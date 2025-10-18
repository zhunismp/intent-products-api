package usecases

import (
	"context"

	"github.com/zhunismp/intent-products-api/internal/common/apperrors"
	"github.com/zhunismp/intent-products-api/internal/core/dtos"
)

type ProductUsecase interface {
	CreateProduct(context.Context, dtos.CreateProductInput) *apperrors.AppError
}