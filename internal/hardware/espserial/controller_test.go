package espserial

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tbe-team/raybot/internal/config"
)

func TestController(t *testing.T) {
	t.Run("Open cargo door successfully", func(t *testing.T) {
		mockPort := &FakeSerialPort{}
		client := NewClient(config.Serial{})
		client.port = mockPort

		err := client.OpenCargoDoor(context.Background(), 10)
		assert.NoError(t, err)

		actual := overrideIDAndRemoveMarkers(t, mockPort.WriteBuffer.Bytes(), "abc")
		expected := `
		{
			"id":"abc",
			"type":0,
			"data":{
				"state":1,
				"speed":10,
				"enable":1
			}
		}`
		assert.JSONEq(t, expected, string(actual))
	})

	t.Run("Close cargo door successfully", func(t *testing.T) {
		mockPort := &FakeSerialPort{}
		client := NewClient(config.Serial{})
		client.port = mockPort

		err := client.CloseCargoDoor(context.Background(), 10)
		assert.NoError(t, err)

		actual := overrideIDAndRemoveMarkers(t, mockPort.WriteBuffer.Bytes(), "abc")
		expected := `
		{
			"id":"abc",
			"type":0,
			"data":{
				"state":0,
				"speed":10,
				"enable":1
			}
		}`
		assert.JSONEq(t, expected, string(actual))
	})
}

// overrideIDAndRemoveMarkers overrides the id of the command, and remove the markers
func overrideIDAndRemoveMarkers(t *testing.T, data []byte, id string) []byte {
	b := data[1 : len(data)-2] // remove > and \r\n

	var m map[string]any
	err := json.Unmarshal(b, &m)
	assert.NoError(t, err)
	m["id"] = id

	b, err = json.Marshal(m)
	assert.NoError(t, err)

	return b
}
