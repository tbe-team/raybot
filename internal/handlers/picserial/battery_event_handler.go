package picserial

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/lithammer/shortuuid/v4"

	"github.com/tbe-team/raybot/internal/events"
)

func (s *Service) HandleUpdateBatteryChargeSettingEvent(_ context.Context, event events.UpdateBatteryChargeSettingEvent) {
	cmd := picCommand{
		ID:   shortuuid.New(),
		Type: picCommandTypeBatteryCharge,
		Data: picCommandBatteryChargeData{
			CurrentLimit: event.CurrentLimit,
			Enable:       event.Enable,
		},
	}

	cmdJSON, err := json.Marshal(cmd)
	if err != nil {
		s.log.Error("failed to marshal battery charge command", slog.Any("error", err))
		return
	}

	s.commandStore.AddCommand(cmd)
	if err := s.client.Write(cmdJSON); err != nil {
		s.log.Error("failed to write battery charge command", slog.Any("error", err))
	}
}

func (s *Service) HandleUpdateBatteryDischargeSettingEvent(_ context.Context, event events.UpdateBatteryDischargeSettingEvent) {
	cmd := picCommand{
		ID:   shortuuid.New(),
		Type: picCommandTypeBatteryDischarge,
		Data: picCommandBatteryDischargeData{
			CurrentLimit: event.CurrentLimit,
			Enable:       event.Enable,
		},
	}

	cmdJSON, err := json.Marshal(cmd)
	if err != nil {
		s.log.Error("failed to marshal battery discharge command", slog.Any("error", err))
		return
	}

	s.commandStore.AddCommand(cmd)
	if err := s.client.Write(cmdJSON); err != nil {
		s.log.Error("failed to write battery discharge command", slog.Any("error", err))
	}
}
