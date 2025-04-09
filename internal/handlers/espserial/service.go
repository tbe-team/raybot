package espserial

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/cargo"
)

type Service struct {
	cfg    config.ESP
	client *client
	log    *slog.Logger

	cargoService cargo.Service
	commandStore *commandStore
}

type CleanupFunc func(context.Context) error

func New(
	cfg config.ESP,
	log *slog.Logger,
	cargoService cargo.Service,
) *Service {
	s := &Service{
		cfg:          cfg,
		client:       newClient(cfg.Serial),
		log:          log.With("service", "espserial"),
		cargoService: cargoService,
		commandStore: newCommandStore(),
	}

	events.CloseCargoDoorSignal.AddListener(s.HandleCloseCargoDoorEvent)
	events.OpenCargoDoorSignal.AddListener(s.HandleOpenCargoDoorEvent)

	return s
}

func (s *Service) Run(ctx context.Context) (CleanupFunc, error) {
	if err := s.client.Open(); err != nil {
		// We don't want to fail the service if the serial client fails to open
		s.log.Error("failed to open ESP serial client",
			slog.Any("serial_cfg", s.client.cfg),
			slog.Any("error", err),
		)
		events.ESPSerialDisconnectedSignal.Emit(ctx, events.ESPSerialDisconnectedEvent{
			Error: err,
		})
		return func(_ context.Context) error { return nil }, nil
	}

	events.ESPSerialConnectedSignal.Emit(ctx, events.ESPSerialConnectedEvent{})

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
				events.ESPSerialDisconnectedSignal.Emit(ctx, events.ESPSerialDisconnectedEvent{
					Error: err,
				})
				return
			}
			s.routeMessage(ctx, msg)
		}
	}
}

func (s *Service) routeMessage(ctx context.Context, msg []byte) {
	s.log.Debug("routing message", slog.Any("message", msg))
	var temp struct {
		Type messageType `json:"type"`
	}
	if err := json.Unmarshal(msg, &temp); err != nil {
		s.log.Error("failed to unmarshal message type", slog.Any("error", err), slog.Any("message", msg))
		return
	}

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
	}
}

// messageType is the type of message received from the ESP
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
