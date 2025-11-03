package product

import (
	"fmt"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
)

func GenerateErrorMap(errs validator.ValidationErrors) map[string]string {
	errMap := make(map[string]string)

	for _, e := range errs {
		errMap[e.Field()] = fmt.Sprintf("failed on '%s' rule", e.Tag())
	}

	return errMap
}

// custom validator
func IsDateAfter(fl validator.FieldLevel) bool {
	otherField := fl.Parent().FieldByName(fl.Param())
	field := fl.Parent().FieldByName(fl.FieldName())

	msg := fmt.Sprintf("type of field = %s, otherField = %s", field.Type(), otherField.Type())
	log.Print(msg)

	// if either itself is nil or comparision field is nil, skip validation.
	if field.IsNil() || otherField.IsNil() {
		return true
	}

	end := field.Interface().(*time.Time)
	start := otherField.Interface().(*time.Time)

	return end.After(*start)
}
