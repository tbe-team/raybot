-- +goose Up
-- +goose StatementBegin
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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE commands;
-- +goose StatementEnd
