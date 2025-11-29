package cause

import "context"

type CauseUsecase interface {
	BulkCreateCauses(ctx context.Context, productID uint, reasons []string) error
	GetCauses(ctx context.Context, productID uint) ([]*Cause, error)
	DeleteCauses(ctx context.Context, productID uint) error
}

type CauseRepository interface {
	BulkSaveCauses(ctx context.Context, productID uint, causes []*Cause) error
	FindByProductID(ctx context.Context, productID uint) ([]*Cause, error)
	DeleteByProductID(ctx context.Context, productID uint) error
}
