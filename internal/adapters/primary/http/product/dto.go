package product

type CreateProductRequest struct {
	Title   string   `json:"title" validate:"required"`
	Price   float64  `json:"price" validate:"min=1"`
	Link    string   `json:"link" validate:"omitempty,url"`
	Reasons []string `json:"reasons" validate:"omitempty,dive,required"`
}

type UpdatePriorityRequest struct {
	ProductID      uint  `json:"productId" validate:"required"`
	ProductIDAfter *uint `json:"productIdAfter"`
}

type GetAllProductsRequest struct {
	Status string `query:"status" validate:"omitempty,oneof=pending installment bought"`
	Page   int    `query:"page" validate:"omitempty,min=1"`
	Size   int    `query:"size" validate:"omitempty,min=1"`
}
