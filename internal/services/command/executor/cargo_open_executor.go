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

type cargoOpenExecutor struct {
	log          *slog.Logger
	subscriber   eventbus.Subscriber
	cargoService cargo.Service
}

func newCargoOpenExecutor(
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	cargoService cargo.Service,
) CommandExecutor[command.CargoOpenInputs, command.CargoOpenOutputs] {
	return cargoOpenExecutor{
		log:          log,
		subscriber:   subscriber,
		cargoService: cargoService,
	}
}

func (e cargoOpenExecutor) Execute(ctx context.Context, _ command.CargoOpenInputs) (command.CargoOpenOutputs, error) {
	cargo, err := e.cargoService.GetCargo(ctx)
	if err != nil {
		return command.CargoOpenOutputs{}, fmt.Errorf("failed to get cargo: %w", err)
	}

	if cargo.IsOpen {
		return command.CargoOpenOutputs{}, nil
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		e.trackingCargoDoorUntilOpen(ctx)
	}()

	if err := e.cargoService.OpenCargoDoor(ctx); err != nil {
		return command.CargoOpenOutputs{}, fmt.Errorf("failed to open cargo: %w", err)
	}

	// wait for tracking to finish
	wg.Wait()

	return command.CargoOpenOutputs{}, nil
}

func (e cargoOpenExecutor) OnCancel(_ context.Context) error {
	return nil
}

func (e cargoOpenExecutor) trackingCargoDoorUntilOpen(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		e.log.Debug("stop tracking cargo door open")
		cancel()
	}()

	doneCh := make(chan struct{})
	e.log.Debug("start tracking cargo door open")
	e.subscriber.Subscribe(ctx, events.CargoDoorUpdatedTopic, func(_ context.Context, msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.CargoDoorUpdatedEvent)
		if !ok {
			e.log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		if ev.IsOpen {
			e.log.Debug("cargo door open completed")
			close(doneCh)
		}
	})

	select {
	case <-doneCh:
	case <-ctx.Done():
	}
}
