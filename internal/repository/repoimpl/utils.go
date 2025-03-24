package repoimpl

import (
	"fmt"
	"time"
)

// Helper function to parse time consistently
func parseTime(timeStr string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse time: %w", err)
	}
	return t, nil
}
