package data

import (
	"math"
	"strings"
)

// Define a Metadata struct for holding the pagination metadata.
type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
}

//	Check that the client-provided Sort field matches one of the entries in our safelist.
//
// And if it does, extract the column name from the Sort field by stripping the leading
// hyphen character (if one exists).
func (f Filters) sortColumn() string {
	for _, safeSort := range f.SortSafeList {
		if f.Sort == safeSort {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	panic("unsafe sort: " + f.Sort)
}

// Return the sort direction ("ASC" or "DESC")
func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}

// Get limit for put in sql query
func (f Filters) limit() int {
	return f.PageSize
}

// Get start point of movies base on query parameters that provided
// by user
func (f Filters) offset() int {
	return (f.Page - 1) * f.limit()
}

// The calculateMetadata() function calculates the appropriate pagination metadata
// values given the total number of records, current page, and page size values.
func calculateMetadata(totalRecords int, page int, pageSize int) Metadata {
	if totalRecords == 0 {
		// Return an empty Metadata struct if there are no records.
		return Metadata{}
	}

	// Create Metadata{} and return it.
	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}
