package product

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/zhunismp/intent-products-api/internal/core/domain/product"
	"github.com/zhunismp/intent-products-api/internal/core/domain/shared/apperrors"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) CreateProduct(ctx context.Context, product *domain.Product) (uint, error) {
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

func (r *productRepository) GetProductByStatus(ctx context.Context, ownerID uint, status string) ([]*domain.Product, error) {
	var models []ProductModel

	err := r.db.WithContext(ctx).
		Where("owner_id = ? AND status = ?", ownerID, status).
		Order("priority ASC").
		Find(&models).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.New(
			apperrors.ErrCodeNotFound,
			fmt.Sprintf("no products found for owner %d with status %s", ownerID, status),
			err,
		)
	}

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

func (r *productRepository) GetLastProductPriority(ctx context.Context, ownerID uint) (int64, error) {
	var maxPriority *int64

	err := r.db.WithContext(ctx).
		Model(&ProductModel{}).
		Where("owner_id = ?", ownerID).
		Select("MAX(priority)").
		Scan(&maxPriority).Error

	if err != nil {
		return 0, apperrors.New(
			apperrors.ErrCodeInternal,
			"failed to get last product priority",
			err,
		)
	}

	if maxPriority == nil {
		return 0, nil
	}

	return *maxPriority, nil
}

func (r *productRepository) BulkGetProducts(ctx context.Context, ownerID uint, productIDs []uint) ([]*domain.Product, error) {
	if len(productIDs) == 0 {
		return []*domain.Product{}, nil
	}

	var models []ProductModel

	err := r.db.WithContext(ctx).
		Where("owner_id = ? AND id IN ?", ownerID, productIDs).
		Find(&models).Error

	if err != nil {
		return nil, apperrors.New(
			apperrors.ErrCodeInternal,
			"failed to bulk get products",
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
