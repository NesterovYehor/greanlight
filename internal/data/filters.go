package data

import (
	"math"
	"strings"

	"greenlight.nesty.net/internal/validator"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

func ValidateFilters(v *validator.Validator, input *Filters) {
	v.Check(input.Page > 0, "page", "must be greater than zero")
	v.Check(input.Page <= 10_000_000, "page", "must be maximum of 10 milion")

	v.Check(input.PageSize > 0, "page_size", "must be greater thn zero")
	v.Check(input.PageSize >= 100, "page_size", "must be maximum of 100")

	v.Check(validator.In(input.Sort, input.SortSafelist...), "sort", "invalid sort value")
}

func (filter *Filters) sortColumn() string {
	for _, safeValue := range filter.SortSafelist {
		if safeValue == filter.Sort {
			return strings.TrimPrefix(filter.Sort, "-")
		}
	}
	panic("unsafe sort parameter: " + filter.Sort)
}

func (filter *Filters) sortDirection() string {
	if strings.HasPrefix(filter.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (filter *Filters) offset() int {
	return (filter.PageSize - 1) * filter.Page
}

func (filter *Filters) limit() int {
	return filter.PageSize
}

func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		TotalRecords: totalRecords,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
	}
}
