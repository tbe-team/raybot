package paging_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tbe-team/raybot/pkg/paging"
)

func TestListCreation(t *testing.T) {
	items := []string{"item1", "item2", "item3"}
	totalItems := int64(100)

	list := paging.NewList(items, totalItems)

	assert.Equal(t, items, list.Items)
	assert.Equal(t, totalItems, list.TotalItems)
}
