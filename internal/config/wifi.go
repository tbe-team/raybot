package config

import (
	"fmt"
	"net"
	"regexp"
)

var ssidRegex = regexp.MustCompile(`^[\w\s\-\.]{1,32}$`)      // only allow alphanumeric, space, -, _, . (max 32 characters)
var passwordRegex = regexp.MustCompile(`^[\x21-\x7E]{8,63}$`) // only allow printable characters (min 8, max 63 characters)

type Wifi struct {
	AP  APConfig  `yaml:"ap"`
	STA STAConfig `yaml:"sta"`
}

func (w *Wifi) Validate() error {
	if w.AP.Enable && w.STA.Enable {
		return fmt.Errorf("only one of AP or STA can be enabled")
	}

	if w.AP.Enable {
		if err := w.AP.Validate(); err != nil {
			return fmt.Errorf("invalid AP config: %w", err)
		}
	}

	if w.STA.Enable {
		if err := w.STA.Validate(); err != nil {
			return fmt.Errorf("invalid STA config: %w", err)
		}
	}

	return nil
}

type APConfig struct {
	Enable   bool   `yaml:"enable"`
	SSID     string `yaml:"ssid"`
	Password string `yaml:"password"`
	IP       string `yaml:"ip"`
}

func (c APConfig) Validate() error {
	if !ssidRegex.MatchString(c.SSID) {
		return fmt.Errorf("invalid ssid: %s", c.SSID)
	}

	if !passwordRegex.MatchString(c.Password) {
		return fmt.Errorf("invalid password: %s", c.Password)
	}

	ip := net.ParseIP(c.IP)
	if ip == nil {
		return fmt.Errorf("invalid ip format: %s", c.IP)
	}

	if ip.To4() == nil {
		return fmt.Errorf("only ipv4 is supported: %s", c.IP)
	}

	return nil
}

type STAConfig struct {
	Enable   bool   `yaml:"enable"`
	SSID     string `yaml:"ssid"`
	Password string `yaml:"password"`
}

func (c STAConfig) Validate() error {
	if !ssidRegex.MatchString(c.SSID) {
		return fmt.Errorf("invalid ssid: %s", c.SSID)
	}

	if !passwordRegex.MatchString(c.Password) {
		return fmt.Errorf("invalid password: %s", c.Password)
	}

	return nil
}
