package product

import (
	"context"

	"github.com/zhunismp/intent-products-api/internal/core/domain/cause"
	"github.com/zhunismp/intent-products-api/internal/core/domain/priority"
	"github.com/zhunismp/intent-products-api/internal/core/domain/shared/apperrors"
)

type productService struct {
	productRepo ProductRepository
	causeSvc    cause.CauseUsecase
	prioritySvc priority.PriorityUsecase
}

func NewProductService(productRepo ProductRepository, causeSvc cause.CauseUsecase, prioritySvc priority.PriorityUsecase) ProductUsecase {
	return &productService{
		productRepo: productRepo,
		causeSvc:    causeSvc,
		prioritySvc: prioritySvc,
	}
}

func (s *productService) CreateProduct(ctx context.Context, cmd CreateProductCmd) error {

	product := &Product{
		OwnerID:  cmd.OwnerID,
		Name:     cmd.Title,
		ImageUrl: nil,
		Link:     cmd.Link,
		Price:    cmd.Price,
		Status:   PENDING,
	}

	lastPriority, err := s.productRepo.GetLastProductPriority(ctx, cmd.OwnerID)
	if err != nil {
		return err
	}

	newPriority, err := s.prioritySvc.CalculateNewPriority(ctx, nil, &lastPriority)
	if err != nil {
		return err
	}
	product.Priority = newPriority

	productID, err := s.productRepo.CreateProduct(ctx, product)
	if err != nil {
		return err
	}

	if err := s.causeSvc.BulkCreateCauses(ctx, productID, cmd.Reasons); err != nil {
		return err
	}

	return nil
}

func (s *productService) GetProduct(ctx context.Context, cmd GetProductCmd) (*Product, error) {
	product, err := s.productRepo.GetProduct(ctx, cmd.OwnerID, cmd.ProductID)
	if err != nil {
		return nil, err
	}

	causes, err := s.causeSvc.GetCauses(ctx, product.ID)
	if err != nil {
		return nil, err
	}

	product.Causes = causes
	return product, nil
}

func (s *productService) GetProductByStatus(ctx context.Context, cmd GetProductByStatusCmd) ([]*Product, error) {
	return s.productRepo.GetProductByStatus(ctx, cmd.OwnerID, cmd.Status)
}

func (s *productService) UpdatePriority(ctx context.Context, cmd UpdatePriorityCmd) error {
	if cmd.ProductIDBefore == nil && cmd.ProductIDAfter == nil {
		return apperrors.New(apperrors.ErrCodeValidation, "before or after product must be set", nil)
	}

	ids := []uint{cmd.ProductID}
	if cmd.ProductIDBefore != nil {
		ids = append(ids, *cmd.ProductIDBefore)
	}
	if cmd.ProductIDAfter != nil {
		ids = append(ids, *cmd.ProductIDAfter)
	}

	products, err := s.productRepo.BulkGetProducts(ctx, cmd.OwnerID, ids)
	if err != nil {
		return err
	}

	productMap := make(map[uint]*Product, len(products))
	for i := range products {
		productMap[products[i].ID] = products[i]
	}

	product, ok := productMap[cmd.ProductID]
	if !ok {
		return apperrors.New(apperrors.ErrCodeNotFound, "product not found", nil)
	}

	var beforePriority, afterPriority *int64
	if cmd.ProductIDBefore != nil {
		if p, ok := productMap[*cmd.ProductIDBefore]; ok {
			beforePriority = &p.Priority
		} else {
			return apperrors.New(apperrors.ErrCodeNotFound, "before product not found", nil)
		}
	}
	if cmd.ProductIDAfter != nil {
		if p, ok := productMap[*cmd.ProductIDAfter]; ok {
			afterPriority = &p.Priority
		} else {
			return apperrors.New(apperrors.ErrCodeNotFound, "after product not found", nil)
		}
	}

	newPriority, err := s.prioritySvc.CalculateNewPriority(ctx, beforePriority, afterPriority)
	if err != nil {
		return err
	}

	product.Priority = newPriority
	_, err = s.productRepo.CreateProduct(ctx, product)
	return err
}

func (s *productService) DeleteProduct(ctx context.Context, cmd DeleteProductCmd) error {
	if err := s.productRepo.DeleteProduct(ctx, cmd.OwnerID, cmd.ProductID); err != nil {
		return err
	}
	if err := s.causeSvc.DeleteCauses(ctx, cmd.ProductID); err != nil {
		return err
	}

	return nil
}
