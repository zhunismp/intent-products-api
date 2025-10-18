package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/zhunismp/intent-products-api/internal/applications/repositories"
	"github.com/zhunismp/intent-products-api/internal/applications/validators"
	"github.com/zhunismp/intent-products-api/internal/common/apperrors"

	"github.com/zhunismp/intent-products-api/internal/core/domainerrors"
	"github.com/zhunismp/intent-products-api/internal/core/dtos"
	"github.com/zhunismp/intent-products-api/internal/core/entities"
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

func (s *ProductService) CreateProduct(ctx context.Context, createProductInput dtos.CreateProductInput) *apperrors.AppError {

	// validate request
	if err := validators.ValidateCreateProductReq(createProductInput); err != nil {
		return apperrors.BadRequest("validation failed", err)
	}

	node, err := snowflake.NewNode(1)
	if err != nil {
		return apperrors.Internal("error while get node from snowflake", err)
	}

	// transform to core model
	causes := make([]entities.Cause, len(createProductInput.Reasons))
	for i, reason := range createProductInput.Reasons {
		causes[i] = entities.Cause{
			Reason: reason,
			Status: true,
		}
	}

	product := entities.Product{
		ID:        node.Generate().Int64(),
		OwnerID:   1,
		Name:      createProductInput.Title,
		ImageUrl:  nil,
		Link:      createProductInput.Link,
		Price:     createProductInput.Price,
		AddedAt:   time.Now(),
		UpdatedAt: time.Now(),
		Status:    entities.STAGING,
		Causes:    causes,
	}

	if err := s.productRepo.CreateProduct(ctx, product); err != nil {
		log.Print(err.Error())
		if errors.Is(err, domainerrors.ErrorDuplicateProduct) {
			return apperrors.BadRequest("can not add duplicate product", err)
		}
		return apperrors.Internal("error occur while saving to db", err)
	}

	return nil
}
