package product

import (
	domain "github.com/zhunismp/intent-products-api/internal/core/domain/product"
	"gorm.io/gorm"
)

type ProductModel struct {
	gorm.Model
	OwnerID  uint    `gorm:"type:bigint;not null"`
	Name     string  `gorm:"type:varchar(255);not null"`
	ImageURL *string `gorm:"type:text"`
	Link     *string `gorm:"type:text"`
	Price    float64 `gorm:"not null;check:price >= 0"`
	Status   string  `gorm:"type:varchar(50);not null;default:'active'"`
	Priority int64   `gorm:"type:bigint;default:-1"`
}

func (ProductModel) TableName() string {
	return "products"
}

func toProductModel(d *domain.Product) ProductModel {
	return ProductModel{
		Model:    gorm.Model{ID: d.ID},
		OwnerID:  d.OwnerID,
		Name:     d.Name,
		ImageURL: d.ImageUrl,
		Link:     d.Link,
		Price:    d.Price,
		Status:   d.Status,
		Priority: d.Priority,
	}
}

func toDomainProduct(m ProductModel) *domain.Product {
	return &domain.Product{
		ID:        m.ID,
		OwnerID:   m.OwnerID,
		Name:      m.Name,
		ImageUrl:  m.ImageURL,
		Link:      m.Link,
		Price:     m.Price,
		Status:    m.Status,
		Priority:  m.Priority,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
