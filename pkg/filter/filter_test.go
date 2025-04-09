package filter

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	f := New("name", TypeEQ, "test")
	assert.Equal(t, "name", f.column)
	assert.Equal(t, TypeEQ, f.ftype)
	assert.Equal(t, "test", f.value)
	assert.Equal(t, OpAnd, f.operator)
	assert.Empty(t, f.filters)
}

func TestSetOperator(t *testing.T) {
	f := New("name", TypeEQ, "test")
	f.SetOperator(OpOr)
	assert.Equal(t, OpOr, f.operator)
}

func TestWithFilters(t *testing.T) {
	f := New("name", TypeEQ, "test")
	f2 := New("age", TypeGT, 18)
	f.WithFilters(f2)
	assert.Len(t, f.filters, 1)
	assert.Equal(t, "age", f.filters[0].column)
}

func TestAdd(t *testing.T) {
	f := New("name", TypeEQ, "test")
	f.Add("age", TypeGT, 18)
	assert.Len(t, f.filters, 1)
	assert.Equal(t, "age", f.filters[0].column)
}

func TestCondition(t *testing.T) {
	tests := []struct {
		name     string
		filter   Filter
		expected sq.Sqlizer
	}{
		{
			name:     "EQ",
			filter:   New("name", TypeEQ, "test"),
			expected: sq.Eq{"name": "test"},
		},
		{
			name:     "NotEQ",
			filter:   New("name", TypeNotEQ, "test"),
			expected: sq.NotEq{"name": "test"},
		},
		{
			name:     "GT",
			filter:   New("age", TypeGT, 18),
			expected: sq.Gt{"age": 18},
		},
		{
			name:     "GTE",
			filter:   New("age", TypeGTE, 18),
			expected: sq.GtOrEq{"age": 18},
		},
		{
			name:     "LT",
			filter:   New("age", TypeLT, 18),
			expected: sq.Lt{"age": 18},
		},
		{
			name:     "LTE",
			filter:   New("age", TypeLTE, 18),
			expected: sq.LtOrEq{"age": 18},
		},
		{
			name:     "Like",
			filter:   New("name", TypeLike, "%test%"),
			expected: sq.Like{"name": "%test%"},
		},
		{
			name:     "NotLike",
			filter:   New("name", TypeNotLike, "%test%"),
			expected: sq.NotLike{"name": "%test%"},
		},
		{
			name:     "ILike",
			filter:   New("name", TypeILike, "%test%"),
			expected: sq.ILike{"name": "%test%"},
		},
		{
			name:     "NotILike",
			filter:   New("name", TypeNotILike, "%test%"),
			expected: sq.NotILike{"name": "%test%"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.filter.condition()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestGetConditions(t *testing.T) {
	tests := []struct {
		name     string
		filter   *Filter
		expected sq.Sqlizer
	}{
		{
			name:     "Single condition",
			filter:   &Filter{column: "name", ftype: TypeEQ, value: "test"},
			expected: sq.Eq{"name": "test"},
		},
		{
			name: "Multiple conditions with AND",
			filter: &Filter{
				column:   "name",
				ftype:    TypeEQ,
				value:    "test",
				operator: OpAnd,
				filters: []Filter{
					{column: "age", ftype: TypeGT, value: 18},
					{column: "active", ftype: TypeEQ, value: true},
				},
			},
			expected: sq.And{
				sq.Eq{"name": "test"},
				sq.Gt{"age": 18},
				sq.Eq{"active": true},
			},
		},
		{
			name: "Multiple conditions with OR",
			filter: &Filter{
				column:   "name",
				ftype:    TypeEQ,
				value:    "test",
				operator: OpOr,
				filters: []Filter{
					{column: "age", ftype: TypeGT, value: 18},
					{column: "active", ftype: TypeEQ, value: true},
				},
			},
			expected: sq.Or{
				sq.Eq{"name": "test"},
				sq.Gt{"age": 18},
				sq.Eq{"active": true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.filter.getConditions()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestAttach(t *testing.T) {
	f := New("name", TypeEQ, "test")
	builder := sq.Select("*").From("users")
	result := f.Attach(builder)

	sql, args, err := result.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT * FROM users WHERE name = ?", sql)
	assert.Equal(t, []interface{}{"test"}, args)
}
