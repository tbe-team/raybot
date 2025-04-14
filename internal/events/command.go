package events

const (
	CommandCreatedTopic = "command:created"
)

type CommandCreatedEvent struct {
	CommandID int64
}
