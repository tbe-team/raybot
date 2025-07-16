package events

import "github.com/tbe-team/raybot/internal/services/battery"

const (
	BatteryUpdatedTopic = "battery:updated"
)

type BatteryUpdatedEvent struct {
	BatteryState battery.BatteryState
}
