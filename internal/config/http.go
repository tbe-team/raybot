package config

import "fmt"

type HTTP struct {
	Swagger bool   `yaml:"swagger"`
	Port    uint32 `yaml:"port"`
}

func (h HTTP) Validate() error {
	if h.Port < 1 || h.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}

	return nil
}
