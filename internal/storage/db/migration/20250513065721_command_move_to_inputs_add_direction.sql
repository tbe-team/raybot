-- +goose Up
-- +goose StatementBegin
UPDATE commands
SET inputs = json_set(inputs, '$.direction', 'FORWARD')
WHERE type = 'MOVE_TO';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE commands
SET inputs = json_remove(inputs, '$.direction')
WHERE type = 'MOVE_TO';
-- +goose StatementEnd
