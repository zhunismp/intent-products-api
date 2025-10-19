package transformer

import (
	"fmt"
	"net/url"
	"time"

	"github.com/zhunismp/intent-products-api/internal/core/dtos"
	"github.com/zhunismp/intent-products-api/internal/interfaces/http/transport"
)

const hardcodedUserID = "101234567890123456789"

func ToCreateProductInput(req transport.CreateProductRequest) (dtos.CreateProductInput, error) {
	return dtos.CreateProductInput{
		OwnerID: hardcodedUserID, // hardcoded google oauth's sub
		Title:   req.Title,
		Price:   req.Price,
		Link:    req.Link,
		Reasons: req.Reasons,
	}, nil
}

func ToQueryProductInput(q url.Values) (dtos.QueryProductInput, error) {
	startStr := q.Get("start")
	endStr := q.Get("end")
	statusStr := q.Get("status")

	// TODO: move to opt logic somewhere.
	var startOpt *time.Time
	var endOpt *time.Time
	var statusOpt *string

	if startStr != "" {
		t, err := time.Parse(time.RFC3339, startStr)
		if err != nil {
			return dtos.QueryProductInput{}, fmt.Errorf("invalid 'start' format: %w", err)
		}
		startOpt = &t
	}

	if endStr != "" {
		t, err := time.Parse(time.RFC3339, endStr)
		if err != nil {
			return dtos.QueryProductInput{}, fmt.Errorf("invalid 'end' format: %w", err)
		}
		endOpt = &t
	}

	if statusStr != "" {
		statusOpt = &statusStr
	}

	queryProductFilters := dtos.QueryProductFilter{
		Start:  startOpt,
		End:    endOpt,
		Status: statusOpt,
	}

	return dtos.QueryProductInput{
		OwnerID: hardcodedUserID,
		Filters: queryProductFilters,
	}, nil
}
