package ptr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tbe-team/raybot/pkg/ptr"
)

func TestNew(t *testing.T) {
	t.Run("string pointer", func(t *testing.T) {
		input := "test"
		got := ptr.New(input)
		assert.NotNil(t, got)
		assert.Equal(t, input, *got)
	})

	t.Run("int pointer", func(t *testing.T) {
		input := 42
		got := ptr.New(input)
		assert.NotNil(t, got)
		assert.Equal(t, input, *got)
	})

	t.Run("bool pointer", func(t *testing.T) {
		input := true
		got := ptr.New(input)
		assert.NotNil(t, got)
		assert.Equal(t, input, *got)
	})

	t.Run("struct pointer", func(t *testing.T) {
		input := struct{ name string }{name: "test"}
		got := ptr.New(input)
		assert.NotNil(t, got)
		assert.Equal(t, input, *got)
	})
}
