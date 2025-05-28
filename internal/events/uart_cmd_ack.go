package events

const (
	PICCmdAckTopic = "pic:cmd:ack"
	ESPCmdAckTopic = "esp:cmd:ack"
)

type PICCmdAckEvent struct {
	ID      string
	Success bool
}

type ESPCmdAckEvent struct {
	ID      string
	Success bool
}
