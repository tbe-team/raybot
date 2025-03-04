package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/tbe-team/raybot/internal/service"
)

// MessageType is the type of message received from the PIC
type MessageType int

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (m *MessageType) UnmarshalText(text []byte) error {
	switch int(text[0]) {
	case 0:
		*m = MessageTypeSyncState
	default:
		return fmt.Errorf("invalid message type: %s", string(text))
	}
	return nil
}

const (
	MessageTypeSyncState MessageType = iota
)

type SyncStateMessage struct {
	StateType int8            `json:"state_type"`
	Data      json.RawMessage `json:"data"`
}

type SyncStateHandler struct {
	robotService service.RobotService
	log          *slog.Logger
}

func NewSyncStateHandler(robotService service.RobotService) *SyncStateHandler {
	return &SyncStateHandler{
		robotService: robotService,
		log: slog.With(
			slog.String("module", "pic"),
			slog.String("handler", "SyncStateHandler"),
		),
	}
}

func (h *SyncStateHandler) Handle(ctx context.Context, msg SyncStateMessage) {
	params := service.UpdateRobotStateParams{}
	switch msg.StateType {
	case 0:
		params.SetBattery = true

	case 1:
	case 2:
	case 3:
	case 4:
	case 5:
	default:
		h.log.Error("unknown state type", "type", msg.StateType)
	}
}
