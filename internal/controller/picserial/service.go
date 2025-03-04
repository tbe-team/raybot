package picserial

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"

	"github.com/tbe-team/raybot/internal/controller/picserial/handler"
	"github.com/tbe-team/raybot/internal/controller/picserial/serial"
	"github.com/tbe-team/raybot/internal/service"
)

type Config struct {
	Serial serial.Config `yaml:"serial"`
}

// RegisterFlags registers flags for the PIC configuration.
func (cfg *Config) RegisterFlags(f *flag.FlagSet) {
	cfg.Serial.RegisterFlagsWithPrefix("pic.", f)
}

// Validate validates the PIC configuration.
func (cfg *Config) Validate() error {
	return cfg.Serial.Validate()
}

type Handlers struct {
	SyncStateHandler *handler.SyncStateHandler
}

//nolint:revive
type PICSerialService struct {
	cfg Config

	serialClient serial.Client
	robotService service.RobotService

	handlers Handlers
	log      *slog.Logger
}

type CleanupFunc func(context.Context) error

func NewPICSerialService(cfg Config, service service.Service) (*PICSerialService, error) {
	serialClient, err := serial.NewClient(cfg.Serial)
	if err != nil {
		return nil, fmt.Errorf("failed to create serial client: %w", err)
	}

	handlers := Handlers{
		SyncStateHandler: handler.NewSyncStateHandler(service.RobotService()),
	}

	return &PICSerialService{
		cfg:          cfg,
		serialClient: serialClient,
		robotService: service.RobotService(),
		handlers:     handlers,
		log: slog.With(
			slog.String("module", "pic"),
			slog.String("service", "PICSerialService"),
		),
	}, nil
}

// Run runs the PIC serial service.
func (s *PICSerialService) Run(ctx context.Context) (CleanupFunc, error) {
	go s.readLoop(ctx)

	cleanup := func(_ context.Context) error {
		return s.serialClient.Stop()
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
		s.log.Error("failed to unmarshal message", "error", err)
		return
	}

	//nolint:gocritic
	switch temp.Type {
	case messageTypeSyncState:
		var syncStateMsg handler.SyncStateMessage
		if err := json.Unmarshal(msg, &syncStateMsg); err != nil {
			s.log.Error("failed to unmarshal sync state message", "error", err)
			return
		}
		s.handlers.SyncStateHandler.Handle(ctx, syncStateMsg)
	}
}

// messageType is the type of message received from the PIC
type messageType int

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (m *messageType) UnmarshalText(text []byte) error {
	length := len(text)
	if length == 0 || length > 1 {
		return fmt.Errorf("invalid message type: %s", string(text))
	}

	switch int(text[0]) {
	case 0:
		*m = messageTypeSyncState
	default:
		return fmt.Errorf("invalid message type: %s", string(text))
	}
	return nil
}

const (
	messageTypeSyncState messageType = iota
)
