package entities

import "time"

type Product struct {
	ID        int
	OwnerID   int
	Name      string
	ImageUrl  *string
	Link      *string
	Price     float64
	AddedAt   time.Time
	UpdatedAt time.Time
	IsBought  bool
	Status    string
	Causes    []Cause
}

type Cause struct {
	Reason string
	Status string
}
