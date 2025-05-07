package espserial

import (
	"bytes"
	"time"

	"github.com/stretchr/testify/mock"
	"go.bug.st/serial"
)

var _ serial.Port = (*FakeSerialPort)(nil)

type FakeSerialPort struct {
	mock.Mock
	ReadBuffer  bytes.Buffer
	WriteBuffer bytes.Buffer
}

func (m *FakeSerialPort) SetMode(mode *serial.Mode) error {
	args := m.Called(mode)
	return args.Error(0)
}

func (m *FakeSerialPort) Read(p []byte) (int, error) {
	return m.ReadBuffer.Read(p)
}

func (m *FakeSerialPort) Write(p []byte) (int, error) {
	n, err := m.WriteBuffer.Write(p)
	return n, err
}

func (m *FakeSerialPort) Drain() error {
	args := m.Called()
	return args.Error(0)
}

func (m *FakeSerialPort) ResetInputBuffer() error {
	args := m.Called()
	return args.Error(0)
}

func (m *FakeSerialPort) ResetOutputBuffer() error {
	args := m.Called()
	return args.Error(0)
}

func (m *FakeSerialPort) SetDTR(dtr bool) error {
	args := m.Called(dtr)
	return args.Error(0)
}

func (m *FakeSerialPort) SetRTS(rts bool) error {
	args := m.Called(rts)
	return args.Error(0)
}

func (m *FakeSerialPort) GetModemStatusBits() (*serial.ModemStatusBits, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*serial.ModemStatusBits), args.Error(1)
}

func (m *FakeSerialPort) SetReadTimeout(t time.Duration) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *FakeSerialPort) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *FakeSerialPort) Break(d time.Duration) error {
	args := m.Called(d)
	return args.Error(0)
}
