-- name: CargoGet :one
SELECT * FROM cargo
WHERE id = 1;

-- name: CargoDoorMotorGet :one
SELECT * FROM cargo_door_motor
WHERE id = 1;

-- name: CargoUpdateIsOpen :one
UPDATE cargo
SET
	is_open = @is_open,
	updated_at = @updated_at
WHERE id = 1
RETURNING *;

-- name: CargoUpdateQRCode :one
UPDATE cargo
SET
	qr_code = @qr_code,
	updated_at = @updated_at
WHERE id = 1
RETURNING *;

-- name: CargoUpdateBottomDistance :one
UPDATE cargo
SET
	bottom_distance = @bottom_distance,
	updated_at = @updated_at
WHERE id = 1
RETURNING *;

-- name: CargoDoorMotorUpdate :one
UPDATE cargo_door_motor
SET
	direction = @direction,
	speed = @speed,
	is_running = @is_running,
	enabled = @enabled,
	updated_at = @updated_at
WHERE id = 1
RETURNING *;
