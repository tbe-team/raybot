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
	drive_motor_state = ?,
	location_state = ?;
