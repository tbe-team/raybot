package serial

import (
	"flag"
	"fmt"
	"time"
)

// Config is the configuration for the PIC serial port.
type Config struct {
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

// RegisterFlagsWithPrefix registers flags with a prefix for the PIC serial port configuration.
func (cfg *Config) RegisterFlagsWithPrefix(prefix string, f *flag.FlagSet) {
	f.StringVar(&cfg.Port, prefix+"serial-port", "/dev/ttyUSB0", "PIC board serial port")
	f.IntVar(&cfg.BaudRate, prefix+"serial-baud-rate", 9600, "PIC board serial baud rate")
	f.IntVar(&cfg.DataBits, prefix+"serial-data-bits", 8, "PIC board serial data bits (8, 7, 6, 5)")
	f.Float64Var(&cfg.StopBits, prefix+"serial-stop-bits", 1, "PIC board serial stop bits (1, 1.5, 2)")
	f.StringVar(&cfg.Parity, prefix+"serial-parity", "none", "PIC board serial parity (none, odd, even)")
	f.DurationVar(&cfg.ReadTimeout, prefix+"serial-read-timeout", 1*time.Second, "PIC board serial read timeout")
}

// Validate verifies the configuration for the PIC serial port.
func (cfg *Config) Validate() error {
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
