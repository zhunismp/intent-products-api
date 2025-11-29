package cause

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type causeService struct {
	causeRepo CauseRepository
	logger    *zap.Logger
}

func NewCauseService(causeRepo CauseRepository, logger *zap.Logger) CauseUsecase {
	return &causeService{
		causeRepo: causeRepo,
		logger:    logger,
	}
}

func (s *causeService) BulkCreateCauses(ctx context.Context, productID uint, reasons []string) error {
	log := s.logger.With(
		zap.Uint("product_id", productID),
	)

	causes := make([]*Cause, 0, len(reasons))
	for _, reason := range reasons {
		causes = append(causes, &Cause{
			Reason: reason,
			Status: true,
		})
	}

	if err := s.causeRepo.BulkSaveCauses(ctx, productID, causes); err != nil {
		log.Error("failed to bulk save causes", zap.Error(err))
		return fmt.Errorf("failed to bulk save causes for product %d: %w", productID, err)
	}

	log.Info("bulk created causes successfully")
	return nil
}

func (s *causeService) GetCauses(ctx context.Context, productID uint) ([]*Cause, error) {
	log := s.logger.With(zap.Uint("product_id", productID))

	causes, err := s.causeRepo.FindByProductID(ctx, productID)
	if err != nil {
		log.Error("failed to fetch causes", zap.Error(err))
		return nil, fmt.Errorf("failed to get causes for product %d: %w", productID, err)
	}

	log.Info("fetched causes successfully", zap.Int("count", len(causes)))
	return causes, nil
}

func (s *causeService) DeleteCauses(ctx context.Context, productID uint) error {
	log := s.logger.With(zap.Uint("product_id", productID))

	if err := s.causeRepo.DeleteByProductID(ctx, productID); err != nil {
		log.Error("failed to delete causes", zap.Error(err))
		return fmt.Errorf("failed to delete causes for product %d: %w", productID, err)
	}

	log.Info("deleted causes successfully")
	return nil
}
