package product

import "time"

// alias type
type OwnerId = string
type ProductId = string
type CauseId = string

// common type
type Sort struct {
	Field     string
	Direction string
}

type UpdateCauseStatus struct {
	CauseID string
	Status  bool
}

// usecase command
type CreateProductCmd struct {
	OwnerID string
	Title   string
	Price   float64
	Link    *string
	Reasons []string
}

type DeleteProductCmd struct {
	OwnerID   string
	ProductID string
}

type GetProductCmd struct {
	OwnerID   string
	ProductID string
}

type UpdateProductCmd struct {
	OwnerID   string
	ProductID string
	Causes    []UpdateCauseStatus
}

type QueryProductCmd struct {
	OwnerID string
	Start   *time.Time
	End     *time.Time
	Status  *string
	Sort    *Sort
}

// repo query
type QueryProductSpec struct {
	OwnerID string
	Start   *time.Time
	End     *time.Time
	Status  *string
	Sort    Sort
}
