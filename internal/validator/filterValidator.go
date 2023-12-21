package validator

import "Meow/internal/data"

// The ValidateFilters() function
// validate page - pageSize - sort query parameters.
func ValidateFilters(v *Validator, f data.Filters) {
	// Check page & page_size parameters in query parameters.
	v.Check(f.Page > 0, "page", "must be greated than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be less than 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be less than 100")

	// Check the sort parameter matches a value from safe list.
	v.Check(In(f.Sort, f.SortSafeList...), "sort", "invalid sort value")
}
