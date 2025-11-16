package product

// alias type
type OwnerId = string
type ProductId = string
type CauseId = string

// common type
type Sort struct {
	Field     string
	Direction string
}

type CauseStatus struct {
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

type UpdateCauseStatusCmd struct {
	OwnerID     string
	ProductID   string
	CauseStatus CauseStatus
}
