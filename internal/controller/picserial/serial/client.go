package serial

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"go.bug.st/serial"
)

const (
	readBufferSize = 64
)

var (
	ErrPortNotOpen = errors.New("port not open")
)

// Client is the interface for the serial client.
type Client interface {
	// Write sends data to the serial port.
	// It is safe to call this method from multiple goroutines.
	Write(data []byte) error

	// Read reads data from the serial port.
	// It returns a channel that receives the data.
	// Data is formatted as JSON.
	Read() <-chan []byte

	// Stop stops the client.
	Stop() error
}

type client struct {
	cfg Config

	writeMu  sync.Mutex
	port     serial.Port
	readChan chan []byte
	log      *slog.Logger
	stop     chan struct{}
}

// NewClient creates a new serial client.
func NewClient(cfg Config, log *slog.Logger) (Client, error) {
	log = log.With(
		slog.Group("serial_client",
			slog.String("port", cfg.Port),
			slog.Int("baud_rate", cfg.BaudRate),
			slog.Int("data_bits", cfg.DataBits),
			slog.String("parity", cfg.Parity),
			slog.Float64("stop_bits", cfg.StopBits),
		),
	)
	mode := &serial.Mode{
		BaudRate: cfg.BaudRate,
		DataBits: cfg.DataBits,
	}

	switch cfg.StopBits {
	case 1:
		mode.StopBits = serial.OneStopBit
	case 1.5:
		mode.StopBits = serial.OnePointFiveStopBits
	case 2:
		mode.StopBits = serial.TwoStopBits
	default:
		return nil, fmt.Errorf("invalid stop bits: %f", cfg.StopBits)
	}

	switch strings.ToLower(cfg.Parity) {
	case "none":
		mode.Parity = serial.NoParity
	case "odd":
		mode.Parity = serial.OddParity
	case "even":
		mode.Parity = serial.EvenParity
	default:
		return nil, fmt.Errorf("invalid parity: %s", cfg.Parity)
	}

	var port serial.Port
	var openErr error

	port, openErr = serial.Open(cfg.Port, mode)
	if openErr != nil {
		// Now we just ignore the error
		log.Error("failed to open serial port", slog.Any("error", openErr))
		// return nil, err
	}

	if openErr == nil {
		if err := port.SetReadTimeout(cfg.ReadTimeout); err != nil {
			return nil, err
		}
	}

	client := &client{
		cfg:      cfg,
		port:     port,
		readChan: make(chan []byte, 1),
		stop:     make(chan struct{}),
		log:      log,
	}

	if openErr == nil {
		go client.readLoop()
	}

	return client, nil
}

// Write sends data to the serial port.
// It prefixes the data with '>' and suffixes it with CR LF (\r\n)
func (c *client) Write(data []byte) error {
	data = append([]byte(">"), data...)
	data = append(data, '\r', '\n')

	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	if c.port == nil {
		return ErrPortNotOpen
	}

	_, err := c.port.Write(data)
	return err
}

// Read reads data from the serial port.
// It returns a channel that receives the data.
// Data is formatted as JSON.
func (c *client) Read() <-chan []byte {
	return c.readChan
}

// Stop stops the client.
func (c *client) Stop() error {
	close(c.stop)
	if c.port != nil {
		return c.port.Close()
	}
	return nil
}

// readLoop reads from the serial port and sends the data to the read channel.
func (c *client) readLoop() {
	for {
		select {
		case <-c.stop:
			return
		default:
			data, err := c.read()
			if err != nil {
				c.log.Error("failed to read from serial port", "error", err)
				continue
			}

			select {
			case c.readChan <- data:
			default:
				c.log.Warn("read channel is blocked, dropping data")
			}
		}
	}
}

// read continuously reads from the port until a complete message is received.
// A complete message starts with '>' and ends with CR LF (\r\n).
// The message is returned without the prefix and suffix
func (c *client) read() ([]byte, error) {
	var res []byte
	messageStarted := false

	for {
		buf := make([]byte, readBufferSize)
		n, err := c.port.Read(buf)
		if err != nil {
			return nil, err
		}

		// Only append the bytes that were actually read
		chunk := buf[:n]

		// If we haven't found the start marker yet, look for it
		if !messageStarted {
			startIdx := bytes.IndexByte(chunk, '>')
			if startIdx >= 0 {
				// Found the start marker, only append from that point
				res = append(res, chunk[startIdx:]...)
				messageStarted = true
			}
		} else {
			// Already found start marker, append the whole chunk
			res = append(res, chunk...)
		}

		// Check if we have a complete message
		if messageStarted && len(res) > 0 && res[0] == '>' && bytes.HasSuffix(res, []byte("\r\n")) {
			// Remove the prefix and suffix
			res = res[1 : len(res)-2]
			// Remove null bytes
			res = bytes.Trim(res, "\x00")
			return res, nil
		}
	}
}
