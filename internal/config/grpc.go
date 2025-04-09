package config

import "fmt"

type GRPC struct {
	Port   uint32 `yaml:"port"`
	Enable bool   `yaml:"enable"`
}

func (g GRPC) Validate() error {
	if g.Port < 1 || g.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}

	return nil
}
