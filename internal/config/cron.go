package config

import (
	"fmt"
	"time"
)

type Cron struct {
	DeleteOldCommand DeleteOldCommand `yaml:"delete_old_command"`
}

func (c Cron) Validate() error {
	if err := c.DeleteOldCommand.Validate(); err != nil {
		return fmt.Errorf("delete_old_command: %w", err)
	}
	return nil
}

type DeleteOldCommand struct {
	Schedule  string        `yaml:"schedule"`
	Threshold time.Duration `yaml:"threshold"`
}

func (c DeleteOldCommand) Validate() error {
	if c.Schedule == "" {
		return fmt.Errorf("schedule is required")
	}
	if c.Threshold.Hours() < 1 {
		return fmt.Errorf("threshold must be greater than 1 hour")
	}

	return nil
}
