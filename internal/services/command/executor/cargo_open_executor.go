package executor

import (
	"context"
	"log/slog"
	"sync"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/cargo"
	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

type cargoOpenExecutor struct {
	log *slog.Logger

	subscriber   eventbus.Subscriber
	cargoService cargo.Service
}

func newCargoOpenExecutor(
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	cargoService cargo.Service,
) cargoOpenExecutor {
	return cargoOpenExecutor{
		log:          log,
		subscriber:   subscriber,
		cargoService: cargoService,
	}
}

func (e cargoOpenExecutor) Execute(ctx context.Context, _ command.CargoOpenInputs) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		e.trackingCargoDoorOpen(ctx)
	}()

	if err := e.cargoService.OpenCargoDoor(ctx); err != nil {
		return NewExecutorError(err, "failed to open cargo")
	}

	// wait for cargo door open to finish
	wg.Wait()

	return nil
}

func (e cargoOpenExecutor) trackingCargoDoorOpen(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		e.log.Debug("stop tracking cargo door open")
		cancel() // cancel the context to stop the subscriber
	}()

	doneCh := make(chan struct{})
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
