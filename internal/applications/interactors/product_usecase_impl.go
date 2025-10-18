package interactors

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/zhunismp/intent-products-api/internal/applications/validators"
	"github.com/zhunismp/intent-products-api/internal/core/repositories"

	"github.com/zhunismp/intent-products-api/internal/core/domainerrors"
	"github.com/zhunismp/intent-products-api/internal/core/dtos"
	"github.com/zhunismp/intent-products-api/internal/core/entities"
	"github.com/zhunismp/intent-products-api/internal/core/usecases"
)

type ProductUsecaseImpl struct {
	productRepo repositories.ProductRepository
}

func NewProductUsecase(productRepo repositories.ProductRepository) usecases.ProductUsecase {
	return &ProductUsecaseImpl{
		productRepo: productRepo,
	}
}

func (s *ProductUsecaseImpl) CreateProduct(ctx context.Context, createProductInput dtos.CreateProductInput) (*entities.Product, error) {

	// validate request
	if err := validators.ValidateCreateProductReq(createProductInput); err != nil {
		return nil, domainerrors.ErrorProducInput
	}

	node, err := snowflake.NewNode(1)
	if err != nil {
		return nil, fmt.Errorf("error while get node from snowflake: %v", err)
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

	createdProdcut, err := s.productRepo.CreateProduct(ctx, product) 
	if err != nil {
		if errors.Is(err, domainerrors.ErrorDuplicateProduct) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to save product to database: %w", err)
	}

	return createdProdcut, nil
}
