package cause

import (
	domain "github.com/zhunismp/intent-products-api/internal/core/domain/cause"
	"gorm.io/gorm"
)

type CauseModel struct {
	gorm.Model
	ProductID uint   `gorm:"type:bigint;not null"`
	Reason    string `gorm:"type:text;not null"`
	Status    bool   `gorm:"not null;default:true"`
}

func (CauseModel) TableName() string {
	return "causes"
}

func (m *CauseModel) ToDomain() *domain.Cause {
	return &domain.Cause{
		ID:     m.ID,
		Reason: m.Reason,
		Status: m.Status,
	}
}

func FromDomain(productID uint, d *domain.Cause) *CauseModel {
	return &CauseModel{
		Model:     gorm.Model{ID: d.ID},
		ProductID: productID,
		Reason:    d.Reason,
		Status:    d.Status,
	}
}
