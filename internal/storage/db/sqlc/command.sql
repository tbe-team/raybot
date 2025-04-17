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
	started_at,
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
	@started_at,
	@created_at,
	@updated_at,
	@completed_at
)
RETURNING id;

-- name: CommandUpdate :one
UPDATE commands
SET
	status = CASE WHEN @set_status = 1 THEN @status ELSE status END,
	error = CASE WHEN @set_error = 1 THEN @error ELSE error END,
	started_at = CASE WHEN @set_started_at = 1 THEN @started_at ELSE started_at END,
	completed_at = CASE WHEN @set_completed_at = 1 THEN @completed_at ELSE completed_at END,
	updated_at = @updated_at
WHERE id = @id
RETURNING *;

-- name: CommandDeleteByIDAndNotProcessing :execrows
DELETE FROM commands
WHERE id = @id
AND status != 'PROCESSING';
