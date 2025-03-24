package sort_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tbe-team/raybot/pkg/sort"
)

func TestNewList(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expected   []sort.Sort
		expectErr  bool
		errMessage string
	}{
		{
			name:  "Single column, ascending",
			input: "created_at",
			expected: []sort.Sort{
				{Col: "created_at", Order: sort.OrderASC},
			},
			expectErr: false,
		},
		{
			name:  "Single column, descending",
			input: "-updated_at",
			expected: []sort.Sort{
				{Col: "updated_at", Order: sort.OrderDESC},
			},
			expectErr: false,
		},
		{
			name:  "Multiple columns",
			input: "name,-created_at,updated_at",
			expected: []sort.Sort{
				{Col: "name", Order: sort.OrderASC},
				{Col: "created_at", Order: sort.OrderDESC},
				{Col: "updated_at", Order: sort.OrderASC},
			},
			expectErr: false,
		},
		{
			name:      "Empty input",
			input:     "",
			expected:  []sort.Sort{},
			expectErr: false,
		},
		{
			name:       "Invalid input with trailing space",
			input:      "name ,created_at",
			expectErr:  true,
			errMessage: "invalid sort column",
		},
		{
			name:       "Invalid input with leading space",
			input:      " -created_at",
			expectErr:  true,
			errMessage: "invalid sort column",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sort.NewListFromString(tt.input)

			if tt.expectErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, sort.ErrInvalidSort, "unexpected error occurred")
			} else {
				require.NoError(t, err, "unexpected error occurred")
				assert.Equal(t, tt.expected, result, "result mismatch for input '%s'", tt.input)
			}
		})
	}
}
