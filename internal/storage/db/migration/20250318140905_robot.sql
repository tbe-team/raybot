-- +goose Up
-- +goose StatementBegin
CREATE TABLE robot (
	id INTEGER PRIMARY KEY CHECK (id = 1)
);

CREATE TABLE battery_charge_setting (
	id INTEGER PRIMARY KEY CHECK (id = 1),
	current_limit INTEGER NOT NULL,
	enabled INTEGER NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE TABLE battery_discharge_setting (
	id INTEGER PRIMARY KEY CHECK (id = 1),
	current_limit INTEGER NOT NULL,
	enabled INTEGER NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE TABLE location (
	id INTEGER PRIMARY KEY CHECK (id = 1),
	current_location TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

INSERT INTO robot (id) VALUES (1);
INSERT INTO battery_charge_setting (id, current_limit, enabled, updated_at) VALUES (1, 0, true, '2025-01-01T00:00:00Z');
INSERT INTO battery_discharge_setting (id, current_limit, enabled, updated_at) VALUES (1, 0, true, '2025-01-01T00:00:00Z');
INSERT INTO location (id, current_location, updated_at) VALUES (1, '', '2025-01-01T00:00:00Z');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE location;
DROP TABLE battery_discharge_setting;
DROP TABLE battery_charge_setting;
DROP TABLE robot;
-- +goose StatementEnd
