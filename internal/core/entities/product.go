package entities

import (
	"time"
)

type Product struct {
	ID        string    `bson:"id"`
	OwnerID   string    `bson:"owner_id"`
	Name      string    `bson:"name"`
	ImageUrl  *string   `bson:"image_url,omitempty"`
	Link      *string   `bson:"link,omitempty"`
	Price     float64   `bson:"price"`
	AddedAt   time.Time `bson:"added_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	Status    Status    `bson:"status"`
	Causes    []Cause   `bson:"causes,omitempty"`
}

type Cause struct {
	Reason string `bson:"reason"`
	Status bool   `bson:"status"`
}
