package validators

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/zhunismp/intent-products-api/internal/core/dtos"
)

func ValidateCreateProductReq(req dtos.CreateProductInput) error {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return errors.New("validation failed")
	}

	return nil
}