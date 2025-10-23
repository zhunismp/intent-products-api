package product

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func GenerateErrorMap(errs validator.ValidationErrors) map[string]string {
	errMap := make(map[string]string)

	for _, e := range errs {
		errMap[e.Field()] = fmt.Sprintf("failed on '%s' rule", e.Tag())
	}

	return errMap
}
