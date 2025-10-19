package interactors

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zhunismp/intent-products-api/internal/applications/utils"
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

func NewProductUsecaseImpl(productRepo repositories.ProductRepository) usecases.ProductUsecase {
	return &ProductUsecaseImpl{
		productRepo: productRepo,
	}
}

func (s *ProductUsecaseImpl) CreateProduct(ctx context.Context, createProductInput dtos.CreateProductInput) (*entities.Product, error) {

	// validate request
	if err := validators.ValidateCreateProductReq(createProductInput); err != nil {
		return nil, domainerrors.ErrorProducInput
	}

	// transform to core model
	causes := make([]entities.Cause, len(createProductInput.Reasons))
	for i, reason := range createProductInput.Reasons {
		causes[i] = entities.Cause{
			Reason: reason,
			Status: true,
		}
	}

	currTime := time.Now()

	product := entities.Product{
		ID:        utils.GenULID(time.Now()),
		OwnerID:   "101234567890123456789", // hardcoded google oauth's sub
		Name:      createProductInput.Title,
		ImageUrl:  nil,
		Link:      createProductInput.Link,
		Price:     createProductInput.Price,
		AddedAt:   currTime,
		UpdatedAt: currTime,
		Status:    entities.STAGING,
		Causes:    causes,
	}

	createdProdcut, err := s.productRepo.CreateProduct(ctx, product)
	if err != nil {
		if errors.Is(err, domainerrors.ErrorDuplicateProduct) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to save product to database: %v", err)
	}

	return createdProdcut, nil
}
