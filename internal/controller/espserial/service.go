package espserial

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/tbe-team/raybot/internal/service"
)

type Config struct {
	Serial SerialConfig `yaml:"serial"`
}

// Validate validates the ESP configuration.
func (cfg *Config) Validate() error {
	return cfg.Serial.Validate()
}

type CleanupFunc func(context.Context) error

type Handlers struct {
	SyncStateHandler  *SyncStateHandler
	CommandACKHandler *CommandACKHandler
}

type Service struct {
	cfg Config

	serialClient Client

	handlers Handlers
	log      *slog.Logger
}

func New(cfg Config, client Client, service service.Service, log *slog.Logger) (*Service, error) {
	handlers := Handlers{
		SyncStateHandler:  NewSyncStateHandler(service.CargoControlService()),
		CommandACKHandler: NewCommandACKHandler(service.CargoControlService()),
	}

	return &Service{
		cfg:          cfg,
		serialClient: client,
		handlers:     handlers,
		log:          log.With("service", "espserial"),
	}, nil
}

func (s Service) Run(ctx context.Context) (CleanupFunc, error) {
	s.log.Info("ESP serial service is running")

	if s.serialClient != nil {
		go s.readLoop(ctx)
	}

	cleanup := func(_ context.Context) error {
		s.log.Debug("ESP serial service shut down complete")
		return nil
	}

	return cleanup, nil
}

func (s Service) readLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := s.serialClient.Read()
			if err != nil {
				s.log.Error("failed to read from serial client", slog.Any("error", err))
				return
			}
			s.routeMessage(ctx, msg)
		}
	}
}

func (s Service) routeMessage(ctx context.Context, msg []byte) {
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

		if err := s.handlers.SyncStateHandler.Handle(ctx, syncStateMsg); err != nil {
			s.log.Error("failed to handle sync state message", slog.Any("error", err), slog.Any("message", msg))
		}

	case messageTypeACK:
		var commandACKMsg commandACKMessage
		if err := json.Unmarshal(msg, &commandACKMsg); err != nil {
			s.log.Error("failed to unmarshal command ack message", slog.Any("error", err), slog.Any("message", msg))
			return
		}

		if err := s.handlers.CommandACKHandler.Handle(ctx, commandACKMsg); err != nil {
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
