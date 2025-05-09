package application

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/hardware/espserial"
	"github.com/tbe-team/raybot/internal/hardware/picserial"
	"github.com/tbe-team/raybot/internal/logging"
	"github.com/tbe-team/raybot/internal/services/apperrorcode"
	"github.com/tbe-team/raybot/internal/services/apperrorcode/apperrorcodeimpl"
	"github.com/tbe-team/raybot/internal/services/appstate"
	"github.com/tbe-team/raybot/internal/services/appstate/appstateimpl"
	"github.com/tbe-team/raybot/internal/services/battery"
	"github.com/tbe-team/raybot/internal/services/battery/batteryimpl"
	"github.com/tbe-team/raybot/internal/services/cargo"
	"github.com/tbe-team/raybot/internal/services/cargo/cargoimpl"
	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/command/commandimpl"
	"github.com/tbe-team/raybot/internal/services/command/executor"
	"github.com/tbe-team/raybot/internal/services/command/processinglockimpl"
	configsvc "github.com/tbe-team/raybot/internal/services/config"
	"github.com/tbe-team/raybot/internal/services/config/configimpl"
	"github.com/tbe-team/raybot/internal/services/dashboarddata"
	"github.com/tbe-team/raybot/internal/services/dashboarddata/dashboarddataimpl"
	"github.com/tbe-team/raybot/internal/services/distancesensor"
	"github.com/tbe-team/raybot/internal/services/distancesensor/distancesensorimpl"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
	"github.com/tbe-team/raybot/internal/services/drivemotor/drivemotorimpl"
	"github.com/tbe-team/raybot/internal/services/liftmotor"
	"github.com/tbe-team/raybot/internal/services/liftmotor/liftmotorimpl"
	"github.com/tbe-team/raybot/internal/services/location"
	"github.com/tbe-team/raybot/internal/services/location/locationimpl"
	"github.com/tbe-team/raybot/internal/services/peripheral"
	"github.com/tbe-team/raybot/internal/services/peripheral/peripheralimpl"
	"github.com/tbe-team/raybot/internal/services/system"
	"github.com/tbe-team/raybot/internal/services/system/systemimpl"
	"github.com/tbe-team/raybot/internal/services/wifi/wifiimpl"
	"github.com/tbe-team/raybot/internal/storage/db"
	"github.com/tbe-team/raybot/internal/storage/db/sqlc"
	"github.com/tbe-team/raybot/internal/storage/file"
	"github.com/tbe-team/raybot/pkg/eventbus"
	"github.com/tbe-team/raybot/pkg/ptr"
	"github.com/tbe-team/raybot/pkg/validator"
)

type Application struct {
	Cfg     *config.Config
	Log     *slog.Logger
	Context context.Context

	EventBus eventbus.EventBus

	ESPSerialClient espserial.Client
	PICSerialClient picserial.Client

	BatteryService        battery.Service
	DistanceSensorService distancesensor.Service
	DriveMotorService     drivemotor.Service
	LiftMotorService      liftmotor.Service
	CargoService          cargo.Service
	LocationService       location.Service
	ConfigService         configsvc.Service
	SystemService         system.Service
	DashboardDataService  dashboarddata.Service
	AppStateService       appstate.Service
	PeripheralService     peripheral.Service
	CommandService        command.Service
	ApperrorcodeService   apperrorcode.Service
}

type CleanupFunc func() error

