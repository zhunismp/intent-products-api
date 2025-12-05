package product

import (
	"time"

	"github.com/zhunismp/intent-products-api/internal/core/domain/cause"
)

// TODO: when logic is complex, should not return domain object directly
type Product struct {
	ID       uint           `json:"id"`
	OwnerID  uint           `json:"ownerId"`
	Name     string         `json:"name"`
	ImageUrl string         `json:"imageUrl"`
	Link     string         `json:"link"`
	Price    float64        `json:"price"`
	Status   string         `json:"status"`
	Position string         `json:"-"`
	Causes   []*cause.Cause `json:"causes,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

const (
	PENDING     string = "pending"
	INSTALLMENT string = "installment"
	BOUGHT      string = "bought"
)
