-- name: DistanceSensorGet :one
SELECT * FROM distance_sensor;

-- name: DistanceSensorUpdate :exec
UPDATE distance_sensor
SET
	front_distance = ?,
	back_distance = ?,
	down_distance = ?,
	updated_at = ?
WHERE id = 1;
