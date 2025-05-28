package picserial

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tbe-team/raybot/internal/config"
)

func TestController(t *testing.T) {
	t.Run("Set cargo position successfully", func(t *testing.T) {
		mockPort := &FakeSerialPort{}
		client := NewClient(config.Serial{})
		client.port = mockPort

		err := client.SetCargoPosition(context.Background(), "abc", 100, 10)
		assert.NoError(t, err)

		actual := removeMarkers(mockPort.WriteBuffer.Bytes())
		expected := `
		{
			"id":"abc",
			"type":2,
			"data":{
				"target_position":10,
				"max_output":100,
				"enable":1
			}
		}`
		assert.JSONEq(t, expected, string(actual))
	})

	t.Run("Stop lift cargo motor successfully", func(t *testing.T) {
		mockPort := &FakeSerialPort{}
		client := NewClient(config.Serial{})
		client.port = mockPort

		err := client.StopLiftCargoMotor(context.Background(), "abc")
		assert.NoError(t, err)

		actual := removeMarkers(mockPort.WriteBuffer.Bytes())
		expected := `
		{
			"id":"abc",
			"type":2,
			"data":{
				"enable":0,
				"max_output":0,
				"target_position":0
			}
		}`
		assert.JSONEq(t, expected, string(actual))
	})

	t.Run("Move forward successfully", func(t *testing.T) {
		mockPort := &FakeSerialPort{}
		client := NewClient(config.Serial{})
		client.port = mockPort

		err := client.MoveForward(context.Background(), "abc", 10)
		assert.NoError(t, err)

		actual := removeMarkers(mockPort.WriteBuffer.Bytes())
		expected := `
		{
			"id":"abc",
			"type":3,
			"data":{
				"direction":0,
				"speed":10,
				"enable":1
			}
		}`
		assert.JSONEq(t, expected, string(actual))
	})

	t.Run("Move backward successfully", func(t *testing.T) {
		mockPort := &FakeSerialPort{}
		client := NewClient(config.Serial{})
		client.port = mockPort

		err := client.MoveBackward(context.Background(), "abc", 10)
		assert.NoError(t, err)

		actual := removeMarkers(mockPort.WriteBuffer.Bytes())
		expected := `
		{
			"id":"abc",
			"type":3,
			"data":{
				"direction":1,
				"speed":10,
				"enable":1
			}
		}`
		assert.JSONEq(t, expected, string(actual))
	})

	t.Run("Stop drive motor successfully", func(t *testing.T) {
		mockPort := &FakeSerialPort{}
		client := NewClient(config.Serial{})
		client.port = mockPort

		err := client.StopDriveMotor(context.Background(), "abc", MoveDirectionForward)
		assert.NoError(t, err)

		actual := removeMarkers(mockPort.WriteBuffer.Bytes())
		expected := `
		{
			"id":"abc",
			"type":3,
			"data":{
				"direction":0,
				"speed":0,
				"enable":0
			}
		}`
		assert.JSONEq(t, expected, string(actual))
	})

	t.Run("Config battery charge successfully", func(t *testing.T) {
		mockPort := &FakeSerialPort{}
		client := NewClient(config.Serial{})
		client.port = mockPort

		err := client.ConfigBatteryCharge(context.Background(), "abc", 10, true)
		assert.NoError(t, err)

		actual := removeMarkers(mockPort.WriteBuffer.Bytes())
		expected := `
		{
			"id":"abc",
			"type":0,
			"data":{
				"current_limit":10,
				"enable":1
			}
		}`
		assert.JSONEq(t, expected, string(actual))
	})

	t.Run("Config battery discharge successfully", func(t *testing.T) {
		mockPort := &FakeSerialPort{}
		client := NewClient(config.Serial{})
		client.port = mockPort

		err := client.ConfigBatteryDischarge(context.Background(), "abc", 10, true)
		assert.NoError(t, err)

		actual := removeMarkers(mockPort.WriteBuffer.Bytes())
		expected := `
		{
			"id":"abc",
			"type":1,
			"data":{
				"current_limit":10,
				"enable":1
			}
		}`
		assert.JSONEq(t, expected, string(actual))
	})
}

// removeMarkers removes the markers
func removeMarkers(data []byte) []byte {
	b := data[1 : len(data)-2] // remove > and \r\n
	return b
}
