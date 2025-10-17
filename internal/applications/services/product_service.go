package services

import (
	"context"

	"github.com/zhunismp/intent-products-api/internal/applications/repositories"
	"github.com/zhunismp/intent-products-api/internal/applications/validators"
	"github.com/zhunismp/intent-products-api/internal/common/errors"

	"github.com/zhunismp/intent-products-api/internal/core/dtos"
	"github.com/zhunismp/intent-products-api/internal/core/usecases"
)

type ProductService struct {
	productRepo repositories.ProductRepository
}

func NewProductService(productRepo repositories.ProductRepository) usecases.ProductUsecase {
	return &ProductService{
		productRepo: productRepo,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, createProductReq dtos.CreateProductInput) *errors.AppError {
	
	// validate request
	if err := validators.ValidateCreateProductReq(createProductReq); err != nil {
		return errors.BadRequest("validation failed", err)
	}

	// transform to core model

	return nil
}