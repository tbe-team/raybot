package config

import "fmt"

type Cloud struct {
	Address string `yaml:"address"`
	Token   string `yaml:"token"`
}

func (c Cloud) Validate() error {
	if c.Address == "" {
		return fmt.Errorf("address is required")
	}

	if c.Token == "" {
		return fmt.Errorf("token is required")
	}

	return nil
}
