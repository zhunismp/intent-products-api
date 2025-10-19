package validators

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/zhunismp/intent-products-api/internal/core/dtos"
)

func ValidateCreateProductInput(input dtos.CreateProductInput) error {
	// initial validator
	validate := validator.New()

	// validate
	if err := validate.Struct(input); err != nil {
		return fmt.Errorf("validation failed: %v", err)
	}

	return nil
}

func ValidateQueryProductInput(input dtos.QueryProductInput) error {
	// initial validator
	validate := validator.New()

	// apply rules
	validate.RegisterStructValidation(queryProductFilterRules, dtos.QueryProductFilter{})

	// validate
	if err := validate.Struct(input); err != nil {
		return fmt.Errorf("validation failed: %v", err)
	}

	return nil
}

func queryProductFilterRules(sl validator.StructLevel) {
	filter := sl.Current().Interface().(dtos.QueryProductFilter)

	now := time.Now()

	if filter.Start != nil {
		if filter.Start.After(now) {
			sl.ReportError(filter.Start, "Start", "start", "lte_now", "")
		}
	}

	if filter.End != nil {
		if filter.End.After(now) {
			sl.ReportError(filter.End, "End", "end", "lte_now", "")
		}
		if filter.Start != nil && filter.End.Before(*filter.Start) {
			sl.ReportError(filter.End, "End", "end", "gtfield=Start", "")
		}
	}
}
