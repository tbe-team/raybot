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

func newCargoCloseExecutor(
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	cargoService cargo.Service,
	commandRepository command.Repository,
) *commandExecutor[command.CargoCloseInputs, command.CargoCloseOutputs] {
	handler := cargoCloseHandler{
		log:          log,
		subscriber:   subscriber,
		cargoService: cargoService,
	}

	return newCommandExecutor(
		handler.Handle,
		Hooks[command.CargoCloseOutputs]{},
		log,
		commandRepository,
	)
}

type cargoCloseHandler struct {
	log          *slog.Logger
	subscriber   eventbus.Subscriber
	cargoService cargo.Service
}

func (h cargoCloseHandler) Handle(ctx context.Context, _ command.CargoCloseInputs) (command.CargoCloseOutputs, error) {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		h.trackingCargoDoorUntilClosed(ctx)
	}()

	if err := h.cargoService.CloseCargoDoor(ctx); err != nil {
		return command.CargoCloseOutputs{}, fmt.Errorf("failed to close cargo door: %w", err)
	}

	// wait for tracking to finish
	wg.Wait()

	return command.CargoCloseOutputs{}, nil
}

func (h cargoCloseHandler) trackingCargoDoorUntilClosed(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		h.log.Debug("stop tracking cargo door")
		cancel()
	}()

	doneCh := make(chan struct{})
	h.log.Debug("start tracking cargo door")
	h.subscriber.Subscribe(ctx, events.CargoDoorUpdatedTopic, func(_ context.Context, msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.CargoDoorUpdatedEvent)
		if !ok {
			h.log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		if !ev.IsOpen {
			h.log.Debug("cargo door closed")
			close(doneCh)
		}
	})

	select {
	case <-doneCh:
	case <-ctx.Done():
	}
}
