-- +goose Up
-- +goose StatementBegin
CREATE TABLE robot_state (
	battery_state TEXT NOT NULL,
	charge_state TEXT NOT NULL,
	discharge_state TEXT NOT NULL,
	distance_sensor_state TEXT NOT NULL,
	lift_motor_state TEXT NOT NULL,
	drive_motor_state TEXT NOT NULL,
	location_state TEXT NOT NULL
);

-- Insert default values
INSERT INTO robot_state (
	battery_state,
	charge_state,
	discharge_state,
	distance_sensor_state,
	lift_motor_state,
	drive_motor_state,
	location_state
) VALUES (
	'{"Current": 0, "Temp": 0, "Voltage": 0, "CellVoltages": [0, 0, 0, 0], "Percent": 0, "Fault": 0, "Health": 0, "UpdatedAt": "2023-01-01T00:00:00Z"}',
	'{"CurrentLimit": 0, "Enabled": true, "UpdatedAt": "2023-01-01T00:00:00Z"}',
	'{"CurrentLimit": 0, "Enabled": true, "UpdatedAt": "2023-01-01T00:00:00Z"}',
	'{"FrontDistance": 0, "BackDistance": 0, "DownDistance": 0, "UpdatedAt": "2023-01-01T00:00:00Z"}',
	'{"CurrentPosition": 0, "TargetPosition": 0, "IsRunning": false, "Enabled": true, "UpdatedAt": "2023-01-01T00:00:00Z"}',
	'{"Direction": 0, "Speed": 0, "IsRunning": false, "Enabled": true, "UpdatedAt": "2023-01-01T00:00:00Z"}',
	'{"CurrentLocation": "", "UpdatedAt": "2023-01-01T00:00:00Z"}'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE robot_state;
-- +goose StatementEnd
