package events

const (
	LocationUpdatedTopic = "location:updated"
)

type UpdateLocationEvent struct {
	Location string
}
