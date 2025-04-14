package eventbus

type Payload any

type Message struct {
	// Metadata is additional information about the message.
	Metadata Metadata

	// Payload is the actual message payload.
	Payload Payload
}

func NewMessage(payload Payload) *Message {
	return &Message{
		Metadata: make(map[string]string),
		Payload:  payload,
	}
}
