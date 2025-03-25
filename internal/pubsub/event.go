package pubsub

const (
	TopicCommandCreated       = "command.created"
	TopicRobotLocationUpdated = "robot.location.updated"
)

type CommandCreatedEvent struct {
	CommandID string `json:"command_id"`
}

type RobotLocationUpdatedEvent struct {
	Location string `json:"location"`
}
