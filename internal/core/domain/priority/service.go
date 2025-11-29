package priority

import (
	"context"

	"github.com/zhunismp/intent-products-api/internal/core/domain/shared/apperrors"
)

type priorityService struct {
}

func NewPriorityService() PriorityUsecase {
	return &priorityService{}
}

func (s *priorityService) CalculateNewPriority(ctx context.Context, priorityBefore *int64, priorityAfter *int64) (int64, error) {
	if priorityBefore == nil && priorityAfter == nil {
		return -1, apperrors.New(apperrors.ErrCodeValidation, "before or after priority must be set", nil)
	}

	switch {
	case priorityBefore != nil && priorityAfter == nil:
		return *priorityBefore / 2, nil
	case priorityBefore == nil && priorityAfter != nil:
		return *priorityAfter + Step, nil
	default:
		return *priorityBefore + (*priorityAfter-*priorityBefore)/2, nil
	}
}
