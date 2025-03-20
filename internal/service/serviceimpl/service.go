package serviceimpl

import (
	"log/slog"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/controller/picserial/serial"
	"github.com/tbe-team/raybot/internal/repository"
	"github.com/tbe-team/raybot/internal/service"
	"github.com/tbe-team/raybot/internal/storage/db"
	"github.com/tbe-team/raybot/pkg/validator"
)

type serviceImpl struct {
	robotStateService *RobotStateService
	systemService     *SystemService
	picService        *PICService
	locationService   *LocationService
	commandService    *CommandService
}

func New(
	cfgManager config.Manager,
	picSerialClient serial.Client,
	repo repository.Repository,
	dbProvider db.Provider,
	validator validator.Validator,
	log *slog.Logger,
) service.Service {
	robotStateService := NewRobotStateService(repo.RobotState(), dbProvider)
	systemService := NewSystemService(cfgManager)
	picService := NewPICService(
		repo.RobotState(),
		repo.PICSerialCommand(),
		repo.Battery(),
		repo.DistanceSensor(),
		repo.LiftMotor(),
		repo.DriveMotor(),
		repo.Location(),
		picSerialClient,
		dbProvider,
		validator,
	)
	locationService := NewLocationService(repo.Location(), dbProvider)
	commandService := NewCommandService(
		repo.Command(),
		repo.Location(),
		picService,
		dbProvider,
		validator,
		log,
	)
	return &serviceImpl{
		robotStateService: robotStateService,
		systemService:     systemService,
		picService:        picService,
		locationService:   locationService,
		commandService:    commandService,
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
