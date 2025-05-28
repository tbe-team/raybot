package espserial

import (
	"fmt"
	"strconv"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

func (s *Service) HandleACK(msg ackMessage) error {
	switch msg.Status {
	case ackStatusError:
		return fmt.Errorf("ack error: %s", msg.ID)

	case ackStatusSuccess:
		s.publisher.Publish(events.ESPCmdAckTopic, eventbus.NewMessage(
			events.ESPCmdAckEvent{
				ID:      msg.ID,
				Success: true,
			},
		))
		return nil
	}

	return nil
}

type ackMessage struct {
	ID     string    `json:"id"`
	Status ackStatus `json:"status"`
}

type ackStatus uint8

func (s *ackStatus) UnmarshalJSON(data []byte) error {
	n, err := strconv.ParseUint(string(data), 10, 8)
	if err != nil {
		return fmt.Errorf("parse uint8: %w", err)
	}

	switch n {
	case 0:
		*s = ackStatusError
	case 1:
		*s = ackStatusSuccess
	default:
		return fmt.Errorf("invalid ack status: %s", string(data))
	}

	return nil
}

const (
	ackStatusError ackStatus = iota
	ackStatusSuccess
)
