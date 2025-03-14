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
		log: log.With(slog.String("service", "RFIDService")),
	}

	rfidClient, err := newClient(log)
	if err != nil {
		s.log.Error("failed to create rfid client", slog.Any("error", err))
		return nil, err
	}

	s.rfidClient = rfidClient
	s.rfidTagHandler = newRFIDTagHandler(service.RobotService(), log)
	return s, nil
}

func (s Service) Run(ctx context.Context) (CleanupFunc, error) {
	s.log.Info("RFID service is running")

	go s.readLoop(ctx)

	cleanup := func(_ context.Context) error {
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
