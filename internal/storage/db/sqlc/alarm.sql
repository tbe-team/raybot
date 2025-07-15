-- name: AlarmListActive :many
SELECT
	*
FROM alarms
WHERE deactivated_at IS NULL
ORDER BY activated_at DESC
LIMIT @limit
OFFSET @offset;

-- name: AlarmCountActive :one
SELECT COUNT(*)
FROM alarms
WHERE deactivated_at IS NULL;

-- name: AlarmListDeactive :many
SELECT
	*
FROM alarms
WHERE deactivated_at IS NOT NULL
ORDER BY deactivated_at DESC
LIMIT @limit
OFFSET @offset;

-- name: AlarmCountDeactive :one
SELECT COUNT(*)
FROM alarms
WHERE deactivated_at IS NOT NULL;

-- name: AlarmCreate :one
INSERT INTO alarms (
	type,
	data,
	activated_at
)
VALUES (
	@type,
	@data,
	@activated_at
)
RETURNING *;

-- name: AlarmDeactivate :exec
UPDATE alarms
SET deactivated_at = @deactivated_at
WHERE id = @id;

-- name: AlarmDeactivateAllActivated :exec
UPDATE alarms
SET deactivated_at = @deactivated_at
WHERE deactivated_at IS NULL;

-- name: AlarmDeleteDeactivated :exec
DELETE FROM alarms
WHERE deactivated_at IS NOT NULL;

-- name: AlarmDeleteDeactivatedByThreshold :exec
DELETE FROM alarms
WHERE deactivated_at < @threshold;
