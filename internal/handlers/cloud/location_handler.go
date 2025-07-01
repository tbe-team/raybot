package cloud

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	locationv1 "github.com/tbe-team/raybot-api/location/v1"
	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/location"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

type locationHandler struct {
	locationv1.UnimplementedLocationServiceServer
	log             *slog.Logger
	subscriber      eventbus.Subscriber
	locationService location.Service
}

func newLocationHandler(
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	locationService location.Service,
) locationv1.LocationServiceServer {
	return &locationHandler{
		log:             log,
		subscriber:      subscriber,
		locationService: locationService,
	}
}

func (h locationHandler) GetLocation(
	ctx context.Context,
	_ *locationv1.GetLocationRequest,
) (*locationv1.GetLocationResponse, error) {
	location, err := h.locationService.GetLocation(ctx)
	if err != nil {
		return nil, fmt.Errorf("get location: %w", err)
	}

	return &locationv1.GetLocationResponse{
		CurrentLocation: location.CurrentLocation,
		UpdatedAt:       timestamppb.New(location.UpdatedAt),
	}, nil
}

func (h locationHandler) StreamLocation(
	_ *locationv1.StreamLocationRequest,
	stream grpc.ServerStreamingServer[locationv1.StreamLocationResponse],
) error {
	h.log.Info("streaming location")

	ctx := stream.Context()
	h.subscriber.Subscribe(
		ctx,
		events.LocationUpdatedTopic,
		func(msg *eventbus.Message) {
			ev, ok := msg.Payload.(events.LocationUpdatedEvent)
			if !ok {
				h.log.Error("received invalid event", slog.Any("event", msg.Payload))
				return
			}

			if err := stream.Send(&locationv1.StreamLocationResponse{
				CurrentLocation: ev.Location,
				UpdatedAt:       timestamppb.New(ev.UpdatedAt),
			}); err != nil {
				h.log.Error("failed to send location update", slog.Any("error", err))
			}
		},
	)

	<-ctx.Done()
	h.log.Info("stopped streaming location")
	return nil
}
