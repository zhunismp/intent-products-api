package product

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	domain "github.com/zhunismp/intent-products-api/internal/core/domain/product"
	"github.com/zhunismp/intent-products-api/internal/core/domain/shared/apperrors"
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

// TODO: clean up logic here
func (r *productRepository) QueryProduct(ctx context.Context, spec domain.QueryProductSpec) ([]domain.Product, error) {
	var models []ProductModel
	tx := r.db.WithContext(ctx).Where("owner_id = ?", spec.OwnerID)

	if spec.Status != nil {
		tx = tx.Where("status = ?", *spec.Status)
	}

	switch {
	case spec.Start != nil && spec.End != nil:
		tx = tx.Where("created_at BETWEEN ? AND ?", spec.Start, spec.End)
	case spec.Start != nil:
		tx = tx.Where("created_at >= ?", spec.Start)
	case spec.End != nil:
		tx = tx.Where("created_at <= ?", spec.End)
	}

	tx = tx.Order(clause.OrderByColumn{Column: clause.Column{Name: spec.Sort.Field}, Desc: spec.Sort.Direction == "desc"})

	if err := tx.Find(&models).Error; err != nil {
		return nil, apperrors.New(apperrors.ErrCodeInternal, "failed to query products", err)
	}

	// Convert to domain entities
	products := make([]domain.Product, 0, len(models))
	for _, m := range models {
		products = append(products, *toDomainProduct(m))
	}
	return products, nil
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

func (r *productRepository) BatchGetProduct(ctx context.Context, ownerID string, productIDs []string) ([]domain.Product, error) {
	var models []ProductModel

	err := r.db.WithContext(ctx).
		Where("owner_id = ? AND id IN ?", ownerID, productIDs).
		Find(&models).Error

	if err != nil {
		return nil, apperrors.New(
			apperrors.ErrCodeInternal,
			"failed to batch get products",
			err,
		)
	}

	// Convert models to domain products
	products := make([]domain.Product, 0, len(models))
	for _, model := range models {
		products = append(products, *toDomainProduct(model))
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
