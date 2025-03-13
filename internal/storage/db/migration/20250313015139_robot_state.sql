-- +goose Up
-- +goose StatementBegin
CREATE TABLE robot_state (
	battery_state TEXT NOT NULL,
	charge_state TEXT NOT NULL,
	discharge_state TEXT NOT NULL,
	distance_sensor_state TEXT NOT NULL,
	lift_motor_state TEXT NOT NULL,
	drive_motor_state TEXT NOT NULL
);

-- Insert default values
INSERT INTO robot_state (
	battery_state,
	charge_state,
	discharge_state,
	distance_sensor_state,
	lift_motor_state,
	drive_motor_state
) VALUES (
	'{"current": 0, "temp": 0, "voltage": 0, "cell_voltages": [0, 0, 0, 0], "percent": 0, "fault": 0, "health": 0, "updated_at": "0000-00-00 00:00:00"}',
	'{"current_limit": 0, "enabled": true, "updated_at": "0000-00-00 00:00:00"}',
	'{"current_limit": 0, "enabled": true, "updated_at": "0000-00-00 00:00:00"}',
	'{"front_distance": 0, "back_distance": 0, "down_distance": 0, "updated_at": "0000-00-00 00:00:00"}',
	'{"current_position": 0, "target_position": 0, "is_running": false, "enabled": true, "updated_at": "0000-00-00 00:00:00"}',
	'{"direction": 0, "speed": 0, "is_running": false, "enabled": true, "updated_at": "0000-00-00 00:00:00"}'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE robot_config;
-- +goose StatementEnd
