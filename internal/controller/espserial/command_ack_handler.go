package espserial

import (
	"context"
	"fmt"
	"strconv"

	"github.com/tbe-team/raybot/internal/service"
)

type CommandACKHandler struct {
	service service.CargoControlService
}

func NewCommandACKHandler(service service.CargoControlService) *CommandACKHandler {
	return &CommandACKHandler{
		service: service,
	}
}

func (h *CommandACKHandler) Handle(ctx context.Context, msg commandACKMessage) error {
	params := service.ProcessESPSerialCommandACKParams{
		ID:      msg.ID,
		Success: msg.Status == ACKStatusSuccess,
	}
	if err := h.service.ProcessESPSerialCommandACK(ctx, params); err != nil {
		return fmt.Errorf("failed to process esp serial command ack: %w", err)
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
