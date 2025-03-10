package handler

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/tbe-team/raybot/internal/service"
)

type ACKStatus uint8

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *ACKStatus) UnmarshalJSON(data []byte) error {
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
	ACKStatusFailure ACKStatus = iota
	ACKStatusSuccess
)

type CommandACKMessage struct {
	ID     string    `json:"id"`
	Status ACKStatus `json:"status"`
}

type CommandACKHandler struct {
	picService service.PICService
	log        *slog.Logger
}

func NewCommandACKHandler(picService service.PICService, log *slog.Logger) *CommandACKHandler {
	return &CommandACKHandler{
		picService: picService,
		log: log.With(
			slog.String("handler", "CommandACKHandler"),
		),
	}
}

func (h CommandACKHandler) Handle(ctx context.Context, msg CommandACKMessage) {
	params := service.ProcessSerialCommandACKParams{
		ID:      msg.ID,
		Success: msg.Status == ACKStatusSuccess,
	}
	if err := h.picService.ProcessSerialCommandACK(ctx, params); err != nil {
		h.log.Error("failed to process serial command ack", slog.Any("error", err))
	}
}
