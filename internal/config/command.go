package config

import "fmt"

type Command struct {
	CargoLift CargoLift `yaml:"cargo_lift"`
}

func (c *Command) Validate() error {
	if err := c.CargoLift.Validate(); err != nil {
		return fmt.Errorf("cargo_lift: %w", err)
	}
	return nil
}

type CargoLift struct {
	// StableReadCount is the number of stable bottom distance readings required to consider the lift position reached
	StableReadCount uint8 `yaml:"stable_read_count"`
}

func (c *CargoLift) Validate() error {
	if c.StableReadCount == 0 {
		c.StableReadCount = 1
	}
	return nil
}
