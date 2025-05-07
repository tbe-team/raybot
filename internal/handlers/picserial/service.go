package picserial

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/hardware/picserial"
	"github.com/tbe-team/raybot/internal/services/appstate"
	"github.com/tbe-team/raybot/internal/services/battery"
	"github.com/tbe-team/raybot/internal/services/distancesensor"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
	"github.com/tbe-team/raybot/internal/services/liftmotor"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

type Service struct {
	cfg    config.PIC
	log    *slog.Logger
	client picserial.Client

	publisher eventbus.Publisher

	batteryService        battery.Service
	distanceSensorService distancesensor.Service
	liftMotorService      liftmotor.Service
	driveMotorService     drivemotor.Service
	appStateService       appstate.Service
}

type CleanupFunc func(context.Context) error

func New(
	cfg config.PIC,
	log *slog.Logger,
	client picserial.Client,
	publisher eventbus.Publisher,
	batteryService battery.Service,
	distanceSensorService distancesensor.Service,
	liftMotorService liftmotor.Service,
	driveMotorService drivemotor.Service,
	appStateService appstate.Service,
) *Service {
	s := &Service{
		cfg:                   cfg,
		client:                client,
		publisher:             publisher,
		log:                   log.With("service", "picserial"),
		batteryService:        batteryService,
		distanceSensorService: distanceSensorService,
		liftMotorService:      liftMotorService,
		driveMotorService:     driveMotorService,
		appStateService:       appStateService,
	}

	return s
}

func (s *Service) Run(ctx context.Context) (CleanupFunc, error) {
	if !s.client.Connected() {
		return func(_ context.Context) error { return nil }, nil
	}

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
			msg, err := s.client.Read(ctx)
			if err != nil {
				s.log.Error("failed to read from serial client", slog.Any("error", err))
				if errors.Is(err, picserial.ErrPICSerialNotConnected) {
					s.publisher.Publish(
						events.PICSerialDisconnectedTopic,
						eventbus.NewMessage(events.PICSerialDisconnectedEvent{
							Error: err,
						}),
					)
				}
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
