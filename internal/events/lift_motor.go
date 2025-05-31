package events

const (
	LiftMotorUpdatedTopic = "lift_motor:updated"
)

type LiftMotorStateUpdatedEvent struct {
	CurrentPosition uint16
	TargetPosition  uint16
	IsRunning       bool
	Enabled         bool
}
