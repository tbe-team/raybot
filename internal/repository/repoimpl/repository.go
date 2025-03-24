package repoimpl

import (
	"github.com/tbe-team/raybot/internal/repository"
	"github.com/tbe-team/raybot/internal/storage/db/sqlc"
)

type repo struct {
	robotStateRepo       repository.RobotStateRepository
	picSerialCommandRepo repository.PICSerialCommandRepository
	distanceSensorRepo   repository.DistanceSensorRepository
	batteryRepo          repository.BatteryRepository
	driveMotorRepo       repository.DriveMotorRepository
	liftMotorRepo        repository.LiftMotorRepository
	locationRepo         repository.LocationRepository
	commandRepo          repository.CommandRepository
}

func New() repository.Repository {
	queries := sqlc.New()
	return &repo{
		robotStateRepo:       NewRobotStateRepository(queries),
		picSerialCommandRepo: NewPICSerialCommandRepository(),
		distanceSensorRepo:   NewDistanceSensorRepository(queries),
		batteryRepo:          NewBatteryRepository(queries),
		driveMotorRepo:       NewDriveMotorRepository(queries),
		liftMotorRepo:        NewLiftMotorRepository(queries),
		locationRepo:         NewLocationRepository(queries),
		commandRepo:          NewCommandRepository(queries),
	}
}

func (r *repo) RobotState() repository.RobotStateRepository {
	return r.robotStateRepo
}

func (r *repo) PICSerialCommand() repository.PICSerialCommandRepository {
	return r.picSerialCommandRepo
}

func (r *repo) DistanceSensor() repository.DistanceSensorRepository {
	return r.distanceSensorRepo
}

func (r *repo) Battery() repository.BatteryRepository {
	return r.batteryRepo
}

func (r *repo) DriveMotor() repository.DriveMotorRepository {
	return r.driveMotorRepo
}

func (r *repo) LiftMotor() repository.LiftMotorRepository {
	return r.liftMotorRepo
}

func (r *repo) Location() repository.LocationRepository {
	return r.locationRepo
}

func (r *repo) Command() repository.CommandRepository {
	return r.commandRepo
}
