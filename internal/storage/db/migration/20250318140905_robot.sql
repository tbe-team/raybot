-- +goose Up
-- +goose StatementBegin
CREATE TABLE robot (
	id INTEGER PRIMARY KEY CHECK (id = 1)
);

CREATE TABLE battery (
	id INTEGER PRIMARY KEY CHECK (id = 1),
	current INTEGER NOT NULL,
	temp INTEGER NOT NULL,
	voltage INTEGER NOT NULL,
	cell_voltages TEXT NOT NULL,
	percent INTEGER NOT NULL,
	fault INTEGER NOT NULL,
	health INTEGER NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE TABLE battery_charge (
	id INTEGER PRIMARY KEY CHECK (id = 1),
	current_limit INTEGER NOT NULL,
	enabled INTEGER NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE TABLE battery_discharge (
	id INTEGER PRIMARY KEY CHECK (id = 1),
	current_limit INTEGER NOT NULL,
	enabled INTEGER NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE TABLE distance_sensor (
	id INTEGER PRIMARY KEY CHECK (id = 1),
	front_distance INTEGER NOT NULL,
	back_distance INTEGER NOT NULL,
	down_distance INTEGER NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE TABLE drive_motor (
	id INTEGER PRIMARY KEY CHECK (id = 1),
	direction INTEGER NOT NULL,
	speed INTEGER NOT NULL,
	is_running INTEGER NOT NULL,
	enabled INTEGER NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE TABLE lift_motor (
	id INTEGER PRIMARY KEY CHECK (id = 1),
	current_position INTEGER NOT NULL,
	target_position INTEGER NOT NULL,
	is_running INTEGER NOT NULL,
	enabled INTEGER NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE TABLE location (
	id INTEGER PRIMARY KEY CHECK (id = 1),
	current_location TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

INSERT INTO robot (id) VALUES (1);
INSERT INTO battery (id, current, temp, voltage, cell_voltages, percent, fault, health, updated_at) VALUES (1, 0, 0, 0, '[0, 0, 0, 0]', 0, 0, 0, '2025-01-01T00:00:00Z');
INSERT INTO battery_charge (id, current_limit, enabled, updated_at) VALUES (1, 0, true, '2025-01-01T00:00:00Z');
INSERT INTO battery_discharge (id, current_limit, enabled, updated_at) VALUES (1, 0, true, '2025-01-01T00:00:00Z');
INSERT INTO distance_sensor (id, front_distance, back_distance, down_distance, updated_at) VALUES (1, 0, 0, 0, '2025-01-01T00:00:00Z');
INSERT INTO drive_motor (id, direction, speed, is_running, enabled, updated_at) VALUES (1, 0, 0, false, false, '2025-01-01T00:00:00Z');
INSERT INTO lift_motor (id, current_position, target_position, is_running, enabled, updated_at) VALUES (1, 0, 0, false, false, '2025-01-01T00:00:00Z');
INSERT INTO location (id, current_location, updated_at) VALUES (1, '', '2025-01-01T00:00:00Z');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE location;
DROP TABLE lift_motor;
DROP TABLE drive_motor;
DROP TABLE distance_sensor;
DROP TABLE battery_discharge;
DROP TABLE battery_charge;
DROP TABLE battery;
DROP TABLE robot;
-- +goose StatementEnd
