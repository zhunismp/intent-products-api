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

func (r *productRepository) CreateProduct(ctx context.Context, product domain.Product) (*domain.Product, error) {
	model := toProductModel(product)

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperrors.New(
				apperrors.ErrCodeForbidden,
				fmt.Sprintf("owner id %s attempted to create existing product name '%s'", product.OwnerID, product.Name),
				err,
			)
		}
		return nil, apperrors.New(apperrors.ErrCodeInternal, "failed to create product", err)
	}

	return toDomainProduct(model), nil
}

func (r *productRepository) GetProduct(ctx context.Context, ownerID string, productID string) (*domain.Product, error) {
	var model ProductModel
	err := r.db.WithContext(ctx).
		Where("id = ? AND owner_id = ?", productID, ownerID).
		First(&model).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.New(
			apperrors.ErrCodeNotFound,
			fmt.Sprintf("owner id %s does not own product id %s", ownerID, productID),
			err,
		)
	}

	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeInternal, "failed to get product", err)
	}

	return toDomainProduct(model), nil
}

func (r *productRepository) GetProductByStatus(ctx context.Context, ownerID string, status string) ([]domain.Product, error) {
	var models []ProductModel

	err := r.db.WithContext(ctx).
		Where("owner_id = ? AND status = ?", ownerID, status).
		Find(&models).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.New(
			apperrors.ErrCodeNotFound,
			fmt.Sprintf("no products found for owner %s with status %s", ownerID, status),
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

	products := make([]domain.Product, 0, len(models))
	for _, m := range models {
		products = append(products, *toDomainProduct(m))
	}

	return products, nil
}

func (r *productRepository) DeleteProduct(ctx context.Context, ownerID string, productID string) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND owner_id = ?", productID, ownerID).
		Delete(&ProductModel{})

	if result.Error != nil {
		return apperrors.New(apperrors.ErrCodeInternal, "failed to delete product", result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.New(
			apperrors.ErrCodeNotFound,
			fmt.Sprintf("owner id %s does not own product id %s", ownerID, productID),
			nil,
		)
	}
	return nil
}
