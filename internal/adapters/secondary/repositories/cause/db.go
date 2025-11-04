package cause

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/zhunismp/intent-products-api/internal/core/domain/product"
	"github.com/zhunismp/intent-products-api/internal/core/domain/shared/apperrors"
	"gorm.io/gorm"
)

type causeRepository struct {
	db *gorm.DB
}

func NewCauseRepository(db *gorm.DB) domain.CauseRepository {
	return &causeRepository{db: db}
}

func (r *causeRepository) CreateCauses(ctx context.Context, productID string, causes []domain.Cause) ([]domain.Cause, error) {
	models := make([]CauseModel, len(causes))
	for i, c := range causes {
		models[i] = toCauseModel(productID, c)
	}

	if err := r.db.WithContext(ctx).Create(&models).Error; err != nil {
		return nil, apperrors.New(apperrors.ErrCodeInternal, "failed to create causes", err)
	}

	return causes, nil
}

func (r *causeRepository) GetCauses(ctx context.Context, productID string) ([]domain.Cause, error) {
	var models []CauseModel

	if err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Find(&models).Error; err != nil {
		return nil, apperrors.New(apperrors.ErrCodeInternal, "failed to get causes", err)
	}

	results := make([]domain.Cause, len(models))
	for i := range models {
		results[i] = *toDomainCause(models[i])
	}
	return results, nil
}

func (r *causeRepository) UpdateCauseStatus(ctx context.Context, productID string, cause domain.CauseStatus) (*domain.Cause, error) {
	tx := r.db.WithContext(ctx)

	var updatedModel CauseModel
	res := tx.Model(&CauseModel{}).
		Where("id = ? AND product_id = ?", cause.CauseID, productID).
		Update("status", cause.Status).
		Scan(&updatedModel)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, apperrors.New(
				apperrors.ErrCodeNotFound,
				fmt.Sprintf("no cause found with id %s for product %s", cause.CauseID, productID),
				nil,
			)
		}

		return nil, apperrors.New(
			apperrors.ErrCodeInternal,
			fmt.Sprintf("failed to update cause %s", cause.CauseID),
			res.Error,
		)
	}

	updatedCause := toDomainCause(updatedModel)
	return updatedCause, nil
}

func (r *causeRepository) DeleteCausesByProductID(ctx context.Context, productID string) error {
	result := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Delete(&CauseModel{})

	if result.Error != nil {
		return apperrors.New(apperrors.ErrCodeInternal, "failed to delete causes", result.Error)
	}

	if result.RowsAffected == 0 {
		return apperrors.New(
			apperrors.ErrCodeNotFound,
			fmt.Sprintf("no causes found for product id %s", productID),
			nil,
		)
	}

	return nil
}
