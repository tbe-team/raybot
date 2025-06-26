-- +goose Up
-- +goose StatementBegin
ALTER TABLE commands
ADD COLUMN request_id TEXT;

CREATE UNIQUE INDEX idx_commands_request_id ON commands(request_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_commands_request_id;

ALTER TABLE commands
DROP COLUMN request_id;
-- +goose StatementEnd
