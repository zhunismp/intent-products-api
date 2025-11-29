package cause

import (
	"context"
	"fmt"

	domain "github.com/zhunismp/intent-products-api/internal/core/domain/cause"
	"github.com/zhunismp/intent-products-api/internal/core/domain/shared/apperrors"
	"gorm.io/gorm"
)

type causeRepository struct {
	db *gorm.DB
}

func NewCauseRepository(db *gorm.DB) domain.CauseRepository {
	return &causeRepository{db: db}
}

func (r *causeRepository) BulkSaveCauses(ctx context.Context, productID uint, causes []*domain.Cause) error {
	if len(causes) == 0 {
		return nil
	}

	models := make([]*CauseModel, len(causes))
	for i, c := range causes {
		models[i] = FromDomain(productID, c)
	}

	err := r.db.WithContext(ctx).Save(models).Error
	if err != nil {
		return apperrors.New(apperrors.ErrCodeInternal, "failed to bulk save causes", err)
	}

	return nil
}

func (r *causeRepository) FindByProductID(ctx context.Context, productID uint) ([]*domain.Cause, error) {
	var models []*CauseModel

	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Find(&models).Error

	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeInternal, "failed to find cause by product id", err)
	}

	result := make([]*domain.Cause, len(models))
	for i, model := range models {
		result[i] = model.ToDomain()
	}

	return result, nil
}

func (r *causeRepository) DeleteByProductID(ctx context.Context, productID uint) error {
	result := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Delete(&CauseModel{})

	if result.Error != nil {
		return apperrors.New(apperrors.ErrCodeInternal, "failed to delete cause", result.Error)
	}

	if result.RowsAffected == 0 {
		return apperrors.New(apperrors.ErrCodeNotFound, "no cause found", fmt.Errorf("no cause found for product %d", productID))
	}

	return nil
}
