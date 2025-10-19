package transformer

import (
	"fmt"
	"strconv"
	"time"

	"github.com/zhunismp/intent-products-api/internal/core/dtos"
)

type QueryProductInputFilterFunc func(*dtos.QueryProductInput) error

func NewQueryProductInput(userID string, filters ...QueryProductInputFilterFunc) (*dtos.QueryProductInput, error) {
	input := &dtos.QueryProductInput{
		OwnerID:    userID,
		Pagination: nil,
		Filters:    dtos.QueryProductFilter{},
	}

	for _, queryFilter := range filters {
		if err := queryFilter(input); err != nil {
			return &dtos.QueryProductInput{}, err
		}
	}

	return input, nil
}

func WithStart(startStr string) QueryProductInputFilterFunc {
	return func(input *dtos.QueryProductInput) error {
		if startStr == "" {
			return nil
		}
		t, err := time.Parse(time.RFC3339, startStr)
		if err != nil {
			return fmt.Errorf("invalid 'start' format: %w", err)
		}
		input.Filters.Start = &t
		return nil
	}
}

func WithEnd(endStr string) QueryProductInputFilterFunc {
	return func(input *dtos.QueryProductInput) error {
		if endStr == "" {
			return nil
		}
		t, err := time.Parse(time.RFC3339, endStr)
		if err != nil {
			return fmt.Errorf("invalid 'end' format: %w", err)
		}
		input.Filters.End = &t
		return nil
	}
}

func WithStatus(statusStr string) QueryProductInputFilterFunc {
	return func(input *dtos.QueryProductInput) error {
		if statusStr == "" {
			return nil
		}
		input.Filters.Status = &statusStr
		return nil
	}
}

func WithPagination(pageNumStr string, pageSizeStr string) QueryProductInputFilterFunc {
	return func(input *dtos.QueryProductInput) error {
		// if either num or size is null, set pagination to null
		if pageNumStr == "" || pageSizeStr == "" {
			input.Pagination = nil
			return nil
		}

		// parse string
		pageNum, pnerr := strconv.ParseInt(pageNumStr, 10, 64)
		pageSize, pserr := strconv.ParseInt(pageSizeStr, 10, 64)
		if pnerr != nil || pserr != nil {
			return fmt.Errorf("error while parsing pagination: %v %v", pnerr, pserr)
		}

		// set pagination is specify and valid
		input.Pagination = &dtos.Pagination{
			Page: int(pageNum),
			Size: int(pageSize),
		}

		return nil
	}
}
