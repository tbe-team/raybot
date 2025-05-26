package config

import (
	"fmt"
	"strings"
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
	scheduleDuration time.Duration
	Schedule         string        `yaml:"schedule"`
	Threshold        time.Duration `yaml:"threshold"`
}

func (c DeleteOldCommand) ScheduleDuration() time.Duration {
	return c.scheduleDuration
}

func (c *DeleteOldCommand) Validate() error {
	d, err := parseDuration(c.Schedule)
	if err != nil {
		return fmt.Errorf("schedule: %w", err)
	}

	c.scheduleDuration = d

	if c.Threshold.Hours() < 1 {
		return fmt.Errorf("threshold must be greater than 1 hour")
	}

	return nil
}

// parseDuration parses a duration string with the format "@every <duration>".
// For example, "@every 1h" will be parsed as 1 hour.
func parseDuration(expr string) (time.Duration, error) {
	const prefix = "@every "
	if !strings.HasPrefix(expr, prefix) {
		return 0, fmt.Errorf("not an @every expression")
	}
	return time.ParseDuration(strings.TrimPrefix(expr, prefix))
}
