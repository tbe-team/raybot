package picserial

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/lithammer/shortuuid/v4"

	"github.com/tbe-team/raybot/internal/events"
)

func (s *Service) HandleUpdateDriveMotorStateEvent(_ context.Context, event events.UpdateDriveMotorStateEvent) {
	var direction moveDirection
	switch event.Direction {
	case events.MoveDirectionForward:
		direction = moveDirectionForward
	case events.MoveDirectionBackward:
		direction = moveDirectionBackward
	}

	cmd := picCommand{
		ID:   shortuuid.New(),
		Type: picCommandTypeDriveMotor,
		Data: picCommandDriveMotorData{
			Direction: direction,
			Speed:     event.Speed,
			Enable:    event.Enable,
		},
	}

	cmdJSON, err := json.Marshal(cmd)
	if err != nil {
		s.log.Error("failed to marshal drive motor command", slog.Any("error", err))
		return
	}

	s.commandStore.AddCommand(cmd)
	if err := s.client.Write(cmdJSON); err != nil {
		s.log.Error("failed to write drive motor command", slog.Any("error", err))
	}
}
