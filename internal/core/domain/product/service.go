package product

import (
	"context"
	"log/slog"

	"github.com/zhunismp/intent-products-api/internal/core/domain/cause"
	"github.com/zhunismp/intent-products-api/internal/core/domain/shared/utils/ordering"
	"go.opentelemetry.io/otel"
)

type productService struct {
	productRepo ProductRepository
	causeSvc    cause.CauseUsecase
	logger      *slog.Logger
}

func NewProductService(
	productRepo ProductRepository,
	causeSvc cause.CauseUsecase,
	logger *slog.Logger,
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

	product := &Product{
		OwnerID:  ownerID,
		Name:     title,
		ImageUrl: "",
		Link:     link,
		Price:    price,
		Status:   PENDING,
	}

	productID, err := s.productRepo.CreateProduct(ctx, product)
	if err != nil {
		return err
	}

	if err := s.causeSvc.BulkCreateCauses(ctx, productID, reasons); err != nil {
		return err
	}

	s.logger.Info("product created successfully")

	return nil
}

func (s *productService) GetProduct(ctx context.Context, ownerID, productID uint) (*Product, error) {

	product, err := s.productRepo.GetProduct(ctx, ownerID, productID)
	if err != nil {
		return nil, err
	}

	causes, err := s.causeSvc.GetCauses(ctx, product.ID)
	if err != nil {
		return nil, err
	}

	product.Causes = causes

	s.logger.Info("product fetched successfully")

	return product, nil
}

func (s *productService) GetAllProducts(ctx context.Context, ownerID uint, filter *Filter) ([]*Product, error) {

	products, err := s.productRepo.FindAllProducts(ctx, ownerID, filter)
	if err != nil {
		return nil, err
	}

	s.logger.Info("products fetched by status")

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
		return err
	}

	if err := s.causeSvc.DeleteCauses(ctx, productID); err != nil {
		return err
	}

	s.logger.Info("product deleted successfully")

	return nil
}
