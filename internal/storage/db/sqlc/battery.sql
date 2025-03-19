-- name: BatteryGet :one
SELECT * FROM battery;

-- name: BatteryUpdate :exec
UPDATE battery
SET
	current = ?,
	temp = ?,
	voltage = ?,
	cell_voltages = ?,
	percent = ?,
	fault = ?,
	health = ?,
	updated_at = ?
WHERE id = 1;

-- name: BatteryChargeGet :one
SELECT * FROM battery_charge;

-- name: BatteryChargeUpdate :exec
UPDATE battery_charge
SET
	current_limit = ?,
	enabled = ?,
	updated_at = ?
WHERE id = 1;

-- name: BatteryDischargeGet :one
SELECT * FROM battery_discharge;

-- name: BatteryDischargeUpdate :exec
UPDATE battery_discharge
SET
	current_limit = ?,
	enabled = ?,
	updated_at = ?
WHERE id = 1;
