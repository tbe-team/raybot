-- +goose Up
-- +goose StatementBegin
ALTER TABLE cargo
ADD COLUMN has_item INTEGER NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE cargo DROP COLUMN has_item;
-- +goose StatementEnd
