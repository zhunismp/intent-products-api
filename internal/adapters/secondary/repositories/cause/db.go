package cause

import (
	"context"
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

func (r *causeRepository) UpdateCauseStatus(ctx context.Context, productID string, causes []domain.UpdateCauseStatus) ([]domain.Cause, error) {
	tx := r.db.WithContext(ctx)
	var updatedIds []string

	for _, c := range causes {
		res := tx.Model(&CauseModel{}).
			Where("id = ? AND product_id = ?", c.CauseID, productID).
			Updates(map[string]any{
				"status": c.Status,
			})

		if res.Error != nil {
			return nil, apperrors.New(apperrors.ErrCodeInternal,
				fmt.Sprintf("failed to update cause %s", c.CauseID), res.Error)
		}

		if res.RowsAffected == 0 {
			return nil, apperrors.New(apperrors.ErrCodeNotFound,
				fmt.Sprintf("no cause found with id %s for product %s", c.CauseID, productID), nil)
		}

		updatedIds = append(updatedIds, c.CauseID)
	}

	// query update entry
	var updatedModels []CauseModel
	if err := tx.Where("id IN ?", updatedIds).Find(&updatedModels).Error; err != nil {
		return nil, apperrors.New(apperrors.ErrCodeInternal, "failed to fetch updated causes", err)
	}

	// transform to domain model
	updatedCauses := make([]domain.Cause, 0, len(updatedModels))
	for _, m := range updatedModels {
		updatedCauses = append(updatedCauses, *toDomainCause(m))
	}

	return updatedCauses, nil
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
