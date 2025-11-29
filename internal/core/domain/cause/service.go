package cause

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type causeService struct {
	causeRepo CauseRepository
	logger *zap.Logger
}

func NewCauseService(causeRepo CauseRepository, logger *zap.Logger) CauseUsecase {
	return &causeService{causeRepo: causeRepo, logger: logger}
}

func (s *causeService) BulkCreateCauses(ctx context.Context, productID uint, reasons []string) error {
	causes := make([]*Cause, 0, len(reasons))
	for _, reason := range reasons {

		c := &Cause{
			Reason: reason,
			Status: true, // TODO: decided whether keep this field or not.
		}

		causes = append(causes, c)
	}

	// Save causes to repository
	if err := s.causeRepo.BulkSaveCauses(ctx, productID, causes); err != nil {
		return fmt.Errorf("failed to bulk save causes for product %d: %w", productID, err)
	}

	return nil
}

func (s *causeService) GetCauses(ctx context.Context, productID uint) ([]*Cause, error) {
	causes, err := s.causeRepo.FindByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get causes for product %d: %w", productID, err)
	}

	return causes, nil
}

// DeleteCauses removes all causes for a product
func (s *causeService) DeleteCauses(ctx context.Context, productID uint) error {
	// Delete causes from repository
	if err := s.causeRepo.DeleteByProductID(ctx, productID); err != nil {
		return fmt.Errorf("failed to delete causes for product %d: %w", productID, err)
	}

	return nil
}
