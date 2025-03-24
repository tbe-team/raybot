package espserial

import (
	"bytes"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"go.bug.st/serial"
)

// SerialConfig is the configuration for the ESP serial port.
type SerialConfig struct {
	// Port is the serial port path (e.g., "/dev/ttyUSB0" or "COM3")
	Port string `yaml:"port"`

	// BaudRate is the communication speed
	BaudRate int `yaml:"baud_rate"`

	// DataBits is the number of data bits (usually 8)
	DataBits int `yaml:"data_bits"`

	// StopBits is the number of stop bits (usually 1)
	StopBits float64 `yaml:"stop_bits"`

	// Parity mode (none, odd, even)
	Parity string `yaml:"parity"`

	// ReadTimeout is the timeout for read operations
	ReadTimeout time.Duration `yaml:"read_timeout"`
}

// Validate verifies the configuration for the ESP serial port.
func (cfg SerialConfig) Validate() error {
	if cfg.BaudRate < 1200 || cfg.BaudRate > 115200 {
		return fmt.Errorf("invalid baud rate: %d", cfg.BaudRate)
	}

	if cfg.DataBits != 5 && cfg.DataBits != 6 && cfg.DataBits != 7 && cfg.DataBits != 8 {
		return fmt.Errorf("invalid data bits: %d", cfg.DataBits)
	}

	if cfg.StopBits != 1 && cfg.StopBits != 1.5 && cfg.StopBits != 2 {
		return fmt.Errorf("invalid stop bits: %f", cfg.StopBits)
	}

	if cfg.Parity != "none" && cfg.Parity != "odd" && cfg.Parity != "even" {
		return fmt.Errorf("invalid parity: %s", cfg.Parity)
	}

	return nil
}

const (
	readBufferSize = 64
)

// Client is the interface for the serial client.
type Client interface {
	// Write sends data to the serial port.
	// It is safe to call this method from multiple goroutines.
	Write(data []byte) error

	// Read reads data from the serial port.
	Read() ([]byte, error)

	// Stop stops the client.
	Stop() error
}

type client struct {
	cfg SerialConfig

	writeMu sync.Mutex
	port    serial.Port
	log     *slog.Logger
}

// NewClient creates a new serial client.
func NewClient(cfg SerialConfig, log *slog.Logger) (Client, error) {
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

	port, err := serial.Open(cfg.Port, mode)
	if err != nil {
		return nil, fmt.Errorf("failed to open serial port: %w", err)
	}

	if err := port.SetReadTimeout(cfg.ReadTimeout); err != nil {
		return nil, fmt.Errorf("failed to set read timeout: %w", err)
	}

	return &client{cfg: cfg, port: port, log: log}, nil
}

// Write sends data to the serial port.
// It prefixes the data with '>' and suffixes it with CR LF (\r\n)
func (c *client) Write(data []byte) error {
	data = append([]byte(">"), data...)
	data = append(data, '\r', '\n')

	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	_, err := c.port.Write(data)
	return err
}

// Read reads data from the serial port.
func (c *client) Read() ([]byte, error) {
	return c.read()
}

// Stop stops the client.
func (c *client) Stop() error {
	return c.port.Close()
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
