-- name: CommandGetByStatusInProgress :one
SELECT * FROM commands WHERE status = 1 LIMIT 1;

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
VALUES (?, ?, ?, ?, ?, ?, ?, ?);
