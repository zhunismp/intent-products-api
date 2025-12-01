package product

import (
	"context"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"github.com/zhunismp/intent-products-api/internal/core/domain/cause"
	"github.com/zhunismp/intent-products-api/internal/core/domain/shared/utils/ordering"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

type productService struct {
	productRepo ProductRepository
	causeSvc    cause.CauseUsecase
	logger      *otelzap.Logger
}

func NewProductService(
	productRepo ProductRepository,
	causeSvc cause.CauseUsecase,
	logger *otelzap.Logger,
) ProductUsecase {
	return &productService{
		productRepo: productRepo,
		causeSvc:    causeSvc,
		logger:      logger,
	}
}

func (s *productService) CreateProduct(
	ctx context.Context, 
	ownerID uint, 
	title string, 
	price float64, 
	link string, 
	reasons []string,
) error {
	// tracer
	tr := otel.Tracer("product-service")
	ctx, span := tr.Start(ctx, "CreateProduct")
	defer span.End()

	newPosition, _ := ordering.KeyBetween("", "")

	product := &Product{
		OwnerID:  ownerID,
		Name:     title,
		ImageUrl: "",
		Link:     link,
		Price:    price,
		Status:   PENDING,
		Position: newPosition,
	}

	productID, err := s.productRepo.CreateProduct(ctx, product)
	if err != nil {
		s.logger.Ctx(ctx).Error("failed to create product in repository",
			zap.Uint("owner_id", ownerID),
			zap.String("title", title),
			zap.Error(err),
		)
		return err
	}

	if err := s.causeSvc.BulkCreateCauses(ctx, productID, reasons); err != nil {
		s.logger.Ctx(ctx).Error("failed to create causes for product",
			zap.Uint("product_id", productID),
			zap.Int("causes_count", len(reasons)),
			zap.Error(err),
		)
		return err
	}

	s.logger.Ctx(ctx).Info("product created successfully",
		zap.Uint("product_id", productID),
		zap.Uint("owner_id", ownerID),
		zap.Int("causes_count", len(reasons)),
	)

	return nil
}

func (s *productService) GetProduct(ctx context.Context, ownerID, productID uint) (*Product, error) {

	product, err := s.productRepo.GetProduct(ctx, ownerID, productID)
	if err != nil {
		s.logger.Ctx(ctx).Error("failed to get product",
			zap.Uint("owner_id", ownerID),
			zap.Uint("product_id", productID),
			zap.Error(err),
		)
		return nil, err
	}

	causes, err := s.causeSvc.GetCauses(ctx, product.ID)
	if err != nil {
		s.logger.Ctx(ctx).Error("failed to get causes for product",
			zap.Uint("product_id", product.ID),
			zap.Error(err),
		)
		return nil, err
	}

	product.Causes = causes

	s.logger.Ctx(ctx).Info("product fetched successfully",
		zap.Uint("product_id", product.ID),
		zap.Int("causes_count", len(causes)),
	)

	return product, nil
}

func (s *productService) GetProductByStatus(ctx context.Context, ownerID uint, status string) ([]*Product, error) {

	products, err := s.productRepo.GetProductByStatus(ctx, ownerID, status)
	if err != nil {
		s.logger.Ctx(ctx).Error("failed to get products by status",
			zap.Uint("owner_id", ownerID),
			zap.String("status", string(status)),
			zap.Error(err),
		)
		return nil, err
	}

	s.logger.Ctx(ctx).Info("products fetched by status",
		zap.Uint("owner_id", ownerID),
		zap.String("status", string(status)),
		zap.Int("count", len(products)),
	)

	return products, nil
}

func (s *productService) Move(ctx context.Context, ownerID uint, productID uint, productAfterID *uint) error {
	var prevPos, nextPos string

	if productAfterID == nil {
		// moving to first product
		np, err := s.productRepo.GetFirstPosition(ctx, ownerID)
		if err != nil {
			// TODO: handle log
			return err
		}

		nextPos = np
	} else {
		// moving after specific item
		pp, err := s.productRepo.GetPositionByProductID(ctx, ownerID, *productAfterID)
		if err != nil {
			// TODO: handle log
			return err
		}
		np, err := s.productRepo.GetNextPosition(ctx, ownerID, pp)
		if err != nil {
			// TODO: handle log
			return err
		}
		prevPos, nextPos = pp, np
	}

	newPos, err := ordering.KeyBetween(prevPos, nextPos)
	if err != nil {
		// TODO: handle log
		return err
	}

	if err := s.productRepo.UpdatePosition(ctx, ownerID, productID, newPos); err != nil {
		// TODO: handle log
		return err
	}

	return nil
}

func (s *productService) DeleteProduct(ctx context.Context, ownerID, productID uint) error {

	if err := s.productRepo.DeleteProduct(ctx, ownerID, productID); err != nil {
		s.logger.Ctx(ctx).Error("failed to delete product from repository",
			zap.Uint("owner_id", ownerID),
			zap.Uint("product_id", productID),
			zap.Error(err),
		)
		return err
	}

	if err := s.causeSvc.DeleteCauses(ctx,productID); err != nil {
		s.logger.Ctx(ctx).Error("failed to delete causes for product",
			zap.Uint("product_id", productID),
			zap.Error(err),
		)
		return err
	}

	s.logger.Ctx(ctx).Info("product deleted successfully",
		zap.Uint("owner_id", ownerID),
		zap.Uint("product_id", productID),
	)

	return nil
}
