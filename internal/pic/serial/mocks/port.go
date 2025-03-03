package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"
	"go.bug.st/serial"
)

var _ serial.Port = (*FakePort)(nil)

// FakePort implements serial.Port for testing
type FakePort struct {
	mock.Mock
}

func (m *FakePort) SetMode(mode *serial.Mode) error {
	args := m.Called(mode)
	return args.Error(0)
}

func (m *FakePort) SetReadTimeout(t time.Duration) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *FakePort) Read(p []byte) (int, error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func (m *FakePort) Write(p []byte) (int, error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func (m *FakePort) ResetInputBuffer() error {
	args := m.Called()
	return args.Error(0)
}

func (m *FakePort) ResetOutputBuffer() error {
	args := m.Called()
	return args.Error(0)
}

func (m *FakePort) SetDTR(dtr bool) error {
	args := m.Called(dtr)
	return args.Error(0)
}

func (m *FakePort) SetRTS(rts bool) error {
	args := m.Called(rts)
	return args.Error(0)
}

func (m *FakePort) GetModemStatusBits() (*serial.ModemStatusBits, error) {
	args := m.Called()
	return args.Get(0).(*serial.ModemStatusBits), args.Error(1)
}

func (m *FakePort) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *FakePort) Break(d time.Duration) error {
	args := m.Called(d)
	return args.Error(0)
}

func (m *FakePort) Drain() error {
	args := m.Called()
	return args.Error(0)
}
