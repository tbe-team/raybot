package repository

type Repository interface {
	RobotState() RobotStateRepository
	PICSerialCommand() PICSerialCommandRepository
	Battery() BatteryRepository
	DistanceSensor() DistanceSensorRepository
	DriveMotor() DriveMotorRepository
	LiftMotor() LiftMotorRepository
	Location() LocationRepository
	Command() CommandRepository
}
