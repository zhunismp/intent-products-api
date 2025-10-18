package validators

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/zhunismp/intent-products-api/internal/core/dtos"
)

func ValidateCreateProductReq(input dtos.CreateProductInput) error {
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return errors.New("validation failed")
	}

	return nil
}