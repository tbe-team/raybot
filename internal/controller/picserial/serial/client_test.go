package serial

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/tbe-team/raybot/internal/controller/picserial/serial/mocks"
)

func TestWrite(t *testing.T) {
	mockPort := new(mocks.FakePort)
	client := &client{
		port: mockPort,
	}

	testData := []byte("test data")
	expectedWriteData := []byte(">test data\r\n")

	mockPort.On("Write", expectedWriteData).Return(len(expectedWriteData), nil)

	err := client.Write(testData)
	assert.NoError(t, err)
	mockPort.AssertExpectations(t)
}

func TestWriteError(t *testing.T) {
	mockPort := new(mocks.FakePort)
	client := &client{
		port: mockPort,
	}

	testData := []byte("test data")
	expectedWriteData := []byte(">test data\r\n")

	mockPort.On("Write", expectedWriteData).Return(0, errors.New("write error"))

	err := client.Write(testData)
	assert.Error(t, err)
	assert.Equal(t, "write error", err.Error())
	mockPort.AssertExpectations(t)
}

func TestRead(t *testing.T) {
	mockPort := new(mocks.FakePort)
	client := &client{
		port: mockPort,
	}

	response := []byte(">response data\r\n")
	mockPort.On("Read", mock.Anything).Return(len(response), nil).Run(func(args mock.Arguments) {
		buf := args.Get(0).([]byte)
		copy(buf, response)
	}).Once()

	result, err := client.read()
	assert.NoError(t, err)
	assert.Equal(t, []byte("response data"), result)
	mockPort.AssertExpectations(t)
}

func TestReadMultipleChunks(t *testing.T) {
	mockPort := new(mocks.FakePort)
	client := &client{
		port: mockPort,
	}

	mockPort.On("Read", mock.Anything).Return(5, nil).Run(func(args mock.Arguments) {
		buf := args.Get(0).([]byte)
		copy(buf, []byte(">resp"))
	}).Once()

	// Second chunk has more data but still no suffix
	mockPort.On("Read", mock.Anything).Return(10, nil).Run(func(args mock.Arguments) {
		buf := args.Get(0).([]byte)
		copy(buf, []byte("onse data"))
	}).Once()

	// Third chunk has the suffix
	mockPort.On("Read", mock.Anything).Return(2, nil).Run(func(args mock.Arguments) {
		buf := args.Get(0).([]byte)
		copy(buf, []byte("\r\n"))
	}).Once()

	result, err := client.read()
	assert.NoError(t, err)
	assert.Equal(t, []byte("response data"), result)
	mockPort.AssertExpectations(t)
}

func TestReadError(t *testing.T) {
	mockPort := new(mocks.FakePort)
	client := &client{
		port: mockPort,
	}

	mockPort.On("Read", mock.Anything).Return(0, errors.New("read error"))

	result, err := client.read()
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "read error", err.Error())
	mockPort.AssertExpectations(t)
}
