package cause

import (
	"time"

	domain "github.com/zhunismp/intent-products-api/internal/core/domain/product"
	"github.com/zhunismp/intent-products-api/internal/core/domain/shared/utils"
	"gorm.io/gorm"
)

type CauseModel struct {
	ID        string `gorm:"type:char(26);primaryKey"`
	ProductID string `gorm:"type:char(26);not null;index"`
	Reason    string `gorm:"type:text;not null"`
	Status    bool   `gorm:"not null;default:true"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (CauseModel) TableName() string {
	return "causes"
}

// BeforeCreate hook — generate ULID if not provided
func (c *CauseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = utils.GenULID(time.Now())
	}
	return
}

func toCauseModel(productID string, d domain.Cause) CauseModel {
	return CauseModel{
		ID:        d.ID,
		ProductID: productID,
		Reason:    d.Reason,
		Status:    d.Status,
	}
}

func toDomainCause(m CauseModel) *domain.Cause {
	return &domain.Cause{
		ID:     m.ID,
		Reason: m.Reason,
		Status: m.Status,
	}
}