func New(configFilePath, dbPath string) (*Application, CleanupFunc, error) {
	ctx := context.Background()

	// Set UTC timezone
	time.Local = time.UTC

	cfg, err := config.NewConfig(configFilePath, dbPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create config: %w", err)
	}

	// Initialize logger
	log, cleanupLogger, err := logging.NewSlogLogger(cfg.Log)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create logger: %w", err)
	}

	// Initialize file client
	fileClient := file.NewLocalFileClient()

	// Initialize db
	db, err := db.NewSQLiteDB(dbPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create db: %w", err)
	}

	// Migrate db
	if err := db.AutoMigrate(); err != nil {
		return nil, nil, fmt.Errorf("failed to migrate db: %w", err)
	}

	// Initialize event bus
	eventBus := eventbus.NewInProcEventBus(log)

	// Initialize repositories
	queries := sqlc.New()
	validator := validator.New()
	batteryStateRepository := batteryimpl.NewBatteryStateRepository()
	batterySettingRepository := batteryimpl.NewBatterySettingRepository(db, queries)
	driveMotorStateRepository := drivemotorimpl.NewDriveMotorStateRepository()
	liftMotorStateRepository := liftmotorimpl.NewLiftMotorStateRepository()
	cargoRepository := cargoimpl.NewCargoRepository(db, queries)
	locationRepository := locationimpl.NewLocationRepository(db, queries)
	distanceSensorStateRepository := distancesensorimpl.NewDistanceSensorStateRepository()
	appStateRepository := appstateimpl.NewAppStateRepository()
	commandRepository := commandimpl.NewCommandRepository(db, queries)

	// Initialize hardware components
	espSerialClient := espserial.NewClient(cfg.Hardware.ESP.Serial)
	if err := espSerialClient.Open(); err != nil {
		log.Error("failed to open ESP serial client",
			slog.Any("serial_cfg", cfg.Hardware.ESP.Serial),
			slog.Any("error", err),
		)

		if err := appStateRepository.UpdateESPSerialConnection(ctx, appstate.UpdateESPSerialConnectionParams{
			Connected:    false,
			SetConnected: true,
			Error:        ptr.New(err.Error()),
			SetError:     true,
		}); err != nil {
			log.Error("failed to update ESP serial connection", slog.Any("error", err))
		}
	} else {
		if err := appStateRepository.UpdateESPSerialConnection(ctx, appstate.UpdateESPSerialConnectionParams{
			Connected:          true,
			SetConnected:       true,
			LastConnectedAt:    ptr.New(time.Now()),
			SetLastConnectedAt: true,
		}); err != nil {
			log.Error("failed to update ESP serial connection", slog.Any("error", err))
		}
	}

	picSerialClient := picserial.NewClient(cfg.Hardware.PIC.Serial)
	if err := picSerialClient.Open(); err != nil {
		log.Error("failed to open PIC serial client",
			slog.Any("serial_cfg", cfg.Hardware.PIC.Serial),
			slog.Any("error", err),
		)

		if err := appStateRepository.UpdatePICSerialConnection(ctx, appstate.UpdatePICSerialConnectionParams{
			Connected:    false,
			SetConnected: true,
			Error:        ptr.New(err.Error()),
			SetError:     true,
		}); err != nil {
			log.Error("failed to update PIC serial connection", slog.Any("error", err))
		}
	} else {
		if err := appStateRepository.UpdatePICSerialConnection(ctx, appstate.UpdatePICSerialConnectionParams{
			Connected:          true,
			SetConnected:       true,
			LastConnectedAt:    ptr.New(time.Now()),
			SetLastConnectedAt: true,
		}); err != nil {
			log.Error("failed to update PIC serial connection", slog.Any("error", err))
		}
	}

	// Initialize services
	batteryService := batteryimpl.NewService(validator, batteryStateRepository, batterySettingRepository)
	distanceSensorService := distancesensorimpl.NewService(validator, eventBus, distanceSensorStateRepository)
	driveMotorService := drivemotorimpl.NewService(validator, eventBus, driveMotorStateRepository, picSerialClient)
	liftMotorService := liftmotorimpl.NewService(validator, liftMotorStateRepository, picSerialClient)
	cargoService := cargoimpl.NewService(validator, eventBus, cargoRepository, espSerialClient)
	locationService := locationimpl.NewService(validator, eventBus, locationRepository)
	configService := configimpl.NewService(cfg, fileClient)
	systemService := systemimpl.NewService(log)
	dashboardDataService := dashboarddataimpl.NewService(
		batteryStateRepository,
		batterySettingRepository,
		distanceSensorStateRepository,
		liftMotorStateRepository,
		driveMotorStateRepository,
		locationRepository,
		cargoRepository,
		appStateRepository,
	)
	appStateService := appstateimpl.NewService(appStateRepository)
	peripheralService := peripheralimpl.NewService()
	commandService := commandimpl.NewService(
		cfg.Cron.DeleteOldCommand,
		log,
		validator,
		eventBus,
		commandRepository,
		appStateRepository,
		processinglockimpl.New(),
		executor.NewRouter(
			cfg.Cargo,
			log,
			eventBus,
			driveMotorService,
			liftMotorService,
			cargoService,
			commandRepository,
		),
	)
	wifiService := wifiimpl.NewService(cfg.Wifi, log)
	if err := wifiService.Run(ctx); err != nil {
		return nil, nil, fmt.Errorf("failed to run wifi service: %w", err)
	}

	apperrorcodeService := apperrorcodeimpl.NewService()

	cleanup := func() error {
		var err error
		if espSerialClient.Connected() {
			if espErr := espSerialClient.Close(); espErr != nil {
				err = fmt.Errorf("failed to close ESP serial client: %w", espErr)
			}
		}

		if picSerialClient.Connected() {
			if picErr := picSerialClient.Close(); picErr != nil {
				err = fmt.Errorf("failed to close PIC serial client: %w", picErr)
			}
		}

		if dbErr := db.Close(); dbErr != nil {
			err = fmt.Errorf("failed to close db: %w", dbErr)
		}

		if err := cleanupLogger(); err != nil {
			return fmt.Errorf("failed to cleanup logger: %w", err)
		}

		return err
	}

	return &Application{
		Cfg:                   cfg,
		Log:                   log,
		Context:               ctx,
		EventBus:              eventBus,
		ESPSerialClient:       espSerialClient,
		PICSerialClient:       picSerialClient,
		BatteryService:        batteryService,
		DistanceSensorService: distanceSensorService,
		DriveMotorService:     driveMotorService,
		LiftMotorService:      liftMotorService,
		CargoService:          cargoService,
		LocationService:       locationService,
		ConfigService:         configService,
		SystemService:         systemService,
		DashboardDataService:  dashboardDataService,
		AppStateService:       appStateService,
		PeripheralService:     peripheralService,
		CommandService:        commandService,
		ApperrorcodeService:   apperrorcodeService,
	}, cleanup, nil
}
