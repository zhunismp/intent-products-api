package product

import (
	domain "github.com/zhunismp/intent-products-api/internal/core/domain/product"
	"gorm.io/gorm"
)

type ProductModel struct {
	gorm.Model
	OwnerID  uint    `gorm:"type:bigint;not null"`
	Name     string  `gorm:"type:varchar(255);not null"`
	ImageURL string  `gorm:"type:text"`
	Link     string  `gorm:"type:text"`
	Price    float64 `gorm:"not null;check:price >= 0"`
	Status   string  `gorm:"type:varchar(50);not null;default:'active'"`
	Position string  `gorm:"type:varchar(255) COLLATE \"C\";not null"` // ensure binary order
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
		Position: d.Position,
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
		Position:  m.Position,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
