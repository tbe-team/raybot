-- +goose Up
-- +goose StatementBegin
ALTER TABLE commands
ADD COLUMN outputs TEXT NOT NULL DEFAULT '{}';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE commands_backup AS
SELECT
    id,
    type,
    status,
    source,
    inputs,
    error,
    completed_at,
    created_at,
    updated_at,
    started_at
FROM
    commands;

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
	updated_at TEXT NOT NULL,
	started_at TEXT
);

INSERT INTO commands SELECT * FROM commands_backup;

DROP TABLE commands_backup;
-- +goose StatementEnd
