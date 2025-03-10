package picserial

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/tbe-team/raybot/internal/controller/picserial/handler"
	"github.com/tbe-team/raybot/internal/controller/picserial/serial"
	"github.com/tbe-team/raybot/internal/service"
)

type Config struct {
	Serial serial.Config `yaml:"serial"`
}

// Validate validates the PIC configuration.
func (cfg *Config) Validate() error {
	return cfg.Serial.Validate()
}

type Handlers struct {
	SyncStateHandler  *handler.SyncStateHandler
	CommandACKHandler *handler.CommandACKHandler
}

//nolint:revive
type PICSerialService struct {
	cfg Config

	serialClient serial.Client

	handlers Handlers
	log      *slog.Logger
}

type CleanupFunc func(context.Context) error

func NewPICSerialService(cfg Config, client serial.Client, service service.Service, log *slog.Logger) (*PICSerialService, error) {
	handlers := Handlers{
		SyncStateHandler:  handler.NewSyncStateHandler(service.RobotService(), log),
		CommandACKHandler: handler.NewCommandACKHandler(service.PICService(), log),
	}

	return &PICSerialService{
		cfg:          cfg,
		serialClient: client,
		handlers:     handlers,
		log:          log.With(slog.String("service", "PICSerialService")),
	}, nil
}

// Run runs the PIC serial service.
func (s *PICSerialService) Run(ctx context.Context) (CleanupFunc, error) {
	s.log.Info("PIC serial service is running")

	go s.readLoop(ctx)

	cleanup := func(_ context.Context) error {
		s.log.Debug("PIC serial service shut down complete")
		return nil
	}

	return cleanup, nil
}

func (s *PICSerialService) readLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-s.serialClient.Read():
			if !ok {
				s.log.Error("serial client read channel closed")
				return
			}
			s.routeMessage(ctx, msg)
		}
	}
}

func (s *PICSerialService) routeMessage(ctx context.Context, msg []byte) {
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
		var syncStateMsg handler.SyncStateMessage
		if err := json.Unmarshal(msg, &syncStateMsg); err != nil {
			s.log.Error("failed to unmarshal sync state message", slog.Any("error", err), slog.Any("message", msg))
			return
		}
		s.handlers.SyncStateHandler.Handle(ctx, syncStateMsg)

	case messageTypeSyncStateACK:
		var commandACKMsg handler.CommandACKMessage
		if err := json.Unmarshal(msg, &commandACKMsg); err != nil {
			s.log.Error("failed to unmarshal command ack message", slog.Any("error", err), slog.Any("message", msg))
			return
		}
		s.handlers.CommandACKHandler.Handle(ctx, commandACKMsg)

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
		*m = messageTypeSyncStateACK
	default:
		return fmt.Errorf("invalid message type: %s", string(data))
	}
	return nil
}

const (
	messageTypeSyncState messageType = iota
	messageTypeSyncStateACK
)
