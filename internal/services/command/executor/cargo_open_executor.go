package executor

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/cargo"
	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

func newCargoOpenExecutor(
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	cargoService cargo.Service,
	commandRepository command.Repository,
) *commandExecutor[command.CargoOpenInputs, command.CargoOpenOutputs] {
	handler := cargoOpenHandler{
		log:          log,
		subscriber:   subscriber,
		cargoService: cargoService,
	}

	return newCommandExecutor(
		handler.Handle,
		Hooks[command.CargoOpenOutputs]{},
		log,
		commandRepository,
	)
}

type cargoOpenHandler struct {
	log          *slog.Logger
	subscriber   eventbus.Subscriber
	cargoService cargo.Service
}

func (h cargoOpenHandler) Handle(ctx context.Context, _ command.CargoOpenInputs) (command.CargoOpenOutputs, error) {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		h.trackingCargoDoorUntilOpen(ctx)
	}()

	if err := h.cargoService.OpenCargoDoor(ctx); err != nil {
		return command.CargoOpenOutputs{}, fmt.Errorf("failed to open cargo: %w", err)
	}

	// wait for tracking to finish
	wg.Wait()

	return command.CargoOpenOutputs{}, nil
}

func (h cargoOpenHandler) trackingCargoDoorUntilOpen(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		h.log.Debug("stop tracking cargo door open")
		cancel()
	}()

	doneCh := make(chan struct{})
	h.log.Debug("start tracking cargo door open")
	h.subscriber.Subscribe(ctx, events.CargoDoorUpdatedTopic, func(_ context.Context, msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.CargoDoorUpdatedEvent)
		if !ok {
			h.log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		if ev.IsOpen {
			h.log.Debug("cargo door open completed")
			close(doneCh)
		}
	})

	select {
	case <-doneCh:
	case <-ctx.Done():
	}
}
