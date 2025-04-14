package events

const (
	CargoDoorUpdatedTopic = "cargo:door:updated"
)

type CargoDoorUpdatedEvent struct {
	IsOpen bool
}
