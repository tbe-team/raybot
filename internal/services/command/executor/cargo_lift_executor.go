package executor

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/liftmotor"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

func newCargoLiftExecutor(
	cfg config.Cargo,
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	liftMotorService liftmotor.Service,
	commandRepository command.Repository,
) *commandExecutor[command.CargoLiftInputs, command.CargoLiftOutputs] {
	handler := cargoLiftHandler{
		cfg:              cfg,
		log:              log,
		subscriber:       subscriber,
		liftMotorService: liftMotorService,
	}

	return newCommandExecutor(
		handler.Handle,
		Hooks[command.CargoLiftOutputs]{},
		log,
		commandRepository,
	)
}

type cargoLiftHandler struct {
	cfg              config.Cargo
	log              *slog.Logger
	subscriber       eventbus.Subscriber
	liftMotorService liftmotor.Service
}

func (h cargoLiftHandler) Handle(ctx context.Context, _ command.CargoLiftInputs) (command.CargoLiftOutputs, error) {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		h.trackingLiftPositionUntilReached(ctx, h.cfg.LiftPosition)
	}()

	if err := h.liftMotorService.SetCargoPosition(ctx, liftmotor.SetCargoPositionParams{
		Position: h.cfg.LiftPosition,
	}); err != nil {
		return command.CargoLiftOutputs{}, fmt.Errorf("failed to set cargo position: %w", err)
	}

	// wait for tracking to finish
	wg.Wait()

	return command.CargoLiftOutputs{}, nil
}

func (h cargoLiftHandler) trackingLiftPositionUntilReached(ctx context.Context, liftPosition uint16) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		h.log.Debug("stop tracking lift position")
		cancel()
	}()

	doneCh := make(chan struct{})
	h.log.Debug("start tracking lift position", slog.Int64("lift_position", int64(liftPosition)))
	h.subscriber.Subscribe(ctx, events.DistanceSensorUpdatedTopic, func(_ context.Context, msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.UpdateDistanceSensorEvent)
		if !ok {
			h.log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		// 10% tolerance
		acceptableDistance := liftPosition + liftPosition*10/100
		if ev.DownDistance <= acceptableDistance {
			h.log.Debug("lift position reached", slog.Int64("lift_position", int64(liftPosition)))
			close(doneCh)
		}
	})

	select {
	case <-doneCh:
	case <-ctx.Done():
	}
}
