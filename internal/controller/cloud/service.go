package cloud

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

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

	// Validate address format
	if _, err := url.Parse(c.Address); err != nil {
		return fmt.Errorf("invalid address: %w", err)
	}

	return nil
}

type CleanupFunc func(context.Context) error

//nolint:revive
type CloudService struct {
	cfg Config

	service service.Service
	log     *slog.Logger
}

func NewCloudService(cfg Config, service service.Service, log *slog.Logger) (*CloudService, error) {
	return &CloudService{
		cfg:     cfg,
		service: service,
		log:     log,
	}, nil
}

func (s CloudService) Run() (CleanupFunc, error) {
	conn, err := grpc.NewClient(
		s.cfg.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc client: %w", err)
	}

	tunnel := tunnelpb.NewTunnelServiceClient(conn)
	reverseTunnel := grpctunnel.NewReverseTunnelServer(tunnel)
	s.registerHandlers(reverseTunnel)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		s.log.Info(fmt.Sprintf("serving reverse tunnel on %s", s.cfg.Address))
		if started, err := reverseTunnel.Serve(ctx); !started {
			s.log.Error("failed to serve reverse tunnel", slog.Any("error", err))
		}
	}()

	return func(_ context.Context) error {
		s.log.Debug("closing cloud service")
		reverseTunnel.GracefulStop()
		cancel()

		if err := conn.Close(); err != nil {
			return fmt.Errorf("failed to close grpc client: %w", err)
		}
		s.log.Debug("cloud service closed")
		return nil
	}, nil
}

func (s CloudService) registerHandlers(server *grpctunnel.ReverseTunnelServer) {
	robotStateHandler := handler.NewRobotStateHandler(s.service.RobotStateService())
	raybotv1grpc.RegisterRobotStateServiceServer(server, robotStateHandler)

	driveMotorHandler := handler.NewDriveMotorHandler(s.service.PICService())
	raybotv1grpc.RegisterDriveMotorServiceServer(server, driveMotorHandler)

	liftMotorHandler := handler.NewLiftMotorHandler(s.service.PICService())
	raybotv1grpc.RegisterLiftMotorServiceServer(server, liftMotorHandler)
}
