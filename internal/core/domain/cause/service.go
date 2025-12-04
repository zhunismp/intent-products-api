package cause

import (
	"context"
	"fmt"
	"log/slog"
)

type causeService struct {
	causeRepo CauseRepository
	logger    *slog.Logger
}

func NewCauseService(causeRepo CauseRepository, logger *slog.Logger) CauseUsecase {
	return &causeService{causeRepo: causeRepo, logger: logger}
}

func (s *causeService) BulkCreateCauses(ctx context.Context, productID uint, reasons []string) error {
	causes := make([]*Cause, 0, len(reasons))
	for _, reason := range reasons {
		causes = append(causes, &Cause{
			Reason: reason,
			Status: true,
		})
	}

	if err := s.causeRepo.BulkSaveCauses(ctx, productID, causes); err != nil {
		return fmt.Errorf("failed to bulk save causes for product %d: %w", productID, err)
	}

	s.logger.Info("bulk created causes successfully")
	return nil
}

func (s *causeService) GetCauses(ctx context.Context, productID uint) ([]*Cause, error) {

	causes, err := s.causeRepo.FindByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get causes for product %d: %w", productID, err)
	}

	s.logger.Info("fetched causes successfully", slog.Int("count", len(causes)))
	return causes, nil
}

func (s *causeService) DeleteCauses(ctx context.Context, productID uint) error {
	if err := s.causeRepo.DeleteByProductID(ctx, productID); err != nil {
		return fmt.Errorf("failed to delete causes for product %d: %w", productID, err)
	}

	s.logger.Info("deleted causes successfully")
	return nil
}
