package product

import "time"

type Product struct {
	ID        string
	OwnerID   string
	Name      string
	ImageUrl  *string
	Link      *string
	Price     float64
	CreatedAt time.Time
	UpdatedAt time.Time
	Status    string
	Causes    []Cause
}

type Cause struct {
	ID        string
	Reason    string
	Status    bool
}

const (
	PENDING     string = "pending"
	INSTALLMENT string = "installment"
	BOUGHT      string = "bought"
)
