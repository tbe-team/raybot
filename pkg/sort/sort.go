package sort

import (
	"strings"

	sq "github.com/Masterminds/squirrel"

	"github.com/tbe-team/raybot/pkg/xerror"
)

var (
	ErrInvalidSort = xerror.BadRequest(nil, "sort.invalid", "Invalid sort")
)

const (
	OrderASC  = "ASC"
	OrderDESC = "DESC"
)

type Sort struct {
	Col   string
	Order string
}

// Attach attaches sort to builder.
func (s Sort) Attach(b sq.SelectBuilder) sq.SelectBuilder {
	return b.OrderBy(s.Col + " " + s.Order)
}

// NewListFromString creates a list of Sort from a string.
func NewListFromString(s string) ([]Sort, error) {
	if s == "" {
		return []Sort{}, nil
	}
	orderBys := strings.Split(s, ",")
	sorts := make([]Sort, len(orderBys))

	for i, r := range orderBys {
		orderBy := r
		if strings.HasPrefix(orderBy, " ") || strings.HasSuffix(orderBy, " ") {
			return nil, ErrInvalidSort
		}

		order := OrderASC
		if strings.HasPrefix(orderBy, "-") {
			order = OrderDESC
			orderBy = orderBy[1:]
		}

		sorts[i] = Sort{
			Col:   orderBy,
			Order: order,
		}
	}

	return sorts, nil
}
