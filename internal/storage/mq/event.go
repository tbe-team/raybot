package mq

const (
	TopicCommandCreated = "command.created"
)

type CommandCreatedEvent struct {
	CommandID string `json:"command_id"`
}
