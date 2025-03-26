package cloud

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"buf.build/gen/go/tbe-team/raybot-api/grpc/go/raybot/v1/raybotv1grpc"
	"github.com/jhump/grpctunnel"
	"github.com/jhump/grpctunnel/tunnelpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/tbe-team/raybot/internal/controller/grpc/handler"
	"github.com/tbe-team/raybot/internal/service"
)

type Config struct {
	Address string `yaml:"address"`
}

func (c *Config) Validate() error {
	// Validate address
	if c.Address == "" {
		return fmt.Errorf("address is required")
	}

	return nil
}

type CleanupFunc func(context.Context) error

type Service struct {
	cfg Config

	service service.Service

	conn                *grpc.ClientConn
	reverseTunnelServer *grpctunnel.ReverseTunnelServer

	log *slog.Logger
}

func NewService(cfg Config, service service.Service, log *slog.Logger) (*Service, error) {
	conn, err := grpc.NewClient(
		cfg.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
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
		service:             service,
		log:                 log,
	}, nil
}

func (s Service) Run(ctx context.Context) (CleanupFunc, error) {
	s.registerHandlers()

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		s.log.Info(fmt.Sprintf("serving reverse tunnel on %s", s.cfg.Address))

		attempts := 0
		for {
			started, err := s.reverseTunnelServer.Serve(ctx)

			if ctx.Err() != nil {
				// Context cancelled, exit loop
				return
			}

			if !started || err != nil {
				s.log.Error("serving reverse tunnel failed, retry after 5 seconds",
					slog.Any("error", err),
					slog.Bool("started", started),
					slog.Int("attempts", attempts),
				)

				time.Sleep(5 * time.Second)
				attempts++
				continue
			}
		}
	}()

	return func(_ context.Context) error {
		s.log.Debug("closing cloud service")
		s.reverseTunnelServer.Stop()
		cancel()

		if err := s.conn.Close(); err != nil {
			return fmt.Errorf("failed to close grpc client: %w", err)
		}
		s.log.Debug("cloud service closed")
		return nil
	}, nil
}

func (s Service) registerHandlers() {
	robotStateHandler := handler.NewRobotStateHandler(s.service.RobotStateService())
	raybotv1grpc.RegisterRobotStateServiceServer(s.reverseTunnelServer, robotStateHandler)

	driveMotorHandler := handler.NewDriveMotorHandler(s.service.PICService())
	raybotv1grpc.RegisterDriveMotorServiceServer(s.reverseTunnelServer, driveMotorHandler)

	liftMotorHandler := handler.NewLiftMotorHandler(s.service.PICService())
	raybotv1grpc.RegisterLiftMotorServiceServer(s.reverseTunnelServer, liftMotorHandler)
}
