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

type cargoCloseExecutor struct {
	log          *slog.Logger
	subscriber   eventbus.Subscriber
	cargoService cargo.Service
}

func newCargoCloseExecutor(
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	cargoService cargo.Service,
) CommandExecutor[command.CargoCloseInputs, command.CargoCloseOutputs] {
	return cargoCloseExecutor{
		log:          log,
		subscriber:   subscriber,
		cargoService: cargoService,
	}
}

func (e cargoCloseExecutor) Execute(ctx context.Context, _ command.CargoCloseInputs) (command.CargoCloseOutputs, error) {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		e.trackingCargoDoorUntilClosed(ctx)
	}()

	if err := e.cargoService.CloseCargoDoor(ctx); err != nil {
		return command.CargoCloseOutputs{}, fmt.Errorf("failed to close cargo door: %w", err)
	}

	// wait for tracking to finish
	wg.Wait()

	return command.CargoCloseOutputs{}, nil
}

func (e cargoCloseExecutor) OnCancel(_ context.Context) error {
	return nil
}

func (e cargoCloseExecutor) trackingCargoDoorUntilClosed(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		e.log.Debug("stop tracking cargo door")
		cancel()
	}()

	doneCh := make(chan struct{})
	e.log.Debug("start tracking cargo door")
	e.subscriber.Subscribe(ctx, events.CargoDoorUpdatedTopic, func(_ context.Context, msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.CargoDoorUpdatedEvent)
		if !ok {
			e.log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		if !ev.IsOpen {
			e.log.Debug("cargo door closed")
			close(doneCh)
		}
	})

	select {
	case <-doneCh:
	case <-ctx.Done():
	}
}
