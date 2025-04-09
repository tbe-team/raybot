package picserial

import (
	"context"
	"fmt"
	"strconv"

	"github.com/tbe-team/raybot/internal/services/battery"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
	"github.com/tbe-team/raybot/internal/services/liftmotor"
)

func (s *Service) HandleCommandACK(ctx context.Context, msg commandACKMessage) error {
	if msg.Status == ackStatusFailure {
		return fmt.Errorf("pic serial command failed: %s", msg.ID)
	}

	cmd, ok := s.commandStore.GetCommand(msg.ID)
	if !ok {
		return fmt.Errorf("command not found: %s", msg.ID)
	}

	//nolint:gocritic
	switch cmd.Type {
	case picCommandTypeBatteryCharge:
		data, ok := cmd.Data.(picCommandBatteryChargeData)
		if !ok {
			return fmt.Errorf("invalid command data: %s", cmd.ID)
		}

		if err := s.batteryService.UpdateChargeSetting(ctx, battery.UpdateChargeSettingParams{
			CurrentLimit: data.CurrentLimit,
			Enabled:      data.Enable,
		}); err != nil {
			return fmt.Errorf("failed to update battery charge: %w", err)
		}

		s.commandStore.RemoveCommand(msg.ID)

	case picCommandTypeBatteryDischarge:
		data, ok := cmd.Data.(picCommandBatteryDischargeData)
		if !ok {
			return fmt.Errorf("invalid command data: %s", cmd.ID)
		}

		if err := s.batteryService.UpdateDischargeSetting(ctx, battery.UpdateDischargeSettingParams{
			CurrentLimit: data.CurrentLimit,
			Enabled:      data.Enable,
		}); err != nil {
			return fmt.Errorf("failed to update battery discharge: %w", err)
		}

		s.commandStore.RemoveCommand(msg.ID)

	case picCommandTypeLiftMotor:
		data, ok := cmd.Data.(picCommandLiftMotorData)
		if !ok {
			return fmt.Errorf("invalid command data: %s", cmd.ID)
		}

		if err := s.liftMotorService.UpdateLiftMotorState(ctx, liftmotor.UpdateLiftMotorStateParams{
			TargetPosition: data.TargetPosition,
			Enabled:        data.Enable,
		}); err != nil {
			return fmt.Errorf("failed to update lift motor: %w", err)
		}

		s.commandStore.RemoveCommand(msg.ID)

	case picCommandTypeDriveMotor:
		data, ok := cmd.Data.(picCommandDriveMotorData)
		if !ok {
			return fmt.Errorf("invalid command data: %s", cmd.ID)
		}

		var direction drivemotor.Direction
		switch data.Direction {
		case moveDirectionForward:
			direction = drivemotor.DirectionForward
		case moveDirectionBackward:
			direction = drivemotor.DirectionBackward
		}

		if err := s.driveMotorService.UpdateDriveMotorState(ctx, drivemotor.UpdateDriveMotorStateParams{
			Direction: direction,
			Speed:     data.Speed,
			Enabled:   data.Enable,
		}); err != nil {
			return fmt.Errorf("failed to update drive motor: %w", err)
		}

		s.commandStore.RemoveCommand(msg.ID)

	default:
		return fmt.Errorf("invalid command type: %d", cmd.Type)
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
		*s = ackStatusFailure
	case 1:
		*s = ackStatusSuccess
	default:
		return fmt.Errorf("invalid ack status: %s", string(data))
	}
	return nil
}

const (
	ackStatusFailure ackStatus = iota
	ackStatusSuccess
)

type commandACKMessage struct {
	ID     string    `json:"id"`
	Status ackStatus `json:"status"`
}
