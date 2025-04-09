package espserial

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/lithammer/shortuuid/v4"

	"github.com/tbe-team/raybot/internal/events"
)

const (
	openCargoDoorSpeed  = 50
	closeCargoDoorSpeed = 50
)

func (s *Service) HandleOpenCargoDoorEvent(_ context.Context, _ events.OpenCargoDoorEvent) {
	cmd := espCommand{
		ID:   shortuuid.New(),
		Type: espCommandTypeCargoDoorMotor,
		Data: espCargoDoorMotorData{
			Direction: doorDirectionOpen,
			Speed:     openCargoDoorSpeed,
			Enable:    true,
		},
	}

	cmdJSON, err := json.Marshal(cmd)
	if err != nil {
		s.log.Error("failed to marshal cargo door motor command", slog.Any("error", err))
		return
	}

	s.commandStore.AddCommand(cmd)
	if err := s.client.Write(cmdJSON); err != nil {
		s.log.Error("failed to write cargo door motor command", slog.Any("error", err))
		return
	}
}

func (s *Service) HandleCloseCargoDoorEvent(_ context.Context, _ events.CloseCargoDoorEvent) {
	cmd := espCommand{
		ID:   shortuuid.New(),
		Type: espCommandTypeCargoDoorMotor,
		Data: espCargoDoorMotorData{
			Direction: doorDirectionClose,
			Speed:     closeCargoDoorSpeed,
			Enable:    true,
		},
	}

	cmdJSON, err := json.Marshal(cmd)
	if err != nil {
		s.log.Error("failed to marshal cargo door motor command", slog.Any("error", err))
		return
	}

	s.commandStore.AddCommand(cmd)
	if err := s.client.Write(cmdJSON); err != nil {
		s.log.Error("failed to write cargo door motor command", slog.Any("error", err))
		return
	}
}
