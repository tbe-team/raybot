package config

import "fmt"

type Cargo struct {
	LiftPosition  uint16 `yaml:"lift_position"`
	LowerPosition uint16 `yaml:"lower_position"`
}

func (c Cargo) Validate() error {
	if c.LiftPosition >= c.LowerPosition {
		return fmt.Errorf("lift position must be less than lower position: %d", c.LowerPosition)
	}

	return nil
}
