package paging_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tbe-team/raybot/pkg/paging"
)

func TestNewParams_CustomValues(t *testing.T) {
	params := paging.NewParams(paging.Page(3), paging.PageSize(20))

	assert.Equal(t, uint(20), params.PageSize)
	assert.Equal(t, uint(3), params.Page)
	assert.Equal(t, uint(40), params.Offset())
	assert.Equal(t, uint(20), params.Limit())
}

func TestNewParams_WithOptions(t *testing.T) {
	params := paging.NewParams(
		paging.Page(50), paging.PageSize(5),
		paging.WithMaxPageSize(3),
	)

	assert.Equal(t, uint(3), params.PageSize)
	assert.Equal(t, uint(50), params.Page)
	assert.Equal(t, uint(147), params.Offset())
}

func TestNewParams_WithZeroPage(t *testing.T) {
	params := paging.NewParams(paging.Page(0), paging.PageSize(15))

	assert.Equal(t, uint(15), params.PageSize)
	assert.Equal(t, uint(0), params.Offset())
	assert.Equal(t, uint(15), params.Limit())
}
