package picserial

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/appconnection"
	"github.com/tbe-team/raybot/internal/services/battery"
	"github.com/tbe-team/raybot/internal/services/distancesensor"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
	"github.com/tbe-team/raybot/internal/services/liftmotor"
)

type Service struct {
	cfg    config.PIC
	client *client
	log    *slog.Logger

	batteryService        battery.Service
	distanceSensorService distancesensor.Service
	liftMotorService      liftmotor.Service
	driveMotorService     drivemotor.Service
	appConnectionService  appconnection.Service
	commandStore          *commandStore
}

type CleanupFunc func(context.Context) error

func New(
	cfg config.PIC,
	log *slog.Logger,
	batteryService battery.Service,
	distanceSensorService distancesensor.Service,
	liftMotorService liftmotor.Service,
	driveMotorService drivemotor.Service,
	appConnectionService appconnection.Service,
) *Service {
	s := &Service{
		cfg:                   cfg,
		client:                newClient(cfg.Serial),
		log:                   log,
		batteryService:        batteryService,
		distanceSensorService: distanceSensorService,
		liftMotorService:      liftMotorService,
		driveMotorService:     driveMotorService,
		appConnectionService:  appConnectionService,
		commandStore:          newCommandStore(),
	}

	events.UpdateBatteryChargeSettingSignal.AddListener(s.HandleUpdateBatteryChargeSettingEvent)
	events.UpdateBatteryDischargeSettingSignal.AddListener(s.HandleUpdateBatteryDischargeSettingEvent)
	events.UpdateLiftMotorStateSignal.AddListener(s.HandleUpdateLiftMotorStateEvent)
	events.UpdateDriveMotorStateSignal.AddListener(s.HandleUpdateDriveMotorStateEvent)

	return s
}

func (s *Service) Run(ctx context.Context) (CleanupFunc, error) {
	if err := s.client.Open(); err != nil {
		// We don't want to fail the service if the serial client fails to open
		s.log.Error("failed to open PIC serial client",
			slog.Any("serial_cfg", s.client.cfg),
			slog.Any("error", err),
		)
		events.PICSerialDisconnectedSignal.Emit(ctx, events.PICSerialDisconnectedEvent{
			Error: err,
		})

		return func(_ context.Context) error { return nil }, nil
	}

	events.PICSerialConnectedSignal.Emit(ctx, events.PICSerialConnectedEvent{})

	ctx, cancel := context.WithCancel(ctx)
	go s.readLoop(ctx)

	cleanup := func(_ context.Context) error {
		// Cancel read loop before closing the serial client
		cancel()
		return s.client.Close()
	}

	return cleanup, nil
}

func (s *Service) readLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := s.client.Read()
			if err != nil {
				s.log.Error("failed to read from serial client", slog.Any("error", err))
				events.PICSerialDisconnectedSignal.Emit(ctx, events.PICSerialDisconnectedEvent{
					Error: err,
				})
				return
			}
			s.routeMessage(ctx, msg)
		}
	}
}

func (s *Service) routeMessage(ctx context.Context, msg []byte) {
	var temp struct {
		Type messageType `json:"type"`
	}
	if err := json.Unmarshal(msg, &temp); err != nil {
		s.log.Error("failed to unmarshal message", slog.Any("error", err), slog.Any("message", msg))
		return
	}

	//nolint:gocritic
	switch temp.Type {
	case messageTypeSyncState:
		var syncStateMsg syncStateMessage
		if err := json.Unmarshal(msg, &syncStateMsg); err != nil {
			s.log.Error("failed to unmarshal sync state message", slog.Any("error", err), slog.Any("message", msg))
			return
		}
		if err := s.HandleSyncState(ctx, syncStateMsg); err != nil {
			s.log.Error("failed to handle sync state message", slog.Any("error", err), slog.Any("message", msg))
		}

	case messageTypeACK:
		var commandACKMsg commandACKMessage
		if err := json.Unmarshal(msg, &commandACKMsg); err != nil {
			s.log.Error("failed to unmarshal command ack message", slog.Any("error", err), slog.Any("message", msg))
			return
		}
		if err := s.HandleCommandACK(ctx, commandACKMsg); err != nil {
			s.log.Error("failed to handle command ack message", slog.Any("error", err), slog.Any("message", msg))
		}

	default:
		s.log.Error("unknown message type", slog.Any("type", temp.Type))
	}
}

// messageType is the type of message received from the PIC
type messageType uint8

// UnmarshalJSON implements the json.Unmarshaler interface.
func (m *messageType) UnmarshalJSON(data []byte) error {
	n, err := strconv.ParseUint(string(data), 10, 8)
	if err != nil {
		return fmt.Errorf("parse uint8: %w", err)
	}

	switch n {
	case 0:
		*m = messageTypeSyncState
	case 1:
		*m = messageTypeACK
	default:
		return fmt.Errorf("invalid message type: %s", string(data))
	}
	return nil
}

const (
	messageTypeSyncState messageType = iota
	messageTypeACK
)
