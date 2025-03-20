-- +goose Up
-- +goose StatementBegin
CREATE TABLE commands (
	id TEXT PRIMARY KEY,
	type INTEGER NOT NULL,
	status INTEGER NOT NULL,
	source INTEGER NOT NULL,
	inputs TEXT NOT NULL DEFAULT '{}',
	error TEXT,
	created_at DATETIME NOT NULL,
	completed_at DATETIME
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE commands;
-- +goose StatementEnd
