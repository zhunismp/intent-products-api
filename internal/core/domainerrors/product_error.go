package domainerrors

import "errors"

var (
	ErrorDuplicateProduct = errors.New("duplicate product")
	ErrorProductNotFound  = errors.New("product not found")
	ErrorProducInput      = errors.New("product input is invalid")
)
