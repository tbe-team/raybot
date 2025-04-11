-- name: CommandGetByID :one
SELECT * FROM commands
WHERE id = @id;

-- name: CommandGetProcessing :one
SELECT * FROM commands
WHERE status = 'PROCESSING'
LIMIT 1;

-- name: CommandGetNextExecutable :one
SELECT * FROM commands
WHERE
	status IN ('QUEUED', 'PROCESSING')
ORDER BY
	CASE status
		WHEN 'PROCESSING' THEN 0
		WHEN 'QUEUED' THEN 1
	END ASC,
	created_at ASC
LIMIT 1;

-- name: CommandProcessingExists :one
SELECT EXISTS (
	SELECT 1 FROM commands
	WHERE status = 'PROCESSING'
);

-- name: CommandCreate :one
INSERT INTO commands (
	type,
	status,
	source,
	inputs,
	error,
	created_at,
	updated_at,
	completed_at
)
VALUES (
	@type,
	@status,
	@source,
	@inputs,
	@error,
	@created_at,
	@updated_at,
	@completed_at
)
RETURNING id;

-- name: CommandUpdate :one
UPDATE commands
SET
	status = CASE WHEN @set_status = 1 THEN @status ELSE status END,
	error = CASE WHEN @set_error IS NOT NULL THEN @error ELSE error END,
	completed_at = CASE WHEN @set_completed_at IS NOT NULL THEN @completed_at ELSE completed_at END,
	updated_at = @updated_at
WHERE id = @id
RETURNING *;
