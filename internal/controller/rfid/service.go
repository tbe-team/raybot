package rfid

import (
	"context"
	"log/slog"

	"github.com/tbe-team/raybot/internal/service"
)

type Service struct {
	rfidClient     *client
	rfidTagHandler *rfidTagHandler
	log            *slog.Logger
}

type CleanupFunc func(context.Context) error

func NewRFIDService(service service.Service, log *slog.Logger) (*Service, error) {
	s := &Service{
		log: log,
	}

	rfidClient, err := newClient(log)
	if err != nil {
		s.log.Error("failed to create rfid client", slog.Any("error", err))
		// Avoid crashing the application if RFID reader is not connected
		// return nil, err
	}

	s.rfidClient = rfidClient
	s.rfidTagHandler = newRFIDTagHandler(service.LocationService(), log)
	return s, nil
}

func (s Service) Run(ctx context.Context) (CleanupFunc, error) {
	// Only start the read loop if we have a client
	if s.rfidClient != nil {
		s.log.Info("RFID service is running")
		go s.readLoop(ctx)
	} else {
		s.log.Error("RFID service running without a reader connected")
	}

	cleanup := func(_ context.Context) error {
		s.log.Debug("RFID service is shutting down")

		if s.rfidClient != nil {
			if err := s.rfidClient.Stop(); err != nil {
				s.log.Error("failed to stop rfid client", slog.Any("error", err))
			}
		}

		s.log.Debug("RFID service shut down complete")
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
			tag, err := s.rfidClient.Read()
			if err != nil {
				s.log.Error("failed to read rfid tag", slog.Any("error", err))
				continue
			}
			s.rfidTagHandler.Handle(ctx, tag)
		}
	}
}
