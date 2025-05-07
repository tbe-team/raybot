package espserial

import (
	"context"
	"fmt"
	"sync"

	"go.bug.st/serial"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/pkg/xerror"
)

var ErrESPSerialNotConnected = xerror.NotFound(nil, "espserial.notConnected", "ESP serial not connected")

const readBufferSize = 64

type Client interface {
	Open() error
	Close() error
	Connected() bool
	Write(ctx context.Context, data []byte) error
	Read(ctx context.Context) ([]byte, error)
}

type DefaultClient struct {
	cfg config.Serial

	port serial.Port
	mode serial.Mode

	writeMu sync.Mutex
}

func NewClient(cfg config.Serial) *DefaultClient {
	mode := serial.Mode{
		BaudRate: cfg.BaudRate,
		DataBits: int(cfg.DataBits),
	}

	switch cfg.StopBits {
	case 1:
		mode.StopBits = serial.OneStopBit
	case 1.5:
		mode.StopBits = serial.OnePointFiveStopBits
	case 2:
		mode.StopBits = serial.TwoStopBits
	}

	switch cfg.Parity {
	case "NONE":
		mode.Parity = serial.NoParity
	case "ODD":
		mode.Parity = serial.OddParity
	case "EVEN":
		mode.Parity = serial.EvenParity
	}

	return &DefaultClient{
		cfg:  cfg,
		mode: mode,
	}
}

func (c *DefaultClient) Open() error {
	port, err := serial.Open(c.cfg.Port, &c.mode)
	if err != nil {
		return fmt.Errorf("failed to open serial port: %w", err)
	}

	if err := port.SetReadTimeout(c.cfg.ReadTimeout); err != nil {
		return fmt.Errorf("failed to set read timeout: %w", err)
	}

	c.port = port
	return nil
}

func (c *DefaultClient) Close() error {
	if c.port == nil {
		return ErrESPSerialNotConnected
	}

	return c.port.Close()
}

func (c *DefaultClient) Connected() bool {
	return c.port != nil
}

func (c *DefaultClient) Write(ctx context.Context, data []byte) error {
	if c.port == nil {
		return ErrESPSerialNotConnected
	}

	data = append([]byte(">"), data...)
	data = append(data, '\r', '\n')

	done := make(chan error, 1)
	go func() {
		c.writeMu.Lock()
		defer c.writeMu.Unlock()

		_, err := c.port.Write(data)
		done <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

// Read reads data from the serial port.
func (c *DefaultClient) Read(ctx context.Context) ([]byte, error) {
	if c.port == nil {
		return nil, ErrESPSerialNotConnected
	}

	return c.read(ctx)
}

// read continuously reads from the port until a complete message is received.
// A complete message starts with '>' and ends with CR LF (\r\n).
// The message is returned without the prefix and suffix
func (c *DefaultClient) read(ctx context.Context) ([]byte, error) {
	var msg []byte
	isMsgStarted := false

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()

		default:
			buf := make([]byte, readBufferSize)
			n, err := c.port.Read(buf)
			if err != nil {
				return nil, err
			}

			chunk := buf[:n]

			for _, b := range chunk {
				if !isMsgStarted {
					if b == '>' {
						isMsgStarted = true
					}
					continue
				}

				msg = append(msg, b)
				msgLength := len(msg)

				// check end marker
				if msgLength >= 2 && msg[msgLength-2] == '\r' && msg[msgLength-1] == '\n' {
					msg = msg[:msgLength-2] // remove \r\n
					return msg, nil
				}
			}
		}
	}
}
