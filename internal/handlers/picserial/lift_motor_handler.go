package picserial

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/lithammer/shortuuid/v4"

	"github.com/tbe-team/raybot/internal/events"
)

func (s *Service) HandleUpdateLiftMotorStateEvent(_ context.Context, event events.UpdateLiftMotorStateEvent) {
	cmd := picCommand{
		ID:   shortuuid.New(),
		Type: picCommandTypeLiftMotor,
		Data: picCommandLiftMotorData{
			TargetPosition: event.TargetPosition,
			Enable:         event.Enable,
		},
	}

	cmdJSON, err := json.Marshal(cmd)
	if err != nil {
		s.log.Error("failed to marshal lift motor command", slog.Any("error", err))
		return
	}

	s.commandStore.AddCommand(cmd)
	if err := s.client.Write(cmdJSON); err != nil {
		s.log.Error("failed to write lift motor command", slog.Any("error", err))
	}
}
