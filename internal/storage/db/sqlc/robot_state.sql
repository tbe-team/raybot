-- name: RobotStateGet :one
SELECT * FROM robot_state LIMIT 1;

-- name: RobotStateUpdate :exec
UPDATE robot_state
SET
	battery_state = ?,
	charge_state = ?,
	discharge_state = ?,
	distance_sensor_state = ?,
	lift_motor_state = ?,
	drive_motor_state = ?;

-- name: BatteryStateUpdate :exec
UPDATE robot_state SET battery_state = ?;

-- name: ChargeStateUpdate :exec
UPDATE robot_state SET charge_state = ?;

-- name: DischargeStateUpdate :exec
UPDATE robot_state SET discharge_state = ?;

-- name: DistanceSensorStateUpdate :exec
UPDATE robot_state SET distance_sensor_state = ?;

-- name: LiftMotorStateUpdate :exec
UPDATE robot_state SET lift_motor_state = ?;

-- name: DriveMotorStateUpdate :exec
UPDATE robot_state SET drive_motor_state = ?;
