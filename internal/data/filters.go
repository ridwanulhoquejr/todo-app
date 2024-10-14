package data

import (
	"strings"
	"time"

	"github.com/ridwanulhoquejr/todo-app/internal/validator"
)

// Pagination holds pagination information
type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// Filters holds filtering criteria
type Filters struct {
	Completed bool      `json:"completed"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// Search holds the search criteria
type Search struct {
	Title string `json:"title"`
}

// Sorts holds sorting information
type Sorts struct {
	Sort     string   `json:"sort"`
	SafeList []string `json:"safe_list"`
}

// Queries wraps all filtering, sorting, pagination, and search options
type Queries struct {
	Pagination Pagination `json:"pagination"`
	Filters    Filters    `json:"filters"`
	Search     Search     `json:"search"`
	Sorts      Sorts      `json:"sorts"`
}

func ValidateQueries(v *validator.Validator, q Queries) {
	v.Check(q.Pagination.Page > 0, "page", "must be greater than zero")
	v.Check(q.Pagination.Page <= 10_100_000, "page", "must be maximum of 10 million")
	v.Check(q.Pagination.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(q.Pagination.PageSize <= 100, "page_size", "must be maximum of 100")

	// chekc the sorts.safelist
	v.Check(validator.In(q.Sorts.Sort, q.Sorts.SafeList...), "sort", "invalid sort value")

	// check the Filters of star and end date
	// v.Check(!(q.Filters.StartDate.After(q.Filters.EndDate)), "start_date", "start_date must be less or equal to the current date")
	// v.Check(q.Filters.StartDate.Before(time.Now().AddDate(-1, 0, -1)), "start_date", "start_date must be in between less than 1 year of current date")
}

// Check that the client-provided Sort field matches one of the entries in our safelist
// and if it does, extract the column name from the Sort field by stripping the leading
// hyphen character (if one exists)
func (s Sorts) sortColumn() string {
	for _, safeValue := range s.SafeList {
		if s.Sort == safeValue {
			return strings.TrimPrefix(s.Sort, "-")
		}
	}
	panic("unsafe sort parameters: " + s.Sort)
}

func (s Sorts) sortDirection() string {
	if strings.HasPrefix(s.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (p Pagination) limit() int {
	return p.PageSize
}

func (p Pagination) offset() int {
	return ((p.Page - 1) * p.PageSize)
}
