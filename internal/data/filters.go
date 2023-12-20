package data

import "strings"

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
