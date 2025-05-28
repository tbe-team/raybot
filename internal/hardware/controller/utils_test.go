package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newShortID(t *testing.T) {
	id := newShortID()
	assert.NotEmpty(t, id)
	assert.Len(t, id, 6)
}

func Test_removeMarkers(t *testing.T) {
	data := []byte(">abc\r\n")
	assert.Equal(t, []byte("abc"), removeMarkers(data))
}

func Test_boolToUint8(t *testing.T) {
	assert.Equal(t, uint8(1), boolToUint8(true))
	assert.Equal(t, uint8(0), boolToUint8(false))
}
