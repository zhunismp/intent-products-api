package product

import (
	"time"

	"github.com/zhunismp/intent-products-api/internal/core/domain/cause"
)

type Product struct {
	ID       uint
	OwnerID  uint
	Name     string
	ImageUrl *string
	Link     *string
	Price    float64
	Status   string
	Priority int64
	Causes   []*cause.Cause

	CreatedAt time.Time
	UpdatedAt time.Time
}

const (
	PENDING     string = "pending"
	INSTALLMENT string = "installment"
	BOUGHT      string = "bought"
)
