package cloud

import (
	"context"
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/jhump/grpctunnel"
	"github.com/jhump/grpctunnel/tunnelpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/handlers/cloud/interceptor"
)

type Service struct {
	cfg config.Cloud

	conn                *grpc.ClientConn
	reverseTunnelServer *grpctunnel.ReverseTunnelServer
	closing             atomic.Bool

	log *slog.Logger
}

type CleanupFunc func(context.Context) error

func New(cfg config.Cloud, log *slog.Logger) (*Service, error) {
	conn, err := grpc.NewClient(
		cfg.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStreamInterceptor(interceptor.ReverseCredentialsInterceptor(cfg.Token)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc client: %w", err)
	}

	tunnel := tunnelpb.NewTunnelServiceClient(conn)
	reverseTunnel := grpctunnel.NewReverseTunnelServer(tunnel)

	return &Service{
		cfg:                 cfg,
		conn:                conn,
		reverseTunnelServer: reverseTunnel,
		log:                 log,
	}, nil
}

func (s *Service) Run(ctx context.Context) (CleanupFunc, error) {
	ctx, cancel := context.WithCancel(ctx)

	go s.runReverseTunnel(ctx)

	return func(_ context.Context) error {
		s.closing.Store(true)

		s.reverseTunnelServer.Stop()
		cancel()

		if err := s.conn.Close(); err != nil {
			return fmt.Errorf("failed to close grpc client: %w", err)
		}
		return nil
	}, nil
}

func (s *Service) runReverseTunnel(ctx context.Context) {
	attempts := 0
	retryDelay := 1 * time.Second

	for {
		if s.isClosing() {
			return
		}

		connectingErrChan := make(chan struct{}, 1)
		go func() {
			select {
			case <-ctx.Done():
				return
			// Because Serve function is blocking, we don't know if it's connected or not
			// so we emit a connected event after 3 seconds or it will emit a disconnected event
			// if it's not connected
			case <-time.After(2 * time.Second):
				events.CloudConnectedSignal.Emit(ctx, events.CloudConnectedEvent{})
			case <-connectingErrChan:
			}
		}()

		started, err := s.reverseTunnelServer.Serve(ctx)
		if ctx.Err() != nil {
			// Context cancelled, exit loop
			return
		}
		if !started || err != nil {
			connectingErrChan <- struct{}{}
			s.log.Error("serving reverse tunnel failed, retrying",
				slog.Bool("started", started),
				slog.Int("attempts", attempts),
				slog.Duration("retryDelay", retryDelay),
				slog.Any("error", err),
			)
			events.CloudDisconnectedSignal.Emit(ctx, events.CloudDisconnectedEvent{
				Error: err,
			})

			time.Sleep(retryDelay)
			attempts++
			retryDelay *= 2
			continue
		}
	}
}

func (s *Service) isClosing() bool {
	return s.closing.Load()
}
