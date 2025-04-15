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

type cargoCloseExecutor struct {
	log *slog.Logger

	subscriber   eventbus.Subscriber
	cargoService cargo.Service
}

func newCargoCloseExecutor(
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	cargoService cargo.Service,
) cargoCloseExecutor {
	return cargoCloseExecutor{
		log:          log,
		subscriber:   subscriber,
		cargoService: cargoService,
	}
}

func (e cargoCloseExecutor) Execute(ctx context.Context, _ command.CargoCloseInputs) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		e.trackingCargoDoorClose(ctx)
	}()

	if err := e.cargoService.CloseCargoDoor(ctx); err != nil {
		return NewExecutorError(err, "failed to close cargo")
	}

	// wait for cargo door close to finish
	wg.Wait()

	return nil
}

func (e cargoCloseExecutor) trackingCargoDoorClose(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		e.log.Debug("stop tracking cargo door close")
		cancel() // cancel the context to stop the subscriber
	}()

	doneCh := make(chan struct{})
	e.subscriber.Subscribe(ctx, events.CargoDoorUpdatedTopic, func(_ context.Context, msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.CargoDoorUpdatedEvent)
		if !ok {
			e.log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		if !ev.IsOpen {
			e.log.Debug("cargo door close completed")
			close(doneCh)
		}
	})

	select {
	case <-doneCh:
	case <-ctx.Done():
	}
}
