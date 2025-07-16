-- +goose Up
-- +goose StatementBegin
CREATE TABLE alarms (
	id INTEGER PRIMARY KEY,
	type TEXT NOT NULL,
	data TEXT NOT NULL,
	activated_at TEXT NOT NULL,
	deactivated_at TEXT
);

CREATE INDEX idx_alarms_deactivated_at ON alarms(deactivated_at DESC);
CREATE UNIQUE INDEX idx_unique_active_alarm_type ON alarms(type) WHERE deactivated_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_unique_active_alarm_type;
DROP INDEX idx_alarms_deactivated_at;
DROP TABLE alarms;
-- +goose StatementEnd
