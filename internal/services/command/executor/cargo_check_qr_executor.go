package executor

import (
	"context"
	"log/slog"
	"sync"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

type cargoCheckQRExecutor struct {
	log *slog.Logger

	subscriber eventbus.Subscriber
}

func newCargoCheckQRExecutor(
	log *slog.Logger,
	subscriber eventbus.Subscriber,
) cargoCheckQRExecutor {
	return cargoCheckQRExecutor{
		log:        log,
		subscriber: subscriber,
	}
}

func (e cargoCheckQRExecutor) Execute(ctx context.Context, inputs command.CargoCheckQRInputs) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		e.trackingCargoQRCode(ctx, inputs.QRCode)
	}()

	wg.Wait()

	return nil
}

func (e cargoCheckQRExecutor) trackingCargoQRCode(ctx context.Context, qrCode string) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		e.log.Debug("stop tracking cargo qr code")
		cancel() // cancel the context to stop the subscriber
	}()

	doneCh := make(chan struct{})
	e.subscriber.Subscribe(ctx, events.CargoQRCodeUpdatedTopic, func(_ context.Context, msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.CargoQRCodeUpdatedEvent)
		if !ok {
			e.log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		if ev.QRCode == qrCode {
			e.log.Debug("cargo qr code matched", slog.Any("qrcode", ev.QRCode))
			close(doneCh)
		}
	})

	select {
	case <-doneCh:
	case <-ctx.Done():
	}
}
