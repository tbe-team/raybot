package events

import "time"

const (
	LocationUpdatedTopic = "location:updated"
)

type LocationUpdatedEvent struct {
	Location  string
	UpdatedAt time.Time
}
