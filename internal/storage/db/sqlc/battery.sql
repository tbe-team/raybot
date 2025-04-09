-- name: BatteryChargeSettingGet :one
SELECT * FROM battery_charge_setting
WHERE id = 1;

-- name: BatteryDischargeSettingGet :one
SELECT * FROM battery_discharge_setting
WHERE id = 1;

-- name: BatteryChargeSettingUpdate :exec
UPDATE battery_charge_setting
SET
	current_limit = @current_limit,
	enabled = @enabled,
	updated_at = @updated_at
WHERE id = 1;

-- name: BatteryDischargeSettingUpdate :exec
UPDATE battery_discharge_setting
SET
	current_limit = @current_limit,
	enabled = @enabled,
	updated_at = @updated_at
WHERE id = 1;
