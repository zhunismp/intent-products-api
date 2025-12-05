package product

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/zhunismp/intent-products-api/internal/core/domain/product"
	"github.com/zhunismp/intent-products-api/internal/core/domain/shared/apperrors"
	"github.com/zhunismp/intent-products-api/internal/core/domain/shared/utils/ordering"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) CreateProduct(ctx context.Context, product *domain.Product) (uint, error) {
	// Get the last position to append the new product at the end
	lastPosition, err := r.GetLastPosition(ctx, product.OwnerID)
	if err != nil {
		return 0, err
	}

	// Generate new position after the last item
	newPosition, err := ordering.KeyBetween(lastPosition, "")
	if err != nil {
		return 0, apperrors.New(
			apperrors.ErrCodeInternal,
			"failed to generate position",
			err,
		)
	}

	product.Position = newPosition
	model := toProductModel(product)

	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return 0, apperrors.New(
				apperrors.ErrCodeForbidden,
				fmt.Sprintf("owner id %d attempted to create existing product name '%s'", product.OwnerID, product.Name),
				err,
			)
		}
		return 0, apperrors.New(apperrors.ErrCodeInternal, "failed to create product", err)
	}

	return model.ID, nil
}

func (r *productRepository) GetProduct(ctx context.Context, ownerID uint, productID uint) (*domain.Product, error) {
	var model ProductModel
	err := r.db.WithContext(ctx).
		Where("id = ? AND owner_id = ?", productID, ownerID).
		First(&model).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.New(
			apperrors.ErrCodeNotFound,
			fmt.Sprintf("owner id %d does not own product id %d", ownerID, productID),
			err,
		)
	}

	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeInternal, "failed to get product", err)
	}

	return toDomainProduct(model), nil
}

func (r *productRepository) FindAllProducts(ctx context.Context, ownerID uint, filter *domain.Filter) ([]*domain.Product, error) {
	q := r.db.WithContext(ctx).
		Where("owner_id = ?", ownerID).
		Order("position")

	// TODO: extract logic away from here
	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}

	offset := (filter.Page - 1) * filter.Size
	q = q.Offset(offset).Limit(filter.Size)

	var models []ProductModel
	err := q.Find(&models).Error

	if err != nil {
		return nil, apperrors.New(
			apperrors.ErrCodeInternal,
			"failed to get products by status",
			err,
		)
	}

	products := make([]*domain.Product, 0, len(models))
	for _, m := range models {
		products = append(products, toDomainProduct(m))
	}

	return products, nil
}

func (r *productRepository) DeleteProduct(ctx context.Context, ownerID uint, productID uint) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND owner_id = ?", productID, ownerID).
		Delete(&ProductModel{})

	if result.Error != nil {
		return apperrors.New(apperrors.ErrCodeInternal, "failed to delete product", result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.New(
			apperrors.ErrCodeNotFound,
			fmt.Sprintf("owner id %d does not own product id %d", ownerID, productID),
			nil,
		)
	}
	return nil
}

func (r *productRepository) GetFirstPosition(ctx context.Context, ownerID uint) (string, error) {
	var position string

	err := r.db.WithContext(ctx).
		Table("products").
		Select("position").
		Where("owner_id = ?", ownerID).
		Order("position").
		Limit(1).
		Pluck("position", &position).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", apperrors.New(
				apperrors.ErrCodeNotFound,
				fmt.Sprintf("no products found for owner id %d", ownerID),
				err,
			)
		}
		return "", apperrors.New(
			apperrors.ErrCodeInternal,
			"failed to get first position",
			err,
		)
	}

	if position == "" {
		return "", apperrors.New(
			apperrors.ErrCodeNotFound,
			fmt.Sprintf("no products found for owner id %d", ownerID),
			nil,
		)
	}

	return position, nil
}

func (r *productRepository) GetLastPosition(ctx context.Context, ownerID uint) (string, error) {
	var position string
	err := r.db.WithContext(ctx).
		Table("products").
		Select("position").
		Where("owner_id = ?", ownerID).
		Order("position DESC").
		Limit(1).
		Pluck("position", &position).
		Error

	// Ignore record not found as it's first product
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", apperrors.New(
			apperrors.ErrCodeInternal,
			"failed to get last position",
			err,
		)
	}

	return position, nil
}

func (r *productRepository) GetPositionByProductID(ctx context.Context, ownerID uint, productID uint) (string, error) {
	var position string

	err := r.db.WithContext(ctx).
		Table("products").
		Select("position").
		Where("id = ? AND owner_id = ?", productID, ownerID).
		Pluck("position", &position).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", apperrors.New(
				apperrors.ErrCodeNotFound,
				fmt.Sprintf("product id %d not found for owner id %d", productID, ownerID),
				err,
			)
		}
		return "", apperrors.New(
			apperrors.ErrCodeInternal,
			"failed to get position by product id",
			err,
		)
	}

	if position == "" {
		return "", apperrors.New(
			apperrors.ErrCodeNotFound,
			fmt.Sprintf("product id %d not found for owner id %d", productID, ownerID),
			nil,
		)
	}

	return position, nil
}

func (r *productRepository) GetNextPosition(ctx context.Context, ownerID uint, position string) (string, error) {
	var nextPosition string

	err := r.db.WithContext(ctx).
		Table("products").
		Select("position").
		Where("owner_id = ? AND position > ?", ownerID, position).
		Order("position COLLATE \"C\" ASC").
		Limit(1).
		Pluck("position", &nextPosition).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No next position means this is the last item - return empty string
			return "", nil
		}
		return "", apperrors.New(
			apperrors.ErrCodeInternal,
			"failed to get next position",
			err,
		)
	}

	// Empty result means no next item (last item in list)
	return nextPosition, nil
}

func (r *productRepository) UpdatePosition(ctx context.Context, ownerID uint, productID uint, position string) error {
	result := r.db.WithContext(ctx).
		Table("products").
		Where("id = ? AND owner_id = ?", productID, ownerID).
		Update("position", position)

	if result.Error != nil {
		return apperrors.New(
			apperrors.ErrCodeInternal,
			"failed to update position",
			result.Error,
		)
	}

	if result.RowsAffected == 0 {
		return apperrors.New(
			apperrors.ErrCodeNotFound,
			fmt.Sprintf("product id %d not found for owner id %d", productID, ownerID),
			nil,
		)
	}

	return nil
}

func (r *productRepository) ValidateOwnership(ctx context.Context, ownerID, productID uint) error {
	var count int64
	err := r.db.WithContext(ctx).
		Table("products").
		Where("id = ? AND owner_id = ?", productID, ownerID).
		Count(&count).
		Error

	if err != nil {
		return apperrors.New(
			apperrors.ErrCodeInternal,
			"failed to validate ownership",
			err,
		)
	}

	if count == 0 {
		return apperrors.New(
			apperrors.ErrCodeNotFound,
			fmt.Sprintf("product id %d not found for owner id %d", productID, ownerID),
			nil,
		)
	}

	return nil
}
