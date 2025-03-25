-- name: CargoUpdate :one
UPDATE cargo
SET
	is_open = CASE WHEN @set_is_open IS NOT NULL THEN @is_open ELSE is_open END,
	qr_code = CASE WHEN @set_qr_code IS NOT NULL THEN @qr_code ELSE qr_code END,
	bottom_distance = CASE WHEN @set_bottom_distance IS NOT NULL THEN @bottom_distance ELSE bottom_distance END,
	updated_at = @updated_at
WHERE id = 1
RETURNING *;

-- name: CargoDoorMotorUpdate :one
UPDATE cargo_door_motor
SET
	direction = CASE WHEN @set_direction IS NOT NULL THEN @direction ELSE direction END,
	speed = CASE WHEN @set_speed IS NOT NULL THEN @speed ELSE speed END,
	is_running = CASE WHEN @set_is_running IS NOT NULL THEN @is_running ELSE is_running END,
	enabled = CASE WHEN @set_enabled IS NOT NULL THEN @enabled ELSE enabled END,
	updated_at = @updated_at
WHERE id = 1
RETURNING *;
