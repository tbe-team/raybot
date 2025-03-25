package paging

// List represents a list of items with total item.
type List[T any] struct {
	Items      []T
	TotalItems int64
}

// NewList creates a new List instance with total item.
func NewList[T any](items []T, totalItems int64) List[T] {
	return List[T]{
		Items:      items,
		TotalItems: totalItems,
	}
}
