package config

import "fmt"

type Command struct {
	CargoLift  CargoLift  `yaml:"cargo_lift"`
	CargoLower CargoLower `yaml:"cargo_lower"`
}

func (c *Command) Validate() error {
	if err := c.CargoLift.Validate(); err != nil {
		return fmt.Errorf("cargo_lift: %w", err)
	}

	if err := c.CargoLower.Validate(); err != nil {
		return fmt.Errorf("cargo_lower: %w", err)
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

type CargoLower struct {
	// StableReadCount is the number of stable bottom distance readings required to consider the lift position reached
	StableReadCount uint8 `yaml:"stable_read_count"`

	// BottomObstacleTracking is the configuration for the bottom obstacle tracking
	BottomObstacleTracking ObstacleTracking `yaml:"bottom_obstacle_tracking"`
}

func (c *CargoLower) Validate() error {
	if c.StableReadCount == 0 {
		c.StableReadCount = 1
	}
	return nil
}

type ObstacleTracking struct {
	// Obstacle is considered present when distance <= EnterDistance (cm)
	EnterDistance uint16 `yaml:"enter_distance"`
	// Obstacle is considered cleared when distance >= ExitDistance (cm)
	ExitDistance uint16 `yaml:"exit_distance"`
}

func (c ObstacleTracking) Validate() error {
	if c.EnterDistance >= c.ExitDistance {
		return fmt.Errorf("enter distance must be less than exit distance")
	}
	return nil
}
