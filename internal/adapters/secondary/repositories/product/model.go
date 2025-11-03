package product

import (
	"time"

	"github.com/zhunismp/intent-products-api/internal/core/domain/shared/utils"
	domain "github.com/zhunismp/intent-products-api/internal/core/domain/product"
	"gorm.io/gorm"
)

type ProductModel struct {
	ID       string  `gorm:"type:char(26);primaryKey"`
	OwnerID  string  `gorm:"type:char(26);not null;index"`
	Name     string  `gorm:"type:varchar(255);not null"`
	ImageURL *string `gorm:"type:text"`
	Link     *string `gorm:"type:text"`
	Price    float64 `gorm:"not null;check:price >= 0"`
	Status   string  `gorm:"type:varchar(50);not null;default:'active'"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (ProductModel) TableName() string {
	return "products"
}

// BeforeCreate hook — auto-generate ULID before insert
func (p *ProductModel) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = utils.GenULID(time.Now())
	}
	return
}

func toProductModel(d domain.Product) ProductModel {
	return ProductModel{
		ID:        d.ID,
		OwnerID:   d.OwnerID,
		Name:      d.Name,
		ImageURL:  d.ImageUrl,
		Link:      d.Link,
		Price:     d.Price,
		CreatedAt:   d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
		Status:    d.Status,
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
		CreatedAt:   m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Status:    m.Status,
	}
}
