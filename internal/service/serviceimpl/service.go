package serviceimpl

import (
	"log/slog"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/controller/espserial"
	"github.com/tbe-team/raybot/internal/controller/picserial/serial"
	"github.com/tbe-team/raybot/internal/pubsub"
	"github.com/tbe-team/raybot/internal/repository"
	"github.com/tbe-team/raybot/internal/service"
	"github.com/tbe-team/raybot/internal/service/serviceimpl/cargocontrol"
	"github.com/tbe-team/raybot/internal/service/serviceimpl/command"
	"github.com/tbe-team/raybot/internal/service/serviceimpl/location"
	"github.com/tbe-team/raybot/internal/service/serviceimpl/pic"
	robotstate "github.com/tbe-team/raybot/internal/service/serviceimpl/robot_state"
	"github.com/tbe-team/raybot/internal/service/serviceimpl/system"
	"github.com/tbe-team/raybot/internal/storage/db"
	"github.com/tbe-team/raybot/pkg/validator"
)

type serviceImpl struct {
	robotStateService   *robotstate.Service
	systemService       *system.Service
	picService          *pic.Service
	locationService     *location.Service
	commandService      *command.Service
	cargoControlService *cargocontrol.Service
}

func New(
	cfgManager config.Manager,
	picSerialClient serial.Client,
	espSerialClient espserial.Client,
	repo repository.Repository,
	pubSub pubsub.PubSub,
	dbProvider db.Provider,
	validator validator.Validator,
	log *slog.Logger,
) service.Service {
	robotStateService := robotstate.NewService(repo.RobotState(), dbProvider)
	systemService := system.NewService(cfgManager)
	picService := pic.NewService(
		repo.RobotState(),
		repo.PICSerialCommand(),
		repo.Battery(),
		repo.DistanceSensor(),
		repo.LiftMotor(),
		repo.DriveMotor(),
		repo.Location(),
		dbProvider,
		picSerialClient,
		validator,
	)
	locationService := location.NewService(repo.Location(), pubSub, dbProvider)
	cargoControlService := cargocontrol.NewService(
		repo.Cargo(),
		repo.ESPSerialCommand(),
		espSerialClient,
		dbProvider,
		validator,
	)
	commandService := command.NewService(
		repo.Command(),
		picService,
		cargoControlService,
		dbProvider,
		pubSub,
		pubSub,
		validator,
		log,
	)
	return &serviceImpl{
		robotStateService:   robotStateService,
		systemService:       systemService,
		picService:          picService,
		locationService:     locationService,
		commandService:      commandService,
		cargoControlService: cargoControlService,
	}
}

func (s serviceImpl) RobotStateService() service.RobotStateService {
	return s.robotStateService
}

func (s serviceImpl) SystemService() service.SystemService {
	return s.systemService
}

func (s serviceImpl) PICService() service.PICService {
	return s.picService
}

func (s serviceImpl) LocationService() service.LocationService {
	return s.locationService
}

func (s serviceImpl) CommandService() service.CommandService {
	return s.commandService
}

func (s serviceImpl) CargoControlService() service.CargoControlService {
	return s.cargoControlService
}
