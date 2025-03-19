-- name: LocationUpdate :exec
UPDATE location
SET
	current_location = ?,
	updated_at = ?
WHERE id = 1;
