package espserial

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tbe-team/raybot/internal/config"
)

func TestClientWrite(t *testing.T) {
	t.Run("Should write successfully", func(t *testing.T) {
		mockPort := &FakeSerialPort{}
		client := NewClient(config.Serial{})
		client.port = mockPort

		data := []byte(`{"cmd":"test"}`)
		err := client.Write(context.Background(), data)

		assert.NoError(t, err)

		// Check the written bytes
		expected := append([]byte(">"), data...)
		expected = append(expected, '\r', '\n')
		assert.Equal(t, expected, mockPort.WriteBuffer.Bytes())
	})

	t.Run("Should return error when not connected", func(t *testing.T) {
		client := NewClient(config.Serial{})
		client.port = nil

		data := []byte(`{"cmd":"test"}`)
		err := client.Write(context.Background(), data)

		assert.Error(t, err)
		assert.Equal(t, ErrESPSerialNotConnected, err)
	})

	t.Run("Should return context cancelled error", func(t *testing.T) {
		mockPort := &FakeSerialPort{}
		client := NewClient(config.Serial{})
		client.port = mockPort

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel the context immediately

		data := []byte(`{"cmd":"test"}`)
		err := client.Write(ctx, data)

		assert.Error(t, err)
		assert.Equal(t, context.Canceled, err)
	})
}

func TestClientRead(t *testing.T) {
	t.Run("Should return simple data successfully", func(t *testing.T) {
		mockPort := &FakeSerialPort{}
		client := NewClient(config.Serial{})
		client.port = mockPort

		mockPort.ReadBuffer.Write([]byte(">data\r\n"))

		res, err := client.Read(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, []byte("data"), res)
	})

	t.Run("Should return double start markers", func(t *testing.T) {
		mockPort := &FakeSerialPort{}
		client := NewClient(config.Serial{})
		client.port = mockPort

		mockPort.ReadBuffer.Write([]byte(">>data\r\n"))

		res, err := client.Read(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, []byte(">data"), res)
	})

	t.Run("Should return null bytes in data", func(t *testing.T) {
		mockPort := &FakeSerialPort{}
		client := NewClient(config.Serial{})
		client.port = mockPort

		mockPort.ReadBuffer.Write([]byte(">data\x00\r\n"))

		res, err := client.Read(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, []byte("data\x00"), res)
	})

	t.Run("Should return start marker in data", func(t *testing.T) {
		mockPort := &FakeSerialPort{}
		client := NewClient(config.Serial{})
		client.port = mockPort

		mockPort.ReadBuffer.Write([]byte(">>data\x00\r\n"))

		res, err := client.Read(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, []byte(">data\x00"), res)
	})

	t.Run("Should return multiple line endings", func(t *testing.T) {
		mockPort := &FakeSerialPort{}
		client := NewClient(config.Serial{})
		client.port = mockPort

		mockPort.ReadBuffer.Write([]byte(">>data\x00\r\n\r\n"))

		res, err := client.Read(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, []byte(">data\x00"), res)
	})

	t.Run("Should return random prefix", func(t *testing.T) {
		mockPort := &FakeSerialPort{}
		client := NewClient(config.Serial{})
		client.port = mockPort

		mockPort.ReadBuffer.Write([]byte("random>message here\r\ngarbage"))

		res, err := client.Read(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, []byte("message here"), res)
	})

	t.Run("Should return null bytes after message", func(t *testing.T) {
		mockPort := &FakeSerialPort{}
		client := NewClient(config.Serial{})
		client.port = mockPort

		mockPort.ReadBuffer.Write([]byte(">null\r\n\x00\x00"))

		res, err := client.Read(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, []byte("null"), res)
	})

	t.Run("Should return empty message", func(t *testing.T) {
		mockPort := &FakeSerialPort{}
		client := NewClient(config.Serial{})
		client.port = mockPort

		mockPort.ReadBuffer.Write([]byte(">\r\n"))

		res, err := client.Read(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, []byte(""), res)
	})

	t.Run("Should return not connected error", func(t *testing.T) {
		client := NewClient(config.Serial{})
		client.port = nil

		res, err := client.Read(context.Background())
		assert.Error(t, err)
		assert.Equal(t, ErrESPSerialNotConnected, err)
		assert.Nil(t, res)
	})

	t.Run("Should return EOF error on missing end marker", func(t *testing.T) {
		mockPort := &FakeSerialPort{}
		client := NewClient(config.Serial{})
		client.port = mockPort

		mockPort.ReadBuffer.Write([]byte(">data without end"))

		res, err := client.Read(context.Background())
		assert.Error(t, err)
		assert.Equal(t, io.EOF, err)
		assert.Nil(t, res)
	})

	t.Run("Should return EOF error on missing start marker", func(t *testing.T) {
		mockPort := &FakeSerialPort{}
		client := NewClient(config.Serial{})
		client.port = mockPort

		mockPort.ReadBuffer.Write([]byte("data only"))

		res, err := client.Read(context.Background())
		assert.Error(t, err)
		assert.Equal(t, io.EOF, err)
		assert.Nil(t, res)
	})

	t.Run("Should return context cancelled error", func(t *testing.T) {
		mockPort := &FakeSerialPort{}
		client := NewClient(config.Serial{})
		client.port = mockPort

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		res, err := client.Read(ctx)
		assert.Error(t, err)
		assert.Equal(t, context.Canceled, err)
		assert.Nil(t, res)
	})
}
