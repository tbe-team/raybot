package config

import "fmt"

type Cargo struct {
	LiftPosition   uint16 `yaml:"lift_position"`
	LowerPosition  uint16 `yaml:"lower_position"`
	LiftMotorSpeed uint8  `yaml:"lift_motor_speed"`

	// BottomDistanceHysteresis is the hysteresis for the bottom distance when executing CARGO_LOWER command.
	BottomDistanceHysteresis CargoBottomDistanceHysteresis `yaml:"bottom_distance_hysteresis"`
}

func (c Cargo) Validate() error {
	if c.LiftPosition >= c.LowerPosition {
		return fmt.Errorf("lift position must be less than lower position: %d", c.LowerPosition)
	}

	if c.LiftMotorSpeed > 100 {
		return fmt.Errorf("lift motor speed must be between 0 and 100: %d", c.LiftMotorSpeed)
	}

	if err := c.BottomDistanceHysteresis.Validate(); err != nil {
		return fmt.Errorf("bottom distance hysteresis: %w", err)
	}

	return nil
}

type CargoBottomDistanceHysteresis struct {
	// LowerThreshold is the lower threshold for the bottom distance. Unit is cm.
	LowerThreshold uint16 `yaml:"lower_threshold"`

	// UpperThreshold is the upper threshold for the bottom distance. Unit is cm.
	UpperThreshold uint16 `yaml:"upper_threshold"`
}

func (c CargoBottomDistanceHysteresis) Validate() error {
	if c.LowerThreshold >= c.UpperThreshold {
		return fmt.Errorf("lower threshold must be less than upper threshold: %d", c.UpperThreshold)
	}

	return nil
}
