package serviceimpl

import (
	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/controller/picserial/serial"
	"github.com/tbe-team/raybot/internal/repository"
	"github.com/tbe-team/raybot/internal/service"
	"github.com/tbe-team/raybot/internal/storage/db"
	"github.com/tbe-team/raybot/pkg/validator"
)

type serviceImpl struct {
	robotService  *RobotService
	systemService *SystemService
	picService    *PICService
}

func New(
	cfgManager config.Manager,
	picSerialClient serial.Client,
	repo repository.Repository,
	dbProvider db.Provider,
	validator validator.Validator,
) service.Service {
	return &serviceImpl{
		robotService:  NewRobotService(repo.RobotState(), dbProvider, validator),
		systemService: NewSystemService(cfgManager),
		picService:    NewPICService(repo.RobotState(), repo.PICSerialCommand(), picSerialClient, dbProvider, validator),
	}
}

func (s serviceImpl) RobotService() service.RobotService {
	return s.robotService
}

func (s serviceImpl) SystemService() service.SystemService {
	return s.systemService
}

func (s serviceImpl) PICService() service.PICService {
	return s.picService
}
