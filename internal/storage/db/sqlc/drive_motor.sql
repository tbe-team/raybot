-- name: DriveMotorGet :one
SELECT * FROM drive_motor;

-- name: DriveMotorUpdate :exec
UPDATE drive_motor
SET
	direction = ?,
	speed = ?,
	is_running = ?,
	enabled = ?,
	updated_at = ?
WHERE id = 1;
