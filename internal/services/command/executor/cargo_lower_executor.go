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

func newCargoLowerExecutor(
	cfg config.Cargo,
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	liftMotorService liftmotor.Service,
	commandRepository command.Repository,
) *commandExecutor[command.CargoLowerInputs, command.CargoLowerOutputs] {
	handler := cargoLowerHandler{
		cfg:              cfg,
		log:              log,
		subscriber:       subscriber,
		liftMotorService: liftMotorService,
	}

	return newCommandExecutor(
		handler.Handle,
		Hooks[command.CargoLowerOutputs]{},
		log,
		commandRepository,
	)
}

type cargoLowerHandler struct {
	cfg              config.Cargo
	log              *slog.Logger
	subscriber       eventbus.Subscriber
	liftMotorService liftmotor.Service
}

func (h cargoLowerHandler) Handle(ctx context.Context, _ command.CargoLowerInputs) (command.CargoLowerOutputs, error) {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		h.trackingLowerPositionUntilReached(ctx, h.cfg.LowerPosition)
	}()

	if err := h.liftMotorService.SetCargoPosition(ctx, liftmotor.SetCargoPositionParams{
		Position: h.cfg.LowerPosition,
	}); err != nil {
		return command.CargoLowerOutputs{}, fmt.Errorf("failed to set cargo position: %w", err)
	}

	// wait for tracking to finish
	wg.Wait()

	return command.CargoLowerOutputs{}, nil
}

func (h cargoLowerHandler) trackingLowerPositionUntilReached(ctx context.Context, lowerPosition uint16) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		h.log.Debug("stop tracking lower position")
		cancel()
	}()

	doneCh := make(chan struct{})
	h.log.Debug("start tracking lower position", slog.Int64("lower_position", int64(lowerPosition)))
	h.subscriber.Subscribe(ctx, events.DistanceSensorUpdatedTopic, func(_ context.Context, msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.UpdateDistanceSensorEvent)
		if !ok {
			h.log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		// 10% tolerance
		acceptableDistance := lowerPosition - lowerPosition*10/100
		if ev.DownDistance >= acceptableDistance {
			h.log.Debug("lower position reached", slog.Int64("lower_position", int64(lowerPosition)))
			close(doneCh)
		}
	})

	select {
	case <-doneCh:
	case <-ctx.Done():
	}
}
