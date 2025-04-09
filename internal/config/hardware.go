package config

import (
	"fmt"
	"strings"
	"time"
)

type Hardware struct {
	ESP ESP `yaml:"esp"`
	PIC PIC `yaml:"pic"`
}

func (h *Hardware) Validate() error {
	if err := h.ESP.Validate(); err != nil {
		return fmt.Errorf("validate esp: %w", err)
	}

	if err := h.PIC.Validate(); err != nil {
		return fmt.Errorf("validate pic: %w", err)
	}

	if h.ESP.Serial.Port == h.PIC.Serial.Port {
		return fmt.Errorf("esp and pic serial ports cannot be the same")
	}

	return nil
}

type ESP struct {
	Serial Serial `yaml:"serial"`
}

func (e ESP) Validate() error {
	if err := e.Serial.Validate(); err != nil {
		return fmt.Errorf("validate esp serial: %w", err)
	}

	return nil
}

type PIC struct {
	Serial Serial `yaml:"serial"`
}

func (p PIC) Validate() error {
	if err := p.Serial.Validate(); err != nil {
		return fmt.Errorf("validate pic serial: %w", err)
	}

	return nil
}

type Serial struct {
	Port        string        `yaml:"port"`
	BaudRate    int           `yaml:"baud_rate"`
	DataBits    uint8         `yaml:"data_bits"`
	StopBits    float32       `yaml:"stop_bits"`
	Parity      string        `yaml:"parity"`
	ReadTimeout time.Duration `yaml:"read_timeout"`
}

func (s *Serial) Validate() error {
	if s.BaudRate < 1200 || s.BaudRate > 115200 {
		return fmt.Errorf("invalid baud rate: %d", s.BaudRate)
	}

	if s.DataBits != 5 && s.DataBits != 6 && s.DataBits != 7 && s.DataBits != 8 {
		return fmt.Errorf("invalid data bits: %d", s.DataBits)
	}

	if s.StopBits != 1 && s.StopBits != 1.5 && s.StopBits != 2 {
		return fmt.Errorf("invalid stop bits: %f", s.StopBits)
	}

	p := strings.ToUpper(s.Parity)
	if p != "NONE" && p != "ODD" && p != "EVEN" {
		return fmt.Errorf("invalid parity: %s", s.Parity)
	}
	s.Parity = p

	return nil
}
