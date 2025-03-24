package espserial

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
)

type Config struct {
	Serial SerialConfig `yaml:"serial"`
}

// Validate validates the ESP configuration.
func (cfg *Config) Validate() error {
	return cfg.Serial.Validate()
}

type CleanupFunc func(context.Context) error

type Service struct {
	cfg Config

	serialClient Client

	log *slog.Logger
}

func New(cfg Config, client Client, log *slog.Logger) (*Service, error) {
	return &Service{
		cfg:          cfg,
		serialClient: client,
		log:          log,
	}, nil
}

func (s Service) Run(ctx context.Context) (CleanupFunc, error) {
	s.log.Info("ESP serial service is running")

	go s.readLoop(ctx)

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
				continue
			}
			s.routeMessage(ctx, msg)
		}
	}
}

func (s Service) routeMessage(ctx context.Context, msg []byte) {
	var temp struct {
		Type messageType `json:"type"`
	}
	if err := json.Unmarshal(msg, &temp); err != nil {
		s.log.Error("failed to unmarshal message", slog.Any("error", err), slog.Any("message", msg))
		return
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
