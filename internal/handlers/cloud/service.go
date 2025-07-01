package cloud

import (
	"context"
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/fullstorydev/grpchan"
	"github.com/jhump/grpctunnel"
	"github.com/jhump/grpctunnel/tunnelpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	batteryv1 "github.com/tbe-team/raybot-api/battery/v1"
	cargov1 "github.com/tbe-team/raybot-api/cargo/v1"
	commandv1 "github.com/tbe-team/raybot-api/command/v1"
	distanceSensorv1 "github.com/tbe-team/raybot-api/distancesensor/v1"
	limitSwitchv1 "github.com/tbe-team/raybot-api/limitswitch/v1"
	locationv1 "github.com/tbe-team/raybot-api/location/v1"
	motorv1 "github.com/tbe-team/raybot-api/motor/v1"
	sysv1 "github.com/tbe-team/raybot-api/sys/v1"
	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/handlers/cloud/interceptor"
	"github.com/tbe-team/raybot/internal/services/battery"
	"github.com/tbe-team/raybot/internal/services/cargo"
	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/distancesensor"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
	"github.com/tbe-team/raybot/internal/services/liftmotor"
	"github.com/tbe-team/raybot/internal/services/limitswitch"
	"github.com/tbe-team/raybot/internal/services/location"
	"github.com/tbe-team/raybot/internal/services/system"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

type Service struct {
	opts Options
	cfg  config.Cloud
	log  *slog.Logger

	publisher  eventbus.Publisher
	subscriber eventbus.Subscriber

	commandService        command.Service
	systemService         system.Service
	batteryService        battery.Service
	cargoService          cargo.Service
	distanceSensorService distancesensor.Service
	limitSwitchService    limitswitch.Service
	locationService       location.Service
	driveMotorService     drivemotor.Service
	liftMotorService      liftmotor.Service

	closing atomic.Bool
}

type Options struct {
	connectTimeout time.Duration
}

var defaultOptions = Options{
	connectTimeout: 2 * time.Second,
}

type OptionFunc func(opts *Options)

func WithConnectTimeout(timeout time.Duration) OptionFunc {
	return func(opts *Options) {
		opts.connectTimeout = timeout
	}
}

type CleanupFunc func() error

func New(
	cfg config.Cloud,
	log *slog.Logger,
	publisher eventbus.Publisher,
	subscriber eventbus.Subscriber,
	commandService command.Service,
	systemService system.Service,
	batteryService battery.Service,
	cargoService cargo.Service,
	distanceSensorService distancesensor.Service,
	limitSwitchService limitswitch.Service,
	locationService location.Service,
	driveMotorService drivemotor.Service,
	liftMotorService liftmotor.Service,
	optFuncs ...OptionFunc,
) *Service {
	opts := defaultOptions
	for _, apply := range optFuncs {
		apply(&opts)
	}

	return &Service{
		opts:                  opts,
		cfg:                   cfg,
		log:                   log.With("service", "cloud"),
		publisher:             publisher,
		subscriber:            subscriber,
		commandService:        commandService,
		systemService:         systemService,
		batteryService:        batteryService,
		cargoService:          cargoService,
		distanceSensorService: distanceSensorService,
		limitSwitchService:    limitSwitchService,
		locationService:       locationService,
		driveMotorService:     driveMotorService,
		liftMotorService:      liftMotorService,
	}
}

func (s *Service) Run(ctx context.Context) (CleanupFunc, error) {
	conn, err := grpc.NewClient(
		s.cfg.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStreamInterceptor(interceptor.ReverseCredentialsInterceptor(s.cfg.Token)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc client: %w", err)
	}

	tunnel := tunnelpb.NewTunnelServiceClient(conn)
	reverseTunnelServer := grpctunnel.NewReverseTunnelServer(tunnel)

	sr := s.wrapInterceptors(reverseTunnelServer)
	s.registerHandlers(sr)

	ctx, cancel := context.WithCancel(ctx)
	go s.runReverseTunnel(ctx, reverseTunnelServer)

	return func() error {
		s.closing.Store(true)

		reverseTunnelServer.Stop()
		cancel()

		if err := conn.Close(); err != nil {
			return fmt.Errorf("failed to close grpc client: %w", err)
		}
		return nil
	}, nil
}

func (s *Service) runReverseTunnel(ctx context.Context, reverseTunnelServer *grpctunnel.ReverseTunnelServer) {
	attempts := 0
	retryDelay := 1 * time.Second
	connected := false

	for {
		if s.isClosing() {
			return
		}

		serveErrChan := make(chan struct{}, 1)
		go func() {
			select {
			// Because Serve function is blocking and we don't have a way to know if it's connected or not
			// so we emit a connected event after [connectTimeout] or it will emit a disconnected event when error occurs
			case <-time.After(s.opts.connectTimeout):
				connected = true
				attempts = 0
				s.publisher.Publish(
					events.CloudConnectedTopic,
					eventbus.NewMessage(events.CloudConnectedEvent{}),
				)

			case <-ctx.Done():
			case <-serveErrChan:
			}
		}()

		started, err := reverseTunnelServer.Serve(ctx)
		if !started || err != nil {
			serveErrChan <- struct{}{}
			s.log.Error("serving reverse tunnel failed, retrying",
				slog.Bool("started", started),
				slog.Int("attempts", attempts),
				slog.Duration("retry_delay", retryDelay),
				slog.Any("error", err),
			)

			// If the last state is connected or this is the first attempt to connect
			// emit a disconnected event
			if connected || attempts == 0 {
				s.publisher.Publish(
					events.CloudDisconnectedTopic,
					eventbus.NewMessage(events.CloudDisconnectedEvent{
						Error: err,
					}),
				)
			}

			time.Sleep(retryDelay)
			attempts++
			retryDelay = min(retryDelay*2, 1*time.Minute)
			continue
		}
	}
}

func (s *Service) isClosing() bool {
	return s.closing.Load()
}

func (s *Service) wrapInterceptors(srv *grpctunnel.ReverseTunnelServer) grpc.ServiceRegistrar {
	sr := grpchan.WithInterceptor(
		srv,
		interceptor.UnaryRecoveryInterceptor(s.log),
		interceptor.StreamRecoveryInterceptor(s.log),
	)
	sr = grpchan.WithInterceptor(
		sr,
		interceptor.UnaryLoggingInterceptor(s.log),
		interceptor.StreamLoggingInterceptor(s.log),
	)
	sr = grpchan.WithInterceptor(
		sr,
		interceptor.UnaryErrorInterceptor(),
		interceptor.StreamErrorInterceptor(),
	)
	return sr
}

func (s *Service) registerHandlers(sr grpc.ServiceRegistrar) {
	commandHandler := newCommandHandler(s.commandService)
	commandv1.RegisterCommandServiceServer(sr, commandHandler)

	systemHandler := newSystemHandler(s.systemService)
	sysv1.RegisterSysServiceServer(sr, systemHandler)

	batteryHandler := newBatteryHandler(s.batteryService)
	batteryv1.RegisterBatteryServiceServer(sr, batteryHandler)

	cargoHandler := newCargoHandler(s.cargoService)
	cargov1.RegisterCargoServiceServer(sr, cargoHandler)

	distanceSensorHandler := newDistanceSensorHandler(s.distanceSensorService)
	distanceSensorv1.RegisterDistanceSensorServiceServer(sr, distanceSensorHandler)

	limitSwitchHandler := newLimitSwitchHandler(s.log, s.subscriber, s.limitSwitchService)
	limitSwitchv1.RegisterLimitSwitchServiceServer(sr, limitSwitchHandler)

	locationHandler := newLocationHandler(s.log, s.subscriber, s.locationService)
	locationv1.RegisterLocationServiceServer(sr, locationHandler)

	motorHandler := newMotorHandler(s.driveMotorService, s.liftMotorService, s.cargoService)
	motorv1.RegisterMotorServiceServer(sr, motorHandler)
}
