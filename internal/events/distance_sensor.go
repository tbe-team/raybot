package events

const (
	DistanceSensorUpdatedTopic = "distance_sensor:updated"
)

type UpdateDistanceSensorEvent struct {
	FrontDistance uint16
	BackDistance  uint16
	DownDistance  uint16
}
