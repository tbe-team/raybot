-- name: RobotStateGet :one
SELECT
	json_object(
		'current', b.current,
		'temp', b.temp,
		'voltage', b.voltage,
		'cell_voltages', b.cell_voltages,
		'percent', b.percent,
		'fault', b.fault,
		'health', b.health,
		'updated_at', b.updated_at
	) AS battery,
	json_object(
		'current_limit', bc.current_limit,
		'enabled', bc.enabled,
		'updated_at', bc.updated_at
	) AS battery_charge,
	json_object(
		'current_limit', bd.current_limit,
		'enabled', bd.enabled,
		'updated_at', bd.updated_at
	) AS battery_discharge,
	json_object(
		'front_distance', ds.front_distance,
		'back_distance', ds.back_distance,
		'down_distance', ds.down_distance,
		'updated_at', ds.updated_at
	) AS distance_sensor,
	json_object(
		'direction', dm.direction,
		'speed', dm.speed,
		'is_running', dm.is_running,
		'enabled', dm.enabled,
		'updated_at', dm.updated_at
	) AS drive_motor,
	json_object(
		'current_position', lm.current_position,
		'target_position', lm.target_position,
		'is_running', lm.is_running,
		'enabled', lm.enabled,
		'updated_at', lm.updated_at
	) AS lift_motor,
	json_object(
		'current_location', l.current_location,
		'updated_at', l.updated_at
	) AS location,
	json_object(
		'is_open', c.is_open,
		'qr_code', c.qr_code,
		'bottom_distance', c.bottom_distance,
		'updated_at', c.updated_at
	) AS cargo,
	json_object(
		'direction', cmd.direction,
		'speed', cmd.speed,
		'is_running', cmd.is_running,
		'enabled', cmd.enabled,
		'updated_at', cmd.updated_at
	) AS cargo_door_motor
FROM robot r
LEFT JOIN battery b ON r.id = b.id
LEFT JOIN battery_charge bc ON r.id = bc.id
LEFT JOIN battery_discharge bd ON r.id = bd.id
LEFT JOIN distance_sensor ds ON r.id = ds.id
LEFT JOIN drive_motor dm ON r.id = dm.id
LEFT JOIN lift_motor lm ON r.id = lm.id
LEFT JOIN location l ON r.id = l.id
LEFT JOIN cargo c ON r.id = c.id
LEFT JOIN cargo_door_motor cmd ON r.id = cmd.id
WHERE r.id = 1;
