package application

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/services/appconnection"
	"github.com/tbe-team/raybot/internal/services/appconnection/appconnectionimpl"
	"github.com/tbe-team/raybot/internal/services/battery"
	"github.com/tbe-team/raybot/internal/services/battery/batteryimpl"
	"github.com/tbe-team/raybot/internal/services/cargo"
	"github.com/tbe-team/raybot/internal/services/cargo/cargoimpl"
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
	"github.com/tbe-team/raybot/internal/storage/db"
	"github.com/tbe-team/raybot/internal/storage/db/sqlc"
	"github.com/tbe-team/raybot/internal/storage/file"
	"github.com/tbe-team/raybot/pkg/log"
	"github.com/tbe-team/raybot/pkg/validator"
)

type Application struct {
	Cfg     *config.Config
	Log     *slog.Logger
	Context context.Context

	BatteryService        battery.Service
	DistanceSensorService distancesensor.Service
	DriveMotorService     drivemotor.Service
	LiftMotorService      liftmotor.Service
	CargoService          cargo.Service
	LocationService       location.Service
	ConfigService         configsvc.Service
	SystemService         system.Service
	DashboardDataService  dashboarddata.Service
	AppConnectionService  appconnection.Service
	PeripheralService     peripheral.Service
}

type CleanupFunc func() error

func New(configFilePath, dbPath string) (*Application, CleanupFunc, error) {
	ctx := context.Background()

	cfg, err := config.NewConfig(configFilePath, dbPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create config: %w", err)
	}

	// Initialize logger
	log := log.NewSlogLogger(log.Config{
		Level:     cfg.Log.Level,
		Format:    cfg.Log.Format,
		AddSource: cfg.Log.AddSource,
	})

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

	queries := sqlc.New()
	validator := validator.New()

	// Initialize repositories
	batteryStateRepository := batteryimpl.NewBatteryStateRepository()
	batterySettingRepository := batteryimpl.NewBatterySettingRepository(db, queries)
	driveMotorStateRepository := drivemotorimpl.NewDriveMotorStateRepository()
	liftMotorStateRepository := liftmotorimpl.NewLiftMotorStateRepository()
	cargoRepository := cargoimpl.NewCargoRepository(db, queries)
	locationRepository := locationimpl.NewLocationRepository(db, queries)
	distanceSensorStateRepository := distancesensorimpl.NewDistanceSensorStateRepository()
	appConnectionRepository := appconnectionimpl.NewAppConnectionRepository()

	// Initialize services
	batteryService := batteryimpl.NewService(validator, batteryStateRepository, batterySettingRepository)
	distanceSensorService := distancesensorimpl.NewService(validator, distanceSensorStateRepository)
	driveMotorService := drivemotorimpl.NewService(validator, driveMotorStateRepository)
	liftMotorService := liftmotorimpl.NewService(validator, liftMotorStateRepository)
	cargoService := cargoimpl.NewService(validator, cargoRepository)
	locationService := locationimpl.NewService(validator, locationRepository)
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
		appConnectionRepository,
	)
	appConnectionService := appconnectionimpl.NewService(appConnectionRepository)
	peripheralService := peripheralimpl.NewService()

	cleanup := func() error {
		return db.Close()
	}

	return &Application{
		Cfg:                   cfg,
		Log:                   log,
		Context:               ctx,
		BatteryService:        batteryService,
		DistanceSensorService: distanceSensorService,
		DriveMotorService:     driveMotorService,
		LiftMotorService:      liftMotorService,
		CargoService:          cargoService,
		LocationService:       locationService,
		ConfigService:         configService,
		SystemService:         systemService,
		DashboardDataService:  dashboardDataService,
		AppConnectionService:  appConnectionService,
		PeripheralService:     peripheralService,
	}, cleanup, nil
}
