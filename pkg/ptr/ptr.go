package ptr

// New creates and returns a pointer to the provided value.
// It's a generic function that works with any type T.
func New[T any](v T) *T { return &v }
