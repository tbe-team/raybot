-- +goose Up
-- +goose StatementBegin
CREATE TABLE cargo (
	id INTEGER PRIMARY KEY CHECK (id = 1),
	is_open INTEGER NOT NULL,
	qr_code TEXT NOT NULL,
	bottom_distance INTEGER NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE TABLE cargo_door_motor (
	id INTEGER PRIMARY KEY CHECK (id = 1),
	direction INTEGER NOT NULL,
	speed INTEGER NOT NULL,
	is_running INTEGER NOT NULL,
	enabled INTEGER NOT NULL,
	updated_at TEXT NOT NULL
);

INSERT INTO cargo (id, is_open, qr_code, bottom_distance, updated_at) VALUES (1, false, '', 0, '2025-01-01T00:00:00Z');
INSERT INTO cargo_door_motor (id, direction, speed, is_running, enabled, updated_at) VALUES (1, 0, 0, false, false, '2025-01-01T00:00:00Z');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE cargo_door_motor;
DROP TABLE cargo;
-- +goose StatementEnd
