-- name: LiftMotorGet :one
SELECT * FROM lift_motor;

-- name: LiftMotorUpdate :exec
UPDATE lift_motor
SET
	current_position = ?,
	target_position = ?,
	is_running = ?,
	enabled = ?,
	updated_at = ?
WHERE id = 1;
