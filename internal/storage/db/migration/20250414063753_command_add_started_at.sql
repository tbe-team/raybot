-- +goose Up
-- +goose StatementBegin
ALTER TABLE commands ADD COLUMN started_at TEXT;

UPDATE commands SET type = 'STOP_MOVEMENT' WHERE type = 'STOP';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE commands_backup AS SELECT id, type, status, source, inputs, error, completed_at, created_at, updated_at FROM commands;

DROP TABLE commands;

CREATE TABLE commands (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	type TEXT NOT NULL,
	status TEXT NOT NULL,
	source TEXT NOT NULL,
	inputs TEXT NOT NULL DEFAULT '{}',
	error TEXT,
	completed_at TEXT,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

INSERT INTO commands SELECT * FROM commands_backup;

UPDATE commands SET type = 'STOP' WHERE type = 'STOP_MOVEMENT';

DROP TABLE commands_backup;
-- +goose StatementEnd
