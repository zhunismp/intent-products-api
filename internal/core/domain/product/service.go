package product

import (
	"context"
	"time"

	"github.com/zhunismp/intent-products-api/internal/core/domain/shared/utils"
)

type productService struct {
	productRepo ProductRepository
	causeRepo   CauseRepository
}

func NewProductService(productRepo ProductRepository, causeRepo CauseRepository) ProductUsecase {
	return &productService{
		productRepo: productRepo,
		causeRepo:   causeRepo,
	}
}

func (s *productService) CreateProduct(ctx context.Context, cmd CreateProductCmd) (*Product, error) {
	currTime := time.Now()
	productID := utils.GenULID(currTime)

	product := Product{
		ID:        productID,
		OwnerID:   cmd.OwnerID,
		Name:      cmd.Title,
		ImageUrl:  nil,
		Link:      cmd.Link,
		Price:     cmd.Price,
		CreatedAt: currTime,
		UpdatedAt: currTime,
		Status:    PENDING,
	}

	causes := make([]Cause, len(cmd.Reasons))
	for i, reason := range cmd.Reasons {
		causes[i] = Cause{
			ID:     utils.GenULID(time.Now()),
			Reason: reason,
			Status: true,
		}
	}

	createdProduct, err := s.productRepo.CreateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	createdCause, err := s.causeRepo.CreateCauses(ctx, productID, causes)
	if err != nil {
		return nil, err
	}

	createdProduct.Causes = createdCause
	return createdProduct, nil
}

func (s *productService) QueryProducts(ctx context.Context, cmd QueryProductCmd) ([]Product, error) {
	// default sorting field
	sort := Sort{
		Field:     "created_at",
		Direction: "asc",
	}

	// apply sorting, if applicable.
	if cmd.Sort != nil {
		sort.Field = cmd.Sort.Field
		sort.Direction = cmd.Sort.Direction
	}

	spec := QueryProductSpec{
		OwnerID: cmd.OwnerID,
		Start:   cmd.Start,
		End:     cmd.End,
		Status:  cmd.Status,
		Sort:    sort,
	}

	return s.productRepo.QueryProduct(ctx, spec)
}

func (s *productService) GetProduct(ctx context.Context, cmd GetProductCmd) (*Product, error) {
	product, err := s.productRepo.GetProduct(ctx, cmd.OwnerID, cmd.ProductID)
	if err != nil {
		return nil, err
	}

	causes, err := s.causeRepo.GetCauses(ctx, product.ID)
	if err != nil {
		return nil, err
	}

	product.Causes = causes
	return product, nil
}

func (s *productService) UpdateCauseStatus(ctx context.Context, cmd UpdateCauseStatusCmd) (*Cause, error) {
	product, err := s.productRepo.GetProduct(ctx, cmd.OwnerID, cmd.ProductID)
	if err != nil {
		return nil, err
	}

	updatedCause, err := s.causeRepo.UpdateCauseStatus(ctx, product.ID, cmd.CauseStatus)
	if err != nil {
		return nil, err
	}

	return updatedCause, nil
}

func (s *productService) DeleteProduct(ctx context.Context, cmd DeleteProductCmd) error {
	if err := s.productRepo.DeleteProduct(ctx, cmd.OwnerID, cmd.ProductID); err != nil {
		return err
	}
	if err := s.causeRepo.DeleteCausesByProductID(ctx, cmd.ProductID); err != nil {
		return err
	}

	return nil
}
