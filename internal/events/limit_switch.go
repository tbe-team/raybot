package events

import "time"

const (
	LimitSwitch1PressedTopic = "limit_switch:1:pressed"
)

type LimitSwitch1PressedEvent struct {
	PressedAt time.Time
}
