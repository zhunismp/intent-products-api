package priority

import (
	"context"

	"github.com/zhunismp/intent-products-api/internal/core/domain/shared/apperrors"
	"go.uber.org/zap"
)

type priorityService struct {
	logger *zap.Logger
}

func NewPriorityService(logger *zap.Logger) PriorityUsecase {
	return &priorityService{logger: logger}
}

func (s *priorityService) CalculateNewPriority(ctx context.Context, priorityBefore *int64, priorityAfter *int64) (int64, error) {
	if priorityBefore == nil && priorityAfter == nil {
		s.logger.Error("invalid priority update - missing before/after reference")
		return -1, apperrors.New(apperrors.ErrCodeValidation, "before or after priority must be set", nil)
	}

	switch {
	case priorityBefore != nil && priorityAfter == nil:
		newPriority := *priorityBefore / 2
		s.logger.Info("calculated new priority (insert at top)", zap.Int64("new_priority", newPriority))
		return newPriority, nil

	case priorityBefore == nil && priorityAfter != nil:
		newPriority := *priorityAfter + Step
		s.logger.Info("calculated new priority (insert at bottom)", zap.Int64("new_priority", newPriority))
		return newPriority, nil

	default:
		newPriority := *priorityBefore + (*priorityAfter-*priorityBefore)/2
		s.logger.Info("calculated new priority (between two priorities)", zap.Int64("new_priority", newPriority))
		return newPriority, nil
	}
}
