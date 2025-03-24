-- name: CommandGetByStatusInProgress :one
SELECT * FROM commands WHERE status = 0 LIMIT 1;

-- name: CommandGetByID :one
SELECT * FROM commands WHERE id = @id;

-- name: CommandCreate :exec
INSERT INTO commands (
	id,
	type,
	status,
	source,
	inputs,
	error,
	created_at,
	completed_at
)
VALUES (@id, @type, @status, @source, @inputs, @error, @created_at, @completed_at);

-- name: CommandUpdate :one
UPDATE commands
SET
	status = CASE WHEN @set_status = 1 THEN @status ELSE status END,
	error = CASE WHEN @set_error IS NOT NULL THEN @error ELSE error END,
	completed_at = CASE WHEN @set_completed_at IS NOT NULL THEN @completed_at ELSE completed_at END
WHERE id = @id
RETURNING *;
