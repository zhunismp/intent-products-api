package services

import (
	"context"

	"github.com/zhunismp/intent-products-api/internal/adapters/http/transport"
	"github.com/zhunismp/intent-products-api/internal/applications/repositories"
	// "github.com/zhunismp/intent-products-api/internal/core/entities"
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

func (s *ProductService) CreateProduct(ctx context.Context, createProductReq transport.CreateProductRequest) error {
	return nil
}