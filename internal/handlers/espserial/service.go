package espserial

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/hardware/espserial"
	"github.com/tbe-team/raybot/internal/services/cargo"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

type Service struct {
	cfg config.ESP
	log *slog.Logger

	publisher eventbus.Publisher

	client espserial.Client

	cargoService cargo.Service
}

type CleanupFunc func(context.Context) error

func New(
	cfg config.ESP,
	log *slog.Logger,
	publisher eventbus.Publisher,
	client espserial.Client,
	cargoService cargo.Service,
) *Service {
	s := &Service{
		cfg:          cfg,
		publisher:    publisher,
		client:       client,
		log:          log.With("service", "espserial"),
		cargoService: cargoService,
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
		return nil
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
				s.publisher.Publish(
					events.ESPSerialDisconnectedTopic,
					eventbus.NewMessage(events.ESPSerialDisconnectedEvent{
						Error: err,
					}),
				)
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
