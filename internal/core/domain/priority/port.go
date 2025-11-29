package priority

import (
	"context"
)

type PriorityUsecase interface {
	CalculateNewPriority(ctx context.Context, priorityBefore *int64, priorityAfter *int64) (int64, error)
}
