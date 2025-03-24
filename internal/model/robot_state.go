package model

type RobotState struct {
	Battery        Battery
	Charge         BatteryCharge
	Discharge      BatteryDischarge
	DistanceSensor DistanceSensor
	LiftMotor      LiftMotor
	DriveMotor     DriveMotor
	Location       Location
	Cargo          Cargo
	CargoDoorMotor CargoDoorMotor
}
