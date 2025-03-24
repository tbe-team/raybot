-- name: LocationUpdate :exec
UPDATE location
SET
	current_location = ?,
	updated_at = ?
WHERE id = 1;

-- name: LocationGetCurrent :one
SELECT * FROM location
WHERE id = 1;
