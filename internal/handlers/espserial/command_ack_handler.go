package espserial

import (
	"context"
	"fmt"
	"strconv"

	"github.com/tbe-team/raybot/internal/services/cargo"
)

func (s *Service) HandleCommandACK(ctx context.Context, msg commandACKMessage) error {
	if msg.Status == ACKStatusFailure {
		return fmt.Errorf("esp serial command failed: %s", msg.ID)
	}

	cmd, ok := s.commandStore.GetCommand(msg.ID)
	if !ok {
		return fmt.Errorf("command not found: %s", msg.ID)
	}

	//nolint:gocritic
	switch cmd.Type {
	case espCommandTypeCargoDoorMotor:
		data, ok := cmd.Data.(espCargoDoorMotorData)
		if !ok {
			return fmt.Errorf("invalid command data: %s", cmd.ID)
		}

		var direction cargo.DoorDirection
		switch data.Direction {
		case doorDirectionOpen:
			direction = cargo.DirectionOpen
		case doorDirectionClose:
			direction = cargo.DirectionClose
		default:
			return fmt.Errorf("invalid direction: %d", data.Direction)
		}

		if err := s.cargoService.UpdateCargoDoorMotorState(ctx, cargo.UpdateCargoDoorMotorStateParams{
			Direction: direction,
			Speed:     data.Speed,
			Enabled:   data.Enable,
		}); err != nil {
			return fmt.Errorf("process cargo door motor command: %w", err)
		}

		s.commandStore.RemoveCommand(msg.ID)
	}

	return nil
}

type ackStatus uint8

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *ackStatus) UnmarshalJSON(data []byte) error {
	n, err := strconv.ParseUint(string(data), 10, 8)
	if err != nil {
		return fmt.Errorf("parse uint8: %w", err)
	}

	switch n {
	case 0:
		*s = ACKStatusFailure
	case 1:
		*s = ACKStatusSuccess
	default:
		return fmt.Errorf("invalid ack status: %s", string(data))
	}
	return nil
}

const (
	ACKStatusFailure ackStatus = iota
	ACKStatusSuccess
)

type commandACKMessage struct {
	ID     string    `json:"id"`
	Status ackStatus `json:"status"`
}
