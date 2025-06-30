package cloud

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	limitSwitchv1 "github.com/tbe-team/raybot-api/limitswitch/v1"
	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/limitswitch"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

type limitSwitchHandler struct {
	limitSwitchv1.UnimplementedLimitSwitchServiceServer
	log                *slog.Logger
	subscriber         eventbus.Subscriber
	limitSwitchService limitswitch.Service
}

func newLimitSwitchHandler(
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	limitSwitchService limitswitch.Service,
) *limitSwitchHandler {
	return &limitSwitchHandler{
		log:                log,
		subscriber:         subscriber,
		limitSwitchService: limitSwitchService,
	}
}

func (h limitSwitchHandler) GetLimitSwitch1(ctx context.Context, _ *limitSwitchv1.GetLimitSwitch1Request) (*limitSwitchv1.GetLimitSwitch1Response, error) {
	state, err := h.limitSwitchService.GetLimitSwitchState(ctx)
	if err != nil {
		return nil, fmt.Errorf("get limit switch state: %w", err)
	}

	return &limitSwitchv1.GetLimitSwitch1Response{
		IsPressed: state.LimitSwitch1.Pressed,
		UpdatedAt: timestamppb.New(state.LimitSwitch1.UpdatedAt),
	}, nil
}

func (h limitSwitchHandler) StreamLimitSwitch1PressEvent(
	_ *limitSwitchv1.StreamLimitSwitch1PressEventRequest,
	stream grpc.ServerStreamingServer[limitSwitchv1.StreamLimitSwitch1PressEventResponse],
) error {
	h.log.Info("streaming limit switch 1 press event")

	ctx := stream.Context()
	h.subscriber.Subscribe(
		ctx,
		events.LimitSwitch1PressedTopic,
		func(msg *eventbus.Message) {
			ev, ok := msg.Payload.(events.LimitSwitch1PressedEvent)
			if !ok {
				h.log.Error("received invalid event", slog.Any("event", msg.Payload))
				return
			}

			if err := stream.Send(&limitSwitchv1.StreamLimitSwitch1PressEventResponse{
				PressedAt: timestamppb.New(ev.PressedAt),
			}); err != nil {
				h.log.Error("failed to stream limit switch 1 press event", slog.Any("error", err))
				return
			}
		},
	)

	<-ctx.Done()
	h.log.Info("stopped streaming limit switch 1 press event")
	return nil
}
