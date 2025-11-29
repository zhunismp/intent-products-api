package product

// usecase command
type CreateProductCmd struct {
	OwnerID uint
	Title   string
	Price   float64
	Link    *string
	Reasons []string
}

type DeleteProductCmd struct {
	OwnerID   uint
	ProductID uint
}

type GetProductCmd struct {
	OwnerID   uint
	ProductID uint
}

type GetProductByStatusCmd struct {
	OwnerID uint
	Status  string
}

type UpdatePriorityCmd struct {
	OwnerID         uint
	ProductID       uint
	ProductIDBefore *uint
	ProductIDAfter  *uint
}
