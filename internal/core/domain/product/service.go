package product

import (
	"context"
	"time"

	"github.com/zhunismp/intent-products-api/internal/applications/utils"
)

type productService struct {
	productRepo ProductRepository
}

func NewProductService(productRepo ProductRepository) ProductUsecase {
	return &productService{productRepo: productRepo}
}

func (s *productService) CreateProduct(ctx context.Context, cmd CreateProductCmd) (*Product, error) {
	causes := make([]Cause, len(cmd.Reasons))
	for i, reason := range cmd.Reasons {
		causes[i] = Cause{
			Reason: reason,
			Status: true,
		}
	}

	currTime := time.Now()

	product := Product{
		ID:        utils.GenULID(currTime),
		OwnerID:   cmd.OwnerID,
		Name:      cmd.Title,
		ImageUrl:  nil,
		Link:      cmd.Link,
		Price:     cmd.Price,
		AddedAt:   currTime,
		UpdatedAt: currTime,
		Status:    STAGING,
		Causes:    causes,
	}

	return s.productRepo.CreateProduct(ctx, product)
}

func (s *productService) QueryProduct(ctx context.Context, cmd QueryProductCmd) ([]Product, error) {
	// TODO: remove hardcoded sorting
	spec := QueryProductSpec{
		OwnerID: cmd.OwnerID,
		Start:   cmd.Filters.Start,
		End:     cmd.Filters.End,
		Status:  cmd.Filters.Status,
		Sort: &Sorting{
			SortField:     "added_at",
			SortDirection: -1,
		},
	}

	return s.productRepo.QueryProduct(ctx, spec)
}

func (s *productService) DeleteProduct(ctx context.Context, cmd DeleteProductCmd) error {
	return s.productRepo.DeleteProduct(ctx, cmd.OwnerID, cmd.ProductID)
}
