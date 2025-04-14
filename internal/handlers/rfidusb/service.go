package rfidusb

import (
	"context"
	"log/slog"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/location"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

type Service struct {
	log    *slog.Logger
	client *client

	publisher       eventbus.Publisher
	locationService location.Service
}

type CleanupFunc func(context.Context) error

func New(
	log *slog.Logger,
	publisher eventbus.Publisher,
	locationService location.Service,
) *Service {
	return &Service{
		log:             log.With("service", "rfidusb"),
		publisher:       publisher,
		client:          newClient(),
		locationService: locationService,
	}
}

func (s *Service) Run(ctx context.Context) (CleanupFunc, error) {
	if err := s.client.Open(); err != nil {
		// We don't want to fail the service if the serial client fails to open
		s.log.Error("failed to open RFID reader", slog.Any("error", err))
		s.publisher.Publish(
			events.RFIDUSBDisconnectedTopic,
			eventbus.NewMessage(events.RFIDUSBDisconnectedEvent{
				Error: err,
			}),
		)
		return func(_ context.Context) error { return nil }, nil
	}

	s.publisher.Publish(
		events.RFIDUSBConnectedTopic,
		eventbus.NewMessage(events.RFIDUSBConnectedEvent{}),
	)

	ctx, cancel := context.WithCancel(ctx)
	go s.readLoop(ctx)

	cleanup := func(_ context.Context) error {
		// Cancel read loop before closing the USB client
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
			tag, err := s.client.Read()
			if err != nil {
				s.log.Error("failed to read rfid tag", slog.Any("error", err))
				s.publisher.Publish(
					events.RFIDUSBDisconnectedTopic,
					eventbus.NewMessage(events.RFIDUSBDisconnectedEvent{
						Error: err,
					}),
				)
				continue
			}
			s.HandleRFIDTag(ctx, tag)
		}
	}
}
