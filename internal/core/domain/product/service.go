package product

import (
	"context"

	"github.com/zhunismp/intent-products-api/internal/core/domain/cause"
	"github.com/zhunismp/intent-products-api/internal/core/domain/priority"
	"github.com/zhunismp/intent-products-api/internal/core/domain/shared/apperrors"
	"go.uber.org/zap"
)

type productService struct {
	productRepo ProductRepository
	causeSvc    cause.CauseUsecase
	prioritySvc priority.PriorityUsecase
	logger      *zap.Logger
}

func NewProductService(
	productRepo ProductRepository,
	causeSvc cause.CauseUsecase,
	prioritySvc priority.PriorityUsecase,
	logger *zap.Logger,
) ProductUsecase {
	return &productService{
		productRepo: productRepo,
		causeSvc:    causeSvc,
		prioritySvc: prioritySvc,
		logger: logger,
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
		s.logger.Error("failed to get last product priority",
			zap.Uint("owner_id", cmd.OwnerID),
			zap.Error(err),
		)
		return err
	}

	newPriority, err := s.prioritySvc.CalculateNewPriority(ctx, nil, &lastPriority)
	if err != nil {
		s.logger.Error("failed to calculate new priority",
			zap.Uint("owner_id", cmd.OwnerID),
			zap.Int64("last_priority", lastPriority),
			zap.Error(err),
		)
		return err
	}
	product.Priority = newPriority

	productID, err := s.productRepo.CreateProduct(ctx, product)
	if err != nil {
		s.logger.Error("failed to create product in repository",
			zap.Uint("owner_id", cmd.OwnerID),
			zap.String("title", cmd.Title),
			zap.Error(err),
		)
		return err
	}

	if err := s.causeSvc.BulkCreateCauses(ctx, productID, cmd.Reasons); err != nil {
		s.logger.Error("failed to create causes for product",
			zap.Uint("product_id", productID),
			zap.Int("causes_count", len(cmd.Reasons)),
			zap.Error(err),
		)
		return err
	}

	s.logger.Info("product created successfully",
		zap.Uint("product_id", productID),
		zap.Uint("owner_id", cmd.OwnerID),
		zap.Int("causes_count", len(cmd.Reasons)),
	)

	return nil
}

func (s *productService) GetProduct(ctx context.Context, cmd GetProductCmd) (*Product, error) {

	product, err := s.productRepo.GetProduct(ctx, cmd.OwnerID, cmd.ProductID)
	if err != nil {
		s.logger.Error("failed to get product",
			zap.Uint("owner_id", cmd.OwnerID),
			zap.Uint("product_id", cmd.ProductID),
			zap.Error(err),
		)
		return nil, err
	}

	causes, err := s.causeSvc.GetCauses(ctx, product.ID)
	if err != nil {
		s.logger.Error("failed to get causes for product",
			zap.Uint("product_id", product.ID),
			zap.Error(err),
		)
		return nil, err
	}

	product.Causes = causes

	s.logger.Info("product fetched successfully",
		zap.Uint("product_id", product.ID),
		zap.Int("causes_count", len(causes)),
	)

	return product, nil
}

func (s *productService) GetProductByStatus(ctx context.Context, cmd GetProductByStatusCmd) ([]*Product, error) {

	products, err := s.productRepo.GetProductByStatus(ctx, cmd.OwnerID, cmd.Status)
	if err != nil {
		s.logger.Error("failed to get products by status",
			zap.Uint("owner_id", cmd.OwnerID),
			zap.String("status", string(cmd.Status)),
			zap.Error(err),
		)
		return nil, err
	}

	s.logger.Info("products fetched by status",
		zap.Uint("owner_id", cmd.OwnerID),
		zap.String("status", string(cmd.Status)),
		zap.Int("count", len(products)),
	)

	return products, nil
}

func (s *productService) UpdatePriority(ctx context.Context, cmd UpdatePriorityCmd) error {

	if cmd.ProductIDBefore == nil && cmd.ProductIDAfter == nil {
		s.logger.Error("invalid priority update request - missing before/after reference",
			zap.Uint("product_id", cmd.ProductID),
		)
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
		s.logger.Error("failed to bulk get products",
			zap.Uint("owner_id", cmd.OwnerID),
			zap.Uints("product_ids", ids),
			zap.Error(err),
		)
		return err
	}

	productMap := make(map[uint]*Product, len(products))
	for i := range products {
		productMap[products[i].ID] = products[i]
	}

	product, ok := productMap[cmd.ProductID]
	if !ok {
		s.logger.Error("product not found for priority update",
			zap.Uint("product_id", cmd.ProductID),
			zap.Uint("owner_id", cmd.OwnerID),
		)
		return apperrors.New(apperrors.ErrCodeNotFound, "product not found", nil)
	}

	var beforePriority, afterPriority *int64
	if cmd.ProductIDBefore != nil {
		if p, ok := productMap[*cmd.ProductIDBefore]; ok {
			beforePriority = &p.Priority
		} else {
			s.logger.Error("before product not found",
				zap.Uint("product_id_before", *cmd.ProductIDBefore),
			)
			return apperrors.New(apperrors.ErrCodeNotFound, "before product not found", nil)
		}
	}
	if cmd.ProductIDAfter != nil {
		if p, ok := productMap[*cmd.ProductIDAfter]; ok {
			afterPriority = &p.Priority
		} else {
			s.logger.Error("after product not found",
				zap.Uint("product_id_after", *cmd.ProductIDAfter),
			)
			return apperrors.New(apperrors.ErrCodeNotFound, "after product not found", nil)
		}
	}

	newPriority, err := s.prioritySvc.CalculateNewPriority(ctx, beforePriority, afterPriority)
	if err != nil {
		s.logger.Error("failed to calculate new priority",
			zap.Uint("product_id", cmd.ProductID),
			zap.Int64p("before_priority", beforePriority),
			zap.Int64p("after_priority", afterPriority),
			zap.Error(err),
		)
		return err
	}

	oldPriority := product.Priority
	product.Priority = newPriority

	_, err = s.productRepo.CreateProduct(ctx, product)
	if err != nil {
		s.logger.Error("failed to update product priority",
			zap.Uint("product_id", cmd.ProductID),
			zap.Int64("old_priority", oldPriority),
			zap.Int64("new_priority", newPriority),
			zap.Error(err),
		)
		return err
	}

	s.logger.Info("product priority updated successfully",
		zap.Uint("product_id", cmd.ProductID),
		zap.Int64("old_priority", oldPriority),
		zap.Int64("new_priority", newPriority),
	)

	return nil
}

func (s *productService) DeleteProduct(ctx context.Context, cmd DeleteProductCmd) error {

	if err := s.productRepo.DeleteProduct(ctx, cmd.OwnerID, cmd.ProductID); err != nil {
		s.logger.Error("failed to delete product from repository",
			zap.Uint("owner_id", cmd.OwnerID),
			zap.Uint("product_id", cmd.ProductID),
			zap.Error(err),
		)
		return err
	}

	if err := s.causeSvc.DeleteCauses(ctx, cmd.ProductID); err != nil {
		s.logger.Error("failed to delete causes for product",
			zap.Uint("product_id", cmd.ProductID),
			zap.Error(err),
		)
		return err
	}

	s.logger.Info("product deleted successfully",
		zap.Uint("owner_id", cmd.OwnerID),
		zap.Uint("product_id", cmd.ProductID),
	)

	return nil
}
